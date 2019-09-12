package deployment

type Options struct {
	Token         string
	Replicas      int
	Version       string
	Namespace     string
	Name          string
	AppsVersion   string
	CoreVersion   string
	Volume        int
	IsNodePort    bool
	IsLoadBalance bool
}

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

func NewStatefulSetConfig(o *Options) StatefulSetConfig {
	return StatefulSetConfig{
		ApiVersion: o.AppsVersion,
		Name:       o.Name,
		Namespace:  o.Namespace,
		Replicas:   o.Replicas,
		Token:      o.Token,
		Version:    o.Version,
		Volume:     o.Volume,
	}
}

func (s StatefulSetConfig) Spec() ([]byte, error) {
	t := NewTemplate(defaultStsTemplate, s)
	return t.Get()
}

func NewServiceConfigs(o *Options) []ServiceConfig {
	list := []ServiceConfig{}
	svc := ServiceConfig{
		ApiVersion:    o.CoreVersion,
		Name:          o.Name,
		Namespace:     o.Namespace,
		AppName:       o.Name,
		Type:          "ClusterIP",
		ContainerPort: 5228,
		Protocol:      "TCP",
		TargetPort:    5228,
		TargetApp:     o.Name,
		PortName:      "cluster-port",
	}

	svcGrpc := ServiceConfig{
		ApiVersion:    o.CoreVersion,
		Name:          o.Name + "-grpc",
		Namespace:     o.Namespace,
		AppName:       o.Name,
		Type:          "ClusterIP",
		ContainerPort: 50000,
		Protocol:      "TCP",
		TargetPort:    50000,
		TargetApp:     o.Name,
		PortName:      "grpc-port",
	}
	if o.IsNodePort {
		svcGrpc.Type = "NodePort"
	}
	if o.IsLoadBalance {
		svcGrpc.Type = "LoadBalancer"
	}

	svcRest := ServiceConfig{
		ApiVersion:    o.CoreVersion,
		Name:          o.Name + "-rest",
		Namespace:     o.Namespace,
		AppName:       o.Name,
		Type:          "ClusterIP",
		ContainerPort: 9090,
		Protocol:      "TCP",
		TargetPort:    9090,
		TargetApp:     o.Name,
		PortName:      "rest-port",
	}
	if o.IsNodePort {
		svcRest.Type = "NodePort"
	}
	if o.IsLoadBalance {
		svcRest.Type = "LoadBalancer"
	}

	svcApi := ServiceConfig{
		ApiVersion:    o.CoreVersion,
		Name:          o.Name + "-api",
		Namespace:     o.Namespace,
		AppName:       o.Name,
		Type:          "ClusterIP",
		ContainerPort: 8080,
		Protocol:      "TCP",
		TargetPort:    8080,
		TargetApp:     o.Name,
		PortName:      "api-port",
	}
	if o.IsNodePort {
		svcApi.Type = "NodePort"
	}
	if o.IsLoadBalance {
		svcApi.Type = "LoadBalancer"
	}

	list = append(list, svc, svcGrpc, svcRest, svcApi)

	return list
}

func (s ServiceConfig) Spec() ([]byte, error) {
	t := NewTemplate(defaultServiceTemplate, s)
	return t.Get()

}
