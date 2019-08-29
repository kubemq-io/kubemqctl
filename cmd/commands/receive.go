package commands

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/kubemq"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"time"
)

type CommandsReceiveOptions struct {
	cfg          *config.Config
	transport    string
	channel      string
	group        string
	autoResponse bool
}

var commandsReceiveExamples = `
	# Receive commands from a commands channel (blocks until next message)
	kubetools commands receive some-channel

	# Receive commands from a commands channel with group(blocks until next message)
	kubetools commands receive some-channel -g G1
`
var commandsReceiveLong = `receive a message from a commands`
var commandsReceiveShort = `receive a message from a commands`

func NewCmdCommandsReceive(cfg *config.Config, opts *CommandsOptions) *cobra.Command {
	o := &CommandsReceiveOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "receive",
		Aliases: []string{"r", "rec"},
		Short:   commandsReceiveShort,
		Long:    commandsReceiveLong,
		Example: commandsReceiveExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args, opts.transport))
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().StringVarP(&o.group, "group", "g", "", "set group")
	cmd.PersistentFlags().BoolVarP(&o.autoResponse, "auto-response", "a", false, "set auto response executed command")
	return cmd
}

func (o *CommandsReceiveOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]
		return nil
	}
	return fmt.Errorf("missing channel argument")
}

func (o *CommandsReceiveOptions) Validate() error {
	return nil
}

func (o *CommandsReceiveOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	errChan := make(chan error, 1)
	commandsChan, err := client.SubscribeToCommands(ctx, o.channel, o.group, errChan)

	if err != nil {
		utils.Println(fmt.Errorf("receive commands messages, %s", err.Error()).Error())
	}
	for {
		utils.Println("waiting for the next command message...")
		select {
		case cmd := <-commandsChan:
			fmt.Fprintf(w, "[channel: %s]\t[id: %s]\t[metadata: %s]\t[body: %s]\n", cmd.Channel, cmd.Id, cmd.Metadata, cmd.Body)
			w.Flush()
			if o.autoResponse {
				err = client.R().SetRequestId(cmd.Id).SetExecutedAt(time.Now()).SetResponseTo(cmd.ResponseTo).SetBody([]byte("executed your command")).Send(ctx)
				if err != nil {
					return err
				}
				utils.Println("auto execution sent executed response ")
				continue
			}
			var isExecuted bool
			prompt := &survey.Confirm{
				Renderer: survey.Renderer{},
				Message:  "Set executed?",
				Help:     "",
			}
			err := survey.AskOne(prompt, &isExecuted)

			if err != nil {
				return err
			}
			if isExecuted {
				err = client.R().SetRequestId(cmd.Id).SetExecutedAt(time.Now()).SetResponseTo(cmd.ResponseTo).SetBody([]byte("executed your command")).Send(ctx)
				if err != nil {
					return err
				}
				continue
			}
			err = client.R().SetRequestId(cmd.Id).SetError(fmt.Errorf("commnad not executed")).SetResponseTo(cmd.ResponseTo).Send(ctx)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}

}
