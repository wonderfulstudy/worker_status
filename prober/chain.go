package prober

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/blackbox_exporter/config"
	"github.com/prometheus/client_golang/prometheus"
)

// 查询所有机器状态
var AllWorkersQuery string = `
	query MyQuery {
		workersConnection(orderBy: id_ASC) {
			pageInfo {
				endCursor
				hasNextPage
				hasPreviousPage
				startCursor
			}
			totalCount
			edges {
				node {
					address
					ownerAddress
					status
					uptime
				}
			}
		}
	}`

// 根据矿池地址列表查询机器在线情况
var PoolQuery string = `
	query MyQuery {
		workers(where: {ownerAddress_in: ["%v"]}) {
			address
			ownerAddress
			status
			uptime
		}
	}`

type TotalCountResponse struct {
	Data struct {
		WorkersConnection struct {
			PageInfo struct {
				EndCursor       string `json:"endCursor"`
				HasNextPage     bool   `json:"hasNextPage"`
				HasPreviousPage bool   `json:"hasPreviousPage"`
				StartCursor     string `json:"startCursor"`
			} `json:"pageInfo"`
			TotalCount int `json:"totalCount"`
			Edges      []struct {
				Node struct {
					Address      string `json:"address"`
					OwnerAddress string `json:"ownerAddress"`
					Status       string `json:"status"`
					Uptime       int    `json:"uptime"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"workersConnection"`
	} `json:"data"`
}

type PoolResponse struct {
	Data struct {
		Workers struct {
			Address      string `json:"address"`
			OwnerAddress string `json:"ownerAddress"`
			Status       string `json:"status"`
			Uptime       string `json:"uptime"`
		} `json:"workers"`
	} `json:"data"`
}

func queryPost(query string, url string, logger *slog.Logger) (statusCode int, body []byte, err error) {
	// 将查询字符串转换为 JSON 格式
	requestPayload := map[string]string{
		"query":         query,
		"operationName": "MyQuery",
		"variables":     "null",
	}
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		logger.Error("加载post请求体负载失败", "err", err)
		return 0, []byte{}, err
	}

	// 发起 POST 请求
	resp, post_err := http.Post(url, "application/json", bytes.NewBuffer([]byte(payloadBytes)))
	if post_err != nil {
		logger.Error("发送post请求失败", "err", err)
		return 0, []byte{}, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, resp_err := io.ReadAll(resp.Body)
	if resp_err != nil {
		logger.Error("读取post请求失败", "err", err)
		return 0, []byte{}, err
	}

	return resp.StatusCode, body, nil
}

// 设置链上所有数据
func setChainAllInfo(target string, registry *prometheus.Registry, logger *slog.Logger) (lookupTime float64, err error) {
	var (
		statusCodeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_chain_all_status_code",
			Help: "请求链上所有数据状态码",
		})

		allInfoGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_chain_all_info",
			Help: "全网数据信息",
		}, []string{"phase"})

		ourAddressesGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_chain_addresses_info",
			Help: "worker 区块状态",
		}, []string{"latest_block_hash", "finalized_block_hash", "phase"})
	)

	registry.MustRegister(allInfoGaugeVec)
	registry.MustRegister(ourAddressesGaugeVec)

	resolvestart := time.Now()

	defer func() {
		lookupTime = time.Since(resolvestart).Seconds()
	}()

	var totalCountResponse TotalCountResponse
	// 请求链上所有数据
	statusCode, body, err := queryPost(AllWorkersQuery, target, logger)
	if err != nil {
		logger.Error("请求链上所有机器数据失败", "err", err)
		return 0.0, err
	}
	statusCodeGauge.Set(float64(statusCode))

	// 解析返回的 JSON 数据
	json_err := json.Unmarshal([]byte(string(body)), &totalCountResponse)
	if json_err != nil {
		fmt.Println("解析 JSON 失败:", json_err)
		return 0.0, err
	}

	allInfoGaugeVec.WithLabelValues("totalCount").Set(float64(totalCountResponse.Data.WorkersConnection.TotalCount))

	return lookupTime, nil
}

// 设置矿池中所有链上数据
func setPoolAllInfo(target string, module config.Module, registry *prometheus.Registry, logger *slog.Logger) (lookupTime float64, err error) {
	var (
		statusCodeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_pool_all_status_code",
			Help: "请求矿池中所有数据状态码",
		})

		poolInfoGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe__all_info",
			Help: "矿池中数据信息",
		}, []string{"phase"})
	)

	registry.MustRegister(statusCodeGauge)
	registry.MustRegister(poolInfoGaugeVec)

	resolvestart := time.Now()

	defer func() {
		lookupTime = time.Since(resolvestart).Seconds()
	}()

	// config.CHAINProbe.MonitorAddresses

	var poolResponse PoolResponse
	var query = fmt.Sprintf(PoolQuery, strings.Join(module.CHAIN.MonitorAddresses, `","`))
	// 请求链上所有数据
	statusCode, body, err := queryPost(query, target, logger)
	if err != nil {
		logger.Error("请求矿池中所有机器信息失败", "err", err)
		return 0.0, err
	}
	statusCodeGauge.Set(float64(statusCode))

	// 解析返回的 JSON 数据
	json_err := json.Unmarshal([]byte(string(body)), &poolResponse)
	if json_err != nil {
		fmt.Println("解析 JSON 失败:", json_err)
		return 0.0, err
	}

	allInfoGaugeVec.WithLabelValues("totalCount").Set(float64(totalCountResponse.Data.WorkersConnection.TotalCount))

	return lookupTime, nil
}

func ProbeCHAIN(ctx context.Context, target string, module config.Module, registry *prometheus.Registry, logger *slog.Logger) (success bool) {
	var (
		durationChainAllGaugeVec = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_chain_all_duration_seconds",
			Help: "探测链上数据耗费的总时间",
		})

		durationPoolAllGaugeVec = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_pool_all_duration_seconds",
			Help: "探测链上数据耗费的总时间",
		})
	)

	registry.MustRegister(durationChainAllGaugeVec)
	registry.MustRegister(durationPoolAllGaugeVec)

	// 如果目标没有加http或者https则加上
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}

	lookuptime, err := setChainAllInfo(target, registry, logger)
	if err != nil {
		logger.Error("设置节点指标失败", "err", err)
		return false
	}
	durationChainAllGaugeVec.Set(lookuptime)

	return true
}
