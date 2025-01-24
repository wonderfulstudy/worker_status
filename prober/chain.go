package prober

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
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
		workersConnection(where: {ownerAddress_in: ["%v"]}, orderBy: address_ASC) {
			totalCount
			pageInfo {
				endCursor
				hasNextPage
				hasPreviousPage
				startCursor
			}
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

// 从指定下标开始查询机器在线情况
var PoolAfterQuery string = `
	query MyQuery($after: String = "%s") {
		workersConnection(where: {ownerAddress_in: ["%v"]}, orderBy: address_ASC, after: $after) {
			totalCount
			pageInfo {
				endCursor
				hasNextPage
				hasPreviousPage
				startCursor
			}
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
		WorkersConnection struct {
			TotalCount int `json:"totalCount"`
			PageInfo   struct {
				EndCursor       string `json:"endCursor"`
				HasNextPage     bool   `json:"hasNextPage"`
				HasPreviousPage bool   `json:"hasPreviousPage"`
				StartCursor     string `json:"startCursor"`
			} `json:"pageInfo"`
			Edges []struct {
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

		chainInfoGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_chain_all_info",
			Help: "全网数据信息",
		}, []string{"phase"})
	)

	registry.MustRegister(statusCodeGauge)
	registry.MustRegister(chainInfoGaugeVec)

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

	chainInfoGaugeVec.WithLabelValues("totalCount").Set(float64(totalCountResponse.Data.WorkersConnection.TotalCount))

	return lookupTime, nil
}

func setPoolWorkerInfo(poolResponse PoolResponse, registry *prometheus.GaugeVec) (err error) {
	for _, edge := range poolResponse.Data.WorkersConnection.Edges {
		registry.WithLabelValues(
			edge.Node.Address,
			edge.Node.OwnerAddress,
			edge.Node.Status,
			"uptime").Set(float64(edge.Node.Uptime))
	}

	return nil
}

// 设置矿池中所有链上数据
func setPoolAllInfo(target string, module config.Module, registry *prometheus.Registry, logger *slog.Logger) (lookupTime float64, err error) {
	var (
		poolAllRegisteredGauge = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_pool_all_registered",
			Help: "矿池中数据信息",
		})

		workerInfoGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "probe_pool_worker_info",
			Help: "矿池中每台机器信息",
		}, []string{"address", "owner_address", "status", "phase"})
	)

	registry.MustRegister(poolAllRegisteredGauge)
	registry.MustRegister(workerInfoGaugeVec)

	resolvestart := time.Now()

	defer func() {
		lookupTime = time.Since(resolvestart).Seconds()
	}()

	var poolResponse PoolResponse

	// 首次请求链上数据，获取矿池中注册机器总数
	var query = fmt.Sprintf(PoolQuery, strings.Join(module.CHAIN.MonitorAddresses, `","`))
	// 请求链上所有数据
	_, body, err := queryPost(query, target, logger)
	if err != nil {
		logger.Error("发起POST请求矿池中所有机器信息失败", "err", err)
		return 0.0, err
	}

	// 解析返回的 JSON 数据
	err = json.Unmarshal([]byte(string(body)), &poolResponse)
	if err != nil {
		logger.Error("解析 JSON 失败:", "err", err)
		return 0.0, err
	}
	poolAllRegisteredGauge.Set(float64(poolResponse.Data.WorkersConnection.TotalCount))
	setPoolWorkerInfo(poolResponse, workerInfoGaugeVec)
	if poolResponse.Data.WorkersConnection.PageInfo.HasNextPage {
		startCursor := 100
		for {
			var tmpPoolResponse PoolResponse
			var query = fmt.Sprintf(PoolAfterQuery, strconv.Itoa(startCursor), strings.Join(module.CHAIN.MonitorAddresses, `","`))
			_, body, err := queryPost(query, target, logger)
			if err != nil {
				logger.Error("发起POST请求矿池中所有机器信息失败", "err", err)
				return 0.0, err
			}
			err = json.Unmarshal([]byte(string(body)), &tmpPoolResponse)
			if err != nil {
				logger.Error("解析 JSON 失败:", "err", err)
				return 0.0, err
			}
			setPoolWorkerInfo(tmpPoolResponse, workerInfoGaugeVec)
			if tmpPoolResponse.Data.WorkersConnection.PageInfo.HasNextPage {
				startCursor = startCursor + 100
				continue
			} else {
				break
			}
		}
	}

	return lookupTime, nil
}

func ProbeCHAIN(ctx context.Context, target string, module config.Module, registry *prometheus.Registry, logger *slog.Logger) (success bool) {
	var (
		allDurationGaugeVec = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_all_duration_seconds",
			Help: "探测链上数据耗费的总时间",
		})
	)

	registry.MustRegister(allDurationGaugeVec)

	// 如果目标没有加http或者https则加上
	if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
		target = "http://" + target
	}

	chainTime, err := setChainAllInfo(target, registry, logger)
	if err != nil {
		logger.Error("设置链上总信息指标失败", "err", err)
		return false
	}

	poolTime, err := setPoolAllInfo(target, module, registry, logger)
	if err != nil {
		logger.Error("设置矿池总信息指标失败", "err", err)
		return false
	}

	allDurationGaugeVec.Set(chainTime + poolTime)

	return true
}
