package cluster

type AuthorizationConfig struct {
	// +optional
	// +kubebuilder:validation:MinLength=1
	Policy string `json:"policy"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
	Url string `json:"url"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	AutoReload int32 `json:"autoReload"`
}
