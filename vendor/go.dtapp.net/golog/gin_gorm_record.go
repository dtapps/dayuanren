package golog

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotrace_id"
	"go.dtapp.net/gourl"
	"log"
	"time"
)

// gormRecord 记录日志
func (gg *GinGorm) gormRecord(ctx context.Context, data ginGormLog) {
	if gg.gormConfig.stats == false {
		return
	}
	data.GoVersion = gg.config.GoVersion                         //【程序】GoVersion
	data.SdkVersion = gg.config.SdkVersion                       //【程序】SdkVersion
	data.SystemInfo = gojson.JsonEncodeNoError(gg.config.system) //【系统】SystemInfo

	err := gg.gormClient.WithContext(ctx).
		Table(gg.gormConfig.tableName).
		Create(&data).Error
	if err != nil {
		log.Printf("记录接口日志错误：%s\n", err)
		log.Printf("记录接口日志数据：%+v\n", data)
	}
}

func (gg *GinGorm) recordJson(ginCtx *gin.Context, requestTime time.Time, requestBody gorequest.Params, responseTime time.Time, responseCode int, responseBody string, costTime int64, requestIp string) {

	data := ginGormLog{
		TraceID:       gotrace_id.GetGinTraceId(ginCtx),                             //【系统】跟踪编号
		RequestTime:   requestTime,                                                  //【请求】时间
		RequestURL:    ginCtx.Request.RequestURI,                                    //【请求】链接
		RequestApi:    gourl.UriFilterExcludeQueryString(ginCtx.Request.RequestURI), //【请求】接口
		RequestMethod: ginCtx.Request.Method,                                        //【请求】方式
		RequestProto:  ginCtx.Request.Proto,                                         //【请求】协议
		RequestBody:   gojson.JsonEncodeNoError(requestBody),                        //【请求】参数
		RequestIP:     requestIp,                                                    //【请求】客户端IP
		RequestHeader: gojson.JsonEncodeNoError(ginCtx.Request.Header),              //【请求】头部
		ResponseTime:  responseTime,                                                 //【返回】时间
		ResponseCode:  responseCode,                                                 //【返回】状态码
		ResponseData:  responseBody,                                                 //【返回】数据
		CostTime:      costTime,                                                     //【系统】花费时间
	}
	if ginCtx.Request.TLS == nil {
		data.RequestUri = "http://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】链接
	} else {
		data.RequestUri = "https://" + ginCtx.Request.Host + ginCtx.Request.RequestURI //【请求】链接
	}

	gg.gormRecord(ginCtx, data)
}
