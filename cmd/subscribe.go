// Copyright © 2018 NAME HERE <EMAIL ADDRESS>

package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/kubemq-io/kubetools/transport"

	"github.com/kubemq-io/kubetools/transport/grpc"
	"github.com/kubemq-io/kubetools/transport/option"
	"github.com/kubemq-io/kubetools/transport/rest"

	"github.com/spf13/cobra"
)

var subscribeTransport string
var subscribeGroup string
var subscribeTimeout int
var subscribeCmd = &cobra.Command{
	Use:     "subscribe",
	Aliases: []string{"sub"},
	Short:   "subscribe to events / events_store / commands / queries",
}

var subscribeEventsCmd = &cobra.Command{
	Use:     "event",
	Aliases: []string{"e"},
	Short:   "subscribe to an events channel",
	Run: func(cmd *cobra.Command, args []string) {
		runSubscribe(args, "events")
	},
}

var subscribeEventsStoreCmd = &cobra.Command{
	Use:     "event_store",
	Aliases: []string{"es"},
	Short:   "subscribe to an event_store channel",
	Run: func(cmd *cobra.Command, args []string) {
		runSubscribe(args, "events_store")
	},
}

var subscribeCommandsCmd = &cobra.Command{
	Use:     "command",
	Aliases: []string{"c"},
	Short:   "subscribe to a command to a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runSubscribe(args, "commands")
	},
}

var subscribeQueriesCmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q"},
	Short:   "subscribe to a query channel",
	Run: func(cmd *cobra.Command, args []string) {
		runSubscribe(args, "queries")
	},
}

func runSubscribe(args []string, kind string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	if len(args) != 1 {
		fmt.Println("invalid args, should be a channel name")
		return
	}
	client, err := getSubscribeClient(ctx)
	if err != nil {
		log.Fatalf("error on create send client: %s", err.Error())
		return
	}
	msgCh := make(chan *transport.Message, 1)
	errCh := make(chan error, 1)

	switch kind {
	case "events":
		log.Printf("subscribe to %s channel\n", args[0])
		err = client.ReceiveEvent(ctx, args[0], subscribeGroup, msgCh, errCh)
	case "events_store":
		log.Printf("subscribe to %s channel\n", args[0])
		err = client.ReceiveEventStore(ctx, args[0], subscribeGroup, msgCh, errCh)
	case "commands":
		log.Printf("subscribe to %s channel\n", args[0])
		err = client.ReceiveCommand(ctx, args[0], subscribeGroup, msgCh, errCh)
	case "queries":
		log.Printf("subscribe to %s channel\n", args[0])
		err = client.ReceiveQuery(ctx, args[0], subscribeGroup, msgCh, errCh)
	}
	for {
		select {
		case msg := <-msgCh:
			log.Printf("Message:\n\t%s\n", msg.Payload)
		case err := <-errCh:
			log.Printf("Error:\n\t%s\n", err.Error())
		case <-ctx.Done():
			return
		}
	}
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
	subscribeCmd.AddCommand(subscribeEventsCmd, subscribeEventsStoreCmd, subscribeCommandsCmd, subscribeQueriesCmd)
	subscribeCmd.PersistentFlags().StringVarP(&subscribeTransport, "subscribeTransport", "t", "grpc", "set transport type, grpc or rest")
	subscribeCmd.PersistentFlags().StringVarP(&subscribeGroup, "subscribeGroup", "g", "", "set optional group for a channel")
	subscribeCommandsCmd.PersistentFlags().IntVarP(&subscribeTimeout, "subscribeTimout", "o", 10000, "set command timeout in MSec")
	subscribeQueriesCmd.PersistentFlags().IntVarP(&subscribeTimeout, "subscribeTimout", "o", 10000, "set query timeout in MSec")

}

func getSubscribeClient(ctx context.Context) (transport.Transport, error) {
	switch subscribeTransport {
	case "grpc":
		for _, conn := range cfg.Connections {
			if conn.Kind == option.ConnectionTypeGrpc {
				return grpc.New(ctx, conn)
			}
		}
	case "rest":
		for _, conn := range cfg.Connections {
			if conn.Kind == option.ConnectionTypeRest {
				return rest.New(ctx, conn)
			}
		}

	}
	return nil, errors.New("invalid transport type")
}
