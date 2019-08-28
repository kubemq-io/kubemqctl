package queue

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/kubemq"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

type QueueStreamOptions struct {
	cfg        *config.Config
	transport  string
	channel    string
	visibility int
	wait       int
}

var queueStreamExamples = `
	# stream queue message in transaction mode
	kubetools queue stream some-channel

	# stream queue message in transaction mode with visibility set to 120 seconds and wait time of 180 seconds
	kubetools queue stream some-channel -v 120 -w 180
`
var queueStreamLong = `receive a message from a queue`
var queueStreamShort = `receive a message from a queue`

func NewCmdQueueStream(cfg *config.Config, opts *QueueOptions) *cobra.Command {
	o := &QueueStreamOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "stream",
		Aliases: []string{"st"},
		Short:   queueStreamShort,
		Long:    queueStreamLong,
		Example: queueStreamExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args, opts.transport))
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().IntVarP(&o.visibility, "visibility", "v", 30, "set initial visibility seconds")
	cmd.PersistentFlags().IntVarP(&o.wait, "wait", "w", 60, "set how many seconds to wait for queue messages")

	return cmd
}

func (o *QueueStreamOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]
		return nil
	}
	return fmt.Errorf("missing channel argument")
}

func (o *QueueStreamOptions) Validate() error {
	return nil
}

func (o *QueueStreamOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()
	for {
		stream := client.NewStreamQueueMessage().SetChannel(o.channel)
		utils.Printlnf("waiting for the message in the queue: (waiting for %d seconds, visibility set to %d seconds)", o.wait, o.visibility)
		msg, err := stream.Next(ctx, int32(o.visibility), int32(o.wait))
		if err != nil {
			return err
		}
		utils.Printlnf("[channel: %s] [client id: %s] -> {id: %s, metadata: %s, body: %s}", msg.Channel, msg.ClientId, msg.Id, msg.Metadata, msg.Body)
		action, result, err := o.prompt()
		if err != nil {
			return err
		}
		switch action {
		case "ack":
			err := msg.Ack()
			if err != nil {
				return err
			}
		case "reject":
			err := msg.Reject()
			if err != nil {
				return err
			}
		case "extend visibility":
			val, err := strconv.Atoi(result)
			if err != nil {
				return err
			}
			err = msg.ExtendVisibility(int32(val))
			if err != nil {
				return err
			}
		case "resend to another queue":
			err = msg.Resend(result)
			if err != nil {
				return err
			}
		case "ack and send new message":
			pair := strings.Split(result, ",")
			if len(pair) != 2 {
				return fmt.Errorf("invalid queue-name,message-body format")
			}
			newMessage := client.QM().SetChannel(pair[0]).SetBody([]byte(pair[1]))
			err := stream.ResendWithNewMessage(newMessage)
			if err != nil {
				return err
			}
		case "done":
			err := msg.Ack()
			if err != nil {
				return err
			}

			return nil
		}
		utils.Println("sent.")
	}

}
func (o *QueueStreamOptions) prompt() (string, string, error) {
	action := ""
	prompt := &survey.Select{
		Message: "What next:",
		Options: []string{"ack", "reject", "extend visibility", "resend to another queue", "ack and send new message", "done"},
	}
	err := survey.AskOne(prompt, &action)
	if err != nil {
		return "", "", err
	}
	switch action {
	case "ack", "reject", "done":
		return action, "", nil
	case "extend visibility":
		visibility := ""
		prompt := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "How long to extend visibility",
			Default:  "60",
			Help:     "in seconds",
		}
		err := survey.AskOne(prompt, &visibility)
		if err != nil {
			return "", "", err
		}
		return action, visibility, nil
	case "resend to another queue":
		queueName := ""
		prompt := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "New queue name:",
			Default:  "new-queue",
			Help:     "",
		}
		err := survey.AskOne(prompt, &queueName, survey.WithValidator(survey.MinLength(1)))
		if err != nil {
			return "", "", err
		}
		return action, queueName, nil
	case "ack and send new message":
		newMessage := ""
		prompt := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "New Message:",
			Default:  "new-queue,new-message",
			Help:     "format queue-name,message-body ",
		}
		err := survey.AskOne(prompt, &newMessage, survey.WithValidator(survey.MinLength(1)))
		if err != nil {
			return "", "", err
		}
		return action, newMessage, nil
	}
	return "", "", fmt.Errorf("invalid input")
}
