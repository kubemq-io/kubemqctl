package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type deploySecurityOptions struct {
	enabled      bool
	certData     string
	certFilename string
	keyData      string
	keyFilename  string
	caData       string
	caFilename   string
}

func defaultSecurityConfig(cmd *cobra.Command) *deploySecurityOptions {
	o := &deploySecurityOptions{
		enabled:      false,
		certData:     "",
		certFilename: "",
		keyData:      "",
		keyFilename:  "",
		caData:       "",
		caFilename:   "",
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "security-enabled", "", false, "enable tls security configuration")
	cmd.PersistentFlags().StringVarP(&o.certData, "security-cert-data", "", "", "set tls certificate data")
	cmd.PersistentFlags().StringVarP(&o.certFilename, "security-cert-file", "", "", "set tls certificate filename")
	cmd.PersistentFlags().StringVarP(&o.keyData, "security-key-data", "", "", "set tls key data")
	cmd.PersistentFlags().StringVarP(&o.keyFilename, "security-key-file", "", "", "set tls key filename")
	cmd.PersistentFlags().StringVarP(&o.caData, "security-ca-data", "", "", "set tls ca certificate data")
	cmd.PersistentFlags().StringVarP(&o.caFilename, "security-ca-file", "", "", "set tls ca certificate filename")
	return o
}

func (o *deploySecurityOptions) validate() error {
	if !o.enabled {
		return nil
	}
	if o.certData == "" {
		return fmt.Errorf("error setting security configuration, missing certifcate data")
	}
	if o.keyData == "" {
		return fmt.Errorf("error setting security configuration, missing key data")
	}
	return nil
}

func (o *deploySecurityOptions) complete() error {
	if !o.enabled {
		return nil
	}
	if o.certFilename != "" {
		data, err := ioutil.ReadFile(o.certFilename)
		if err != nil {
			return fmt.Errorf("error loading security certifcate data: %s", err.Error())
		}
		o.certData = string(data)
	}
	if o.keyFilename != "" {
		data, err := ioutil.ReadFile(o.keyFilename)
		if err != nil {
			return fmt.Errorf("error loading security key data: %s", err.Error())
		}
		o.keyData = string(data)
	}

	if o.caFilename != "" {
		data, err := ioutil.ReadFile(o.caFilename)
		if err != nil {
			return fmt.Errorf("error loading security ca security data: %s", err.Error())
		}
		o.caData = string(data)
	}

	return nil
}

func (o *deploySecurityOptions) setConfig(config *deployment.KubeMQManifestConfig) *deploySecurityOptions {
	if !o.enabled {
		return o
	}
	secConfig, ok := config.Secrets[config.Name]
	if ok {
		secConfig.SetVariable("SECURITY_CERT_DATA", o.certData).
			SetVariable("SECURITY_KEY_DATA", o.keyData).
			SetVariable("SECURITY_CA_DATA", o.caData)
	}
	return o
}
