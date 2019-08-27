package cmd

//
//import (
//	"io/ioutil"
//	"log"
//	"net/http"
//
//	"github.com/spf13/cobra"
//)
//
//var metricsCmd = &cobra.Command{
//	Use:     "metrics",
//	Short:   "Call kubemq metrics endpoint",
//	Aliases: []string{"m"},
//	Long:    `Return prometheus metrics for kubemq`,
//	Run: func(cmd *cobra.Command, args []string) {
//		runMetrics()
//	},
//}
//
//func runMetrics() {
//	resp, err := http.Get(cfg.MetricsAddress)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	defer resp.Body.Close()
//
//	data, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Println(string(data))
//
//}
//
//func init() {
//	rootCmd.AddCommand(metricsCmd)
//}
