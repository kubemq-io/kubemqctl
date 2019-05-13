package cmd

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/go-resty/resty"

	"log"

	"github.com/kubemq-io/kubetools/transport/rest"

	"github.com/spf13/cobra"
)

type Health []struct {
	Host    string    `json:"host"`
	UtcTime time.Time `json:"utc_time"`
	Grpc    struct {
		Connections struct {
			Total                int `json:"total"`
			EventsSenders        int `json:"events_senders"`
			EventsStreamSenders  int `json:"events_stream_senders"`
			EventsReceivers      int `json:"events_receivers"`
			EventsStoreReceivers int `json:"events_store_receivers"`
			RequestsSenders      int `json:"requests_senders"`
			ResponsesSenders     int `json:"responses_senders"`
			CommandsReceivers    int `json:"commands_receivers"`
			QueriesReceivers     int `json:"queries_receivers"`
		} `json:"connections"`
		Traffic struct {
			SentEvents          int `json:"sent_events"`
			ReceivedEvents      int `json:"received_events"`
			SentRequests        int `json:"sent_requests"`
			SentError           int `json:"sent_error"`
			SentResponses       int `json:"sent_responses"`
			ReceivedRequests    int `json:"received_requests"`
			SentEventsVol       int `json:"sent_events_vol"`
			ReceivedEventsVol   int `json:"received_events_vol"`
			SentRequestsVol     int `json:"sent_requests_vol"`
			SentErrorsVol       int `json:"sent_errors_vol"`
			SentResponsesVol    int `json:"sent_responses_vol"`
			ReceivedRequestsVol int `json:"received_requests_vol"`
			TotalMessages       int `json:"total_messages"`
			TotalVolume         int `json:"total_volume"`
		} `json:"traffic"`
	} `json:"grpc"`
	Rest struct {
		Connections struct {
			Total                int `json:"total"`
			EventsSenders        int `json:"events_senders"`
			EventsStreamSenders  int `json:"events_stream_senders"`
			EventsReceivers      int `json:"events_receivers"`
			EventsStoreReceivers int `json:"events_store_receivers"`
			RequestsSenders      int `json:"requests_senders"`
			ResponsesSenders     int `json:"responses_senders"`
			CommandsReceivers    int `json:"commands_receivers"`
			QueriesReceivers     int `json:"queries_receivers"`
		} `json:"connections"`
		Traffic struct {
			SentEvents          int `json:"sent_events"`
			ReceivedEvents      int `json:"received_events"`
			SentRequests        int `json:"sent_requests"`
			SentError           int `json:"sent_error"`
			SentResponses       int `json:"sent_responses"`
			ReceivedRequests    int `json:"received_requests"`
			SentEventsVol       int `json:"sent_events_vol"`
			ReceivedEventsVol   int `json:"received_events_vol"`
			SentRequestsVol     int `json:"sent_requests_vol"`
			SentErrorsVol       int `json:"sent_errors_vol"`
			SentResponsesVol    int `json:"sent_responses_vol"`
			ReceivedRequestsVol int `json:"received_requests_vol"`
			TotalMessages       int `json:"total_messages"`
			TotalVolume         int `json:"total_volume"`
		} `json:"traffic"`
	} `json:"rest"`
}

func (h *Health) String() string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(h)
	if err != nil {
		return err.Error()
	}
	return buffer.String()
}

var healthCmd = &cobra.Command{
	Use:     "health",
	Short:   "Call kubemq health endpoint",
	Aliases: []string{"h"},
	Long:    `Return health stats for kubemq`,
	Run: func(cmd *cobra.Command, args []string) {
		runHealth()
	},
}

func runHealth() {
	resp := &rest.Response{}
	h := &Health{}
	r, err := resty.R().SetResult(h).SetError(resp).Get(cfg.HealthAddress)
	if err != nil {
		log.Fatal(err)
	}
	if !r.IsSuccess() {
		log.Fatal(r.Status())
	}
	if resp.Error {
		log.Fatal(resp.ErrorString)
	}
	log.Println(h.String())
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
