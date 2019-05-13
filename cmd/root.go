package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kubemq-io/kubetools/transport/option"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Control struct {
	Address string
}
type Config struct {
	Connections    []*option.Options
	HealthAddress  string
	MetricsAddress string
}

var cfg *Config
var rootCmd = &cobra.Command{
	Use:   "kubetools",
	Short: "Set of tools for kubemq",
	Long: `Set of tools for kubemq:
			1. test - test kubemq installation
			2. health - call kubemq health endpoint
			3. metrics - call kubemq prometheus metrics endpoint
			`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
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
}
