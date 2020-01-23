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
	Ingress         map[string]*IngressConfig
}

func NewKubeMQManifestConfig(id, name, namespace string) *KubeMQManifestConfig {
	return &KubeMQManifestConfig{
		Id:              id,
		Name:            name,
		Namespace:       namespace,
		NamespaceConfig: NewNamespaceConfig(id, namespace),
		StatefulSet:     NewStatefulSetConfig(id, name, namespace),
		Services:        make(map[string]*ServiceConfig),
		ConfigMaps:      make(map[string]*ConfigMapConfig),
		Secrets:         make(map[string]*SecretConfig),
		Ingress:         map[string]*IngressConfig{},
	}
}

func DefaultKubeMQManifestConfig(id, name, namespace string) *KubeMQManifestConfig {
	return &KubeMQManifestConfig{
		Id:              id,
		Name:            name,
		Namespace:       namespace,
		NamespaceConfig: DefaultNamespaceConfig(id, namespace),
		StatefulSet:     DefaultStatefulSetConfig(id, name, namespace),
		Services:        DefaultServiceConfig(id, namespace, name),
		ConfigMaps:      DefaultConfigMap(id, name, namespace),
		Secrets:         DefaultSecretConfig(id, name, namespace),
		Ingress:         DefaultIngressConfig(),
	}
}
func (c *KubeMQManifestConfig) SetConfigMapValues(cmName, key, value string) {
	cm, ok := c.ConfigMaps[cmName]
	if ok {
		cm.SetVariable(key, value)
	}
}
func (c *KubeMQManifestConfig) SetSecretStringValues(secName, key, value string) {
	sec, ok := c.Secrets[secName]
	if ok {
		sec.SetStringVariable(key, value)
	}
}
func (c *KubeMQManifestConfig) SetSecretDataValues(secName, key, value string) {
	sec, ok := c.Secrets[secName]
	if ok {
		sec.SetDataVariable(key, value)
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
	for _, ing := range c.Ingress {
		ingresSpec, err := ing.Spec()
		if err != nil {
			return nil, fmt.Errorf("error on ingres spec rendring: %s", err.Error())
		}
		manifest = append(manifest, string(ingresSpec))
	}
	return []byte(strings.Join(manifest, "\n---\n")), nil
}
