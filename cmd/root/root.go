package root

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/get"
	"github.com/kubemq-io/kubemqctl/cmd/scale"
	"github.com/kubemq-io/kubemqctl/cmd/set"

	"github.com/kubemq-io/kubemqctl/cmd/commands"
	"github.com/kubemq-io/kubemqctl/cmd/create"
	deleteCmd "github.com/kubemq-io/kubemqctl/cmd/delete"
	"github.com/kubemq-io/kubemqctl/cmd/events_store"
	"github.com/kubemq-io/kubemqctl/cmd/queries"
	"github.com/kubemq-io/kubemqctl/cmd/queue"

	configCmd "github.com/kubemq-io/kubemqctl/cmd/config"
	"github.com/kubemq-io/kubemqctl/cmd/events"

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
	ValidArgs: []string{"config", "commands", "queries", "queues", "events", "events_store", "create", "get", "delete", "scale"},
}

func Execute(version string) {

	rootCmd.Version = version
	defer utils.CheckErr(cfg.Save())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rootCmd.AddCommand(queue.NewCmdQueue(ctx, cfg))
	rootCmd.AddCommand(events.NewCmdEvents(ctx, cfg))
	rootCmd.AddCommand(events_store.NewCmdEventsStore(ctx, cfg))
	rootCmd.AddCommand(commands.NewCmdCommands(ctx, cfg))
	rootCmd.AddCommand(queries.NewCmdQueries(ctx, cfg))
	rootCmd.AddCommand(configCmd.NewCmdConfig(ctx, cfg))
	rootCmd.AddCommand(create.NewCmdCreate(ctx, cfg))
	rootCmd.AddCommand(deleteCmd.NewCmdDelete(ctx, cfg))
	rootCmd.AddCommand(get.NewCmdGet(ctx, cfg))
	rootCmd.AddCommand(scale.NewCmdScale(ctx, cfg))
	rootCmd.AddCommand(set.NewCmdSet(ctx, cfg))

	//_ = doc.GenMarkdownTree(rootCmd, "./docs")
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
