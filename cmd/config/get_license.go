package config

import (
	"context"
	"errors"
	"github.com/kubemq-io/kubemqctl/pkg/api"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type GetLicenseOptions struct {
	token    string
	fileName string
}

var getLicenseExamples = `
	# Run Kubemqctl get license data and save to file
	# kubemqctl config get_license  -t b3d360ssc-93ef-4395-bba3-13131eb27d5e -f kubemq.lic
`
var getLicenseLong = `Get License command get license data for offline use of KubeMQ`
var getLicenseShort = `Run Kubemqctl config get_license for retrieving license file for offline use`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdGetLicense(ctx context.Context) *cobra.Command {
	o := &GetLicenseOptions{}
	cmd := &cobra.Command{
		Use:     "get_license",
		Aliases: []string{"lic"},
		Short:   getLicenseShort,
		Long:    getLicenseLong,
		Example: getLicenseExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args))
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringVarP(&o.token, "token", "t", "", "set KubeMQ Token")
	cmd.PersistentFlags().StringVarP(&o.fileName, "file", "f", "", "set export file name")
	return cmd
}

func (o *GetLicenseOptions) Complete(args []string) error {

	return nil
}

func (o *GetLicenseOptions) Validate() error {
	if o.token == "" {
		return errors.New("missing KubeMQ token")
	}
	return nil
}

func (o *GetLicenseOptions) Run(ctx context.Context) error {
	utils.Printlnf("fetching license data for token: %s", o.token)
	lic, err := api.GetLicenseData(o.token, "")
	if err != nil {
		return err
	}
	if o.fileName != "" {
		err := api.SaveToFile(o.fileName, lic)
		if err != nil {
			return err
		}
		utils.PrintlnfNoTitle("License data saved to %s.", o.fileName)
		return nil
	}
	utils.PrintlnfNoTitle("License data:\n%s", lic)
	return nil
}
