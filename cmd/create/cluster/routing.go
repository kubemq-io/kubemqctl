package cluster

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var defaultRoutingConfig = &deployRoutingOptions{
	routingData:     "",
	routingFilename: "",
	url:             "",
	autoReload:      0,
}

type deployRoutingOptions struct {
	routingData     string
	routingFilename string
	url             string
	autoReload      int32
}

func setRoutingConfig(cmd *cobra.Command) *deployRoutingOptions {
	o := &deployRoutingOptions{
		routingData:     "",
		routingFilename: "",
		url:             "",
		autoReload:      0,
	}
	cmd.PersistentFlags().StringVarP(&o.routingData, "routing-data", "", "", "set routing data")
	cmd.PersistentFlags().StringVarP(&o.routingFilename, "routing-filename", "", "", "set routing filename")
	cmd.PersistentFlags().StringVarP(&o.url, "routing-url", "", "", "set routing loading url")
	cmd.PersistentFlags().Int32VarP(&o.autoReload, "routing-auto-reload", "", 0, "set routing auto loading time interval in minutes")
	return o
}

func (o *deployRoutingOptions) validate() error {
	return nil
}

func (o *deployRoutingOptions) complete() error {

	if o.routingFilename != "" {
		data, err := ioutil.ReadFile(o.routingFilename)
		if err != nil {
			return fmt.Errorf("error loading routing public key data: %s", err.Error())
		}
		o.routingData = string(data)
	}
	return nil
}

func (o *deployRoutingOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployRoutingOptions {
	if isDefault(o, defaultRoutingConfig) {
		return o
	}
	deployment.Spec.Routing = &kubemqcluster.RoutingConfig{
		Data:       o.routingData,
		Url:        o.url,
		AutoReload: o.autoReload,
	}
	return o
}
