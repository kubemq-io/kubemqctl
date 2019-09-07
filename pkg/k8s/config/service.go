package config

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
