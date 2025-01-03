/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 配置文件对应结构体
type Config struct {
	Addresses []string `yaml:"addresses"`
}

// 定义全局配置变量
var GlobalConfig Config
var configFile string

// config命令的执行逻辑代码
func configRun(filename string) int {
	fmt.Printf("%s\n", filename)
	return 0
}

func initConfig() {
	if configFile != "" {
		// 从指定的配置文件中读取配置
		viper.SetConfigFile(configFile)
	} else {
		// 从当前目录下读取配置文件
		fmt.Printf("Using config file: %s\n", "config.yaml")
	}
}

// 配置config命令的介绍信息，以及命令执行时默认执行的函数
var configCmd = &cobra.Command{
	Use:   "config [filename]",
	Short: "从配置文件中读取配置",
	Long:  `从配置文件中读取配置(默认配置文件为当前目录下的 config.yaml)`,
	Run: func(cmd *cobra.Command, args []string) {
		var filename string
		if len(args) > 1 {
			cmd.Help()
			os.Exit(1)
		} else if len(args) == 1 {
			filename = args[0]
		} else {
			filename = "config.yaml"
		}
		configRun(filename)
	},
}

// config命令的初始化函数
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", configCmd.Short)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
