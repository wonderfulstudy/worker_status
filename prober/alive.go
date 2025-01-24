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
func setWrokerInfo(target string, registry *prometheus.Registry, logger *slog.Logger) (lookupTime float64, err error) {
	var (
		statusCodeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_alive_status_code",
			Help: "发送请求的返回状态码",
		})

		workerInfoGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_alive_worker_info",
			Help: "worker状态",
		}, []string{
			"worker_status",
			"worker_public_key",
			"worker_address",
			"app_version",
			"phase"})

		workerBlockInfoGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_alive_worker_block_info",
			Help: "worker 区块状态",
		}, []string{"latest_block_hash", "finalized_block_hash", "phase"})
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
		return 0.0, err
	}
	defer resp.Body.Close()
	statusCodeGauge.Set(float64(resp.StatusCode))

	body, resp_err := io.ReadAll(resp.Body)
	if resp_err != nil {
		logger.Error("读取节点信息失败", "err", resp_err)
		return 0.0, err
	}

	var workerInfo WorkerInfo
	json_err := json.Unmarshal([]byte(string(body)), &workerInfo)
	if json_err != nil {
		logger.Error("解析 JSON 失败:", "err", json_err)
		return 0.0, err
	}

	workerInfoGaugeVec.WithLabelValues(
		workerInfo.WorkerStatus,
		workerInfo.WorkerPublicKey,
		workerInfo.WorkerAddress,
		workerInfo.AppVersion,
		"attestedAt").Set(float64(workerInfo.AtestedAt))
	workerInfoGaugeVec.WithLabelValues(
		workerInfo.WorkerStatus,
		workerInfo.WorkerPublicKey,
		workerInfo.WorkerAddress,
		workerInfo.AppVersion,
		"implVersion").Set(float64(workerInfo.ImplVersion))
	workerBlockInfoGaugeVec.WithLabelValues(
		workerInfo.LatestBlockHash,
		workerInfo.FinalizedBlockHash,
		"latestBlockNumber").Set(float64(workerInfo.LatestBlockNumber))
	workerBlockInfoGaugeVec.WithLabelValues(
		workerInfo.LatestBlockHash,
		workerInfo.FinalizedBlockHash,
		"finalizedBlockNumber").Set(float64(workerInfo.FinalizedBlockNumber))

	return lookupTime, nil
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

	lookuptime, err := setWrokerInfo(target, registry, logger)
	if err != nil {
		logger.Error("设置节点指标失败", "err", err)
		return false
	}
	durationGaugeVec.Set(lookuptime)

	return true
}
