package golog

// 结构体
type ginSLogCustom struct {
	TraceID         string `json:"trace_id,omitempty"`          //【系统】跟踪编号
	RequestUri      string `json:"request_uri,omitempty"`       //【请求】请求链接 域名+路径+参数
	RequestUrl      string `json:"request_url,omitempty"`       //【请求】请求链接 域名+路径
	RequestApi      string `json:"request_api,omitempty"`       //【请求】请求接口 路径
	RequestMethod   string `json:"request_method,omitempty"`    //【请求】请求方式
	RequestProto    string `json:"request_proto,omitempty"`     //【请求】请求协议
	RequestUa       string `json:"request_ua,omitempty"`        //【请求】请求UA
	RequestReferer  string `json:"request_referer,omitempty"`   //【请求】请求referer
	RequestUrlQuery string `json:"request_url_query,omitempty"` //【请求】请求URL参数
	RequestHeader   string `json:"request_header,omitempty"`    //【请求】请求头
	RequestIP       string `json:"request_ip,omitempty"`        //【请求】请求客户端IP
	CustomID        string `json:"custom_id,omitempty"`         //【日志】自定义编号
	CustomType      string `json:"custom_type,omitempty"`       //【日志】自定义类型
	CustomContent   string `json:"custom_content,omitempty"`    //【日志】自定义内容
}
