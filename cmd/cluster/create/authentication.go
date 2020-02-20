package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type deployAuthenticationOptions struct {
	enabled           bool
	publicKeyData     string
	publicKeyFilename string
	publicKeyType     string
}

func defaultAuthenticationOptions(cmd *cobra.Command) *deployAuthenticationOptions {
	o := &deployAuthenticationOptions{
		enabled:           false,
		publicKeyData:     "",
		publicKeyFilename: "",
		publicKeyType:     "",
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "authentication-enabled", "", false, "enable authentication configuration")
	cmd.PersistentFlags().StringVarP(&o.publicKeyData, "authentication-public-key-data", "", "", "set authentication public key data")
	cmd.PersistentFlags().StringVarP(&o.publicKeyFilename, "authentication-public-key-file", "", "", "set authentication public key filename")
	cmd.PersistentFlags().StringVarP(&o.publicKeyType, "authentication-public-key-type", "", "", "set authentication public key type")
	return o
}

func (o *deployAuthenticationOptions) validate() error {
	if !o.enabled {
		return nil
	}
	if o.publicKeyData == "" {
		return fmt.Errorf("error setting authentication configuration, missing publilc key data")
	}

	if o.publicKeyType == "" {
		return fmt.Errorf("error setting authentication configuration, missing public key type")
	}
	return nil
}

func (o *deployAuthenticationOptions) complete() error {
	if !o.enabled {
		return nil
	}
	if o.publicKeyFilename != "" {
		data, err := ioutil.ReadFile(o.publicKeyFilename)
		if err != nil {
			return fmt.Errorf("error loading authentication public key data: %s", err.Error())
		}
		o.publicKeyData = string(data)
	}
	return nil
}

func (o *deployAuthenticationOptions) setConfig(deployment *cluster.KubemqCluster) *deployAuthenticationOptions {
	if !o.enabled {
		return o
	}
	deployment.Spec.Authentication = &cluster.AuthenticationConfig{
		Key:  o.publicKeyData,
		Type: o.publicKeyType,
	}
	return o
}
