package dayuanren

import (
	"context"
	"go.dtapp.net/gorequest"
)

// 请求接口
func (c *Client) request(ctx context.Context, url string, param gorequest.Params) (gorequest.Response, error) {

	// 签名
	param.Set("sign", c.sign(param))

	// 创建请求
	client := gorequest.NewHttp()

	// 设置请求地址
	client.SetUri(c.config.apiURL + url)

	// 设置FORM格式
	client.SetContentTypeForm()

	// 设置参数
	client.SetParams(param)

	// 发起请求
	request, err := client.Post(ctx)
	if err != nil {
		return gorequest.Response{}, err
	}

	// 日志
	if c.gormLog.status {
		go c.gormLog.client.Middleware(ctx, request)
	}
	if c.mongoLog.status {
		go c.mongoLog.client.Middleware(ctx, request)
	}

	return request, err
}
