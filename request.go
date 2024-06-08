package dayuanren

import (
	"context"
	"go.dtapp.net/gorequest"
)

// 请求接口
func (c *Client) request(ctx context.Context, url string, param gorequest.Params) (gorequest.Response, error) {

	// 签名
	param.Set("sign", c.sign(param))

	// 设置请求地址
	c.httpClient.SetUri(c.config.apiURL + url)

	// 设置FORM格式
	c.httpClient.SetContentTypeForm()

	// 设置参数
	c.httpClient.SetParams(param)

	// 发起请求
	request, err := c.httpClient.Post(ctx)
	if err != nil {
		return gorequest.Response{}, err
	}

	return request, err
}
