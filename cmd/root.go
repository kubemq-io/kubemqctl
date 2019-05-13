package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kubemq-io/kubetools/transport/option"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Connections []*option.Options
}

var cfg *Config
var rootCmd = &cobra.Command{
	Use:   "kubetools",
	Short: "Set of tools for kubemq",
	Long: `Set of tools for kubemq:
			1. kubetest - test kubemq installation
			2. kiubemon - monitor channel traffic
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
