package dayuanren

import (
	"context"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/codes"
)

type TypecateResponse struct {
	Errno  int64  `json:"errno"`  // 错误码，0代表成功，非0代表失败
	Errmsg string `json:"errmsg"` // 错误描述
	Data   []struct {
		Id       int64  `json:"id"`        // 产品类型id
		TypeName string `json:"type_name"` // 产品类型名称
		Cate     []struct {
			Id   int64  `json:"id"`        // 分类ID
			Cate string `json:"type_name"` // 分类名称
			Type int64  `json:"type"`      // 产品类型ID
		} `json:"cate"` // 分类列表
	} `json:"data,omitempty"`
}

type TypecateResult struct {
	Result TypecateResponse   // 结果
	Body   []byte             // 内容
	Http   gorequest.Response // 请求
}

func newTypecateResult(result TypecateResponse, body []byte, http gorequest.Response) *TypecateResult {
	return &TypecateResult{Result: result, Body: body, Http: http}
}

// Typecate 获取产品类型和产品分类
// https://www.showdoc.com.cn/dyr/9227005390454727
// https://www.kancloud.cn/boyanyun/boyanyun_huafei/3097252
func (c *Client) Typecate(ctx context.Context, notMustParams ...gorequest.Params) (*TypecateResult, error) {

	// OpenTelemetry链路追踪
	ctx = c.TraceStartSpan(ctx, "index/typecate")
	defer c.TraceEndSpan()

	// 参数
	params := gorequest.NewParamsWith(notMustParams...)
	params.Set("userid", c.config.userID) // 商户ID

	// 请求
	request, err := c.request(ctx, "index/typecate", params)
	if err != nil {
		if c.trace {
			c.span.SetStatus(codes.Error, err.Error())
		}
		return newTypecateResult(TypecateResponse{}, request.ResponseBody, request), err
	}

	// 定义
	var response TypecateResponse
	err = gojson.Unmarshal(request.ResponseBody, &response)
	if err != nil && c.trace {
		c.span.SetStatus(codes.Error, err.Error())
	}
	return newTypecateResult(response, request.ResponseBody, request), err
}
