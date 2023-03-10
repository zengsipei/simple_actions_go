package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func d(values ...interface{}) {
	for _, v := range values {
		fmt.Printf("%T - %v\n", v, v)
	}
}

func de(values ...interface{}) {
	d(values...)
	os.Exit(0)
}

// Request 发送 HTTP 请求并返回响应体的内容
func Request(method, remoteUrl string, body io.Reader, headers map[string]string) []byte {
	d(remoteUrl)
	tr := &http.Transport{
		//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// 创建一个 HTTP 客户端
	client := &http.Client{Transport: tr}

	// 创建一个 HTTP 请求
	req, err := http.NewRequest(method, remoteUrl, body)
	if err != nil {
		fmt.Println("Error occurred:", err)
		return nil
	}

	// 设置请求的 header
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error occurred:", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Request failed with status code %d\n", resp.StatusCode)
		return nil
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}
	if len(content) == 0 {
		fmt.Println("Content is empty")
		return nil
	}

	return content
}
