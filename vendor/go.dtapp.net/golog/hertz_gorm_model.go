package golog

import (
	"context"
	"log"
	"time"
)

// 结构体模型
type hertzGormLog struct {
	RequestID          string    `gorm:"index;comment:【日志】ID" json:"request_id,omitempty"`         //【日志】ID
	RequestTime        time.Time `gorm:"index;comment:【请求】Time" json:"request_time,omitempty"`     //【请求】Time
	RequestHost        string    `json:"request_host,omitempty"`                                   //【请求】Host
	RequestPath        string    `gorm:"index;comment:【请求】Path" json:"request_path,omitempty"`     //【请求】Path
	RequestQuery       string    `json:"request_query,omitempty"`                                  //【请求】Query Json
	RequestMethod      string    `gorm:"index;comment:【请求】Method" json:"request_method,omitempty"` //【请求】Method
	RequestScheme      string    `json:"request_scheme,omitempty"`                                 //【请求】Scheme
	RequestContentType string    `json:"request_content_type,omitempty"`                           //【请求】Content-Type
	RequestBody        string    `json:"request_body,omitempty"`                                   //【请求】Body Json
	RequestClientIP    string    `json:"request_client_ip,omitempty"`                              //【请求】ClientIP
	RequestUserAgent   string    `json:"request_user_agent,omitempty"`                             //【请求】User-Agent
	RequestHeader      string    `json:"request_header,omitempty"`                                 //【请求】Header Json
	RequestCostTime    int64     `json:"request_cost_time,omitempty"`                              //【请求】Cost
	ResponseTime       time.Time `json:"response_time,omitempty"`                                  //【响应】Time
	ResponseHeader     string    `json:"response_header,omitempty"`                                //【响应】Header Json
	ResponseStatusCode int       `json:"response_status_code,omitempty"`                           //【响应】StatusCode
	ResponseBody       string    `json:"response_data,omitempty"`                                  //【响应】Body Json
	GoVersion          string    `gorm:"comment:【程序】GoVersion" json:"go_version,omitempty"`        //【程序】GoVersion
	SdkVersion         string    `gorm:"comment:【程序】SdkVersion" json:"sdk_version,omitempty"`      //【程序】SdkVersion
	SystemInfo         string    `gorm:"comment:【系统】SystemInfo" json:"system_info,omitempty"`      //【系统】SystemInfo
}

// 创建模型
func (hg *HertzGorm) gormAutoMigrate(ctx context.Context) {
	if hg.gormConfig.stats == false {
		return
	}

	err := hg.gormClient.WithContext(ctx).
		Table(hg.gormConfig.tableName).
		AutoMigrate(&hertzGormLog{})
	if err != nil {
		log.Printf("创建模型：%s\n", err)
	}
}
