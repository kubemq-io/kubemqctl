package create

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
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
	tollerations    *deployTolerationOptions
	affinity        *deployAffinityOptions
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
		tollerations:    defaultTolerationOptions(cmd),
		affinity:        defaultAffinityOptions(cmd),
	}
	cmd.PersistentFlags().StringVarP(&o.configFilename, "config-file", "", "", "set kubemq config file")
	cmd.PersistentFlags().StringVarP(&o.name, "name", "", "kubemq-cluster", "set kubemq cluster name")
	cmd.PersistentFlags().StringVarP(&o.namespace, "namespace", "", "kubemq", "set kubemq cluster namespace")
	cmd.PersistentFlags().StringVarP(&o.token, "token", "", "", "set kubemq token")
	cmd.PersistentFlags().StringVarP(&o.licenseData, "license-data", "", "", "set license data")
	cmd.PersistentFlags().StringVarP(&o.licenseDataFile, "license-data-file", "", "", "set license data filename")
	cmd.PersistentFlags().StringVarP(&o.tag, "tag", "", "latest", "set kubemq docker image tag")
	cmd.PersistentFlags().UintVarP(&o.volume, "volume", "", 0, "set persistence volume claim size")
	cmd.PersistentFlags().UintVarP(&o.volume, "replicas", "", 3, "set replicas")
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

	if err := o.tollerations.validate(); err != nil {
		return err
	}
	if err := o.affinity.validate(); err != nil {
		return err
	}
	return nil
}

func (o *deployOptions) complete() error {
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

	if err := o.tollerations.complete(); err != nil {
		return err
	}
	if err := o.affinity.complete(); err != nil {
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
	o.service.setConfig(config)
	o.security.setConfig(config)
	o.authentication.setConfig(config)
	o.authorization.setConfig(config)
	o.gateway.setConfig(config)
	o.resources.setConfig(config)
	o.nodeSelectors.setConfig(config)
	o.tollerations.setConfig(config)
	o.affinity.setConfig(config)
	return config
}
