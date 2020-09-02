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
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	v1a1 "github.com/jmckind/conjur-operator/api/v1alpha1"
	"github.com/jmckind/conjur-operator/common"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// reconcileCASecret ensures the CA Secret is created for the given Conjur resource.
func (r *ConjurReconciler) reconcileCASecret(cr *v1a1.Conjur) error {
	secret := newSecret(cr.Namespace, nameWithSuffix(cr.Name, "conjur-ca"))
	if r.isResourceFound(cr.Namespace, secret.Name, secret) {
		return nil // Secret found, do nothing
	}

	key, err := common.NewPrivateKey()
	if err != nil {
		return err
	}

	cert, err := common.NewSelfSignedCACertificate(key)
	if err != nil {
		return err
	}

	secret.Type = corev1.SecretTypeTLS
	secret.Data = map[string][]byte{
		corev1.TLSCertKey:       common.EncodeCertificatePEM(cert),
		corev1.TLSPrivateKeyKey: common.EncodePrivateKeyPEM(key),
	}

	ctrl.SetControllerReference(cr, secret, r.Scheme)
	return r.Client.Create(context.TODO(), secret)
}

// reconcileTLS will ensure that the TLS resources are present and owned by
// the given Conjur resource.
func (r *ConjurReconciler) reconcileTLS(cr *v1a1.Conjur) error {
	if err := r.reconcileCASecret(cr); err != nil {
		return err
	}
	return nil
}

// newCertificateSecret creates a new secret using the given Conjur name, suffix and CA certificate/key.
func newCertificateSecret(cr *v1a1.Conjur, suffix string, caCert *x509.Certificate, caKey *rsa.PrivateKey) (*corev1.Secret, error) {
	secret := newSecret(cr.Namespace, nameWithSuffix(cr.Name, suffix))

	key, err := common.NewPrivateKey()
	if err != nil {
		return nil, err
	}

	cfg := &common.CertConfig{
		CertName:     secret.Name,
		CertType:     common.ClientAndServingCert,
		CommonName:   secret.Name,
		Organization: []string{cr.Namespace},
	}

	dnsNames := []string{
		cr.Name,
		fmt.Sprintf("%s.%s.svc.cluster.local", cr.Name, cr.Namespace),
	}

	cert, err := common.NewSignedCertificate(cfg, dnsNames, key, caCert, caKey)
	if err != nil {
		return nil, err
	}

	secret.Type = corev1.SecretTypeTLS
	secret.Data = map[string][]byte{
		corev1.TLSCertKey:       common.EncodeCertificatePEM(cert),
		corev1.TLSPrivateKeyKey: common.EncodePrivateKeyPEM(key),
	}

	return secret, nil
}
