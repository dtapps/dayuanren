package golog

import (
	"context"
	"log"
	"time"
)

// 结构体模型
type apiGormLog struct {
	LogID              int64     `gorm:"primaryKey;comment:【记录】编号" json:"log_id,omitempty"`     //【记录】编号
	TraceID            string    `gorm:"index;comment:【系统】跟踪编号" json:"trace_id,omitempty"`      //【系统】跟踪编号
	RequestTime        time.Time `gorm:"index;comment:【请求】时间" json:"request_time,omitempty"`    //【请求】时间
	RequestUri         string    `gorm:"comment:【请求】链接" json:"request_uri,omitempty"`           //【请求】链接
	RequestUrl         string    `gorm:"comment:【请求】链接" json:"request_url,omitempty"`           //【请求】链接
	RequestApi         string    `gorm:"index;comment:【请求】接口" json:"request_api,omitempty"`     //【请求】接口
	RequestMethod      string    `gorm:"index;comment:【请求】方式" json:"request_method,omitempty"`  //【请求】方式
	RequestParams      string    `gorm:"comment:【请求】参数" json:"request_params,omitempty"`        //【请求】参数
	RequestHeader      string    `gorm:"comment:【请求】头部" json:"request_header,omitempty"`        //【请求】头部
	RequestIP          string    `gorm:"comment:【请求】请求IP" json:"request_ip,omitempty"`          //【请求】请求IP
	ResponseHeader     string    `gorm:"comment:【返回】头部" json:"response_header,omitempty"`       //【返回】头部
	ResponseStatusCode int       `gorm:"comment:【返回】状态码" json:"response_status_code,omitempty"` //【返回】状态码
	ResponseBody       string    `gorm:"comment:【返回】数据" json:"response_body,omitempty"`         //【返回】数据
	ResponseTime       time.Time `gorm:"index;comment:【返回】时间" json:"response_time,omitempty"`   //【返回】时间
	GoVersion          string    `gorm:"comment:【程序】GoVersion" json:"go_version,omitempty"`     //【程序】GoVersion
	SdkVersion         string    `gorm:"comment:【程序】SdkVersion" json:"sdk_version,omitempty"`   //【程序】SdkVersion
	SystemInfo         string    `gorm:"comment:【系统】SystemInfo" json:"system_info,omitempty"`   //【系统】SystemInfo
}

// 创建模型
func (ag *ApiGorm) gormAutoMigrate(ctx context.Context) {
	if ag.gormConfig.stats == false {
		return
	}

	err := ag.gormClient.WithContext(ctx).
		Table(ag.gormConfig.tableName).
		AutoMigrate(&apiGormLog{})
	if err != nil {
		log.Printf("创建模型：%s\n", err)
	}
}
