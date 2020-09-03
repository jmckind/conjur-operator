// Copyright 2020 The Conjur Operator Developers

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"
	"io/ioutil"

	v1a1 "github.com/jmckind/conjur-operator/api/v1alpha1"
	"github.com/jmckind/conjur-operator/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

// getServerCount will return the replica count for the Conjur server
// Deployment.
func getServerCount(cr *v1a1.Conjur) *int32 {
	count := common.ConjurDefaultServerCount
	if cr.Spec.Server.Count != nil {
		count = *cr.Spec.Server.Count
	}
	return &count
}

// getServerImage will return the container image for the Conjur server.
func getServerImage(cr *v1a1.Conjur) string {
	img := common.ConjurDefaultServerImage
	tag := common.ConjurDefaultServerVersion
	return common.CombineImageTag(img, tag)
}

// getServerLabels will return the labels to use for the Conjur server.
func getServerLabels(cr *v1a1.Conjur) map[string]string {
	return map[string]string{
		common.ConjurKeyApp:          common.ConjurAppName,
		common.ConjurKeyAppComponent: "service",
		common.ConjurKeyAppName:      nameWithSuffix(cr.Name, common.ConjurAppName),
	}
}

// reconcileServer will ensure that all server resources are present owned
// by the given Conjur resource.
func (r *ConjurReconciler) reconcileServer(cr *v1a1.Conjur) error {
	if err := r.reconcileServerTLS(cr); err != nil {
		return err
	}

	if err := r.reconcileAuthenticators(cr); err != nil {
		return err
	}

	if err := r.reconcileNginxConfig(cr); err != nil {
		return err
	}

	if err := r.reconcileServerService(cr); err != nil {
		return err
	}
	return nil
}

// reconcileAuthenticators will ensure that the Conjur server authenticators
// Secret is present and owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileAuthenticators(cr *v1a1.Conjur) error {
	secret := newSecret(cr.Namespace, nameWithSuffix(cr.Name, "conjur-authenticators"))
	if r.isResourceFound(cr.Namespace, secret.Name, secret) {
		return nil
	}

	secret.Labels = getServerLabels(cr)
	secret.Data = map[string][]byte{
		"key": []byte(common.ConjurDefaultAuthenticators),
	}

	ctrl.SetControllerReference(cr, secret, r.Scheme)
	return r.Create(context.TODO(), secret)
}

// reconcileNginxConfig will ensure that the nginx ConfigMap is present and
// owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileNginxConfig(cr *v1a1.Conjur) error {
	cm := newConfigMap(cr.Namespace, nameWithSuffix(cr.Name, "conjur-nginx-configmap"))
	if r.isResourceFound(cr.Namespace, cm.Name, cm) {
		return nil
	}

	conjurSite, err := ioutil.ReadFile("config/files/conjur.conf")
	if err != nil {
		return err
	}

	dhParams, err := ioutil.ReadFile("config/files/dhparams.pem")
	if err != nil {
		return err
	}

	mimeTypes, err := ioutil.ReadFile("config/files/mime.types")
	if err != nil {
		return err
	}

	nginxCfg, err := ioutil.ReadFile("config/files/nginx.conf")
	if err != nil {
		return err
	}

	cm.Labels = getServerLabels(cr)
	cm.Data = map[string]string{
		"conjur_site": string(conjurSite),
		"dhparams":    string(dhParams),
		"mime_types":  string(mimeTypes),
		"nginx_conf":  string(nginxCfg),
	}

	ctrl.SetControllerReference(cr, cm, r.Scheme)
	return r.Create(context.TODO(), cm)
}

// reconcileServerDeployment will ensure that the Conjur server Deployment is
// present and owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileServerDeployment(cr *v1a1.Conjur) error {
	deploy := newDeployment(cr.Namespace, nameWithSuffix(cr.Name, "conjur-oss"))
	if r.isResourceFound(cr.Namespace, deploy.Name, deploy) {
		return nil
	}

	deploy.Labels = getServerLabels(cr)

	ctrl.SetControllerReference(cr, deploy, r.Scheme)
	return r.Create(context.TODO(), deploy)
}

// reconcileDatabaseTLS will ensure that the database TLS Secret is present and
// owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileServerTLS(cr *v1a1.Conjur) error {
	secret := newSecret(cr.Namespace, nameWithSuffix(cr.Name, "conjur-ssl-cert"))
	if r.isResourceFound(cr.Namespace, secret.Name, secret) {
		return nil // Database TLS Secret found, do nothing
	}

	caSecret := newSecret(cr.Namespace, nameWithSuffix(cr.Name, "conjur-ssl-ca-cert"))
	if !r.isResourceFound(cr.Namespace, caSecret.Name, caSecret) {
		r.Log.Info("ca secret not found, waiting to create database tls secret")
		return nil
	}

	caCert, err := common.ParsePEMEncodedCert(caSecret.Data[corev1.TLSCertKey])
	if err != nil {
		return err
	}

	caKey, err := common.ParsePEMEncodedPrivateKey(caSecret.Data[corev1.TLSPrivateKeyKey])
	if err != nil {
		return err
	}

	secret, err = newCertificateSecret(cr, "conjur-ssl-cert", caCert, caKey)
	if err != nil {
		return err
	}

	secret.Labels = getServerLabels(cr)
	ctrl.SetControllerReference(cr, secret, r.Scheme)
	return r.Client.Create(context.TODO(), secret)
}

// reconcileServerService will ensure that the Conjur server Service is present
// and owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileServerService(cr *v1a1.Conjur) error {
	svc := newService(cr.Namespace, nameWithSuffix(cr.Name, "conjur-oss"))
	if r.isResourceFound(cr.Namespace, svc.Name, svc) {
		return nil // Service found, do nothing...
	}

	svc.Labels = getServerLabels(cr)
	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "https",
			Port:       443,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(common.ConjurDefaultProxySecurePort),
		},
	}

	svc.Spec.Selector = getServerLabels(cr)
	svc.Spec.Type = corev1.ServiceTypeClusterIP

	ctrl.SetControllerReference(cr, svc, r.Scheme)
	return r.Create(context.TODO(), svc)
}
