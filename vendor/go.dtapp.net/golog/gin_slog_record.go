package golog

import (
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"time"
)

// record 记录日志
func (gl *GinSLog) record(msg string, data ginSLog) {
	gl.slog.client.WithTraceIdStr(data.TraceID).Info(msg,
		"request_time", data.RequestTime,
		"request_uri", data.RequestUri,
		"request_url", data.RequestUrl,
		"request_api", data.RequestApi,
		"request_method", data.RequestMethod,
		"request_proto", data.RequestProto,
		"request_body", data.RequestBody,
		"request_ip", data.RequestIP,
		"request_header", data.RequestHeader,
		"response_time", data.ResponseTime,
		"response_code", data.ResponseCode,
		"response_data", data.ResponseData,
		"cost_time", data.CostTime,
	)
}

func (gl *GinSLog) recordJson(ginCtx *gin.Context, requestTime time.Time, requestBody gorequest.Params, responseTime time.Time, responseCode int, responseBody string, costTime int64, requestIp string) {
	data := ginSLog{
		TraceID:       gotrace_id.GetGinTraceId(ginCtx),                             //【系统】跟踪编号
		RequestTime:   requestTime,                                                  //【请求】时间
		RequestUrl:    ginCtx.Request.RequestURI,                                    //【请求】请求链接
		RequestApi:    gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】请求接口
		RequestMethod: ginCtx.Request.Method,                                        //【请求】请求方式
		RequestProto:  ginCtx.Request.Proto,                                         //【请求】请求协议
		RequestBody:   requestBody,                                                  //【请求】请求参数
		RequestIP:     requestIp,                                                    //【请求】请求客户端IP
		RequestHeader: ginCtx.Request.Header,                                        //【请求】请求头
		ResponseTime:  responseTime,                                                 //【返回】时间
		ResponseCode:  responseCode,                                                 //【返回】状态码
		ResponseData:  responseBody,                                                 //【返回】数据
		CostTime:      costTime,                                                     //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】请求链接
	}
	gl.record("json", data)
}
