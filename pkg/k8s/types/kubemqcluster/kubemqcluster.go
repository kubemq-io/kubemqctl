package kubemqcluster

import (
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubemqClusterSpec struct {

	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty"`

	// +optional
	License string `json:"license,omitempty"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	ConfigData string `json:"configData,omitempty"`

	// +optional
	Volume *VolumeConfig `json:"volume,omitempty"`

	// +optional
	Image *ImageConfig `json:"image,omitempty"`

	// +optional
	Api *ApiConfig `json:"api,omitempty"`

	// +optional
	Rest *RestConfig `json:"rest,omitempty"`

	// +optional
	Grpc *GrpcConfig `json:"grpc,omitempty"`

	// +optional
	Tls *TlsConfig `json:"tls,omitempty"`

	// +optional
	Resources *ResourceConfig `json:"resources,omitempty"`

	// +optional
	NodeSelectors *NodeSelectorConfig `json:"nodeSelectors,omitempty"`

	// +optional
	Authentication *AuthenticationConfig `json:"authentication,omitempty"`

	// +optional
	Authorization *AuthorizationConfig `json:"authorization,omitempty"`

	// +optional
	Health *HealthConfig `json:"health,omitempty"`

	// +optional
	Routing *RoutingConfig `json:"routing,omitempty"`

	// +optional
	Log *LogConfig `json:"log,omitempty"`

	// +optional
	Notification *NotificationConfig `json:"notification,omitempty"`

	// +optional
	Store *StoreConfig `json:"store,omitempty"`

	// +optional
	Queue *QueueConfig `json:"queue,omitempty"`

	// +optional
	Gateways *GatewayConfig `json:"gateways,omitempty"`
}

// KubemqClusterStatus defines the observed state of KubemqCluster
type KubemqClusterStatus struct {
	Replicas *int32 `json:"replicas"`

	Version string `json:"version"`

	Ready int32 `json:"ready"`

	Grpc string `json:"grpc"`

	Rest string `json:"rest"`

	Api string `json:"api"`

	Selector string `json:"selector"`

	LicenseType string `json:"license_type"`

	LicenseTo string `json:"license_to"`

	LicenseExpire string `json:"license_expire"`

	Status string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// KubemqCluster is the Schema for the kubemqclusters API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubemqclusters,scope=Namespaced
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:printcolumn:JSONPath=".status.version",name=Version,type=string
// +kubebuilder:printcolumn:JSONPath=".status.status",name=Status,type=string
// +kubebuilder:printcolumn:JSONPath=".status.replicas",name=Replicas,type=string
// +kubebuilder:printcolumn:JSONPath=".status.ready",name=Ready,type=string
// +kubebuilder:printcolumn:JSONPath=".status.grpc",name=gRPC,type=string
// +kubebuilder:printcolumn:JSONPath=".status.rest",name=Rest,type=string
// +kubebuilder:printcolumn:JSONPath=".status.api",name=API,type=string
// +kubebuilder:printcolumn:JSONPath=".status.license_type",name=License-type,type=string
// +kubebuilder:printcolumn:JSONPath=".status.license_to",name=License-To,type=string
// +kubebuilder:printcolumn:JSONPath=".status.license_expire",name=License-Expire,type=string
type KubemqCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubemqClusterSpec   `json:"spec,omitempty"`
	Status KubemqClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// KubemqClusterList contains a list of KubemqCluster
type KubemqClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubemqCluster `json:"items"`
}

func (k *KubemqCluster) String() string {
	data, err := yaml.Marshal(k)
	if err != nil {
		return ""
	}
	return string(data)
}
