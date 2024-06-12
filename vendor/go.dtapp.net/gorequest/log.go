package gorequest

import (
	"time"
)

type LogResponse struct {
	SpanID             string    `json:"span_id"`    // 跨度编号
	TraceID            string    `json:"trace_id"`   // 追踪编号
	RequestID          string    `json:"request_id"` // 请求编号
	RequestTime        time.Time // 请求时间
	RequestUri         string    // 请求链接
	RequestUrl         string    // 请求链接
	RequestApi         string    // 请求接口
	RequestMethod      string    // 请求方式
	RequestParams      string    // 请求参数
	RequestHeader      string    // 请求头部
	RequestIP          string    // 请求请求IP
	ResponseHeader     string    // 返回头部
	ResponseStatusCode int       // 返回状态码
	ResponseBody       string    // 返回Json数据
	ResponseBodyJson   string    // 返回Json数据
	ResponseBodyXml    string    // 返回Xml数据
	ResponseTime       time.Time // 返回时间
	GoVersion          string    // 程序GoVersion
	SdkVersion         string    // 程序SdkVersion
}
