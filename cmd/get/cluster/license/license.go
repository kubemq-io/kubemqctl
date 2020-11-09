package license

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/api"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	prefix = "-----BEGIN KUBEMQ KEY-----"
	suffix = "-----END KUBEMQ KEY-----"
)

type getLicenseOptions struct {
	licenseToken string
}

var licenseExamples = `
 	# Get license with token 
	kubemqctl get cluster license -t some-token
`
var licenseLong = `License command allows to fetach KubeMQ license file from token`
var licenseShort = `License command allows to fetach KubeMQ license file from token`

func NewCmdLicense(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &getLicenseOptions{
		licenseToken: "",
	}
	cmd := &cobra.Command{

		Use:     "license",
		Aliases: []string{"lic"},
		Short:   licenseShort,
		Long:    licenseLong,
		Example: licenseExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringVarP(&o.licenseToken, "token", "t", "", "set license token")

	return cmd
}

func (o *getLicenseOptions) Validate() error {
	return nil
}
func (o *getLicenseOptions) Complete(args []string) error {
	return nil
}
func split(v string, n int) string {
	newStr := ""
	buf := []byte(v)
	for i := 0; i < len(v); i++ {
		if i > 0 && i%n == 0 {
			newStr = newStr + "\n" + string(buf[i])
		} else {
			newStr = newStr + string(buf[i])
		}
	}
	return newStr
}

func (o *getLicenseOptions) Run(ctx context.Context) error {
	if o.licenseToken != "" {
		utils.Printlnf("fetching license for token %s :", o.licenseToken)
		data, err := api.GetLicenseDataByToken(o.licenseToken)
		if err != nil {
			utils.PrintlnfNoTitle(", error: %s ", err.Error())
		} else {
			key := fmt.Sprintf("%s\n%s\n%s", prefix, split(data, 64), suffix)
			utils.PrintlnfNoTitle(key)

		}

	} else {
		return fmt.Errorf("no token was provided")
	}
	return nil
}
