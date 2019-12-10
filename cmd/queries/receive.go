package queries

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"time"
)

type QueriesReceiveOptions struct {
	cfg          *config.Config
	transport    string
	channel      string
	group        string
	autoResponse bool
}

var queriesReceiveExamples = `
	# Receive 'queries'  from a 'queries' channel (blocks until next message)
	kubemqctl queries receive some-channel

	# Receive 'queries' from a 'queries' channel with group(blocks until next message)
	kubemqctl queries receive some-channel -g G1
`
var queriesReceiveLong = `Receive (Subscribe) command allows to receive a message from a 'queries' channel and response with appropriate reply`
var queriesReceiveShort = `Receive a message from a 'queries' channel`

func NewCmdQueriesReceive(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &QueriesReceiveOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "receive",
		Aliases: []string{"r", "rec", "subscribe", "sub"},
		Short:   queriesReceiveShort,
		Long:    queriesReceiveLong,
		Example: queriesReceiveExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().StringVarP(&o.group, "group", "g", "", "set 'queries' channel consumer group (load balancing)")
	cmd.PersistentFlags().BoolVarP(&o.autoResponse, "auto-response", "a", false, "set auto response executed query")
	return cmd
}

func (o *QueriesReceiveOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]
		return nil
	}
	return fmt.Errorf("missing channel argument")
}

func (o *QueriesReceiveOptions) Validate() error {
	return nil
}

func (o *QueriesReceiveOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	errChan := make(chan error, 1)
	queriesChan, err := client.SubscribeToQueries(ctx, o.channel, o.group, errChan)

	if err != nil {
		utils.Println(fmt.Errorf("receive 'queries' messages, %s", err.Error()).Error())
	}
	for {
		utils.Println("waiting for the next query message...")
		select {
		case err := <-errChan:
			return fmt.Errorf("server disconnected with error: %s", err.Error())

		case query, opened := <-queriesChan:
			if !opened {
				utils.Println("server disconnected")
				return nil
			}
			fmt.Fprintf(w, "[channel: %s]\t[id: %s]\t[metadata: %s]\t[body: %s]\n", query.Channel, query.Id, query.Metadata, query.Body)
			w.Flush()
			if o.autoResponse {
				err = client.R().SetRequestId(query.Id).SetExecutedAt(time.Now()).SetResponseTo(query.ResponseTo).SetBody([]byte("executed your query")).Send(ctx)
				if err != nil {
					return err
				}
				utils.Println("auto execution sent executed response ")
				continue
			}
			var isExecuted bool
			prompt := &survey.Confirm{
				Renderer: survey.Renderer{},
				Message:  "Set executed ?",
				Help:     "",
			}
			err := survey.AskOne(prompt, &isExecuted)

			if err != nil {
				return err
			}
			if isExecuted {

				respBody := ""
				prompt := &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "Set response message",
					Default:  "response-to",
					Help:     "",
				}
				err := survey.AskOne(prompt, &respBody)
				if err != nil {
					return err
				}
				err = client.R().SetRequestId(query.Id).SetExecutedAt(time.Now()).SetResponseTo(query.ResponseTo).SetBody([]byte(respBody)).Send(ctx)
				if err != nil {
					return err
				}
				continue
			}
			err = client.R().SetRequestId(query.Id).SetError(fmt.Errorf("query not executed")).SetResponseTo(query.ResponseTo).Send(ctx)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}

}
