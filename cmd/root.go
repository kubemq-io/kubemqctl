package cmd

import (
	"github.com/kubemq-io/kubetools/transport/option"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Control struct {
	Address string
}
type Config struct {
	Connections    []*option.Options
	StatsAddress   string
	HealthAddress  string
	MetricsAddress string
	MonitorAddress string
}

var cfg *Config
var version string
var rootCmd = &cobra.Command{
	Use:   "kubetools",
	Short: "Set of tools for kubemq",
	Long: `Set of tools for kubemq:
			1. test - test kubemq installation
			2. health - call kubemq health endpoint
			3. metrics - call kubemq prometheus metrics endpoint
			4. monitor - call kubemq monitor points to watch channel content
			5. pubsub - sending and receiving Pub/Sub messages
			6. queue - sending and receiving Queue messages
			7. rpc - sending and receiving RPC messages
			8. get - get list of resources
			`,
}

func Execute(ver string) {
	version = ver
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	checkConfigFile()
	cfg = &Config{}
	viper.AddConfigPath("./")
	viper.SetConfigName(".config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	//this is for backwards compatible in case yaml config for stats address is not exist
	if cfg.StatsAddress == "" {
		cfg.StatsAddress = strings.Replace(cfg.MonitorAddress, "ws", "http", 1)
	}
}
