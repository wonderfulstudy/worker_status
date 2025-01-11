package cmd

import (
	"fmt"
	"worker_status/tools"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollector struct {
	allRegisteredWrokers     *prometheus.Desc
	registeredWrokers        *prometheus.Desc
	deregisteredWorkers      *prometheus.Desc
	onlineWorkers            *prometheus.Desc
	unresponsiveWorkers      *prometheus.Desc
	nodeStatus               *prometheus.Desc
	requestingOfflineWrokers *prometheus.Desc
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		allRegisteredWrokers: prometheus.NewDesc(
			"all_registered_workers",
			"全网注册机器数",
			nil, nil),
		registeredWrokers: prometheus.NewDesc(
			"registered_workers",
			"矿池中注册机器总数",
			nil, nil),
		deregisteredWorkers: prometheus.NewDesc(
			"deregistered_workers",
			"矿池中未注册成功机器总数",
			nil, nil),
		onlineWorkers: prometheus.NewDesc(
			"online_workers",
			"矿池中在线机器总数",
			nil, nil),
		unresponsiveWorkers: prometheus.NewDesc(
			"unresponsive_workers",
			"矿池中掉线机器总数",
			nil, nil),
		nodeStatus: prometheus.NewDesc(
			"unresponsive_workers_list",
			"矿池中掉线机器总数",
			[]string{"address", "owner_address", "status"}, nil),
		requestingOfflineWrokers: prometheus.NewDesc(
			"requesting_offline_workers",
			"矿池中请求离线的机器总数",
			nil, nil),
	}
}

func (c *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	// 添加自定义指标描述
	ch <- c.allRegisteredWrokers
	ch <- c.registeredWrokers
	ch <- c.deregisteredWorkers
	ch <- c.onlineWorkers
	ch <- c.unresponsiveWorkers
	ch <- c.nodeStatus
	ch <- c.requestingOfflineWrokers
}

func (c *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	r, d, o, u, ro := 0, 0, 0, 0, 0
	data, err := tools.GetWorkersFromOwnerAddresses(config.Pools.Addresses)
	if err != nil {
		fmt.Println("请求链上机器数据失败：", err)
		return
	}
	for _, worker := range data.Data.Workers {
		if worker.Status == "Registered" {
			r++
		} else if worker.Status == "Deregistered" {
			d++
		} else if worker.Status == "Online" {
			o++
		} else if worker.Status == "Unresponsive" {
			u++
			ch <- prometheus.MustNewConstMetric(
				c.nodeStatus,
				prometheus.GaugeValue,
				float64(1),
				worker.Address,
				worker.OwnerAddress,
				worker.Status)
		} else if worker.Status == "RequestingOffline" {
			ro++
		}
	}

	allRegistredWorkers, err := tools.GetAllWorkersCount()
	if err != nil {
		fmt.Println("请求链上注册机器总数数据失败：", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		c.allRegisteredWrokers,
		prometheus.GaugeValue,
		float64(allRegistredWorkers))
	ch <- prometheus.MustNewConstMetric(
		c.registeredWrokers,
		prometheus.GaugeValue,
		float64(r+o+u+ro))
	ch <- prometheus.MustNewConstMetric(
		c.deregisteredWorkers,
		prometheus.GaugeValue,
		float64(d))
	ch <- prometheus.MustNewConstMetric(
		c.onlineWorkers,
		prometheus.GaugeValue,
		float64(o))
	ch <- prometheus.MustNewConstMetric(
		c.unresponsiveWorkers,
		prometheus.GaugeValue,
		float64(u))
	ch <- prometheus.MustNewConstMetric(
		c.requestingOfflineWrokers,
		prometheus.GaugeValue,
		float64(ro))
}
