package commands

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/targets"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"time"
)

type CommandsSendOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	body      string
	metadata  string
	timeout   int
	fileName  bool
	build     bool
}

var commandsSendExamples = `
	# Send command to a 'commands' channel
	kubemqctl commands send some-channel some-command
	
	# Send command to a 'commands' channel with metadata
	kubemqctl commands send some-channel some-body -m some-metadata
	
	# Send command to a 'commands' channel with 120 seconds timeout
	kubemqctl commands send some-channel some-body -o 120
`
var commandsSendLong = `Send command allow to send messages to 'commands' channel with an option to set command time-out`
var commandsSendShort = `Send messages to 'commands' channel command`

func NewCmdCommandsSend(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &CommandsSendOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "send",
		Aliases: []string{"s"},
		Short:   commandsSendShort,
		Long:    commandsSendLong,
		Example: commandsSendExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringVarP(&o.metadata, "metadata", "m", "", "Set metadata body")
	cmd.PersistentFlags().IntVarP(&o.timeout, "timeout", "o", 30, "Set command timeout")
	cmd.PersistentFlags().BoolVarP(&o.fileName, "file", "f", false, "set load body from file")
	cmd.PersistentFlags().BoolVarP(&o.build, "build", "b", false, "build kubemq targets request")
	return cmd
}

func (o *CommandsSendOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]

	} else {
		return fmt.Errorf("missing channel argument")
	}
	if o.build {
		data, err := targets.BuildRequest()
		if err != nil {
			return err
		}
		o.body = string(data)
		return nil
	}
	if o.fileName {
		data, err := targets.BuildFile()
		if err != nil {
			return err
		}
		o.body = string(data)
	} else {
		if len(args) >= 2 {
			o.body = args[1]
		} else {
			return fmt.Errorf("missing body argument")
		}
	}
	return nil
}

func (o *CommandsSendOptions) Validate() error {
	return nil
}

func (o *CommandsSendOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubemqClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())
	}

	defer func() {
		client.Close()
	}()

	msg := client.C().
		SetChannel(o.channel).
		SetId(uuid.New().String()).
		SetBody([]byte(o.body)).
		SetMetadata(o.metadata).
		SetTimeout(time.Duration(o.timeout) * time.Second)
	fmt.Println("Sending Command:")
	printCommand(msg)
	res, err := msg.Send(ctx)
	if err != nil {
		return fmt.Errorf("sending commands body, %s", err.Error())
	}
	fmt.Println("Getting  Command Response:")
	printCommandResponse(res)
	return nil
}
