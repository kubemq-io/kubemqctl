package builder

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
)

type BuilderRequest struct {
	Id    int64  `json:"id"`
	Type  string `json:"type"`
	Model struct {
		Name      string `json:"name,omitempty"`
		Namespace string `json:"namespace,omitempty"`
		Expose    struct {
			Mode     string `json:"mode"`
			NodePort int    `json:"nodePort,omitempty"`
		}
		SetInit     bool `json:"setInit,omitempty"`
		SetOperator bool `json:"setOperator,omitempty"`
	} `json:"model,omitempty"`
	Configuration string `json:"configuration"`
	configuration interface{}
}

func (b *BuilderRequest) Connector() *Connector {
	c := NewConnector().
		SetName(b.Model.Name).
		SetNamespace(b.Model.Namespace).
		SetType(b.Type).
		SetServiceType(b.Model.Expose.Mode).
		SetNodePort(int32(b.Model.Expose.NodePort))

	data, err := yaml.Marshal(b.configuration)
	if err == nil {
		c.SetConfig(string(data))
	}

	return c
}
func (b *BuilderRequest) Clusters() (clusters []*Cluster, namespaces map[string]string) {
	clustersDTO, ok := b.configuration.([]*ClusterDTO)
	if !ok {
		return nil, nil
	}
	namespaces = map[string]string{}
	for _, dto := range clustersDTO {
		cluster := NewCluster().FromDTO(dto)
		clusters = append(clusters, cluster)
		namespaces[cluster.Metadata.Namespace] = cluster.Metadata.Namespace
	}
	return
}

func (b *BuilderRequest) Validate() error {
	if b.Id < 0 {
		return fmt.Errorf("invalid request id")
	}
	switch b.Type {
	case "bridges", "targets", "sources":
		err := json.Unmarshal([]byte(b.Configuration), &b.configuration)
		if err != nil {
			return fmt.Errorf("invalid request connectors configuration: %s", err.Error())
		}
	case "clusters":
		var clusters []*ClusterDTO
		err := json.Unmarshal([]byte(b.Configuration), &clusters)
		if err != nil {
			return fmt.Errorf("invalid request clusters configuration: %s", err.Error())
		}
		b.configuration = clusters
	default:
		return fmt.Errorf("invalid request type")
	}
	return nil
}
