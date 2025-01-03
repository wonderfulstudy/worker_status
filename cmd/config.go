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
	"gopkg.in/yaml.v3"
)

// 配置文件对应结构体
type Config struct {
	Addresses []string `yaml:"addresses,omitempty" json:"addresses,omitempty"`
}

// 定义全局配置变量
var configFile string

// 配置config命令的介绍信息，以及命令执行时默认执行的函数
var configCmd = &cobra.Command{
	Use:   "config [filename]",
	Short: "从配置文件中读取配置",
	Long:  `从配置文件中读取配置(默认配置文件为当前目录下的 config.yaml)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// 使用 filepath.FromSlash 确保路径格式正确
			configFile = filepath.FromSlash(args[0])
		} else {
			// 设置默认配置文件路径为当前目录下的 config.yaml
			var err error
			configFile, err = filepath.Abs("config.yaml")
			if err != nil {
				fmt.Printf("Error getting absolute path for config.yaml: %s\n", err)
				os.Exit(1)
			}
		}
		configRun()
	},
}

// config命令的执行逻辑代码
func configRun() int {
	fmt.Printf("Using config file: %s\n", configFile)
	if configFile != "" {
		// 从指定的配置文件中读取配置
		viper.SetConfigFile(configFile)
	} else {
		// 从当前目录下读取配置文件
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	// 读取 Addresses 字段并打印
	addresses := viper.GetStringSlice("pools.addresses")
	fmt.Printf("Addresses: %v\n", addresses)

	// 遍历 Addresses 数组并打印每个地址
	for _, address := range addresses {
		fmt.Printf("Address: %s\n", address)
	}

	fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	return 0
}

// 读取配置文件
func readConfig() *Config {
	config := &Config{}

	// 1.读取文件
	b, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	// 2. 绑定参数
	err2 := yaml.Unmarshal(b, config)
	if err2 != nil {
		panic(err)
	}

	return config
}

// config命令的初始化函数
func init() {
	rootCmd.AddCommand(configCmd)
}
