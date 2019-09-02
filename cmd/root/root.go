package root

import (
	"context"
	"github.com/kubemq-io/kubetools/cmd/commands"
	config2 "github.com/kubemq-io/kubetools/cmd/config"
	"github.com/kubemq-io/kubetools/cmd/dashboard"
	"github.com/kubemq-io/kubetools/cmd/delete"
	"github.com/kubemq-io/kubetools/cmd/deploy"
	"github.com/kubemq-io/kubetools/cmd/events"
	"github.com/kubemq-io/kubetools/cmd/events_store"
	"github.com/kubemq-io/kubetools/cmd/logs"
	"github.com/kubemq-io/kubetools/cmd/proxy"
	"github.com/kubemq-io/kubetools/cmd/queries"
	"github.com/kubemq-io/kubetools/cmd/queue"
	"github.com/kubemq-io/kubetools/cmd/scale"
	"github.com/kubemq-io/kubetools/cmd/status"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfg *config.Config
var version string
var rootCmd = &cobra.Command{
	Use: "kubetools",
}

func Execute(ver string) {
	version = ver
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
	if !exists(".kubetools.yaml") {
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
		viper.SetConfigName(".kubetools")
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
	rootCmd.AddCommand(proxy.NewCmdProxy(cfg))
	rootCmd.AddCommand(logs.NewCmdLogs(cfg))
	rootCmd.AddCommand(deploy.NewCmdDeploy(cfg))
	rootCmd.AddCommand(delete.NewCmdDelete(cfg))
	rootCmd.AddCommand(scale.NewCmdScale(cfg))
	rootCmd.AddCommand(status.NewCmdStatus(cfg))

}
