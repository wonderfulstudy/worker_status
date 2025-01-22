package config

type Config struct {
	Ipaddr string
	User   struct {
		Name   string
		Passwd string
		Pass   string
	}
}

var configFile string
var config Config

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
