package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// CreateJsonFile 创建 json 文件
func CreateJsonFile(filename string, jsonData interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	err = encoder.Encode(jsonData)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
