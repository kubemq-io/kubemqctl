package build

import "fmt"

type resources struct {
	clusters   map[string]string
	connectors map[string]string
}

func newResources() *resources {
	return &resources{
		clusters:   map[string]string{},
		connectors: map[string]string{},
	}
}

func (r *resources) updateClusterName(namespace, name string) error {
	key := fmt.Sprintf("%s/%s", namespace, name)
	_, ok := r.clusters[key]
	if ok {
		return fmt.Errorf("cluster with the name %s already exists in namesapce %s", name, namespace)
	}
	r.clusters[key] = key

	return nil
}

func (r *resources) updateConnector(namespace, name string) error {
	key := fmt.Sprintf("%s/%s", namespace, name)
	_, ok := r.connectors[key]
	if ok {
		return fmt.Errorf("connector with the name %s already exists in namesapce %s", name, namespace)
	}
	r.connectors[key] = key
	return nil
}
