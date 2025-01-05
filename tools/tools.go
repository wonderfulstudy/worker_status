/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetWorkersFromOwnerAddress(ownerAddress string) (WorkersResponse, error) {
	// 生成 GraphQL 查询字符串
	AddressQuery = fmt.Sprintf(AddressQuery, ownerAddress)

	// 将查询字符串转换为 JSON 格式
	requestPayload := map[string]string{
		"query":         AddressQuery,
		"operationName": "MyQuery",
		"variables":     "null",
	}
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("将查询字符串转换为 JSON 失败:")
		return WorkersResponse{}, err
	}

	// 发起 POST 请求
	resp, post_err := http.Post(URL, "application/json", bytes.NewBuffer([]byte(payloadBytes)))
	if post_err != nil {
		fmt.Println("发起请求失败:", post_err)
		return WorkersResponse{}, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		fmt.Println("读取响应失败:", resp_err)
		return WorkersResponse{}, err
	}

	// 定义目标结构体变量
	var response WorkersResponse
	// 解析 JSON 数据
	json_err := json.Unmarshal([]byte(string(body)), &response)
	if json_err != nil {
		fmt.Println("解析 JSON 失败:", json_err)
		return WorkersResponse{}, err
	}

	return response, nil
}

func GetWorkersFromOwnerAddresses(ownerAddresses []string) (WorkersResponse, error) {
	// 生成 GraphQL 查询字符串
	query := fmt.Sprintf(AddressesQuery, strings.Join(ownerAddresses, `","`))
	fmt.Println(query)
	// 将查询字符串转换为 JSON 格式
	requestPayload := map[string]string{
		"query":         query,
		"operationName": "MyQuery",
		"variables":     "null",
	}
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Println("将查询字符串转换为 JSON 失败:")
		return WorkersResponse{}, err
	}

	// 发起 POST 请求
	resp, post_err := http.Post(URL, "application/json", bytes.NewBuffer([]byte(payloadBytes)))
	if post_err != nil {
		fmt.Println("发起请求失败:", post_err)
		return WorkersResponse{}, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		fmt.Println("读取响应失败:", resp_err)
		return WorkersResponse{}, err
	}

	// 定义目标结构体变量
	var response WorkersResponse
	// 解析 JSON 数据
	json_err := json.Unmarshal([]byte(string(body)), &response)
	if json_err != nil {
		fmt.Println("解析 JSON 失败:", json_err)
		return WorkersResponse{}, err
	}

	return response, nil
}
