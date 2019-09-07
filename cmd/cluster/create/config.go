package create

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
)

type StatefulSetConfig struct {
	ApiVersion string
	Name       string
	Namespace  string
	Replicas   int
	Token      string
	Version    string
	Volume     int
}

type ServiceConfig struct {
	ApiVersion    string
	Name          string
	Namespace     string
	AppName       string
	Type          string
	ContainerPort int
	Protocol      string
	TargetPort    int
	TargetApp     string
	PortName      string
}

type StatefulSetDeployment struct {
	Namespace   *apiv1.Namespace
	StatefulSet *appsv1.StatefulSet
	Services    map[string]*apiv1.Service
}

func NewStatefulSetConfig(o *CreateOptions) StatefulSetConfig {
	return StatefulSetConfig{
		ApiVersion: o.appsVersion,
		Name:       o.name,
		Namespace:  o.namespace,
		Replicas:   o.replicas,
		Token:      o.token,
		Version:    o.version,
		Volume:     o.volume,
	}
}

func (s StatefulSetConfig) Spec() ([]byte, error) {
	t := NewTemplate(defaultStsTemplate, s)
	return t.Get()
}

func NewServiceConfigs(o *CreateOptions) []ServiceConfig {
	list := []ServiceConfig{}
	svc := ServiceConfig{
		ApiVersion:    o.coreVersion,
		Name:          o.name,
		Namespace:     o.namespace,
		AppName:       o.name,
		Type:          "ClusterIP",
		ContainerPort: 5228,
		Protocol:      "TCP",
		TargetPort:    5228,
		TargetApp:     o.name,
		PortName:      "cluster-port",
	}

	svcGrpc := ServiceConfig{
		ApiVersion:    o.coreVersion,
		Name:          o.name + "-grpc",
		Namespace:     o.namespace,
		AppName:       o.name,
		Type:          "ClusterIP",
		ContainerPort: 50000,
		Protocol:      "TCP",
		TargetPort:    50000,
		TargetApp:     o.name,
		PortName:      "grpc-port",
	}
	if o.isNodePort {
		svcGrpc.Type = "NodePort"
	}
	if o.isLoadBalance {
		svcGrpc.Type = "LoadBalancer"
	}

	svcRest := ServiceConfig{
		ApiVersion:    o.coreVersion,
		Name:          o.name + "-rest",
		Namespace:     o.namespace,
		AppName:       o.name,
		Type:          "ClusterIP",
		ContainerPort: 9090,
		Protocol:      "TCP",
		TargetPort:    9090,
		TargetApp:     o.name,
		PortName:      "rest-port",
	}
	if o.isNodePort {
		svcRest.Type = "NodePort"
	}
	if o.isLoadBalance {
		svcRest.Type = "LoadBalancer"
	}

	svcApi := ServiceConfig{
		ApiVersion:    o.coreVersion,
		Name:          o.name + "-api",
		Namespace:     o.namespace,
		AppName:       o.name,
		Type:          "ClusterIP",
		ContainerPort: 8080,
		Protocol:      "TCP",
		TargetPort:    8080,
		TargetApp:     o.name,
		PortName:      "api-port",
	}
	if o.isNodePort {
		svcApi.Type = "NodePort"
	}
	if o.isLoadBalance {
		svcApi.Type = "LoadBalancer"
	}

	list = append(list, svc, svcGrpc, svcRest, svcApi)

	return list
}

func (s ServiceConfig) Spec() ([]byte, error) {
	t := NewTemplate(defaultServiceTemplate, s)
	return t.Get()

}
