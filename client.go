package dayuanren

import (
	"errors"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/trace"
)

// ClientConfig 实例配置
type ClientConfig struct {
	ClientIP string // 客户端IP
	ApiURL   string // 接口地址
	UserID   int64  // 商户ID
	ApiKey   string // 秘钥
}

// Client 实例
type Client struct {
	httpClient *gorequest.App
	config     struct {
		clientIP string // 客户端IP
		apiURL   string // 接口地址
		userID   int64  // 商户ID
		apiKey   string // 秘钥
	}
	trace bool       // OpenTelemetry链路追踪
	span  trace.Span // OpenTelemetry链路追踪
}

// NewClient 创建实例化
func NewClient(config *ClientConfig) (*Client, error) {

	c := &Client{}
	if config.ApiURL == "" {
		return nil, errors.New("需要配置ApiURL")
	}

	c.httpClient = gorequest.NewHttp()

	c.config.clientIP = config.ClientIP
	c.config.apiURL = config.ApiURL
	c.config.userID = config.UserID
	c.config.apiKey = config.ApiKey

	c.trace = true
	return c, nil
}
