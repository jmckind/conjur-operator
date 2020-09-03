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
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	v1a1 "github.com/jmckind/conjur-operator/api/v1alpha1"
	"github.com/jmckind/conjur-operator/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

// getDatabaseCount will return the replica count for the Conjur database
// StatefulSet.
func getDatabaseCount(cr *v1a1.Conjur) *int32 {
	count := common.ConjurDefaultDatabaseCount
	if cr.Spec.Database.Count != nil {
		count = *cr.Spec.Database.Count
	}
	return &count
}

// getDatabaseImage will return the container image for the Conjur database.
func getDatabaseImage(cr *v1a1.Conjur) string {
	img := common.ConjurDefaultDatabaseImage
	tag := common.ConjurDefaultDatabaseVersion
	return common.CombineImageTag(img, tag)
}

// getDatabaseLabels will return the labels to use for the database component.
func getDatabaseLabels(cr *v1a1.Conjur) map[string]string {
	return map[string]string{
		common.ConjurKeyApp:          common.ConjurAppName,
		common.ConjurKeyAppComponent: common.ConjurComponentPostgres,
		common.ConjurKeyAppName:      nameWithSuffix(cr.Name, common.ConjurAppName),
	}
}

// reconcileDatabase will ensure that all database resources are present owned
// by the given Conjur resource.
func (r *ConjurReconciler) reconcileDatabase(cr *v1a1.Conjur) error {
	if err := r.reconcileDataKey(cr); err != nil {
		return err
	}

	if err := r.reconcileDatabaseTLS(cr); err != nil {
		return err
	}

	if err := r.reconcileDatabasePassword(cr); err != nil {
		return err
	}

	if err := r.reconcileDatabaseURL(cr); err != nil {
		return err
	}

	if err := r.reconcileDatabasePVC(cr); err != nil {
		return err
	}

	if err := r.reconcileDatabaseStatefulSet(cr); err != nil {
		return err
	}

	if err := r.reconcileDatabaseService(cr); err != nil {
		return err
	}
	return nil
}

// reconcileDatabasePassword will ensure that the database password Secret is
// present and owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileDatabasePassword(cr *v1a1.Conjur) error {
	secret := newSecret(cr.Namespace, nameWithSuffix(cr.Name, "conjur-database-password"))
	if r.isResourceFound(cr.Namespace, secret.Name, secret) {
		return nil
	}

	if len(cr.Spec.Database.URL) > 0 {
		return nil // Database URL is provided, no need to generate database root password
	}

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}

	secret.Labels = getDatabaseLabels(cr)
	secret.Data = map[string][]byte{
		"key": []byte(hex.EncodeToString(key)),
	}

	ctrl.SetControllerReference(cr, secret, r.Scheme)
	return r.Create(context.TODO(), secret)
}

// reconcileDatabasePVC will ensure that the database PersistentVolumeClaim is
// present and owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileDatabasePVC(cr *v1a1.Conjur) error {
	pvc := newPersistentVolumeClaim(cr.Namespace, nameWithSuffix(cr.Name, "conjur-postgres"))
	if r.isResourceFound(cr.Namespace, pvc.Name, pvc) {
		return nil // PVC found, do nothing...
	}

	pvc.Labels = getDatabaseLabels(cr)
	pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
	pvc.Spec.Resources = corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			"storage": resource.MustParse(common.ConjurDefaultPVCCapicity),
		},
	}

	ctrl.SetControllerReference(cr, pvc, r.Scheme)
	return r.Client.Create(context.TODO(), pvc)
}

// reconcileDatabaseService will ensure that the database Service is present
// and owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileDatabaseService(cr *v1a1.Conjur) error {
	svc := newService(cr.Namespace, nameWithSuffix(cr.Name, "conjur-postgres"))
	if r.isResourceFound(cr.Namespace, svc.Name, svc) {
		return nil // Service found, do nothing...
	}

	svc.Labels = getDatabaseLabels(cr)
	svc.Spec.Ports = []corev1.ServicePort{
		{
			Name:       common.ConjurComponentPostgres,
			Port:       common.ConjurDefaultDatabasePort,
			Protocol:   corev1.ProtocolTCP,
			TargetPort: intstr.FromInt(common.ConjurDefaultDatabasePort),
		},
	}

	svc.Spec.Selector = getDatabaseLabels(cr)
	svc.Spec.Type = corev1.ServiceTypeClusterIP

	ctrl.SetControllerReference(cr, svc, r.Scheme)
	return r.Create(context.TODO(), svc)
}

// reconcileDatabaseStatefulSet will ensure that the database StatefulSet is
// present and owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileDatabaseStatefulSet(cr *v1a1.Conjur) error {
	ss := newStatefulSet(cr.Namespace, nameWithSuffix(cr.Name, "conjur-postgres"))
	if r.isResourceFound(cr.Namespace, ss.Name, ss) {
		return nil // StatefulSet found, do nothing...
	}

	ss.Labels = getDatabaseLabels(cr)
	ss.Spec.Replicas = getDatabaseCount(cr)

	ss.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: getDatabaseLabels(cr),
	}

	ss.Spec.ServiceName = nameWithSuffix(cr.Name, "conjur-postgres")
	ss.Spec.Template.Labels = getDatabaseLabels(cr)

	ss.Spec.Template.Spec.Containers = []corev1.Container{
		{
			Args: []string{
				"/usr/bin/run-postgresql",
				"-c",
				"ssl=on",
				"-c",
				"ssl_cert_file=/etc/certs/tls.crt",
				"-c",
				"ssl_key_file=/etc/certs/tls.key",
			},
			Env: []corev1.EnvVar{
				{
					Name: "POSTGRESQL_ADMIN_PASSWORD",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: nameWithSuffix(cr.Name, "conjur-database-password"),
							},
							Key: "key",
						},
					},
				},
			},
			Image: getDatabaseImage(cr),
			Name:  common.ConjurComponentPostgres,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "postgres-data",
					MountPath: "/var/lib/postgresql/data",
				}, {
					Name:      "ssl-certs",
					MountPath: "/etc/certs",
					ReadOnly:  true,
				},
			},
		},
	}

	tlsMode := common.ConjurDefaultTLSFileMode
	ss.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: "postgres-data",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: nameWithSuffix(cr.Name, "conjur-postgres"),
				},
			},
		}, {
			Name: "ssl-certs",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					DefaultMode: &tlsMode,
					SecretName:  nameWithSuffix(cr.Name, "conjur-database-ssl"),
				},
			},
		},
	}

	ctrl.SetControllerReference(cr, ss, r.Scheme)
	return r.Create(context.TODO(), ss)
}

// reconcileDatabaseTLS will ensure that the database TLS Secret is present and
// owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileDatabaseTLS(cr *v1a1.Conjur) error {
	secret := newSecret(cr.Namespace, nameWithSuffix(cr.Name, "conjur-database-ssl"))
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

	secret, err = newCertificateSecret(cr, "conjur-database-ssl", caCert, caKey)
	if err != nil {
		return err
	}

	secret.Labels = getDatabaseLabels(cr)
	ctrl.SetControllerReference(cr, secret, r.Scheme)
	return r.Client.Create(context.TODO(), secret)
}

// reconcileDatabaseURL will ensure that the Secret for the database URL is
// present and owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileDatabaseURL(cr *v1a1.Conjur) error {
	secret := newSecret(cr.Namespace, nameWithSuffix(cr.Name, "conjur-database-url"))
	if r.isResourceFound(cr.Namespace, secret.Name, secret) {
		return nil
	}

	dbURL := cr.Spec.Database.URL
	if len(dbURL) <= 0 {
		dbPwdSecret := newSecret(cr.Namespace, nameWithSuffix(cr.Name, "conjur-database-password"))
		if !r.isResourceFound(cr.Namespace, dbPwdSecret.Name, dbPwdSecret) {
			r.Log.Info("database password secret not found, waiting to create database url secret")
			return nil
		}
		dbURL = fmt.Sprintf(
			"postgres://postgres:%s@%s/postgres?sslmode=require",
			dbPwdSecret.Data["key"],
			nameWithSuffix(cr.Name, "conjur-postgres"),
		)
	}

	secret.Labels = getDatabaseLabels(cr)
	secret.Data = map[string][]byte{
		"key": []byte(dbURL),
	}

	ctrl.SetControllerReference(cr, secret, r.Scheme)
	return r.Create(context.TODO(), secret)
}

// reconcileDataKey will ensure that the Secret for the data key is present and
// owned by the given Conjur resource.
func (r *ConjurReconciler) reconcileDataKey(cr *v1a1.Conjur) error {
	secret := newSecret(cr.Namespace, nameWithSuffix(cr.Name, "conjur-data-key"))
	if r.isResourceFound(cr.Namespace, secret.Name, secret) {
		return nil
	}

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}

	secret.Data = map[string][]byte{
		"key": []byte(base64.StdEncoding.EncodeToString(key)),
	}

	ctrl.SetControllerReference(cr, secret, r.Scheme)
	return r.Create(context.TODO(), secret)
}
