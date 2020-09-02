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
	"fmt"

	"github.com/jmckind/conjur-operator/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// defaultLabels returns the default set of labels based on the given name.
func defaultLabels(name string) map[string]string {
	return map[string]string{
		common.ConjurKeyApp:     common.ConjurAppName,
		common.ConjurKeyAppName: nameWithSuffix(name, "conjur-oss"),
	}
}

// defaultMeta returns the default ObjectMeta based on the given namespace and name.
func defaultMeta(namespace string, name string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    defaultLabels(name),
	}
}

// getResource will retrieve the object with the given namespace and name using the Kubernetes API.
// The result will be stored in the given object.
func (r *ConjurReconciler) getResource(namespace string, name string, obj runtime.Object) error {
	return r.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, obj)
}

// isResourceFound will perform a basic check that the given object exists via the Kubernetes API.
// If an error occurs as part of the check, the function will return false.
func (r *ConjurReconciler) isResourceFound(namespace string, name string, obj runtime.Object) bool {
	if err := r.getResource(namespace, name, obj); err != nil {
		return false
	}
	return true
}

// nameWithSuffix is a shortcut function to return a string using the given
// name and suffix appended using a hyphen '-' as a separator. For example, if
// the given name is "example" and the suffix is "test", then the result will
// be "example-test".
func nameWithSuffix(name string, suffix string) string {
	return fmt.Sprintf("%s-%s", name, suffix)
}

// newConfigMap will return a new ConfigMap with the given namespace and name.
func newConfigMap(namespace string, name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: defaultMeta(namespace, name),
	}
}

// newDeployment will return a new Deployment with the given namespace and name.
func newDeployment(namespace string, name string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: defaultMeta(namespace, name),
	}
}

// newPersistentVolumeClaim will return a new PersistentVolumeClaim with the
// given namespace and name.
func newPersistentVolumeClaim(namespace string, name string) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: defaultMeta(namespace, name),
	}
}

// newSecret will return a new Secret with the given namespace and name.
func newSecret(namespace string, name string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: defaultMeta(namespace, name),
	}
}

// newService will return a new Service with the given namespace and name.
func newService(namespace string, name string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: defaultMeta(namespace, name),
	}
}

// newStatefulSet will return a new StatefulSet with the given namespace and
// name.
func newStatefulSet(namespace string, name string) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: defaultMeta(namespace, name),
	}
}
