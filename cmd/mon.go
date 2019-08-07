// Copyright © 2018 NAME HERE <EMAIL ADDRESS>

package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/gorilla/websocket"

	"github.com/spf13/cobra"
)

// monCmd represents the pub command
var monCmd = &cobra.Command{
	Use:     "mon",
	Aliases: []string{"m"},
	Short:   "monitor messages/requests channels",
}

// monCmd represents the pub command
var monEventsCmd = &cobra.Command{
	Use:     "events",
	Aliases: []string{"e"},
	Short:   "monitor events channels",
	Run: func(cmd *cobra.Command, args []string) {
		runMon(args, "events")
	},
}

// monCmd represents the pub command
var monEventsStoreCmd = &cobra.Command{
	Use:     "events_store",
	Aliases: []string{"es"},
	Short:   "monitor events store channels",
	Run: func(cmd *cobra.Command, args []string) {
		runMon(args, "events_store")
	},
}

var monCommandsCmd = &cobra.Command{
	Use:     "commands",
	Aliases: []string{"c"},
	Short:   "monitor commands channels",
	Run: func(cmd *cobra.Command, args []string) {
		runMon(args, "commands")
	},
}

var monQueriesCmd = &cobra.Command{
	Use:     "queries",
	Aliases: []string{"q"},
	Short:   "monitor query channels",
	Run: func(cmd *cobra.Command, args []string) {
		runMon(args, "queries")
	},
}

var monQueuesCmd = &cobra.Command{
	Use:     "queue",
	Aliases: []string{"qu"},
	Short:   "monitor queue channels",
	Run: func(cmd *cobra.Command, args []string) {
		runMon(args, "queue")
	},
}

func runMon(args []string, kind string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	if len(args) != 1 {
		log.Error("invalid args, should be channel name")
		return
	}

	uri := fmt.Sprintf("%s/attach?node=some-node&channel=%s&kind=%s", cfg.MonitorAddress, args[0], kind)
	rxChan := make(chan string, 10)
	txChan := make(chan string, 10)
	ready := make(chan struct{})
	errCh := make(chan error, 10)
	log.Print(fmt.Sprintf("connecting to %s ...", uri))
	go runWebsocketClientReaderWriter(ctx, uri, rxChan, txChan, ready, errCh)
	<-ready
	txChan <- "start"

	for {
		select {
		case msg := <-rxChan:
			log.Print(msg)
		case <-ctx.Done():
			return
		case <-errCh:
			return
		default:
			time.Sleep(1 * time.Millisecond)
		}
	}

}

func init() {
	rootCmd.AddCommand(monCmd)
	monCmd.AddCommand(monEventsCmd, monEventsStoreCmd, monCommandsCmd, monQueriesCmd, monQueuesCmd)

}

var retries = 10
var retryInterval = 100 * time.Millisecond

func runWebsocketClientReaderWriter(ctx context.Context, uri string, chRead chan string, chWrite chan string, ready chan struct{}, errCh chan error) {
	var c *websocket.Conn
	for i := 0; i < retries; i++ {
		conn, res, err := websocket.DefaultDialer.Dial(uri, nil)
		if err != nil {
			buf := make([]byte, 1024)
			if res != nil {
				n, _ := res.Body.Read(buf)
				//	errCh <- errors.New(string(buf[:n]))
				log.Print(string(buf[:n]))
			} else {
				log.Error(err)
			}
			time.Sleep(1 * time.Second)
		} else {
			c = conn
			break
		}
	}
	if c == nil {
		os.Exit(1)
	} else {
		defer c.Close()
	}

	ready <- struct{}{}
	go func() {
		for {
			select {
			case msg := <-chWrite:
				err := c.WriteMessage(1, []byte(msg))
				if err != nil {
					log.Error(err)
					errCh <- err
					return
				}
			case <-ctx.Done():
				c.Close()
				return
				//default:
				//	time.Sleep(100 * time.Millisecond)
			}

		}

	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Error(err)
			errCh <- err
			return
		} else {
			chRead <- string(message)
		}
	}

}
