package cmd

import (
	"fmt"
	"worker_status/tools"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollector struct {
	allRegisteredWrokers *prometheus.Desc
	myRegisteredWrokers  *prometheus.Desc
	workerUptime         *prometheus.Desc
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		allRegisteredWrokers: prometheus.NewDesc(
			"all_registered_workers",
			"全网注册机器数",
			nil, nil),
		myRegisteredWrokers: prometheus.NewDesc(
			"registered_workers",
			"矿池中注册机器总数",
			nil, nil),
		workerUptime: prometheus.NewDesc(
			"worker_uptime",
			"矿池中机器运行时间",
			[]string{"worker_address", "owner_address", "status"}, nil),
	}
}

func (c *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	// 添加自定义指标描述
	ch <- c.allRegisteredWrokers
	ch <- c.myRegisteredWrokers
	ch <- c.workerUptime
}

func (c *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	data, err := tools.GetWorkersFromOwnerAddresses(config.Pools.Addresses)
	if err != nil {
		fmt.Println("请求链上机器数据失败：", err)
		return
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
		c.myRegisteredWrokers,
		prometheus.GaugeValue,
		float64(len(data.Data.Workers)))

	for _, worker := range data.Data.Workers {
		ch <- prometheus.MustNewConstMetric(
			c.workerUptime,
			prometheus.GaugeValue,
			float64(worker.Uptime),
			worker.Address,
			worker.OwnerAddress,
			worker.Status)
	}
}
