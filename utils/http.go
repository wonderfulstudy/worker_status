/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package utils

import (
	"bytes"
	"net/http"
)

var url = "https://dev-squid.cybros.network/graphql"

func postRquest(data []byte) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
