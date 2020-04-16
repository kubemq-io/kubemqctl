package kubemqcluster

type VolumeConfig struct {
	// +optional
	Size string `json:"size"`

	// +optional
	StorageClass string `json:"storageClass"`
}
