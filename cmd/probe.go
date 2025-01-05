package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// probeHandler 处理 /probe 请求
func probeHandler(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数
	target := r.URL.Query().Get("target")

	fmt.Println("target:", target)

	// 检查参数是否为空
	if target == "" {
		http.Error(w, "target and module parameters are required", http.StatusBadRequest)
		return
	}

	resp, get_err := http.Get(target)
	if get_err != nil {
		http.Error(w, "Failed to fetch target", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		http.Error(w, "Failed to fetch target", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(body))

	var workerInfo WorkerInfo
	json_err := json.Unmarshal([]byte(string(body)), &workerInfo)
	if json_err != nil {
		http.Error(w, "解析 JSON 失败:", http.StatusInternalServerError)
		return
	}

	// 创建一个新的注册表
	workerRegistry := prometheus.NewRegistry()

	// 注册机器总数指标
	workerStatus := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "workerStatus",
		Help: "运行状态",
		ConstLabels: prometheus.Labels{
			"value": workerInfo.WorkerStatus,
		},
	})
	workerAddress := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "workerAddress",
		Help: "机器地址",
		ConstLabels: prometheus.Labels{
			"value": workerInfo.WorkerAddress,
		},
	})
	workerPublicKey := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "workerPublicKey",
		Help: "机器公钥",
		ConstLabels: prometheus.Labels{
			"value": workerInfo.WorkerPublicKey,
		},
	})
	latestBlockNumber := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "latestBlockNumber",
		Help: "机器当前区块高度",
	})
	finalizedBlockNumber := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "finalizedBlockNumber",
		Help: "全网链上区块高度",
	})

	workerRegistry.MustRegister(workerStatus)
	workerRegistry.MustRegister(workerAddress)
	workerRegistry.MustRegister(workerPublicKey)
	workerRegistry.MustRegister(latestBlockNumber)
	workerRegistry.MustRegister(finalizedBlockNumber)

	workerStatus.Set(float64(1))
	workerAddress.Set(float64(1))
	workerPublicKey.Set(float64(1))
	latestBlockNumber.Set(float64(workerInfo.LatestBlockNumber))
	finalizedBlockNumber.Set(float64(workerInfo.FinalizedBlockNumber))

	// 使用 promhttp.HandlerFor 返回指标数据
	handler := promhttp.HandlerFor(workerRegistry, promhttp.HandlerOpts{})
	handler.ServeHTTP(w, r)
}
