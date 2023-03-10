package utils

import (
	"fmt"
	"reflect"
)

// FillStruct 使用 map 填充结构体
func FillStruct(m map[string]interface{}, s interface{}) error {
	// 检查 s 是否为结构体变量的指针
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("FillStruct requires a pointer to a struct")
	}

	// 遍历 map，将值填充到结构体中
	for field, value := range m {
		fieldValue := v.Elem().FieldByName(field)
		if fieldValue.IsValid() && fieldValue.CanSet() {
			fieldType := fieldValue.Type()

			// 将值转换为结构体字段的类型
			newVal := reflect.ValueOf(value).Convert(fieldType)

			// 将值设置到结构体字段中
			fieldValue.Set(newVal)
		}
	}

	return nil
}
