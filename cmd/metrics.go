package cmd

import (
	"fmt"
	"worker_status/tools"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollector struct {
	allRegisteredWrokers       *prometheus.Desc
	myRegisteredWrokers        *prometheus.Desc
	myDeregisteredWorkers      *prometheus.Desc
	myOnlineWorkers            *prometheus.Desc
	myUnresponsiveWorkers      *prometheus.Desc
	myUnresponsiveList         *prometheus.Desc
	myRequestingOfflineWrokers *prometheus.Desc
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
		myDeregisteredWorkers: prometheus.NewDesc(
			"deregistered_workers",
			"矿池中未注册成功机器总数",
			nil, nil),
		myOnlineWorkers: prometheus.NewDesc(
			"online_workers",
			"矿池中在线机器总数",
			nil, nil),
		myUnresponsiveWorkers: prometheus.NewDesc(
			"unresponsive_workers",
			"矿池中掉线机器总数",
			nil, nil),
		myUnresponsiveList: prometheus.NewDesc(
			"unresponsive_workers_list",
			"矿池中掉线机器总数",
			[]string{"address", "owner_address"}, nil),
		myRequestingOfflineWrokers: prometheus.NewDesc(
			"requesting_offline_workers",
			"矿池中请求离线的机器总数",
			nil, nil),
	}
}

func (c *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	// 添加自定义指标描述
	ch <- c.allRegisteredWrokers
	ch <- c.myRegisteredWrokers
	ch <- c.myDeregisteredWorkers
	ch <- c.myOnlineWorkers
	ch <- c.myUnresponsiveWorkers
	ch <- c.myUnresponsiveList
	ch <- c.myRequestingOfflineWrokers
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
				c.myUnresponsiveList,
				prometheus.GaugeValue,
				float64(1),
				worker.Address,
				worker.OwnerAddress)
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
		c.myRegisteredWrokers,
		prometheus.GaugeValue,
		float64(r+o+u+ro))
	ch <- prometheus.MustNewConstMetric(
		c.myDeregisteredWorkers,
		prometheus.GaugeValue,
		float64(d))
	ch <- prometheus.MustNewConstMetric(
		c.myOnlineWorkers,
		prometheus.GaugeValue,
		float64(o))
	ch <- prometheus.MustNewConstMetric(
		c.myUnresponsiveWorkers,
		prometheus.GaugeValue,
		float64(u))
	ch <- prometheus.MustNewConstMetric(
		c.myRequestingOfflineWrokers,
		prometheus.GaugeValue,
		float64(ro))
}
