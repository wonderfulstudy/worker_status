/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GetWorkersFromOwnerAddresses(ownerAddresses []string) (WorkersResponse, error) {
	// 定义目标结构体变量
	var response WorkersResponse
	// 生成 GraphQL 查询字符串
	query := fmt.Sprintf(AddressesQuery, strings.Join(ownerAddresses, `","`))
	body, err := QueryPost(query)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return WorkersResponse{}, err
	}
	// 解析 JSON 数据
	json_err := json.Unmarshal([]byte(string(body)), &response)
	if json_err != nil {
		fmt.Println("解析 JSON 失败:", json_err)
		return WorkersResponse{}, err
	}

	return response, nil
}

// 查询链上所有机器
func GetAllWorkersCount() (int, error) {
	var response AllWorkersResponse

	// 第一次发起请求，获取到全网上链机器总数
	// 获取全网机器状态功能待开发 todo
	// 生成 GraphQL 查询字符串
	query := fmt.Sprintf(AllWorkersQuery, "")
	body, err := QueryPost(query)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return 0, err
	}
	// 解析 JSON 数据
	json_err := json.Unmarshal([]byte(string(body)), &response)
	if json_err != nil {
		fmt.Println("解析 JSON 失败:", json_err)
		return 0, err
	}

	return response.Data.WorkersConnection.TotalCount, nil

}
