package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type deployAuthorizationOptions struct {
	enabled        bool
	policyData     string
	policyFilename string
	url            string
	autoReload     int32
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
	cmd.PersistentFlags().Int32VarP(&o.autoReload, "authorization-auto-reload", "", 0, "set authorization auto policy loading time interval in minutes")
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
			return fmt.Errorf("error loading authorization public key data: %s", err.Error())
		}
		o.policyData = string(data)
	}
	return nil
}

func (o *deployAuthorizationOptions) setConfig(deployment *cluster.KubemqCluster) *deployAuthorizationOptions {
	if !o.enabled {
		return o
	}
	deployment.Spec.Authorization = &cluster.AuthorizationConfig{
		Policy:     o.policyData,
		Url:        o.url,
		AutoReload: o.autoReload,
	}
	return o
}
