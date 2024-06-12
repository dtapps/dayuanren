package dayuanren

import (
	"context"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/codes"
)

type UserResponse struct {
	Errno  int64  `json:"errno"`  // 错误码，0代表成功，非0代表失败
	Errmsg string `json:"errmsg"` // 错误描述
	Data   struct {
		Id       int64  `json:"id"`       // userid
		Username string `json:"username"` // 名称
		Balance  string `json:"balance"`  // 余额
	} `json:"data,omitempty"`
}

type UserResult struct {
	Result UserResponse       // 结果
	Body   []byte             // 内容
	Http   gorequest.Response // 请求
}

func newUserResult(result UserResponse, body []byte, http gorequest.Response) *UserResult {
	return &UserResult{Result: result, Body: body, Http: http}
}

// User 查询用户信息
// https://www.showdoc.com.cn/dyr/9227004018562421
// https://www.kancloud.cn/boyanyun/boyanyun_huafei/3097251
func (c *Client) User(ctx context.Context, notMustParams ...gorequest.Params) (*UserResult, error) {

	// OpenTelemetry链路追踪
	ctx = c.TraceStartSpan(ctx, "index/user")
	defer c.TraceEndSpan()

	// 参数
	params := gorequest.NewParamsWith(notMustParams...)
	params.Set("userid", c.config.userID) // 账号ID

	// 请求
	request, err := c.request(ctx, "index/user", params)
	if err != nil {
		c.TraceSetStatus(codes.Error, err.Error())
		return newUserResult(UserResponse{}, request.ResponseBody, request), err
	}

	// 定义
	var response UserResponse
	err = gojson.Unmarshal(request.ResponseBody, &response)
	if err != nil {
		c.TraceSetStatus(codes.Error, err.Error())
	}
	return newUserResult(response, request.ResponseBody, request), err
}
