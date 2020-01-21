package deployment

import (
	"fmt"
	"strings"
)

type KubeMQManifestConfig struct {
	Id              string
	Name            string
	Namespace       string
	NamespaceConfig *NamespaceConfig
	StatefulSet     *StatefulSetConfig
	Services        map[string]*ServiceConfig
	ConfigMaps      map[string]*ConfigMapConfig
	Secrets         map[string]*SecretConfig
}

func NewKubeMQManifestConfig(id, name, namespace string) *KubeMQManifestConfig {
	return &KubeMQManifestConfig{
		Id:              id,
		Name:            name,
		Namespace:       namespace,
		NamespaceConfig: NewNamespaceConfig(id, name),
		StatefulSet:     NewStatefulSetConfig(id, name, namespace),
		Services:        make(map[string]*ServiceConfig),
		ConfigMaps:      make(map[string]*ConfigMapConfig),
		Secrets:         make(map[string]*SecretConfig),
	}
}

func DefaultKubeMQManifestConfig(id, name, namespace string) *KubeMQManifestConfig {
	return &KubeMQManifestConfig{
		Id:              id,
		Name:            name,
		Namespace:       namespace,
		NamespaceConfig: DefaultNamespaceConfig(id, name),
		StatefulSet:     DefaultStatefulSetConfig(id, name, namespace),
		Services:        DefaultServiceConfig(id, namespace, name),
		ConfigMaps:      DefaultConfigMap(id, name, namespace),
		Secrets:         DefaultSecretConfig(id, name, namespace),
	}
}
func (c *KubeMQManifestConfig) SetConfigMapValues(cmName, key, value string) {
	cm, ok := c.ConfigMaps[cmName]
	if ok {
		cm.SetVariable(key, value)
	}
}
func (c *KubeMQManifestConfig) SetSecretValues(secName, key, value string) {
	sec, ok := c.Secrets[secName]
	if ok {
		sec.SetVariable(key, value)
	}
}
func (c *KubeMQManifestConfig) Spec() ([]byte, error) {
	var manifest []string
	nsSpec, err := c.NamespaceConfig.Spec()
	if err != nil {
		return nil, fmt.Errorf("error on namespace spec rendring: %s", err.Error())
	}
	manifest = append(manifest, string(nsSpec))
	stsSpec, err := c.StatefulSet.Spec()
	if err != nil {
		return nil, fmt.Errorf("error on statefull spec rendring: %s", err.Error())
	}

	manifest = append(manifest, string(stsSpec))

	for name, svc := range c.Services {
		svcSpec, err := svc.Spec()
		if err != nil {
			return nil, fmt.Errorf("error on service %s spec rendring: %s", name, err.Error())
		}
		manifest = append(manifest, string(svcSpec))
	}
	for _, cm := range c.ConfigMaps {
		configMapSpec, err := cm.Spec()
		if err != nil {
			return nil, fmt.Errorf("error on config map spec rendring: %s", err.Error())
		}
		manifest = append(manifest, string(configMapSpec))
	}

	for _, sec := range c.Secrets {
		secretSpec, err := sec.Spec()
		if err != nil {
			return nil, fmt.Errorf("error on secret spec rendring: %s", err.Error())
		}
		manifest = append(manifest, string(secretSpec))
	}

	return []byte(strings.Join(manifest, "\n---\n")), nil
}
