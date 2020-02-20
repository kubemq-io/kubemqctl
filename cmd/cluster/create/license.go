package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type deployLicenseOptions struct {
	licenseData     string
	licenseFilename string
	licenseToken    string
}

func defaultLicenseConfig(cmd *cobra.Command) *deployLicenseOptions {
	o := &deployLicenseOptions{
		licenseData:     "",
		licenseFilename: "",
		licenseToken:    "",
	}
	cmd.PersistentFlags().StringVarP(&o.licenseData, "license-data", "", "", "set license data")
	cmd.PersistentFlags().StringVarP(&o.licenseFilename, "license-filename", "", "", "set license filename")
	cmd.PersistentFlags().StringVarP(&o.licenseToken, "license-token", "", "", "set license token")

	return o
}

func (o *deployLicenseOptions) validate() error {
	return nil
}
func (o *deployLicenseOptions) complete() error {
	if o.licenseFilename != "" {
		data, err := ioutil.ReadFile(o.licenseFilename)
		if err != nil {
			return fmt.Errorf("error loading license file data: %s", err.Error())
		}
		o.licenseData = string(data)
	}
	return nil
}

func (o *deployLicenseOptions) setConfig(deployment *cluster.KubemqCluster) *deployLicenseOptions {
	deployment.Spec.License = &cluster.LicenseConfig{
		Data:  o.licenseData,
		Token: o.licenseToken,
	}
	return o
}
