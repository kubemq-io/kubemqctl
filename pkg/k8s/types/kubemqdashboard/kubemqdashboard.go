package kubemqdashboard

import (
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubemqDashboardSpec defines the desired state of KubemqDashboard
type KubemqDashboardSpec struct {
	// +optional
	Port int32 `json:"port,omitempty"`

	// +optional
	Prometheus *PrometheusConfig `json:"prometheus,omitempty"`

	// +optional
	Grafana *GrafanaConfig `json:"grafana,omitempty"`
}

// KubemqDashboardStatus defines the observed state of KubemqDashboard
type KubemqDashboardStatus struct {
	Status            string `json:"status"`
	Address           string `json:"address"`
	PrometheusVersion string `json:"prometheus_version"`
	GrafanaVersion    string `json:"grafana_version"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubemqDashboard is the Schema for the kubemqdashboards API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubemqdashboards,scope=Namespaced
// +kubebuilder:printcolumn:JSONPath=".status.status",name=Status,type=string
// +kubebuilder:printcolumn:JSONPath=".status.address",name=Address,type=string
// +kubebuilder:printcolumn:JSONPath=".status.prometheus_version",name=Prometheus-Version,type=string
// +kubebuilder:printcolumn:JSONPath=".status.grafana_version",name=Grafana-Version,type=string
type KubemqDashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubemqDashboardSpec   `json:"spec,omitempty"`
	Status KubemqDashboardStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// KubemqDashboardList contains a list of KubemqDashboard
type KubemqDashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubemqDashboard `json:"items"`
}

type PrometheusConfig struct {
	// +optional
	NodePort int32 `json:"nodePort,omitempty"`
	// +optional
	Image string `json:"image,omitempty"`
}

type GrafanaConfig struct {
	// +optional
	DashboardUrl string `json:"dashboardUrl,omitempty"`
	// +optional
	Image string `json:"image,omitempty"`
}

func (k *KubemqDashboard) String() string {
	data, err := yaml.Marshal(k)
	if err != nil {
		return ""
	}
	return string(data)
}
