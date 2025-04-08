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

type OwnerStruct struct {
	SolanaAddress string `json:"solana_address"`
}

type TestNet2NodeInfo struct {
	DeviceStatus string `json:"device_status"`
	Owner        OwnerStruct
	Pubkey       string `json:"pubkey"`
	Registry_api string `json:"registry_api"`
	RelayStatus  string `json:"relay_status"`
	Version      string `json:"version"`
}

// 获取节点的json数据返回
func setTestNet2NodeInfo(target string, registry *prometheus.Registry, logger *slog.Logger) (lookupTime float64, err error) {
	var (
		statusCodeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_testnet2_status_code",
			Help: "发送请求的返回状态码",
		})

		workerInfoGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_testnet2_worker_info",
			Help: "worker状态",
		}, []string{
			"device_status",
			"solana_address",
			"pubkey",
			"registry_api",
			"relay_status",
			"version",
		})
	)

	registry.MustRegister(statusCodeGauge)
	registry.MustRegister(workerInfoGaugeVec)

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

	var workerInfo TestNet2NodeInfo
	json_err := json.Unmarshal([]byte(string(body)), &workerInfo)
	if json_err != nil {
		logger.Error("解析 JSON 失败:", "err", json_err)
		return 0.0, err
	}

	workerInfoGaugeVec.WithLabelValues(
		workerInfo.DeviceStatus,
		workerInfo.Owner.SolanaAddress,
		workerInfo.Pubkey,
		workerInfo.Registry_api,
		workerInfo.RelayStatus,
		workerInfo.Version,
	).Set(float64(1))

	return lookupTime, nil
}

func ProbeTestNet2(ctx context.Context, target string, module config.Module, registry *prometheus.Registry, logger *slog.Logger) (success bool) {
	var (
		durationGaugeVec = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_testnet2_duration_seconds",
			Help: "探测节点耗费的总时间",
		})
	)
	// 如果目标没有加http或者https则加上
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}

	lookuptime, err := setTestNet2NodeInfo(target, registry, logger)
	if err != nil {
		logger.Error("设置节点指标失败", "err", err)
		return false
	}
	durationGaugeVec.Set(lookuptime)

	return true
}
