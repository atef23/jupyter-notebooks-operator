package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JupyterNotebooksSpec defines the desired state of JupyterNotebooks
type JupyterNotebooksSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Size int32 `json:"size"`
}

// JupyterNotebooksStatus defines the observed state of JupyterNotebooks
type JupyterNotebooksStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Nodes []string `json:"nodes"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JupyterNotebooks is the Schema for the jupyternotebooks API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=jupyternotebooks,scope=Namespaced
type JupyterNotebooks struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JupyterNotebooksSpec   `json:"spec,omitempty"`
	Status JupyterNotebooksStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JupyterNotebooksList contains a list of JupyterNotebooks
type JupyterNotebooksList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JupyterNotebooks `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JupyterNotebooks{}, &JupyterNotebooksList{})
}
