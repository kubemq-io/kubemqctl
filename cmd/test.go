package cmd

import (
	"context"
	"time"

	"github.com/kubemq-io/kubetools/cmd/kubetest"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t"},
	Short:   "test your kubemq installation",
	Long:    `A testing tool for kubemq installation which simulate all the messaging patterns for gRPC and Rest interfaces`,
	Example: "kubetools test",
	Run: func(cmd *cobra.Command, args []string) {
		runTest()
	},
}

func runTest() {
	eventsTest := &kubetest.Scenario{
		Name:           "testing events pattern",
		Type:           kubetest.ScenarioTypeEvents,
		Producers:      10,
		Consumers:      10,
		MessageSize:    100,
		MessagesToSend: 1,
		Timeout:        30000,
	}
	eventsStoreTest := &kubetest.Scenario{
		Name:           "testing events store pattern",
		Type:           kubetest.ScenarioTypeEventsStore,
		Producers:      10,
		Consumers:      10,
		MessageSize:    100,
		MessagesToSend: 1,
		Timeout:        30000,
	}
	commandsTest := &kubetest.Scenario{
		Name:           "testing commands pattern",
		Type:           kubetest.ScenarioTypeCommands,
		Producers:      10,
		Consumers:      10,
		MessageSize:    100,
		MessagesToSend: 1,
		Timeout:        30000,
	}
	queriesTest := &kubetest.Scenario{
		Name:           "testing queries pattern",
		Type:           kubetest.ScenarioTypeQueries,
		Producers:      10,
		Consumers:      10,
		MessageSize:    100,
		MessagesToSend: 1,
		Timeout:        30000,
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	for _, conn := range cfg.Connections {
		kubetest.Execute(ctx, eventsTest, conn)
		kubetest.Execute(ctx, eventsStoreTest, conn)
		kubetest.Execute(ctx, commandsTest, conn)
		kubetest.Execute(ctx, queriesTest, conn)
	}
}

func init() {
	rootCmd.AddCommand(testCmd)

}
