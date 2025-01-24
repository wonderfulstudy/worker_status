package prober

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/blackbox_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

type WorkerInfo struct {
	LatestBlockNumber    int    `json:"latestBlockNumber"`
	LatestBlockHash      string `json:"latestBlockHash"`
	FinalizedBlockNumber int    `json:"finalizedBlockNumber"`
	FinalizedBlockHash   string `json:"finalizedBlockHash"`
	WorkerAddress        string `json:"workerAddress"`
	WorkerPublicKey      string `json:"workerPublicKey"`
	WorkerStatus         string `json:"workerStatus"`
	AtestedAt            int    `json:"attestedAt"`
	AppVersion           string `json:"appVersion"`
	ImplVersion          int    `json:"implVersion"`
}

// 获取节点的json数据返回
func getWrokerInfo(ctx context.Context, target string, registry *prometheus.Registry, logger *slog.Logger) (info *WorkerInfo, lookupTime float64, err error) {
	var (
		statusCodeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_alive_status_code",
			Help: "发送请求的返回状态码",
		})

		workerInfoGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_alive_worker_info",
			Help: "worker状态",
		}, []string{})

		workerBlockInfoGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_alive_worker_block_info",
			Help: "worker状态",
		}, []string{})
	)

	registry.MustRegister(statusCodeGauge)
	registry.MustRegister(workerInfoGaugeVec)
	registry.MustRegister(workerBlockInfoGaugeVec)

	resolvestart := time.Now()

	defer func() {
		lookupTime = time.Since(resolvestart).Seconds()
	}()

	resp, err := http.Get(target)
	if err != nil {
		logger.Error("获取节点信息失败", "err", err)
		return nil, 0.0, err
	}
	defer resp.Body.Close()
	statusCodeGauge.Set(float64(resp.StatusCode))

	body, resp_err := io.ReadAll(resp.Body)
	if resp_err != nil {
		logger.Error("读取节点信息失败", http.StatusInternalServerError)
		return nil, 0.0, err
	}

	var workerInfo WorkerInfo
	json_err := json.Unmarshal([]byte(string(body)), &workerInfo)
	if json_err != nil {
		logger.Error("解析 JSON 失败:", http.StatusInternalServerError)
		return nil, 0.0, err
	}

	for _, ph := range []string{"attestedAt", "appVersion", "implVersion"} {
		workerInfoGaugeVec.WithLabelValues(ph).set
	}

}

func ProbeALIVE(ctx context.Context, target string, module config.Module, registry *prometheus.Registry, logger *slog.Logger) (success bool) {
	var (
		durationGaugeVec = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_alive_duration_seconds",
			Help: "探测节点耗费的总时间",
		})
	)
	// 如果目标没有加http或者https则加上
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}

}
