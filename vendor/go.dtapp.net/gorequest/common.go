package gorequest

import (
	"go.dtapp.net/gojson"
	"go.dtapp.net/gostring"
)

// GetParamsString 获取参数字符串
func GetParamsString(src interface{}) string {
	switch src.(type) {
	case string:
		return src.(string)
	case int, int8, int32, int64:
	case uint8, uint16, uint32, uint64:
	case float32, float64:
		return gostring.ToString(src)
	}
	data, _ := gojson.Marshal(src)
	return string(data)
}
