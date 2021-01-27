package cluster

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var defaultLicenseConfig = &deployLicenseOptions{
	licenseData:     "",
	licenseFilename: "",
}

type deployLicenseOptions struct {
	licenseData     string
	licenseFilename string
}

func setLicenseConfig(cmd *cobra.Command) *deployLicenseOptions {
	o := &deployLicenseOptions{
		licenseData:     "",
		licenseFilename: "",
	}
	cmd.PersistentFlags().StringVarP(&o.licenseData, "license-data", "", "", "set license data")
	cmd.PersistentFlags().StringVarP(&o.licenseFilename, "license-file", "", "", "set license file")

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

func (o *deployLicenseOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployLicenseOptions {
	if isDefault(o, defaultLicenseConfig) {
		return o
	}
	deployment.Spec.License = o.licenseData
	return o
}
