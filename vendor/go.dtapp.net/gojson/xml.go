package gojson

import (
	xj "github.com/basgys/goxml2json"
	"strings"
)

// XmlDecodeNoError xml字符串转结构体，不报错
func XmlDecodeNoError(b []byte) map[string]interface{} {
	xtj := strings.NewReader(string(b))
	jtx, _ := xj.Convert(xtj)
	var data map[string]interface{}
	_ = Unmarshal(jtx.Bytes(), &data)
	return data
}

// XmlEncodeNoError 结构体转json字符串，不报错
func XmlEncodeNoError(data interface{}) string {
	return JsonEncodeNoError(data)
}
