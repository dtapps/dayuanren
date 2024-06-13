package dayuanren

import (
	"context"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/codes"
)

type CancelResponse struct {
	Errno  int64    `json:"errno"`  // 错误码，0代表成功，非0代表失败
	Errmsg string   `json:"errmsg"` // 错误描述
	Data   struct{} `json:"data"`
}

type CancelResult struct {
	Result CancelResponse     // 结果
	Body   []byte             // 内容
	Http   gorequest.Response // 请求
}

func newCancelResult(result CancelResponse, body []byte, http gorequest.Response) *CancelResult {
	return &CancelResult{Result: result, Body: body, Http: http}
}

// Cancel 退单申请
// out_trade_num = 商户订单号；多个用英文,分割
// https://www.kancloud.cn/boyanyun/boyanyun_huafei/3182909
func (c *Client) Cancel(ctx context.Context, outTradeNums string, notMustParams ...gorequest.Params) (*CancelResult, error) {

	// OpenTelemetry链路追踪
	ctx = c.TraceStartSpan(ctx, "index/cancel")
	defer c.TraceEndSpan()

	// 参数
	params := gorequest.NewParamsWith(notMustParams...)
	params.Set("userid", c.config.userID)      // 账户ID
	params.Set("out_trade_nums", outTradeNums) // 商户订单号；多个用英文,分割

	// 请求
	request, err := c.request(ctx, "index/cancel", params)
	if err != nil {
		c.TraceSetStatus(codes.Error, err.Error())
		c.TraceRecordError(err)
		return newCancelResult(CancelResponse{}, request.ResponseBody, request), err
	}

	// 定义
	var response CancelResponse
	err = gojson.Unmarshal(request.ResponseBody, &response)
	if err != nil {
		c.TraceSetStatus(codes.Error, err.Error())
		c.TraceRecordError(err)
	}
	return newCancelResult(response, request.ResponseBody, request), err
}
