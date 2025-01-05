package cmd

import (
	"fmt"
	"worker_status/tools"

	"github.com/prometheus/client_golang/prometheus"
)

func NewWorkerCollector() *WorkerCollector {
	return &WorkerCollector{
		registeredWrokers: prometheus.NewDesc(
			"registered_workers",
			"矿池中注册机器总数",
			nil, nil),
		workerUptime: prometheus.NewDesc(
			"worker_uptime",
			"矿池中机器运行时间",
			[]string{"worker_address", "owner_address", "status"}, nil),
	}
}

func (c *WorkerCollector) Describe(ch chan<- *prometheus.Desc) {
	// 添加自定义指标描述
	ch <- c.registeredWrokers
	ch <- c.workerUptime
}

func (c *WorkerCollector) Collect(ch chan<- prometheus.Metric) {
	data, err := tools.GetWorkersFromOwnerAddresses(config.Pools.Addresses)
	if err != nil {
		fmt.Println("请求链上机器数据失败：", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.registeredWrokers,
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
