package builder

import (
	"fmt"
)

var deployTemplate = `
apiVersion: core.k8s.kubemq.io/v1beta1
kind: KubemqCluster
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  replicas: 3
  key: {{.Key}}
`

type DeployManifest struct {
	Name      string
	Namespace string
	Key       string
}

func NewDeployManifest() *DeployManifest {
	return &DeployManifest{
		Name:      "kubemq-cluster",
		Namespace: "kubemq",
		Key:       "",
	}
}

func (d *DeployManifest) SetName(value string) *DeployManifest {
	d.Name = value
	return d
}

func (d *DeployManifest) SetNamespace(value string) *DeployManifest {
	d.Namespace = value
	return d
}

func (d *DeployManifest) SetKey(value string) *DeployManifest {
	d.Key = value
	return d
}

func (d *DeployManifest) Manifest() (string, error) {
	if d.Key == "" {
		return "", fmt.Errorf("deployment manifest must have a key, please register at https://kubemq.io")
	}
	om, _ := NewOperatorManifest().SetNamespace(d.Namespace).Manifest()
	t := NewTemplate(deployTemplate, d)
	data, err := t.Get()
	if err != nil {
		return "", fmt.Errorf("error during deployment manifest generation, error: %s", err.Error())
	}
	cluster := string(data)
	manifest := fmt.Sprintf("%s\n---\n%s", om, cluster)
	return manifest, nil
}
