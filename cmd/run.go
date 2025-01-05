/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "作为服务运行",
	Long:  `作为服务运行，默认端口为8080，默认为exporter模式`,
	Run: func(cmd *cobra.Command, args []string) {
		onlineRegistry := prometheus.NewRegistry()
		collector := NewWorkerCollector()
		onlineRegistry.MustRegister(collector)
		http.Handle("/metrics", promhttp.HandlerFor(onlineRegistry, promhttp.HandlerOpts{}))
		http.HandleFunc("/probe", probeHandler)

		// // 启动Web服务器
		http.ListenAndServe(":8080", nil)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
