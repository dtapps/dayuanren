package golog

import "time"

// 结构体
type ginSLog struct {
	TraceID       string              `json:"trace_id,omitempty"`       //【系统】跟踪编号
	RequestTime   time.Time           `json:"request_time,omitempty"`   //【请求】时间
	RequestUri    string              `json:"request_uri,omitempty"`    //【请求】请求链接 域名+路径+参数
	RequestUrl    string              `json:"request_url,omitempty"`    //【请求】请求链接 域名+路径
	RequestApi    string              `json:"request_api,omitempty"`    //【请求】请求接口 路径
	RequestMethod string              `json:"request_method,omitempty"` //【请求】请求方式
	RequestProto  string              `json:"request_proto,omitempty"`  //【请求】请求协议
	RequestBody   map[string]any      `json:"request_body,omitempty"`   //【请求】请求参数
	RequestIP     string              `json:"request_ip,omitempty"`     //【请求】请求客户端IP
	RequestHeader map[string][]string `json:"request_header,omitempty"` //【请求】请求头
	ResponseTime  time.Time           `json:"response_time,omitempty"`  //【返回】时间
	ResponseCode  int                 `json:"response_code,omitempty"`  //【返回】状态码
	ResponseData  string              `json:"response_data,omitempty"`  //【返回】数据
	CostTime      int64               `json:"cost_time,omitempty"`      //【系统】花费时间
}
