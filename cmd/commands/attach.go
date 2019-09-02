package commands

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/attach"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
)

type CommandsAttachOptions struct {
	cfg       *config.Config
	include   []string
	exclude   []string
	resources []string
}

var commandsAttachExamples = `
	# attach to all commands channels and output running messages
	kubetools commands attach *
	
	# attach to some-commands commands channel and output running messages
	kubetools commands attach some-commands

	# attach to some-commands1 and some-commands2 commands channels and output running messages
	kubetools commands attach some-commands1 some-commands2 

	# attach to some-commands commands channel and output running messages filter by include regex (some*)
	kubetools commands attach some-commands -i some*

	# attach to some-commands commands channel and output running messages filter by exclude regex (not-some*)
	kubetools commands attach some-commands -e not-some*
`
var commandsAttachLong = `Attach to commands channels and display channel content`
var commandsAttachShort = `Attach to commands channels`

func NewCmdCommandsAttach(cfg *config.Config) *cobra.Command {
	o := &CommandsAttachOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "attach",
		Aliases: []string{"a", "att", "at"},
		Short:   commandsAttachShort,
		Long:    commandsAttachLong,
		Example: commandsAttachExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().StringArrayVarP(&o.include, "include", "i", []string{}, "Set (regex) strings to include")
	cmd.PersistentFlags().StringArrayVarP(&o.exclude, "exclude", "e", []string{}, "Set (regex) strings to exclude")
	return cmd
}

func (o *CommandsAttachOptions) Complete(args []string, transport string) error {

	if len(args) == 0 {
		return fmt.Errorf("missing channel argument")
	}

	for _, a := range args {
		rsc := fmt.Sprintf("commands/%s", a)
		o.resources = append(o.resources, rsc)
		utils.Printlnf("adding '%s' to attach list", a)
	}
	return nil
}

func (o *CommandsAttachOptions) Validate() error {
	return nil
}

func (o *CommandsAttachOptions) Run(ctx context.Context) error {
	err := attach.Run(ctx, o.cfg, o.resources, o.include, o.exclude)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
