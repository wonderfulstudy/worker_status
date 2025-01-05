package cmd

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// config.yaml 配置文件结构体
type Config struct {
	Pools struct {
		UpgradeTime time.Duration `yaml:"upgradeTime,omitempty"`
		Addresses   []string      `yaml:"addresses,omitempty"`
	} `yaml:"pools"`
}

// 请求本地节点8080端口返回数据结构体
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

type WorkerCollector struct {
	registeredWrokers *prometheus.Desc
	workerUptime      *prometheus.Desc
}
