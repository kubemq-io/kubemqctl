package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "print kubemq version",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(fmt.Sprintf("kubetools version %s", version))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
