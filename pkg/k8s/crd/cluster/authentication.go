package cluster

type AuthenticationConfig struct {
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key"`

	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=(HS256|HS384|HS512|RS256|RS384|RS512|ES256|ES384|ES512)
	Type string `json:"type"`
}
