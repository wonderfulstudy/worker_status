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

type NodeInfo struct {
	DeviceStatus string `json:"device_status"`
	Pubkey       string
	RegistryApi  string `json:"registry_api"`
	RelayStatus  string `json:"relay_status"`
	Version      string
}

func ProbeDephy(ctx context.Context, target string, module config.Module, registry *prometheus.Registry, logger *slog.Logger) (success bool) {
	startT := time.Now()

	var (
		durationGaugeVec = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_dephy_duration_seconds",
			Help: "探测节点耗费的总时间",
		})

		statusCodeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_dephy_status_code",
			Help: "发送请求的返回状态码",
		})

		nodeInfoGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_dephy_node_info",
			Help: "worker状态",
		}, []string{
			"device_status",
			"pubkey",
			"registry_api",
			"relay_status",
			"version",
		})
	)
	defer func() {
		durationGaugeVec.Set(float64(time.Since(startT).Seconds()))
	}()

	registry.MustRegister(statusCodeGauge)
	registry.MustRegister(nodeInfoGaugeVec)

	// 如果目标没有加http或者https则加上
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}

	res, err := http.Get(target)
	if err != nil {
		logger.Error("对节点地址发起http请求失败", "err", err)
		return false
	}
	defer res.Body.Close()
	defer statusCodeGauge.Set(float64(res.StatusCode))

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error("获取节点返回的状态数据失败", "err", err)
		return false
	}

	node := NodeInfo{}
	if err := json.Unmarshal(body, &node); err != nil {
		logger.Error("解析节点状态数据失败", "err", err)
		return false
	}

	nodeInfoGaugeVec.WithLabelValues(
		node.DeviceStatus,
		node.Pubkey,
		node.RegistryApi,
		node.RelayStatus,
		node.Version,
	).Set(float64(1))

	return true
}
