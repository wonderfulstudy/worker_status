package prober

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// LoadPrivateKey 加载私钥
func loadPrivateKey(file string) (ssh.Signer, error) {
	key, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.New("无法读取私钥文件: " + err.Error())
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, errors.New("解析私钥失败: " + err.Error())
	}
	return signer, nil
}

func getTime(target string) (string, error) {
	key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		return "", errors.New("无法读取私钥文件: " + err.Error())
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", errors.New("解析私钥失败: " + err.Error())
	}

	parsedURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	var host string
	if strings.Contains(parsedURL.Host, ":") {
		host = strings.Split(parsedURL.Host, ":")[0]
	}

	sshConfig := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 不验证主机密钥（生产环境请替换）
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, 22)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return "", errors.New("SSH连接失败: " + err.Error())
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", errors.New("创建会话失败: " + err.Error())
	}
	defer session.Close()

	output, err := session.Output("date +%s")
	if err != nil {
		return "", errors.New("执行命令失败: " + err.Error())
	}

	var timestamp int64
	_, err = fmt.Sscanf(string(output), "%d", &timestamp)
	if err != nil {
		return "", errors.New("解析时间失败: " + err.Error())
	}

	result := strconv.FormatInt(timestamp, 10)

	return result, nil
}
