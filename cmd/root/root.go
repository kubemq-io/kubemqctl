package root

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/cluster"
	"github.com/kubemq-io/kubemqctl/cmd/commands"
	config2 "github.com/kubemq-io/kubemqctl/cmd/config"

	"github.com/kubemq-io/kubemqctl/cmd/dashboard"

	"github.com/kubemq-io/kubemqctl/cmd/events"
	"github.com/kubemq-io/kubemqctl/cmd/events_store"

	"github.com/kubemq-io/kubemqctl/cmd/queries"
	"github.com/kubemq-io/kubemqctl/cmd/queue"

	version2 "github.com/kubemq-io/kubemqctl/cmd/version"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfg *config.Config
var Version string
var rootCmd = &cobra.Command{
	Use: "kubemqctl",
}

func Execute(version string) {
	Version = version
	defer utils.CheckErr(cfg.Save())
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
		utils.Println("No configuration found, initialize first time configuration:")
		cfgOpts := &config2.ConfigOptions{Cfg: config.DefaultConfig}
		err := cfgOpts.Run(context.Background())
		utils.CheckErr(err)
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
	rootCmd.AddCommand(queue.NewCmdQueue(cfg))
	rootCmd.AddCommand(events.NewCmdEvents(cfg))
	rootCmd.AddCommand(events_store.NewCmdEventsStore(cfg))
	rootCmd.AddCommand(commands.NewCmdCommands(cfg))
	rootCmd.AddCommand(queries.NewCmdQueries(cfg))
	rootCmd.AddCommand(config2.NewCmdConfig(cfg))
	rootCmd.AddCommand(dashboard.NewCmdDashboard(cfg))
	rootCmd.AddCommand(version2.NewCmdVersion(&Version))
	rootCmd.AddCommand(cluster.NewCmdCluster(cfg))

}
