package cluster

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/api"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var defaultLicenseConfig = &deployLicenseOptions{
	licenseData:     "",
	licenseFilename: "",
	licenseToken:    "",
}

type deployLicenseOptions struct {
	licenseData     string
	licenseFilename string
	licenseToken    string
}

func setLicenseConfig(cmd *cobra.Command) *deployLicenseOptions {
	o := &deployLicenseOptions{
		licenseData:     "",
		licenseFilename: "",
		licenseToken:    "",
	}
	cmd.PersistentFlags().StringVarP(&o.licenseData, "license-data", "", "", "set license data")
	cmd.PersistentFlags().StringVarP(&o.licenseFilename, "license-filename", "", "", "set license filename")
	cmd.PersistentFlags().StringVarP(&o.licenseToken, "license-token", "t", "", "set license token")

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
	} else {
		if o.licenseToken != "" {
			utils.Printf("fetching license for token %s ", o.licenseToken)
			data, err := api.GetLicenseDataByToken(o.licenseToken)
			if err != nil {
				utils.PrintlnfNoTitle(", error: %s ", err.Error())
			} else {
				utils.PrintlnfNoTitle(" completed")
				o.licenseData = data
			}
			if o.licenseData == "" {
				utils.Printlnf("no valid license data received")
			}
		}
	}
	return nil
}

func (o *deployLicenseOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployLicenseOptions {
	if isDefault(o, defaultLicenseConfig) {
		return o
	}
	deployment.Spec.License = &kubemqcluster.LicenseConfig{
		Data:  o.licenseData,
		Token: o.licenseToken,
	}
	return o
}
