package version

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type VersionOptions struct {
	version *string
}

var versionExamples = `
 	# Show Kubemqctl version
	kubemqctl version
`
var versionLong = `Show Kubemqctl version`
var versionShort = `Show Kubemqctl version`

func NewCmdVersion(version *string) *cobra.Command {
	o := VersionOptions{
		version: version,
	}
	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver", "v"},
		Short:   versionShort,
		Long:    versionLong,
		Example: versionExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *VersionOptions) Complete(args []string) error {
	return nil
}

func (o *VersionOptions) Validate() error {
	return nil
}

func (o *VersionOptions) Run(ctx context.Context) error {

	utils.Printlnf("Kubemqctl version %s", *o.version)
	return nil
}
