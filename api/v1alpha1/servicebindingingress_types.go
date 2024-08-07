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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServiceBindingIngressSpec defines the desired state of ServiceBindingIngress
type ServiceBindingIngressSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ServiceBindingIngress. Edit servicebindingingress_types.go to remove/update
	//Foo string `json:"foo,omitempty"`

	// EnableIngress is a flag to enable ingress,it can be empty
	EnableIngress bool `json:"enable_ingress,omitempty"`
	// EnableService is a flag to enable service,requeired
	EnableService bool `json:"enable_service"`
	// Replicas is the number of replicas,requeired
	Replicas int32 `json:"replicas"`
	// Image is the image to use,requeired
	Image string `json:"image"`
}

// ServiceBindingIngressStatus defines the observed state of ServiceBindingIngress
type ServiceBindingIngressStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ServiceBindingIngress is the Schema for the servicebindingingresses API
type ServiceBindingIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceBindingIngressSpec   `json:"spec,omitempty"`
	Status ServiceBindingIngressStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ServiceBindingIngressList contains a list of ServiceBindingIngress
type ServiceBindingIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceBindingIngress `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceBindingIngress{}, &ServiceBindingIngressList{})
}
