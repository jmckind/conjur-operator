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
	v1a1 "github.com/jmckind/conjur-operator/api/v1alpha1"
)

// reconcileServer will ensure that all server resources are present owned
// by the given Conjur resource.
func (r *ConjurReconciler) reconcileServer(cr *v1a1.Conjur) error {
	return nil
}
