package create

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type deployOptions struct {
	configFilename  string
	name            string
	namespace       string
	token           string
	licenseData     string
	licenseDataFile string
	tag             string
	volume          uint
	replicas        uint
	service         *deployServiceOptions
	security        *deploySecurityOptions
	authentication  *deployAuthenticationOptions
	authorization   *deployAuthorizationOptions
	gateway         *deployGatewayOptions
	resources       *deployResourceOptions
	nodeSelectors   *deployNodeSelectorOptions
	healthProbe     *deployHealthOptions
}

func defaultDeployOptions(cmd *cobra.Command) *deployOptions {
	o := &deployOptions{
		configFilename:  "",
		name:            "",
		namespace:       "",
		token:           "",
		licenseData:     "",
		licenseDataFile: "",
		tag:             "",
		volume:          0,
		replicas:        0,
		service:         defaultServiceConfig(cmd),
		security:        defaultSecurityConfig(cmd),
		authentication:  defaultAuthenticationOptions(cmd),
		authorization:   defaultAuthorizationConfig(cmd),
		gateway:         defaultGatewayOptions(cmd),
		resources:       defaultResourceOptions(cmd),
		nodeSelectors:   defaultNodeSelectorOptions(cmd),
		healthProbe:     defaultHealthOptions(cmd),
	}
	cmd.PersistentFlags().StringVarP(&o.configFilename, "config-file", "c", "", "set kubemq config file")
	cmd.PersistentFlags().StringVarP(&o.name, "name", "n", "kubemq-cluster", "set kubemq cluster name")
	cmd.PersistentFlags().StringVarP(&o.namespace, "namespace", "", "kubemq", "set kubemq cluster namespace")
	cmd.PersistentFlags().StringVarP(&o.token, "token", "t", "", "set kubemq token")
	cmd.PersistentFlags().StringVarP(&o.licenseData, "license-data", "d", "", "set license data")
	cmd.PersistentFlags().StringVarP(&o.licenseDataFile, "license-data-file", "l", "", "set license data filename")
	cmd.PersistentFlags().StringVarP(&o.tag, "tag", "T", "latest", "set kubemq docker image tag")
	cmd.PersistentFlags().UintVarP(&o.volume, "volume", "v", 0, "set persistence volume claim size")
	cmd.PersistentFlags().UintVarP(&o.replicas, "replicas", "r", 3, "set replicas")
	return o
}

func (o *deployOptions) validate() error {
	if o.name == "" {
		return fmt.Errorf("error setting deploy configuration, missing kubemq cluster name")
	}
	if o.namespace == "" {
		return fmt.Errorf("error setting deploy configuration, missing kubemq cluster namespace")
	}
	if o.tag == "" {
		return fmt.Errorf("error setting deploy configuration, missing kubemq cluster docker image tag")
	}

	if err := o.service.validate(); err != nil {
		return err
	}
	if err := o.security.validate(); err != nil {
		return err
	}

	if err := o.authentication.validate(); err != nil {
		return err
	}
	if err := o.authorization.validate(); err != nil {
		return err
	}
	if err := o.gateway.validate(); err != nil {
		return err
	}
	if err := o.resources.validate(); err != nil {
		return err
	}

	if err := o.nodeSelectors.validate(); err != nil {
		return err
	}
	if err := o.healthProbe.validate(); err != nil {
		return err
	}
	return nil
}

func (o *deployOptions) complete() error {
	if o.licenseDataFile != "" {
		data, err := ioutil.ReadFile(o.licenseDataFile)
		if err != nil {
			return fmt.Errorf("error loading license file data: %s", err.Error())
		}
		o.licenseData = string(data)
	}
	if err := o.service.complete(); err != nil {
		return err
	}
	if err := o.security.complete(); err != nil {
		return err
	}

	if err := o.authentication.complete(); err != nil {
		return err
	}
	if err := o.authorization.complete(); err != nil {
		return err
	}
	if err := o.gateway.complete(); err != nil {
		return err
	}
	if err := o.resources.complete(); err != nil {
		return err
	}

	if err := o.nodeSelectors.complete(); err != nil {
		return err
	}
	if err := o.healthProbe.complete(); err != nil {
		return err
	}
	return nil
}

func (o *deployOptions) getConfig() *deployment.KubeMQManifestConfig {
	id := uuid.New().String()
	config := deployment.DefaultKubeMQManifestConfig(id, o.name, o.namespace)
	config.StatefulSet.SetImageTag(o.tag)
	config.StatefulSet.SetReplicas(int(o.replicas))
	config.StatefulSet.SetVolume(int(o.volume))
	if o.token != "" {
		config.SetConfigMapValues(o.name, "KUBEMQ_TOKEN", o.token)
	}
	if o.licenseData != "" {
		config.SetSecretStringValues(o.name, "LICENSE_KEY_DATA", o.licenseData)
	}

	o.service.setConfig(config)
	o.security.setConfig(config)
	o.authentication.setConfig(config)
	o.authorization.setConfig(config)
	o.gateway.setConfig(config)
	o.resources.setConfig(config)
	o.nodeSelectors.setConfig(config)
	o.healthProbe.setConfig(config)
	return config
}
