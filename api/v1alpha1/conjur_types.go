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

// Important: Run "make" to regenerate code after modifying this file
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConjurDatabaseSpec defines the desired state for the Conjur database.
type ConjurDatabaseSpec struct {
	// Count is the replica count for the Conjur database StatefulSet.
	Count *int32 `json:"count,omitempty"`

	// URL is the PostgreSQL connection string to be used.
	// If not provided, an integrated PostgreSQL deployment will be created.
	URL string `json:"url,omitempty"`
}

// ConjurServerSpec defines the desired state for the Conjur server.
type ConjurServerSpec struct {
	// Count is the replica count for the Conjur server Deployment.
	Count *int32 `json:"count,omitempty"`
}

// ConjurSpec defines the desired state of Conjur
type ConjurSpec struct {
	// Account is the name of the Conjur account to be used.
	Account string `json:"account,omitempty"`

	// Database is the desired state for the Conjur database.
	Database ConjurDatabaseSpec `json:"database,omitempty"`

	// Server is the desired state for the Conjur sever.
	Server ConjurServerSpec `json:"server,omitempty"`
}

// ConjurStatus defines the observed state of Conjur
type ConjurStatus struct {
	// Phase is a simple, high-level summary of where Conjur is in its lifecycle.
	// There are five possible phase values:
	// Pending: The Conjur resource has been accepted by the Kubernetes system.
	// Available: All of the resources for Conjur are ready.
	// Failed: At least one resource has experienced a failure.
	// Unknown: For some reason the state of the Conjur phase could not be obtained.
	Phase string `json:"phase,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Conjur is the Schema for the conjurs API
type Conjur struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConjurSpec   `json:"spec,omitempty"`
	Status ConjurStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ConjurList contains a list of Conjur
type ConjurList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Conjur `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Conjur{}, &ConjurList{})
}
