// Copyright © 2018 NAME HERE <EMAIL ADDRESS>

package cmd

import (
	"context"
	"errors"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/kubemq-io/kubetools/transport/option"
	"github.com/spf13/cobra"
)

var rpcTransport string
var rpcRecGroup string
var rpcSendTimeout int
var rpcCmd = &cobra.Command{
	Use:     "rpc",
	Aliases: []string{"r"},
	Short:   "send and receive rpc commands and queries",
}

var rpcRecCmd = &cobra.Command{
	Use:     "receive",
	Aliases: []string{"rec", "r"},
	Short:   "receive commands or queries",
}
var rpcSendCmd = &cobra.Command{
	Use:     "send",
	Aliases: []string{"s"},
	Short:   "send commands and queries",
}

var rpcRecCommandsCmd = &cobra.Command{
	Use:     "command",
	Aliases: []string{"c"},
	Short:   "subscribe to receive commands from a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runRpc(args, "sub_commands")
	},
}

var rpcRecQueriesCmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q"},
	Short:   "subscribe to receive queries from channel",
	Run: func(cmd *cobra.Command, args []string) {
		runRpc(args, "sub_queries")
	},
}

var rpcSendCommandsCmd = &cobra.Command{
	Use:     "command",
	Aliases: []string{"c"},
	Short:   "send command to a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runRpc(args, "send_command")
	},
}

var rpcSendQueriesCmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q"},
	Short:   "send query to a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runRpc(args, "send_query")
	},
}

func runRpc(args []string, kind string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	client, err := getRpcClient(ctx)
	if err != nil {
		log.Fatalf("error on create send client: %s", err.Error())
		return
	}
	errCh := make(chan error, 1)
	switch kind {
	case "send_command":
		if len(args) != 2 {
			log.Fatal("invalid args, should be a channel name and message body")
			return
		}

		res, err := client.C().SetChannel(args[0]).SetBody([]byte(args[1])).SetTimeout(time.Duration(rpcSendTimeout) * 1000).Send(ctx)
		if err != nil {
			log.Printf("error sending command: %s", err.Error())
			return
		}
		log.Printf("response received:\n executed at: %s\n", res.ExecutedAt.String())
		return

	case "send_query":
		if len(args) != 2 {
			log.Fatal("invalid args, should be a channel name and message body")
			return
		}

		res, err := client.Q().SetChannel(args[0]).SetBody([]byte(args[1])).SetTimeout(time.Duration(rpcSendTimeout) * 1000).Send(ctx)
		if err != nil {
			log.Printf("error sending query: %s", err.Error())
			return
		}

		if res != nil {
			if res.Error != "" {
				log.Printf("query error: %s", res.Error)
				return
			}
			log.Printf("response body: %s", res.Body)
		}

	case "sub_commands":
		if len(args) != 1 {
			log.Fatal("invalid args, should be a channel name")
			return
		}
		log.Printf("subscribe to %s channel\n", args[0])
		commandsCh, err := client.SubscribeToCommands(ctx, args[0], rpcRecGroup, errCh)
		if err != nil {
			log.Printf("Error:\n\t%s\n", err.Error())
			return
		}
		for {
			select {
			case command, more := <-commandsCh:
				if !more {
					return
				}

				log.Printf("Command Received:\n\t%s\n", command.Body)
				log.Println("Sending Response")
				err = client.R().SetRequestId(command.Id).
					SetExecutedAt(time.Now()).
					SetResponseTo(command.ResponseTo).
					Send(ctx)
			case err := <-errCh:
				log.Printf("Error:\n\t%s\n", err.Error())
			case <-ctx.Done():
				return
			}
		}
	case "sub_queries":
		if len(args) != 1 {
			log.Fatal("invalid args, should be a channel name")
			return
		}
		log.Printf("subscribe to %s channel\n", args[0])
		queriesCh, err := client.SubscribeToQueries(ctx, args[0], rpcRecGroup, errCh)
		if err != nil {
			log.Printf("Error:\n\t%s\n", err.Error())
			return
		}
		for {
			select {
			case query, more := <-queriesCh:
				if !more {
					return
				}
				log.Printf("Query Received:\n\t%s\n", query.Body)
				log.Println("Sending Response")
				err = client.R().SetRequestId(query.Id).
					SetExecutedAt(time.Now()).
					SetResponseTo(query.ResponseTo).
					SetBody([]byte("query received and handled")).
					Send(ctx)
			case err := <-errCh:
				log.Printf("Error:\n\t%s\n", err.Error())
			case <-ctx.Done():
				return
			}
		}

	}

}

func init() {
	rootCmd.AddCommand(rpcCmd)
	rpcCmd.AddCommand(rpcSendCmd, rpcRecCmd)
	rpcSendCmd.AddCommand(rpcSendCommandsCmd, rpcSendQueriesCmd)
	rpcRecCmd.AddCommand(rpcRecCommandsCmd, rpcRecQueriesCmd)
	rpcCmd.PersistentFlags().StringVarP(&rpcTransport, "rpcTransport", "t", "grpc", "set transport type, grpc or rest")
	rpcRecCommandsCmd.PersistentFlags().StringVarP(&rpcRecGroup, "rpcRecGroup", "g", "", "set optional group for a channel")
	rpcRecQueriesCmd.PersistentFlags().StringVarP(&rpcRecGroup, "rpcRecGroup", "g", "", "set optional group for a channel")
	rpcSendCommandsCmd.PersistentFlags().IntVarP(&rpcSendTimeout, "rpcSendTimout", "o", 10000, "set command timeout in msec")
	rpcSendQueriesCmd.PersistentFlags().IntVarP(&rpcSendTimeout, "rpcSendTimout", "o", 10000, "set query timeout in msec")
}

func getRpcClient(ctx context.Context) (*kubemq.Client, error) {
	switch rpcTransport {
	case "grpc":
		for _, conn := range cfg.Connections {
			if conn.Kind == option.ConnectionTypeGrpc {
				client, err := kubemq.NewClient(ctx,
					kubemq.WithAddress(conn.Host, conn.Port),
					kubemq.WithClientId(uuid.New().String()),
					kubemq.WithTransportType(kubemq.TransportTypeGRPC))

				return client, err
			}
		}
	case "rest":
		for _, conn := range cfg.Connections {
			if conn.Kind == option.ConnectionTypeRest {
				client, err := kubemq.NewClient(ctx,
					kubemq.WithUri(conn.Uri()),
					kubemq.WithClientId(uuid.New().String()),
					kubemq.WithTransportType(kubemq.TransportTypeRest))
				return client, err
			}
		}

	}
	return nil, errors.New("invalid transport type")
}
