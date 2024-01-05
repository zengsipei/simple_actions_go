package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"github.com/uniplaces/carbon"
	"os"
	"simple_actions/utils"
	"strconv"
	"strings"
	"io"
	"net/http"
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

func main() {
	viper.SetConfigName("base")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	proxyGroups := viper.Get("proxy-groups").([]interface{})
	proxy := proxyGroups[0].(map[string]interface{})
	proxies := proxy["proxies"].([]interface{})
	now := carbon.Now()
	d(now)
	for {
		// v2rayshare.com
		filename := now.Format("20060102") + ".yaml"
		urlArr := []string{
			"https://v2rayshare.com/wp-content/uploads",
			strconv.Itoa(now.Year()),
			now.Format("01"),
			filename,
		}
		reqUrl := strings.Join(urlArr, "/")

		content := Request("GET", reqUrl, nil, nil)
		if content == nil {
			break
		}

		v2rayViper := viper.New()
		v2rayViper.SetConfigType("yaml")
		err = v2rayViper.ReadConfig(bytes.NewBuffer(content))
		if err != nil {
			d(err)
			v2rayViper.SetConfigName("v2rayshare")
			v2rayViper.AddConfigPath(".")
			err = v2rayViper.ReadInConfig()
			if err != nil {
				d(err)
				break
			}
		}

		err = v2rayViper.WriteConfigAs("v2rayshare.yml")
		if err != nil {
			d(err)
		}

		v2rayProxies := v2rayViper.Get("proxies").([]interface{})
		proxies = append(proxies, v2rayProxies...)
		break
	}
	for {
		isExpire := true
		// MrdoorVpn
		file, err := os.ReadFile("mrdoor.json")
		jsonData := make(map[string]interface{})
		if err == nil {
			err = json.Unmarshal(file, &jsonData)
			if err == nil {
				c, _ := carbon.CreateFromTimestamp(int64(jsonData["expire"].(float64))/1000, "Asia/Shanghai")
				d(c)
				isExpire = c.LessThan(now)
			}
		}
		if isExpire {
			// autosign
			body := []byte("device_id=C5" + utils.RandomStringNumber(11111, 99999) + "-3ED4-449C-806F-16167220F497")
			headers := map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			}
			body = Request("POST", "http://192.53.161.228/api/v2/autosignup", bytes.NewBuffer(body), headers)
			if body == nil {
				break
			}

			bodyJson := make(map[string]interface{})
			err = json.Unmarshal(body, &bodyJson)
			if err != nil {
				break
			}

			// account
			bodyJsonData := bodyJson["data"].(map[string]interface{})
			jsonData = utils.MergeMaps(jsonData, bodyJsonData)
			headers = map[string]string{
				"uuid": jsonData["uuid"].(string),
			}
			body = Request("GET", "http://192.53.161.228/api/v2/account", nil, headers)
			if body == nil {
				break
			}

			bodyJson = make(map[string]interface{})
			err = json.Unmarshal(body, &bodyJson)
			if err != nil {
				break
			}

			bodyJsonData = bodyJson["data"].(map[string]interface{})
			jsonData = utils.MergeMaps(jsonData, bodyJsonData)

			err = utils.CreateJsonFile("mrdoor.json", jsonData)
			if err != nil {
				//
			}
		}
		if jsonData["uuid"] == "" {
			break
		}

		//servers
		headers := map[string]string{
			"uuid": jsonData["uuid"].(string),
		}
		body := Request("GET", "http://192.53.161.228/api/v3/servers", nil, headers)
		if body == nil {
			break
		}

		bodyJson := make(map[string]interface{})
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			break
		}

		bodyJsonData := bodyJson["data"].([]interface{})
		file, err = os.ReadFile("type.json")
		typeData := make(map[string]interface{})
		if err == nil {
			err = json.Unmarshal(file, &typeData)
			if err == nil {
				//
			}
		}
		for _, v := range bodyJsonData {
			vMap := v.(map[string]interface{})
			if vMap["type"] == "Shadowsocks" {
				tempV := make(map[string]interface{})
				tempV["name"] = vMap["name"]
				tempV["server"] = vMap["host"]
				tempV["port"] = vMap["port"]
				tempV["type"] = vMap["ss"]
				tempV["cipher"] = vMap["method"]
				tempV["password"] = vMap["password"]

				proxies = append(proxies, tempV)
				continue
			}

			typeData[vMap["type"].(string)] = v
		}
		if typeData != nil {
			err := utils.CreateJsonFile("type.json", typeData)
			if err != nil {
			}
		}
		break
	}

	viper.Set("proxies", proxies)
	viper.Set("proxy-groups.0.proxies", proxies)
	err = viper.WriteConfigAs("3Q.yml")
	if err != nil {
		fmt.Println(err)
	}

	de("end")
}
