/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// registeredCmd represents the registered command
var registeredCmd = &cobra.Command{
	Use:   "registered",
	Short: "查询账号下已注册的机器信息",
	Long:  `查询账号下已注册的机器信息`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("registered called")
	},
}

func init() {
	getCmd.AddCommand(registeredCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registeredCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registeredCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
