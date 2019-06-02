// Copyright © 2018 NAME HERE <EMAIL ADDRESS>

package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kubemq-io/kubetools/transport"

	"github.com/kubemq-io/kubetools/transport/grpc"
	"github.com/kubemq-io/kubetools/transport/option"
	"github.com/kubemq-io/kubetools/transport/rest"

	"github.com/spf13/cobra"
)

var sendTransport string
var sendTimeout int
var sendCmd = &cobra.Command{
	Use:     "send",
	Aliases: []string{"s"},
	Short:   "send event / event_store / command / query",
}

var sendEventsCmd = &cobra.Command{
	Use:     "event",
	Aliases: []string{"e"},
	Short:   "send event to a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runSend(args, "events")
	},
}

var sendEventsStoreCmd = &cobra.Command{
	Use:     "event_store",
	Aliases: []string{"es"},
	Short:   "send event_store to a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runSend(args, "events_store")
	},
}

var sendCommandsCmd = &cobra.Command{
	Use:     "command",
	Aliases: []string{"c"},
	Short:   "send command to a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runSend(args, "commands")
	},
}

var sendQueriesCmd = &cobra.Command{
	Use:     "query",
	Aliases: []string{"q"},
	Short:   "send query to a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runSend(args, "queries")
	},
}

func runSend(args []string, kind string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	if len(args) != 2 {
		fmt.Println("invalid args, should be channel name and body payload")
		return
	}
	client, err := getSendClient(ctx)
	if err != nil {
		log.Fatalf("error on create send client: %s", err.Error())
		return
	}
	msg := &transport.Message{
		Id:       uuid.New().String(),
		SendTime: time.Now().Unix(),
		Payload:  []byte(args[1]),
	}

	switch kind {
	case "events":
		err = client.SendEvent(ctx, args[0], msg)
		if err != nil {
			log.Printf("error sending event: %s", err.Error())
		}

	case "events_store":
		err = client.SendEventStore(ctx, args[0], msg)
		if err != nil {
			log.Printf("error sending event_store: %s", err.Error())
		}

	case "commands":
		err = client.SendCommand(ctx, args[0], msg, time.Duration(sendTimeout)*1000)
		if err != nil {
			log.Printf("error sending command: %s", err.Error())
		}
		return

	case "queries":
		resp, err := client.SendQuery(ctx, args[0], msg, time.Duration(sendTimeout)*1000)
		if err != nil {
			log.Printf("error sending query: %s", err.Error())
			return
		}
		if resp != nil {
			log.Printf("response body: %s", resp.Payload)
		}
	}
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.AddCommand(sendEventsCmd, sendEventsStoreCmd, sendCommandsCmd, sendQueriesCmd)
	sendCmd.PersistentFlags().StringVarP(&sendTransport, "sendTransport", "t", "grpc", "set transport type, grpc or rest")
	sendCommandsCmd.PersistentFlags().IntVarP(&sendTimeout, "sendTimout", "o", 1000, "set command timeout in MSec")
	sendQueriesCmd.PersistentFlags().IntVarP(&sendTimeout, "sendTimout", "o", 1000, "set query timeout in MSec")

}

func getSendClient(ctx context.Context) (transport.Transport, error) {
	switch sendTransport {
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
