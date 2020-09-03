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

	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1a1 "github.com/jmckind/conjur-operator/api/v1alpha1"
)

// ConjurReconciler reconciles a Conjur object
type ConjurReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:namespace=conjur,groups=apps,resources=deployments/finalizers,verbs=update
// +kubebuilder:rbac:namespace=conjur,groups=apps,resources=deployments;statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:namespace=conjur,groups=core,resources=configmaps;events;persistentvolumeclaims;secrets;serviceaccounts;services;services/finalizers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:namespace=conjur,groups=oss.cyberark.com,resources=conjurs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:namespace=conjur,groups=oss.cyberark.com,resources=conjurs/status,verbs=get;update;patch
// +kubebuilder:rbac:namespace=conjur,groups=route.openshift.io,resources=routes;routes/custom-host,verbs=get;list;watch;create;update;patch;delete

// Reconcile the actual vs desired state for a given Conjur custom resource.
func (r *ConjurReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("conjur", req.NamespacedName)

	cr := &v1a1.Conjur{}
	if err := r.Get(context.TODO(), req.NamespacedName, cr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.reconcileConjur(cr); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager configures the owned resources for the Conjur custom resource.
func (r *ConjurReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1a1.Conjur{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Owns(&routev1.Route{}).
		Complete(r)
}

// reconcileConjur will ensure that all of the required resources are
// present for the given Conjur instance.
func (r *ConjurReconciler) reconcileConjur(cr *v1a1.Conjur) error {
	if err := r.reconcileStatus(cr); err != nil {
		return err
	}

	if err := r.reconcileTLS(cr); err != nil {
		return err
	}

	if err := r.reconcileDatabase(cr); err != nil {
		return err
	}

	if err := r.reconcileServer(cr); err != nil {
		return err
	}
	return nil
}
