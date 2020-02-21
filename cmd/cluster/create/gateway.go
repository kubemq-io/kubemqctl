package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type deployGatewayOptions struct {
	enabled      bool
	remotes      []string
	port         int32
	certData     string
	certFilename string
	keyData      string
	keyFilename  string
	caData       string
	caFilename   string
}

func setGatewayOptions(cmd *cobra.Command) *deployGatewayOptions {
	o := &deployGatewayOptions{
		enabled:      false,
		remotes:      nil,
		port:         7000,
		certData:     "",
		certFilename: "",
		keyData:      "",
		keyFilename:  "",
		caData:       "",
		caFilename:   "",
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "gateway-enabled", "", false, "enable gateway configuration")
	cmd.PersistentFlags().StringArrayVarP(&o.remotes, "gateway-remotes", "", nil, "set tls certificate data for remote gateway")
	cmd.PersistentFlags().Int32VarP(&o.port, "gateway-port", "", 7000, "set gateway listen port value")
	cmd.PersistentFlags().StringVarP(&o.certData, "gateway-cert-data", "", "", "set tls certificate data for remote gateway")
	cmd.PersistentFlags().StringVarP(&o.certFilename, "gateway-cert-file", "", "", "set tls certificate filename for remote gateway")
	cmd.PersistentFlags().StringVarP(&o.keyData, "gateway-key-data", "", "", "set tls key data for remote gateway")
	cmd.PersistentFlags().StringVarP(&o.keyFilename, "gateway-key-file", "", "", "set tls key filename for remote gateway ")
	cmd.PersistentFlags().StringVarP(&o.caData, "gateway-ca-data", "", "", "set tls ca certificate data for remote gateway")
	cmd.PersistentFlags().StringVarP(&o.caFilename, "gateway-ca-file", "", "", "set tls ca certificate filename for remote gateway")
	return o
}

func (o *deployGatewayOptions) validate() error {
	if !o.enabled {
		return nil
	}
	if o.remotes == nil {
		return fmt.Errorf("error setting gateway configuration, missing remotes gateway data")
	}
	if o.port == 0 || o.port > 65535 {
		return fmt.Errorf("invalid gateway value: %d", o.port)
	}
	return nil
}

func (o *deployGatewayOptions) complete() error {
	if !o.enabled {
		return nil
	}
	if !o.enabled {
		return nil
	}
	if o.certFilename != "" {
		data, err := ioutil.ReadFile(o.certFilename)
		if err != nil {
			return fmt.Errorf("error loading gateway certifcate data: %s", err.Error())
		}
		o.certData = string(data)
	}
	if o.keyFilename != "" {
		data, err := ioutil.ReadFile(o.keyFilename)
		if err != nil {
			return fmt.Errorf("error loading gateway key data: %s", err.Error())
		}
		o.keyData = string(data)
	}

	if o.caFilename != "" {
		data, err := ioutil.ReadFile(o.caFilename)
		if err != nil {
			return fmt.Errorf("error loading gateway ca security data: %s", err.Error())
		}
		o.caData = string(data)
	}

	return nil
}

func (o *deployGatewayOptions) setConfig(deployment *cluster.KubemqCluster) *deployGatewayOptions {
	if !o.enabled {
		return o
	}
	deployment.Spec.Gateways = &cluster.GatewayConfig{
		Remotes: o.remotes,
		Port:    o.port,
		Cert:    o.certData,
		Key:     o.keyData,
		Ca:      o.caData,
	}
	return o
}
