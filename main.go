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
)

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
	flag := true
	for flag {
		flag = false
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
	}

	flag = true
	for flag {
		flag = false
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
	}

	viper.Set("proxies", proxies)
	viper.Set("proxy-groups.0.proxies", proxies)
	err = viper.WriteConfigAs("3Q.yml")
	if err != nil {
		fmt.Println(err)
	}

	de("end")
}
