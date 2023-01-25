/*
Copyright 2022.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OpenStackDataPlaneNodeTemplateSpec defines the desired state of OpenStackDataPlaneNodeTemplate
type OpenStackDataPlaneNodeTemplateSpec struct {
	// +kubebuilder:validation:Optional
	// Node - node section defining attributes for this template
	Node NodeSection `json:"node,omitempty"`
}

// OpenStackDataPlaneNodeTemplateStatus defines the observed state of OpenStackDataPlaneNodeTemplate
type OpenStackDataPlaneNodeTemplateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OpenStackDataPlaneNodeTemplate is the Schema for the openstackdataplanenodetemplates API
type OpenStackDataPlaneNodeTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackDataPlaneNodeTemplateSpec   `json:"spec,omitempty"`
	Status OpenStackDataPlaneNodeTemplateStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpenStackDataPlaneNodeTemplateList contains a list of OpenStackDataPlaneNodeTemplate
type OpenStackDataPlaneNodeTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackDataPlaneNodeTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenStackDataPlaneNodeTemplate{}, &OpenStackDataPlaneNodeTemplateList{})
}
