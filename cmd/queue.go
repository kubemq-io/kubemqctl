// Copyright © 2018 NAME HERE <EMAIL ADDRESS>

package cmd

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubetools/transport"
	"github.com/kubemq-io/kubetools/transport/option"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var queueTransport string
var sendQueueExpiration int
var sendQueueDelay int
var receiveQueueMessages int
var receiveQueueWaitSecond int

var queueCmd = &cobra.Command{
	Use:     "queue",
	Aliases: []string{"q"},
	Short:   "send and receive queue messages",
}

var queueSendCmd = &cobra.Command{
	Use:     "send",
	Aliases: []string{"s"},
	Short:   "send message to a queue",
	Run: func(cmd *cobra.Command, args []string) {
		runQueue(args, "send")
	},
}

var queueReceiveCmd = &cobra.Command{
	Use:     "receive",
	Aliases: []string{"r"},
	Short:   "receive messages from a queue",
	Run: func(cmd *cobra.Command, args []string) {
		runQueue(args, "receive")
	},
}

var queuePeakCmd = &cobra.Command{
	Use:     "peak",
	Aliases: []string{"p"},
	Short:   "peak messages from a queue",
	Run: func(cmd *cobra.Command, args []string) {
		runQueue(args, "peak")
	},
}

var queueAckAllCmd = &cobra.Command{
	Use:     "ack",
	Aliases: []string{"a"},
	Short:   "acl all messages in a queue",
	Run: func(cmd *cobra.Command, args []string) {
		runQueue(args, "ack")
	},
}

func runQueue(args []string, kind string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	client, err := getQueueClient(ctx)
	if err != nil {
		log.Fatalf("error on create send client: %s", err.Error())
		return
	}

	switch kind {
	case "send":
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
		res, err := client.QM().
			SetChannel(args[0]).
			SetBody(msg.Marshal()).
			SetPolicyExpirationSeconds(sendQueueExpiration).
			SetPolicyDelaySeconds(sendQueueDelay).
			SetId(msg.Id).
			Send(ctx)
		if err != nil {
			log.Printf("error sending queue message: %s", err.Error())
			return
		}

		if res != nil {
			if res.IsError {
				log.Printf("send queue message error: %s", res.Error)
				return
			}
			log.Printf("queue message sent at: %s", time.Unix(0, res.SentAt))
		}

	case "receive":
		if len(args) != 1 {
			log.Fatal("invalid args, should be a channel name")
			return
		}
		log.Printf("receiving queue message from %s channel\n", args[0])
		res, err := client.RQM().
			SetChannel(args[0]).
			SetWaitTimeSeconds(receiveQueueWaitSecond).
			SetMaxNumberOfMessages(receiveQueueMessages).
			Send(ctx)
		if err != nil {
			log.Printf("Error:\n\t%s\n", err.Error())
			return
		}
		if res.IsError {
			log.Printf("Error:\n\t%s\n", res.Error)
			return
		}
		log.Printf("received %d messages, %d messages Expired \n", res.MessagesReceived, res.MessagesExpired)
		for _, item := range res.Messages {
			msg, err := transport.Unmarshal(item.Body)
			if err != nil {
				log.Printf("Error:\n\t%s\n", err.Error())
				return
			}
			log.Printf("queue message received:\n\t%s\n", msg.Payload)
		}

	case "peak":
		if len(args) != 1 {
			log.Fatal("invalid args, should be a channel name")
			return
		}
		log.Printf("peaking queue message from %s channel\n", args[0])
		res, err := client.RQM().
			SetChannel(args[0]).
			SetIsPeak(true).
			SetWaitTimeSeconds(receiveQueueWaitSecond).
			SetMaxNumberOfMessages(receiveQueueMessages).
			Send(ctx)
		if err != nil {
			log.Printf("Error:\n\t%s\n", err.Error())
			return
		}
		if res.IsError {
			log.Printf("Error:\n\t%s\n", res.Error)
			return
		}
		log.Printf("peaked %d messages, %d messages Expired \n", res.MessagesReceived, res.MessagesExpired)
		for _, item := range res.Messages {
			msg, err := transport.Unmarshal(item.Body)
			if err != nil {
				log.Printf("Error:\n\t%s\n", err.Error())
				return
			}
			log.Printf("queue message peaked:\n\t%s\n", msg.Payload)
		}
	case "ack":
		if len(args) != 1 {
			log.Fatal("invalid args, should be a channel name")
			return
		}
		log.Printf("accl all messages in queue %s\n", args[0])
		res, err := client.NewAckAllQueueMessagesRequest().
			SetChannel(args[0]).
			SetWaitTimeSeconds(receiveQueueWaitSecond).
			Send(ctx)
		if err != nil {
			log.Printf("Error:\n\t%s\n", err.Error())
			return
		}
		if res.IsError {
			log.Printf("Error:\n\t%s\n", res.Error)
			return
		}
		log.Printf("ack all messages in a qeueu, affected messages %d", res.AffectedMessages)
	}

}

func init() {
	rootCmd.AddCommand(queueCmd)
	queueCmd.AddCommand(queueSendCmd, queueReceiveCmd, queuePeakCmd, queueAckAllCmd)
	queueCmd.PersistentFlags().StringVarP(&queueTransport, "queueTransport", "t", "grpc", "set transport type, grpc or rest")
	queueSendCmd.PersistentFlags().IntVarP(&sendQueueExpiration, "sendExpiration", "e", 0, "set queue message expiration seconds")
	queueSendCmd.PersistentFlags().IntVarP(&sendQueueDelay, "sendDelay", "d", 0, "set queue message send delay seconds")
	queueReceiveCmd.PersistentFlags().IntVarP(&receiveQueueMessages, "receiveMessages", "i", 1, "set how many messages we want to get from queue")
	queueReceiveCmd.PersistentFlags().IntVarP(&receiveQueueWaitSecond, "receiveWait", "w", 2, "set how many seconds to wait for queue messages")
	queuePeakCmd.PersistentFlags().IntVarP(&receiveQueueMessages, "receiveMessages", "i", 1, "set how many messages we peak to get from queue")
	queuePeakCmd.PersistentFlags().IntVarP(&receiveQueueWaitSecond, "receiveWait", "w", 2, "set how many seconds to peak for queue messages")
	queueAckAllCmd.PersistentFlags().IntVarP(&receiveQueueWaitSecond, "receiveWait", "w", 2, "set how many seconds wait to ack all messages in queue")
}

func getQueueClient(ctx context.Context) (*kubemq.Client, error) {
	switch queueTransport {
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
