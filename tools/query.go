package tools

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var URL string = "https://dev-squid.cybros.network/graphql"
var RESPONSE WorkersResponse

type WorkersResponse struct {
	Data struct {
		Workers []struct {
			Address      string `json:"address"`
			OwnerAddress string `json:"ownerAddress"`
			Status       string `json:"status"`
			Uptime       int    `json:"uptime"`
		} `json:"workers"`
	} `json:"data"`
}

type AllWorkersResponse struct {
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

// 查询所有机器状态
var AllWorkersQuery string = `
	query MyQuery {
		workersConnection(orderBy: id_ASC, %s) {
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

// 根据矿池地址查询机器在线情况
var AddressQuery string = `
	query MyQuery {
		workers(where: {ownerAddress_eq: "%s"}) {
			address
			ownerAddress
			status
			uptime
		}
	}`

// 根据矿池地址列表查询机器在线情况
var AddressesQuery string = `
	query MyQuery {
		workers(where: {ownerAddress_in: ["%v"]}) {
			address
			ownerAddress
			status
			uptime
		}
	}`

func QueryPost(query string) ([]byte, error) {
	// 将查询字符串转换为 JSON 格式
	requestPayload := map[string]string{
		"query":         query,
		"operationName": "MyQuery",
		"variables":     "null",
	}
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return []byte{}, err
	}

	// 发起 POST 请求
	resp, post_err := http.Post(URL, "application/json", bytes.NewBuffer([]byte(payloadBytes)))
	if post_err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, resp_err := io.ReadAll(resp.Body)
	if resp_err != nil {
		return []byte{}, err
	}

	return body, nil
}
