package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type deployAuthorizationOptions struct {
	enabled        bool
	policyData     string
	policyFilename string
	url            string
	autoReload     int
}

func defaultAuthorizationConfig(cmd *cobra.Command) *deployAuthorizationOptions {
	o := &deployAuthorizationOptions{
		enabled:        false,
		policyData:     "",
		policyFilename: "",
		url:            "",
		autoReload:     0,
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "authorization-enabled", "", false, "enable authorization configuration")
	cmd.PersistentFlags().StringVarP(&o.policyData, "authorization-policy-data", "", "", "set authorization policy data")
	cmd.PersistentFlags().StringVarP(&o.policyFilename, "authorization-policy-file", "", "", "set authorization policy filename")
	cmd.PersistentFlags().StringVarP(&o.url, "authorization-url", "", "", "set authorization policy loading url")
	cmd.PersistentFlags().IntVarP(&o.autoReload, "authorization-auto-reload", "", 0, "set authorization auto policy loading time interval in minutes")
	return o
}

func (o *deployAuthorizationOptions) validate() error {
	if !o.enabled {
		return nil
	}
	if o.policyData == "" && o.url == "" {
		return fmt.Errorf("error setting authorization configuration, no ploicy data or policy url was set")
	}
	return nil
}

func (o *deployAuthorizationOptions) complete() error {
	if !o.enabled {
		return nil
	}
	if o.policyFilename != "" {
		data, err := ioutil.ReadFile(o.policyFilename)
		if err != nil {
			return fmt.Errorf("error loading authentication public key data: %s", err.Error())
		}
		o.policyData = string(data)
	}
	return nil
}

func (o *deployAuthorizationOptions) setConfig(config *deployment.KubeMQManifestConfig) *deployAuthorizationOptions {
	if !o.enabled {
		return o
	}
	cmConfig, ok := config.ConfigMaps[config.Name]
	if ok {
		cmConfig.SetVariable("AUTHORIZATION_ENABLE", "true").
			SetVariable("AUTHORIZATION_POLICY_DATA", o.policyData).
			SetVariable("AUTHORIZATION_URL", o.url).
			SetVariable("AUTHORIZATION_AUTO_RELOAD", fmt.Sprintf("%d", o.autoReload))
	}
	return o
}
