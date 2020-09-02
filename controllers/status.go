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

	v1a1 "github.com/jmckind/conjur-operator/api/v1alpha1"
)

// reconcileStatus will ensure that all of the Status properties are updated for the given ArgoCD.
func (r *ConjurReconciler) reconcileStatus(cr *v1a1.Conjur) error {
	if err := r.reconcileStatusPhase(cr); err != nil {
		return err
	}
	return nil
}

// reconcileStatusPhase will ensure that the Status Phase is updated for the given ArgoCD.
func (r *ConjurReconciler) reconcileStatusPhase(cr *v1a1.Conjur) error {
	phase := "Unknown"

	// if cr.Status.Server == "Running" {
	// 	phase = "Available"
	// } else {
	// 	phase = "Pending"
	// }
	phase = "Pending"

	if cr.Status.Phase != phase {
		cr.Status.Phase = phase
		return r.Client.Status().Update(context.TODO(), cr)
	}
	return nil
}
