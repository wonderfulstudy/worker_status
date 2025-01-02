package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

// 读取配置文件，功能待开发
type Config struct { // 待定
	Addresses []string `yaml:"Addresses"`
}

// 全局配置变量
var GlobalConfig Config

func InitConfig() error {
	// 定义命令行参数
	configPath := flag.String("config", "config.yml", "path to config file")
	flag.Parse()

	// 读取配置文件
	config, err := readConfig(*configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// 设置全局配置
	GlobalConfig = *config

	fmt.Println("Addresses from config:", GlobalConfig.Addresses)
	return nil
}

func readConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// 定义单台矿机结构体
type Worker struct {
	Status  string `json:"status"`
	Address string `json:"address"`
}

// 定义矿机集合结构体
type Data struct {
	Workers []Worker `json:"workers"`
}

// // 定义post请求获取到的数据负载结构体
// type Response struct {
// 	Data Data `json:"data"`
// }

// // 定义post请求负载数据
// type GraphQLQuery struct {
// 	OperationName string                 `json:"operationName"`
// 	Query         string                 `json:"query"`
// 	Variables     map[string]interface{} `json:"variables"`
// }

// 读取链上数据
func get_registered_workers() int {
	url := "https://dev-squid.cybros.network/graphql"
	data := []byte(`
		{
			"operationName": "MyQuery",
			"query": "query MyQuery {\n  workers(where: {ownerAddress_eq: \"5ES6Mu3uUUB8QRarGB4J73RnoeTxADWV8HKhnKfrASpoLSEP\"}) {\n    status\n    address\n  }\n}\n",
			"variables":null
		}`)

	// 发起 POST 请求
	resp, post_err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if post_err != nil {
		fmt.Println("发起请求失败:", post_err)
		return -1
	}
	defer resp.Body.Close()

	// 读取响应
	body, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		fmt.Println("读取响应失败:", resp_err)
		return -1
	}

	// 定义目标结构体变量
	var response Response
	// 解析 JSON 数据
	json_err := json.Unmarshal([]byte(string(body)), &response)
	if json_err != nil {
		fmt.Println("解析 JSON 失败:", json_err)
		return -1
	}

	// 打印解析结果
	count := 0
	fmt.Println("解析成功，workers 信息如下：")
	for _, worker := range response.Data.Workers {
		fmt.Printf("Status: %s, Address: %s\n", worker.Status, worker.Address)
		count++
	}

	return count
}

func main() {
	// 初始化配置
	err := InitConfig()
	if err != nil {
		fmt.Printf("Failed to initialize config: %v\n", err)
		return
	}

	fmt.Print("%s", GlobalConfig.Addresses)

	// 创建一个新的注册表
	registry := prometheus.NewRegistry()

	// 注册机器总数指标
	registeredWorkers := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "registeredWorkers", // 指标名称
		Help: "在链上已经注册的机器总数",      // 指标帮助信息
	})

	// 注册指标
	registry.MustRegister(registeredWorkers)

	// 增加指标值
	go func() {
		for {
			count := get_registered_workers()
			registeredWorkers.Set(float64(count))
			time.Sleep(10 * time.Second)
		}
	}()

	// 创建一个HTTP处理器来暴露指标
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	// 启动Web服务器
	http.ListenAndServe(":8080", nil)
}
