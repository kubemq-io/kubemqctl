package root

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/cluster"
	"github.com/kubemq-io/kubemqctl/cmd/commands"
	configCmd "github.com/kubemq-io/kubemqctl/cmd/config"
	"github.com/kubemq-io/kubemqctl/cmd/events"
	"github.com/kubemq-io/kubemqctl/cmd/events_store"
	"github.com/kubemq-io/kubemqctl/cmd/queries"
	"github.com/kubemq-io/kubemqctl/cmd/queue"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfg *config.Config
var Version string

var rootCmd = &cobra.Command{
	Use:       "kubemqctl",
	ValidArgs: []string{"cluster", "config", "commands", "queries", "queues", "events", "events_store"},
}

func Execute(version string) {
	rootCmd.Version = version
	defer utils.CheckErr(cfg.Save())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rootCmd.AddCommand(queue.NewCmdQueue(cfg))
	rootCmd.AddCommand(events.NewCmdEvents(cfg))
	rootCmd.AddCommand(events_store.NewCmdEventsStore(cfg))
	rootCmd.AddCommand(commands.NewCmdCommands(cfg))
	rootCmd.AddCommand(queries.NewCmdQueries(cfg))
	rootCmd.AddCommand(configCmd.NewCmdConfig(cfg))
	rootCmd.AddCommand(cluster.NewCmdCluster(ctx, cfg))

	utils.CheckErr(rootCmd.Execute())

}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func init() {
	cfg = &config.Config{}
	if !exists(".kubemqctl.yaml") {
		utils.Println("No configuration found, initialize first time default configuration. Run 'kubemqctl config' to run expert configuration wizard.")
	}

	defaultCfg, err := config.CheckConfigFile()
	if err != nil && defaultCfg != nil {
		cfg = defaultCfg
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName(".kubemqctl")
		err := viper.ReadInConfig()
		utils.CheckErr(err)
		err = viper.Unmarshal(cfg)
		utils.CheckErr(err)
	}

}
