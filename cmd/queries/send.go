package queries

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/targets"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"time"
)

type QueriesSendOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	body      string
	metadata  string
	timeout   int
	cacheKey  string
	cacheTTL  time.Duration
	fileName  string
	build     bool
}

var queriesSendExamples = `
	# Send query to a 'queries' channel
	kubemqctl queries send some-channel some-query
	
	# Send query to a 'queries' channel with metadata
	kubemqctl queries send some-channel some-body -m some-metadata
	
	# Send query to a 'queries' channel with 120 seconds timeout
	kubemqctl queries send some-channel some-body -o 120
	
	# Send query to a 'queries' channel with cache-key and cache duration of 1m
	kubemqctl queries send some-channel some-body -c cache-key -d 1m
`
var queriesSendLong = `Send command allow to send messages to 'queries' channel with an option to set query time-out and caching parameters`
var queriesSendShort = `Send messages to a 'queries' channel command`

func NewCmdQueriesSend(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &QueriesSendOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "send",
		Aliases: []string{"s"},
		Short:   queriesSendShort,
		Long:    queriesSendLong,
		Example: queriesSendExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringVarP(&o.metadata, "metadata", "m", "", "set query body metadata field")
	cmd.PersistentFlags().StringVarP(&o.cacheKey, "cache-key", "c", "", "set query cache key")
	cmd.PersistentFlags().IntVarP(&o.timeout, "timeout", "o", 30, "set query timeout")
	cmd.PersistentFlags().DurationVarP(&o.cacheTTL, "cache-duration", "d", 10*time.Minute, "set cache duration timeout")
	cmd.PersistentFlags().StringVarP(&o.fileName, "file", "f", "", "set load body from file")
	cmd.PersistentFlags().BoolVarP(&o.build, "build", "b", false, "build kubemq targets request")
	return cmd
}

func (o *QueriesSendOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]

	} else {
		return fmt.Errorf("missing channel argument")
	}
	if o.build {
		data, err := targets.BuildRequest()
		if err != nil {
			return err
		}
		o.body = string(data)
		return nil
	}
	if o.fileName != "" {
		data, err := ioutil.ReadFile(o.fileName)
		if err != nil {
			return err
		}
		o.body = string(data)
	} else {
		if len(args) >= 2 {
			o.body = args[1]
		} else {
			return fmt.Errorf("missing body argument")
		}
	}
	return nil
}
func (o *QueriesSendOptions) Validate() error {
	return nil
}

func (o *QueriesSendOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubemqClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())
	}

	defer func() {
		client.Close()
	}()
	fmt.Println("Sending Query:")
	msg := client.Q().
		SetChannel(o.channel).
		SetId(uuid.New().String()).
		SetBody([]byte(o.body)).
		SetMetadata(o.metadata).
		SetTimeout(time.Duration(o.timeout) * time.Second).
		SetCacheKey(o.cacheKey).
		SetCacheTTL(o.cacheTTL)
	printQuery(msg)
	res, err := msg.Send(ctx)
	if err != nil {
		return fmt.Errorf("sending query body, %s", err.Error())
	}
	fmt.Println("Getting Query Response:")
	printQueryResponse(res)
	return nil
}
