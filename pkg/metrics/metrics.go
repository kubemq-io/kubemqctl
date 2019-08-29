package metrics

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kubemq-io/kubetools/pkg/config"
	"io"
	"text/tabwriter"
	"time"
)

type Metric struct {
	Error       bool   `json:"error"`
	ErrorString string `json:"error_string"`
	Data        struct {
		KindID      int    `json:"kind_id"`
		KindName    string `json:"kind_name"`
		SubKindID   int    `json:"sub_kind_id"`
		SubKindName string `json:"sub_kind_name"`
		ChartsData  []struct {
			ID             string    `json:"id"`
			Node           string    `json:"node"`
			Time           time.Time `json:"time"`
			TimeUnix       int64     `json:"time_unix"`
			KindID         int       `json:"kind_id"`
			SubKindID      int       `json:"sub_kind_id"`
			Channel        string    `json:"channel"`
			Group          string    `json:"group"`
			ClientID       string    `json:"client_id"`
			Volume         float64   `json:"volume"`
			EgressVolume   float64   `json:"egress_volume"`
			IngressVolume  float64   `json:"ingress_volume"`
			Errors         int64     `json:"errors"`
			Events         int64     `json:"events"`
			Commands       int64     `json:"commands"`
			Queries        int64     `json:"queries"`
			Responses      int64     `json:"responses"`
			QueueMessages  int64     `json:"queue_messages"`
			CacheHits      int64     `json:"cache_hits"`
			CacheMiss      int64     `json:"cache_miss"`
			Traffic        float64   `json:"traffic"`
			EgressTraffic  float64   `json:"egress_traffic"`
			IngressTraffic float64   `json:"ingress_traffic"`
			ErrorRate      float64   `json:"error_rate"`
			Latency        float64   `json:"latency"`
			Active         bool      `json:"active"`
			CacheRatio     float64   `json:"cache_ratio"`
			MeanLatency    float64   `json:"mean_latency"`
		} `json:"charts_data"`
	} `json:"data"`
	Channels struct {
		Headers  []string `json:"headers"`
		Channels []struct {
			ID            string  `json:"id"`
			Active        bool    `json:"active"`
			Channel       string  `json:"channel"`
			Group         string  `json:"group"`
			Clients       int64   `json:"clients"`
			Traffic       float64 `json:"traffic"`
			Volume        float64 `json:"volume"`
			Events        int64   `json:"events"`
			Errors        int64   `json:"errors"`
			Commands      int64   `json:"commands"`
			Queries       int64   `json:"queries"`
			QueueMessages int64   `json:"queue_messages"`
			Responses     int64   `json:"responses"`
			CacheHits     int64   `json:"cache_hits"`
			CacheMiss     int64   `json:"cache_miss"`
			ErrorRate     float64 `json:"error_rate"`
			CacheRatio    float64 `json:"cache_ratio"`
			MeanLatency   float64 `json:"mean_latency"`
		} `json:"channels"`
	} `json:"channels"`
	LastUpdate int `json:"last_update"`
}

func PrintMetrics(ctx context.Context, out io.Writer, cfg *config.Config, kind, subkind string, timeframe string, top int) error {
	metrics := &Metric{}
	r, err := resty.New().R().SetResult(metrics).SetError(metrics).SetQueryParam("kind_id", kind).SetQueryParam("sub_kind_id", subkind).SetQueryParam("time_frame", timeframe).Get(fmt.Sprintf("%s/v1/stats/sub_kind", cfg.GetApiHttpURI()))
	if err != nil {
		return nil
	}
	if !r.IsSuccess() {
		return fmt.Errorf(r.Status())
	}
	if metrics.Error {
		return fmt.Errorf(metrics.ErrorString)
	}
	switch subkind {
	case "1", "2":
		header := "NODE\tTIME\tCHANNEL\tVOLUME\tEVENTS\tERRORS"
		row := "%s\t%s\t%s\t%2.f\t%d\t%d\t\n"
		tw := tabwriter.NewWriter(out, 0, 0, 1, ' ', tabwriter.TabIndent)
		_, _ = fmt.Fprintln(tw, header)
		for i := 0; i < len(metrics.Data.ChartsData) && i < top; i++ {
			d := metrics.Data.ChartsData[i]
			_, _ = fmt.Fprintf(tw, row,
				d.Node,
				d.Time.Format("2006-01-02 15:04:05"),
				d.Channel,
				d.Volume,
				d.Events,
				d.Errors,
			)
		}
		tw.Flush()
	}
	return nil
}
