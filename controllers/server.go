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
	ctrl "sigs.k8s.io/controller-runtime"
)

// reconcileServer will ensure that all server resources are present owned
// by the given Conjur resource.
func (r *ConjurReconciler) reconcileServer(cr *v1a1.Conjur) error {
	if err := r.reconcileNginxConfig(cr); err != nil {
		return err
	}
	return nil
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

	cm.Data = map[string]string{
		"conjur_site": string(conjurSite),
		"dhparams":    string(dhParams),
		"mime_types":  string(mimeTypes),
		"nginx_conf":  string(nginxCfg),
	}

	ctrl.SetControllerReference(cr, cm, r.Scheme)
	return r.Create(context.TODO(), cm)
}
