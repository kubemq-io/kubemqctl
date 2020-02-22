package cluster

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type deployTlsOptions struct {
	enabled      bool
	certData     string
	certFilename string
	keyData      string
	keyFilename  string
	caData       string
	caFilename   string
}

func setTolsConfig(cmd *cobra.Command) *deployTlsOptions {
	o := &deployTlsOptions{
		enabled:      false,
		certData:     "",
		certFilename: "",
		keyData:      "",
		keyFilename:  "",
		caData:       "",
		caFilename:   "",
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "tls-enabled", "", false, "enable tls tls configuration")
	cmd.PersistentFlags().StringVarP(&o.certData, "tls-cert-data", "", "", "set tls certificate data")
	cmd.PersistentFlags().StringVarP(&o.certFilename, "tls-cert-file", "", "", "set tls certificate filename")
	cmd.PersistentFlags().StringVarP(&o.keyData, "tls-key-data", "", "", "set tls key data")
	cmd.PersistentFlags().StringVarP(&o.keyFilename, "tls-key-file", "", "", "set tls key filename")
	cmd.PersistentFlags().StringVarP(&o.caData, "tls-ca-data", "", "", "set tls ca certificate data")
	cmd.PersistentFlags().StringVarP(&o.caFilename, "tls-ca-file", "", "", "set tls ca certificate filename")
	return o
}

func (o *deployTlsOptions) validate() error {
	if !o.enabled {
		return nil
	}
	if o.certData == "" {
		return fmt.Errorf("error setting tls configuration, missing certifcate data")
	}
	if o.keyData == "" {
		return fmt.Errorf("error setting tls configuration, missing key data")
	}
	return nil
}

func (o *deployTlsOptions) complete() error {
	if !o.enabled {
		return nil
	}
	if o.certFilename != "" {
		data, err := ioutil.ReadFile(o.certFilename)
		if err != nil {
			return fmt.Errorf("error loading tls certifcate data: %s", err.Error())
		}
		o.certData = string(data)
	}
	if o.keyFilename != "" {
		data, err := ioutil.ReadFile(o.keyFilename)
		if err != nil {
			return fmt.Errorf("error loading tls key data: %s", err.Error())
		}
		o.keyData = string(data)
	}

	if o.caFilename != "" {
		data, err := ioutil.ReadFile(o.caFilename)
		if err != nil {
			return fmt.Errorf("error loading tls ca tls data: %s", err.Error())
		}
		o.caData = string(data)
	}

	return nil
}

func (o *deployTlsOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployTlsOptions {
	if !o.enabled {
		return o
	}
	deployment.Spec.Tls = &kubemqcluster.TlsConfig{
		Cert: o.certData,
		Key:  o.keyData,
		Ca:   o.caData,
	}

	return o
}
