package golog

import "time"

// 结构体
type apiSLog struct {
	TraceID            string              `json:"trace_id,omitempty"`
	RequestTime        time.Time           `json:"request_time,omitempty"`
	RequestUri         string              `json:"request_uri,omitempty"`
	RequestUrl         string              `json:"request_url,omitempty"`
	RequestApi         string              `json:"request_api,omitempty"`
	RequestMethod      string              `json:"request_method,omitempty"`
	RequestParams      map[string]any      `json:"request_params,omitempty"`
	RequestHeader      map[string]string   `json:"request_header,omitempty"`
	ResponseHeader     map[string][]string `json:"response_header,omitempty"`
	ResponseStatusCode int                 `json:"response_status_code,omitempty"`
	ResponseBody       map[string]any      `json:"response_body,omitempty"`
	ResponseTime       time.Time           `json:"response_time,omitempty,omitempty"`
}
