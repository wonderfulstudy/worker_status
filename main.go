package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 读取配置文件，功能待开发
type pool struct { // 待定
	Id      int
	Name    string
	Address string
}

func get_config() {
	fmt.Println("todolist")
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

// 定义post请求获取到的数据负载结构体
type Response struct {
	Data Data `json:"data"`
}

// 读取链上数据
func get_workers() {
	url := "https://dev-squid.cybros.network/graphql"
	data := []byte(`{"operationName": "MyQuery", "query": "query MyQuery {\n  workers(where: {ownerAddress_eq: \"5ES6Mu3uUUB8QRarGB4J73RnoeTxADWV8HKhnKfrASpoLSEP\"}) {\n    status\n    address\n  }\n}\n", "variables":null}`)

	// 发起 POST 请求
	resp, post_err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if post_err != nil {
		fmt.Println("发起请求失败:", post_err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		fmt.Println("读取响应失败:", resp_err)
		return
	}

	// 定义目标结构体变量
	var response Response
	// 解析 JSON 数据
	json_err := json.Unmarshal([]byte(string(body)), &response)
	if json_err != nil {
		fmt.Println("解析 JSON 失败:", json_err)
		return
	}

	// 打印解析结果
	fmt.Println("解析成功，workers 信息如下：")
	for _, worker := range response.Data.Workers {
		fmt.Printf("Status: %s, Address: %s\n", worker.Status, worker.Address)
	}
}

func recordResigters() {
	go func ()  {
		for {
			
		}
		
	}
}

func main() {
	// 注册机器总数指标
	counterMetric := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "register_workers", // 指标名称
		Help: "在链上已经注册的机器总数",     // 指标帮助信息
	})

	// 注册指标
	prometheus.MustRegister(counterMetric)

	// 增加指标值
	counterMetric.Inc()

	// 创建一个HTTP处理器来暴露指标
	http.Handle("/metrics", promhttp.Handler())

	// 启动Web服务器
	http.ListenAndServe(":8080", nil)
	get_workers()
}
