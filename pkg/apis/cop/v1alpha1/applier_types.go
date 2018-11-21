package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ApplierSpec defines the desired state of Applier
type ApplierSpec struct {
	Source         ApplierSource  `json:"source"`
	Webhook        ApplierWebhook `json:"webhook,omitempty"`
	ServiceAccount string         `json:"serviceAccount,omitempty"`
}

// ApplierStatus defines the observed state of Applier
type ApplierStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Applier is the Schema for the appliers API
// +k8s:openapi-gen=true
type Applier struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplierSpec   `json:"spec,omitempty"`
	Status ApplierStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApplierList contains a list of Applier
type ApplierList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Applier `json:"items"`
}

// ApplierSource contains a source for the Applier
type ApplierSource struct {
	Git ApplierGit `json:"git"`
}

// ApplierGit contains the git asset for the Applier Source
type ApplierGit struct {
	URI          string `json:"uri"`
	Ref          string `json:"ref,omitempty"`
	InventoryDir string `json:"inventoryDir,omitempty"`
	HTTPProxy    string `json:"httpProxy,omitempty"`
	HTTPSProxy   string `json:"httpsProxy,omitempty"`
	NoProxy      string `json:"noProxy,omitempty"`
	SecretName   string `json:"secretName,omitempty"`
}

// ApplierWebhook contains the webhook associated with the Applier
type ApplierWebhook struct {
	Token string `json:"token"`
}

func init() {
	SchemeBuilder.Register(&Applier{}, &ApplierList{})
}
