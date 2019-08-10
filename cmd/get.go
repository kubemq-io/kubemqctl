package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"os"
	"text/tabwriter"
	"time"

	"log"

	"github.com/kubemq-io/kubetools/transport/rest"

	"github.com/spf13/cobra"
)

type Queues struct {
	Now    time.Time `json:"now"`
	Total  int       `json:"total"`
	Queues []*Queue  `json:"queues"`
}

type Queue struct {
	Name          string    `json:"name"`
	Messages      int64     `json:"messages"`
	Bytes         int64     `json:"bytes"`
	FirstSequence int64     `json:"first_sequence"`
	LastSequence  int64     `json:"last_sequence"`
	Clients       []*Client `json:"clients"`
}

type Client struct {
	ClientId         string `json:"client_id"`
	Active           bool   `json:"active"`
	LastSequenceSent int64  `json:"last_sequence_sent"`
	IsStalled        bool   `json:"is_stalled"`
	Pending          int64  `json:"pending"`
}

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   "Call kubemq get resources endpoint",
	Aliases: []string{"g"},
	Long:    `Return list of resources`,
}

var getQueuesCmd = &cobra.Command{
	Use:     "queues",
	Short:   "Call kubemq get list of queues",
	Aliases: []string{"qu", "queue"},
	Long:    `Return list of Queues`,
	Run: func(cmd *cobra.Command, args []string) {
		runGetQueues()
	},
}

var getEventsStoreCmd = &cobra.Command{
	Use:     "events_stores",
	Short:   "Call kubemq get list of events store channels",
	Aliases: []string{"es", "events_store", "event_stores", "event_store"},
	Long:    `Return list of Events Store channels`,
	Run: func(cmd *cobra.Command, args []string) {
		runGetEventsStore()
	},
}

func (q *Queues) String() string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(q)
	if err != nil {
		return err.Error()
	}
	return buffer.String()
}

func (q *Queues) printChannelsTab() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "CHANNELS:\n")
	fmt.Fprintln(w, "NAME\tCLIENTS\tMESSAGES\tBYTES\tFIRST_SEQUENCE\tLAST_SEQUENCE")
	for _, q := range q.Queues {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%d\n", q.Name, len(q.Clients), q.Messages, q.Bytes, q.FirstSequence, q.LastSequence)
	}
	fmt.Fprintf(w, "\nTOTAL CHANNELS:\t%d\n", q.Total)
	w.Flush()
}
func (q *Queues) printClientsTab() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "\nCLIENTS:\n")
	fmt.Fprintln(w, "CLIENT_ID\tCHANNEL\tACTIVE\tLAST_SENT\tPENDING\tSTALLED")
	cnt := 0
	for _, q := range q.Queues {
		for _, c := range q.Clients {
			if c.ClientId == "" {
				c.ClientId = "N/A"
			}
			cnt++
			fmt.Fprintf(w, "%s\t%s\t%t\t%d\t%d\t%t\n", c.ClientId, q.Name, c.Active, c.LastSequenceSent, c.Pending, c.IsStalled)
		}

	}
	fmt.Fprintf(w, "\nTOTAL CLIENTS:\t%d\n", cnt)
	w.Flush()
}

func runGetQueues() {
	resp := &rest.Response{}
	q := &Queues{}

	r, err := resty.R().SetResult(resp).SetError(resp).Get(cfg.StatsAddress + "/queues")
	if err != nil {
		log.Fatal(err)
	}
	if !r.IsSuccess() {
		log.Fatal("not available in current KubeMQ version, consider upgrade KubeMQ version")
	}
	if resp.Error {
		log.Fatal(resp.ErrorString)
	}
	err = json.Unmarshal(resp.Data, q)
	if err != nil {
		log.Fatal(err)
	}
	q.printChannelsTab()
	q.printClientsTab()

}

func runGetEventsStore() {
	resp := &rest.Response{}
	q := &Queues{}

	r, err := resty.R().SetResult(resp).SetError(resp).Get(cfg.StatsAddress + "/events_stores")
	if err != nil {
		log.Fatal(err)
	}
	if !r.IsSuccess() {
		log.Fatal("not available in current KubeMQ version, consider upgrade KubeMQ version")
	}
	if resp.Error {
		log.Fatal(resp.ErrorString)
	}
	err = json.Unmarshal(resp.Data, q)
	if err != nil {
		log.Fatal(err)
	}
	q.printChannelsTab()
	q.printClientsTab()

}
func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getQueuesCmd)
	getCmd.AddCommand(getEventsStoreCmd)
}
