package dayuanren

import "go.dtapp.net/gorequest"

// ConfigClient 配置
func (c *Client) ConfigClient(config *ClientConfig) {
	c.config.apiURL = config.ApiURL
	c.config.userID = config.UserID
	c.config.apiKey = config.ApiKey
}

func (c *Client) SetUserID(userID int64) *Client {
	c.config.userID = userID
	return c
}

func (c *Client) SetApiKey(apiKey string) *Client {
	c.config.apiKey = apiKey
	return c
}

// SetClientIP 配置
func (c *Client) SetClientIP(clientIP string) *Client {
	c.config.clientIP = clientIP
	if c.httpClient != nil {
		c.httpClient.SetClientIP(clientIP)
	}
	return c
}

// SetLogFun 设置日志记录函数
func (c *Client) SetLogFun(logFun gorequest.LogFunc) {
	if c.httpClient != nil {
		c.httpClient.SetLogFunc(logFun)
	}
}
