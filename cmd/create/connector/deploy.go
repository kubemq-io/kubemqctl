package connector

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqconnector"
	"github.com/spf13/cobra"
	"io/ioutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deployOptions struct {
	name          string
	namespace     string
	port          int32
	replicas      int32
	connectorType string
	image         string
	serviceType   string
	configFile    string
	configData    string
}

func defaultDeployOptions(cmd *cobra.Command) *deployOptions {
	o := &deployOptions{
		port: 0,
	}
	cmd.PersistentFlags().StringVarP(&o.name, "name", "", "kubemq-connector", "set kubemq connector name")
	cmd.PersistentFlags().StringVarP(&o.namespace, "namespace", "n", "kubemq", "set kubemq connector namespace")
	cmd.PersistentFlags().Int32VarP(&o.port, "port", "p", 0, "set kubemq connector api port")
	cmd.PersistentFlags().Int32VarP(&o.replicas, "replicas", "r", 1, "set kubemq connector replicase")
	cmd.PersistentFlags().StringVarP(&o.connectorType, "type", "t", "", "set kubemq connector type: targets/sources/bridges")
	cmd.PersistentFlags().StringVarP(&o.image, "image", "", "", "set kubemq connector docker image")
	cmd.PersistentFlags().StringVarP(&o.serviceType, "service-type", "", "ClusterIP", "set kubemq connector api service type, default ClusterIP")
	cmd.PersistentFlags().StringVarP(&o.configFile, "config", "c", "", "set kubemq connector configFile file name")
	return o
}

func (o *deployOptions) validate() error {
	if o.replicas < 0 {
		return fmt.Errorf("invalid replicas value, must be greater than 0")
	}
	if o.port < 0 {
		return fmt.Errorf("invalid port value, must be greater than 0")
	}

	if o.connectorType != "targets" && o.connectorType != "source" && o.connectorType != "bridges" {
		return fmt.Errorf("invalid connector type, must be one of targets/sources/bridges")
	}

	if o.serviceType != "ClusterIP" && o.serviceType != "NodePort" && o.serviceType != "LoadBalancer" {
		return fmt.Errorf("invalid service type, must be one of ClusterIP/NodePort/LoadBalancer")
	}

	if o.configData == "" {
		return fmt.Errorf("invalid configuration data, cannot be empty configuration")
	}

	return nil
}

func (o *deployOptions) complete() error {
	if o.connectorType == "" {
		prompt := &survey.Select{
			Message: "Choose Connector type:",
			Options: []string{"targets", "sources", "bridges"},
		}
		_ = survey.AskOne(prompt, &o.connectorType)
	}
	if o.configFile != "" {
		data, err := ioutil.ReadFile(o.configFile)
		if err != nil {
			return fmt.Errorf("error reading config file data: %s", err.Error())
		}
		o.configData = string(data)
	}
	if o.configData == "" {
		prompt := &survey.Editor{
			Message:  "Config file",
			FileName: "*.yaml",
		}
		_ = survey.AskOne(prompt, &o.configData)
	}
	return nil
}

func (o *deployOptions) getConnectorDeployment() *kubemqconnector.KubemqConnector {

	deployment := &kubemqconnector.KubemqConnector{
		TypeMeta: v1.TypeMeta{
			Kind:       "KubemqConnector",
			APIVersion: "core.k8s.kubemq.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", o.name, o.connectorType),
			Namespace: o.namespace,
		},
		Spec: kubemqconnector.KubemqConnectorSpec{
			Replicas:    new(int32),
			Type:        o.connectorType,
			Image:       o.image,
			Config:      o.configData,
			NodePort:    o.port,
			ServiceType: o.serviceType,
		},
		Status: kubemqconnector.KubemqConnectorStatus{},
	}

	return deployment
}
