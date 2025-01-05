/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 定义全局配置变量
var configFile string
var config Config

// 配置config命令的介绍信息，以及命令执行时默认执行的函数
var configCmd = &cobra.Command{
	Use:   "config [filename]",
	Short: "从配置文件中读取配置",
	Long:  `从配置文件中读取配置(默认配置文件为当前目录下的 config.yaml)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			cmd.Help()
			os.Exit(1)
		} else if len(args) == 1 {
			// 使用 filepath.FromSlash 确保路径格式正确
			configFile = filepath.FromSlash(args[0])
			// 从指定的配置文件中读取配置
			viper.SetConfigFile(configFile)
		}

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("读取配置文件失败: %s\n", err)
			os.Exit(1)
		}

		// 将配置文件内容绑定到 config 结构体
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Printf("解析文件内容并绑定到config结构体失败: %s\n", err)
			os.Exit(1)
		}
	},
}

// config命令的初始化函数
func init() {
	rootCmd.AddCommand(configCmd)
	// 设置默认配置文件路径为当前目录下的 config.yaml
	var err error
	configFile, err = filepath.Abs("config.yaml")
	if err != nil {
		fmt.Printf("Error getting absolute path for config.yaml: %s\n", err)
		os.Exit(1)
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	// 将配置文件内容绑定到 config 结构体
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshalling config file into struct: %s\n", err)
		os.Exit(1)
	}
}
