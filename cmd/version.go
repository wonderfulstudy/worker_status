package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func versionRun() {
	fmt.Println("v0.0.1")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "输出脚本工具版本号",
	Long:  `输出脚本工具版本号`,
	Run: func(cmd *cobra.Command, args []string) {
		versionRun()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd) // 添加 versionCmd 作为子命令
	rootCmd.Flags().BoolP("version", "v", false, versionCmd.Short)
	rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if version, _ := cmd.Flags().GetBool(versionCmd.Use); version {
			versionRun()
			os.Exit(0)
		}
	}
}
