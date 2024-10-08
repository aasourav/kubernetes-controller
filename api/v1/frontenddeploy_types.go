/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FrontendDeploySpec defines the desired state of FrontendDeploy
type EnvironmentVariable struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type FrontendDeploySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of FrontendDeploy. Edit frontenddeploy_types.go to remove/update
	ImageName             string                `json:"imageName"`
	Port                  int32                 `json:"port"`
	Replicas              int32                 `json:"replicas,omitempty"`
	NodeName              string                `json:"nodeName,omitempty"`
	IsHost                bool                  `json:"isHost,omitempty"`
	EnvironmentVarialbles []EnvironmentVariable `json:"environmentVariables,omitempty"`
}

// FrontendDeployStatus defines the observed state of FrontendDeploy
type FrontendDeployStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// FrontendDeploy is the Schema for the frontenddeploys API
type FrontendDeploy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FrontendDeploySpec   `json:"spec,omitempty"`
	Status FrontendDeployStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FrontendDeployList contains a list of FrontendDeploy
type FrontendDeployList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FrontendDeploy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FrontendDeploy{}, &FrontendDeployList{})
}
