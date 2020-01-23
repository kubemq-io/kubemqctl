package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
	"strings"
)

var serviceTypes = map[string]string{"clusterip": "ClusterIP", "nodeport": "NodePort", "loadbalancer": "LoadBalancer", "ingress": "Ingress"}

type deployServiceOptions struct {
	apiPortValue  uint
	apiPortType   string
	grpcPortValue uint
	grpcPortType  string
	restPortValue uint
	restPortType  string
}

func defaultServiceConfig(cmd *cobra.Command) *deployServiceOptions {
	o := &deployServiceOptions{
		apiPortValue:  8080,
		apiPortType:   "ClusterIP",
		grpcPortValue: 50000,
		grpcPortType:  "ClusterIP",
		restPortValue: 9090,
		restPortType:  "ClusterIP",
	}
	cmd.PersistentFlags().UintVarP(&o.apiPortValue, "api-port", "", 8080, "set api port value")
	cmd.PersistentFlags().StringVarP(&o.apiPortType, "api-port-type", "", "ClusterIP", "set api port service type (ClusterIP,NodePort,LoadBalancer)")
	cmd.PersistentFlags().UintVarP(&o.grpcPortValue, "grpc-port", "", 50000, "set grpc port value")
	cmd.PersistentFlags().StringVarP(&o.grpcPortType, "grpc-port-type", "", "ClusterIP", "set grpc port service type (ClusterIP,NodePort,LoadBalancer)")
	cmd.PersistentFlags().UintVarP(&o.restPortValue, "rest-port", "", 9090, "set grpc port value")
	cmd.PersistentFlags().StringVarP(&o.restPortType, "rest-port-type", "", "ClusterIP", "set rest port service type (ClusterIP,NodePort,LoadBalancer)")
	return o
}

func (o *deployServiceOptions) validate() error {
	if o.apiPortValue == 0 || o.apiPortValue > 65535 {
		return fmt.Errorf("invalid api-port value: %d", o.apiPortValue)
	}
	if _, ok := serviceTypes[strings.ToLower(o.apiPortType)]; !ok {
		return fmt.Errorf("invalid api-port-type value: %s", o.apiPortType)
	}
	o.apiPortType = serviceTypes[strings.ToLower(o.apiPortType)]

	if o.grpcPortValue == 0 || o.grpcPortValue > 65535 {
		return fmt.Errorf("invalid grpc-port value: %d", o.grpcPortValue)
	}

	if _, ok := serviceTypes[strings.ToLower(o.grpcPortType)]; !ok {
		return fmt.Errorf("invalid grpc-port-type value: %s", o.grpcPortType)
	}
	o.grpcPortType = serviceTypes[strings.ToLower(o.grpcPortType)]

	if o.restPortValue == 0 || o.restPortValue > 65535 {
		return fmt.Errorf("invalid rest-port value: %d", o.restPortValue)
	}
	if _, ok := serviceTypes[strings.ToLower(o.restPortType)]; !ok {
		return fmt.Errorf("invalid rest-port-type value: %s", o.restPortType)
	}
	o.restPortType = serviceTypes[strings.ToLower(o.restPortType)]
	return nil
}

func (o *deployServiceOptions) complete() error {
	return nil
}
func (o *deployServiceOptions) addIngressConfig(config *deployment.KubeMQManifestConfig, svcName string, svcPort int) *deployServiceOptions {
	ingressName := fmt.Sprintf("%s-ingress", svcName)
	config.Ingress[ingressName] = deployment.NewIngressConfig(config.Id, svcName+"-ingress", config.Namespace, config.Name).
		SetServiceName(svcName).
		SetServicePort(svcPort)

	return o
}

func (o *deployServiceOptions) setConfig(config *deployment.KubeMQManifestConfig) *deployServiceOptions {
	svc, ok := config.Services["api"]
	if ok {
		svc.TargetPort = 8080
		svc.ContainerPort = int(o.apiPortValue)
		svc.Type = o.apiPortType
		if svc.Type == "Ingress" {
			svc.Type = "ClusterIP"
			o.addIngressConfig(config, svc.Name, int(o.apiPortValue))
		}
	}
	svc, ok = config.Services["grpc"]
	if ok {
		svc.TargetPort = 50000
		svc.ContainerPort = int(o.grpcPortValue)
		svc.Type = o.grpcPortType
		if svc.Type == "Ingress" {
			svc.Type = "ClusterIP"
			o.addIngressConfig(config, svc.Name, int(o.grpcPortValue))
		}
	}
	svc, ok = config.Services["rest"]
	if ok {
		svc.TargetPort = 9090
		svc.ContainerPort = int(o.restPortValue)
		svc.Type = o.restPortType
		if svc.Type == "Ingress" {
			svc.Type = "ClusterIP"
			o.addIngressConfig(config, svc.Name, int(o.restPortValue))
		}
	}
	return o
}
