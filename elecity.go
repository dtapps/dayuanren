package dayuanren

import (
	"context"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/codes"
)

type ElecityResponse struct {
	Errno  int64  `json:"errno"`  // 错误码，0代表成功，非0代表失败
	Errmsg string `json:"errmsg"` // 错误描述
	Data   []struct {
		Id        int64  `json:"id,omitempty"`
		CityName  string `json:"city_name"`  // 地区名称
		Sort      int64  `json:"sort"`       // 排序
		Initial   string `json:"initial"`    // 首字母
		NeedYtype int64  `json:"need_ytype"` // 是否三要素认证
		NeedCity  int64  `json:"need_city"`  // 是否需要选择城市（当此开关打开以后才有下面的城市列表）
		City      []struct {
			Id       int64  `json:"id,omitempty"`
			CityName string `json:"city_name"` // 城市名称
			Initial  string `json:"initial"`   // 首字母
		} `json:"city"` // 支持的地级市
	} `json:"data,omitempty"`
}

type ElecityResult struct {
	Result ElecityResponse    // 结果
	Body   []byte             // 内容
	Http   gorequest.Response // 请求
}

func newElecityResult(result ElecityResponse, body []byte, http gorequest.Response) *ElecityResult {
	return &ElecityResult{Result: result, Body: body, Http: http}
}

// Elecity 电费支持地区查询
// https://www.showdoc.com.cn/dyr/9227008514209156
// https://www.kancloud.cn/boyanyun/boyanyun_huafei/3097256
func (c *Client) Elecity(ctx context.Context, notMustParams ...gorequest.Params) (*ElecityResult, error) {

	// OpenTelemetry链路追踪
	ctx = c.TraceStartSpan(ctx, "index/elecity")
	defer c.TraceEndSpan()

	// 参数
	params := gorequest.NewParamsWith(notMustParams...)
	params.Set("userid", c.config.userID) // 账号ID

	// 请求
	request, err := c.request(ctx, "index/elecity", params)
	if err != nil {
		c.TraceSetStatus(codes.Error, err.Error())
		return newElecityResult(ElecityResponse{}, request.ResponseBody, request), err
	}

	// 定义
	var response ElecityResponse
	err = gojson.Unmarshal(request.ResponseBody, &response)
	if err != nil {
		c.TraceSetStatus(codes.Error, err.Error())
	}
	return newElecityResult(response, request.ResponseBody, request), err
}
