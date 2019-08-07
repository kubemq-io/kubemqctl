// Copyright © 2018 NAME HERE <EMAIL ADDRESS>

package cmd

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go"
	"log"
	"time"

	"github.com/kubemq-io/kubetools/transport"
	"github.com/kubemq-io/kubetools/transport/option"
	"github.com/spf13/cobra"
)

var pubsubTransport string
var pubsubRecGroup string
var pubsubSendTimeout int
var pubsubCmd = &cobra.Command{
	Use:     "pubsub",
	Aliases: []string{"p", "ps"},
	Short:   "send and receive pub/sub real-time and persistent events",
}

var pubsubRecCmd = &cobra.Command{
	Use:     "receive",
	Aliases: []string{"rec", "r"},
	Short:   "receive pub/sub real-time and persistent events",
}

var pubsubSendCmd = &cobra.Command{
	Use:     "send",
	Aliases: []string{"s"},
	Short:   "send pub/sub real-time and persistent events",
}

var pubsubRecEventsCmd = &cobra.Command{
	Use:     "events",
	Aliases: []string{"e"},
	Short:   "subscribe to receive real-time events from a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runPubSub(args, "sub_events")
	},
}

var pubsubRecEventsStoreCmd = &cobra.Command{
	Use:     "events_store",
	Aliases: []string{"es"},
	Short:   "subscribe to receive persistent events from channel",
	Run: func(cmd *cobra.Command, args []string) {
		runPubSub(args, "sub_events_store")
	},
}

var pubsubSendEventsCmd = &cobra.Command{
	Use:     "events",
	Aliases: []string{"e"},
	Short:   "send real-time event to a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runPubSub(args, "send_event")
	},
}

var pubsubSendEventsStoreCmd = &cobra.Command{
	Use:     "events_store",
	Aliases: []string{"es"},
	Short:   "send persistent event to a channel",
	Run: func(cmd *cobra.Command, args []string) {
		runPubSub(args, "send_event_store")
	},
}

func runPubSub(args []string, kind string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	client, err := getpubsubClient(ctx)
	if err != nil {
		log.Fatalf("error on create send client: %s", err.Error())
		return
	}
	errCh := make(chan error, 1)
	switch kind {
	case "send_event":
		if len(args) != 2 {
			log.Fatal("invalid args, should be a channel name and message body")
			return
		}
		msg := &transport.Message{
			Id:       uuid.New().String(),
			SendTime: time.Now().Unix(),
			Payload:  []byte(args[1]),
		}
		msg.SetSendTime()

		err := client.E().SetChannel(args[0]).SetBody(msg.Marshal()).SetId(msg.Id).Send(ctx)
		if err != nil {
			log.Printf("error sending event: %s", err.Error())
		}
	case "send_event_store":
		if len(args) != 2 {
			log.Fatal("invalid args, should be a channel name and message body")
			return
		}
		msg := &transport.Message{
			Id:       uuid.New().String(),
			SendTime: time.Now().Unix(),
			Payload:  []byte(args[1]),
		}
		msg.SetSendTime()
		res, err := client.ES().SetChannel(args[0]).SetBody(msg.Marshal()).SetId(msg.Id).Send(ctx)

		if err != nil {
			log.Printf("error sending event_store: %s", err.Error())
		}
		if res.Err != nil {
			log.Printf("error sending event_store: %s", res.Err)
		}
	case "sub_events":
		if len(args) != 1 {
			log.Fatal("invalid args, should be a channel name")
			return
		}
		log.Printf("subscribe to %s channel\n", args[0])
		eventsCh, err := client.SubscribeToEvents(ctx, args[0], pubsubRecGroup, errCh)
		if err != nil {
			log.Printf("Error:\n\t%s\n", err.Error())
			return
		}
		for {
			select {
			case event, more := <-eventsCh:
				if !more {
					return
				}
				msg, err := transport.Unmarshal(event.Body)
				if err != nil {
					log.Printf("Error:\n\t%s\n", err.Error())
					return
				}
				msg.SetReceiveTime()
				log.Printf("Message:\n\t%s\n", msg.Payload)
			case err := <-errCh:
				log.Printf("Error:\n\t%s\n", err.Error())
			case <-ctx.Done():
				return
			}
		}

	case "sub_events_store":
		if len(args) != 1 {
			log.Fatal("invalid args, should be a channel name")
			return
		}
		log.Printf("subscribe to %s channel\n", args[0])
		eventsStoreCh, err := client.SubscribeToEventsStore(ctx, args[0], pubsubRecGroup, errCh, kubemq.StartFromNewEvents())
		if err != nil {
			log.Printf("Error:\n\t%s\n", err.Error())
			return
		}
		for {
			select {
			case eventStore, more := <-eventsStoreCh:
				if !more {
					return
				}
				msg, err := transport.Unmarshal(eventStore.Body)
				if err != nil {
					log.Printf("Error:\n\t%s\n", err.Error())
					return
				}
				msg.SetReceiveTime()
				log.Printf("Message:\n\t%s\n", msg.Payload)
			case err := <-errCh:
				log.Printf("Error:\n\t%s\n", err.Error())
			case <-ctx.Done():
				return
			}
		}

	}

}

func init() {
	rootCmd.AddCommand(pubsubCmd)
	pubsubCmd.AddCommand(pubsubSendCmd, pubsubRecCmd)
	pubsubSendCmd.AddCommand(pubsubSendEventsCmd, pubsubSendEventsStoreCmd)
	pubsubRecCmd.AddCommand(pubsubRecEventsCmd, pubsubRecEventsStoreCmd)
	pubsubCmd.PersistentFlags().StringVarP(&pubsubTransport, "pubsubTransport", "t", "grpc", "set transport type, grpc or rest")
	pubsubRecEventsCmd.PersistentFlags().StringVarP(&pubsubRecGroup, "pubsubRecGroup", "g", "", "set optional group for a channel")
	pubsubRecEventsStoreCmd.PersistentFlags().StringVarP(&pubsubRecGroup, "pubsubRecGroup", "g", "", "set optional group for a channel")
}

func getpubsubClient(ctx context.Context) (*kubemq.Client, error) {
	switch pubsubTransport {
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
