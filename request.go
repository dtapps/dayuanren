package dayuanren

import (
	"context"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/attribute"
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

	// OpenTelemetry链路追踪
	c.TraceSetAttributes(attribute.String("http.url", c.config.apiURL+url))
	c.TraceSetAttributes(attribute.String("http.params", gojson.JsonEncodeNoError(param)))

	// 发起请求
	request, err := c.httpClient.Post(ctx)
	if err != nil {
		return gorequest.Response{}, err
	}

	return request, err
}
