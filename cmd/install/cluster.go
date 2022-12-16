package install

import (
	"fmt"

	"github.com/kubemq-io/kubemqctl/pkg/template"
)

var clusterTemplate = `
apiVersion: core.k8s.kubemq.io/v1beta1
kind: KubemqCluster
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  replicas: {{.Replicas}}
{{if not (eq .Key "")}}
  key: {{.Key}}
{{end}}
{{if not (eq .License "")}}
  license: |-
{{.License | indent 4}}
{{end}}
`

type ClusterManifest struct {
	Name      string
	Namespace string
	Key       string
	License   string
	Replicas  int
}

func NewClusterManifest() *ClusterManifest {
	return &ClusterManifest{
		Name:      "kubemq-cluster",
		Namespace: "kubemq",
		Key:       "",
		License:   "",
		Replicas:  3,
	}
}

func (d *ClusterManifest) SetName(value string) *ClusterManifest {
	d.Name = value
	return d
}

func (d *ClusterManifest) SetNamespace(value string) *ClusterManifest {
	d.Namespace = value
	return d
}

func (d *ClusterManifest) SetKey(value string) *ClusterManifest {
	d.Key = value
	return d
}

func (d *ClusterManifest) SetLicense(value string) *ClusterManifest {
	d.License = value
	return d
}

func (d *ClusterManifest) SetReplicas(value int) *ClusterManifest {
	d.Replicas = value
	return d
}

func (d *ClusterManifest) Manifest() (string, error) {
	if d.Key == "" && d.License == "" {
		return "", fmt.Errorf("cluster manifest must have a key or license, please register at https://kubemq.io")
	}
	if d.Replicas < 3 {
		return "", fmt.Errorf("cluster manifest must have at least 3 replicas")
	}

	t := template.NewTemplate(clusterTemplate, d)
	data, err := t.Get()
	if err != nil {
		return "", fmt.Errorf("error during cluster manifest generation, error: %s", err.Error())
	}
	return string(data), nil
}
