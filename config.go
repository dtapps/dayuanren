package dayuanren

import (
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/trace"
)

// ConfigClient 配置
func (c *Client) ConfigClient(config *ClientConfig) {
	c.config.apiURL = config.ApiURL
	c.config.userID = config.UserID
	c.config.apiKey = config.ApiKey
}

// SetClientIP 配置
func (c *Client) SetClientIP(clientIP string) {
	if clientIP == "" {
		return
	}
	c.config.clientIP = clientIP
	if c.httpClient != nil {
		c.httpClient.SetClientIP(clientIP)
	}
}

// SetTracer 设置链路追踪
func (c *Client) SetTracer(tr trace.Tracer) {
	if c.httpClient != nil {
		c.httpClient.SetTracer(tr)
	}
}

// SetLogFun 设置日志记录函数
func (c *Client) SetLogFun(logFun gorequest.LogFunc) {
	if c.httpClient != nil {
		c.httpClient.SetLogFunc(logFun)
	}
}
