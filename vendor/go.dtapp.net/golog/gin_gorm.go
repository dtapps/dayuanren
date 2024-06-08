package golog

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"gorm.io/gorm"
	"io/ioutil"
)

// GinGorm 框架日志
type GinGorm struct {
	gormClient *gorm.DB // 数据库驱动
	config     struct {
		GoVersion  string // go版本
		SdkVersion string // sdk版本
		system     struct {
			SystemVersion  string  `json:"system_version"`   // 系统版本
			SystemOs       string  `json:"system_os"`        // 系统类型
			SystemArch     string  `json:"system_arch"`      // 系统内核
			SystemInsideIP string  `json:"system_inside_ip"` // 内网IP
			SystemCpuModel string  `json:"system_cpu_model"` // CPU型号
			SystemCpuCores int     `json:"system_cpu_cores"` // CPU核数
			SystemCpuMhz   float64 `json:"system_cpu_mhz"`   // CPU兆赫
		}
	}
	gormConfig struct {
		stats     bool   // 状态
		tableName string // 表名
	}
}

// GinGormFun *GinGorm 框架日志驱动
type GinGormFun func() *GinGorm

// NewGinGorm 创建框架实例化
func NewGinGorm(ctx context.Context, gormClient *gorm.DB, gormTableName string) (*GinGorm, error) {

	gg := &GinGorm{}
	gg.setConfig(ctx)

	if gormClient == nil {
		gg.gormConfig.stats = false
	} else {

		gg.gormClient = gormClient

		if gormTableName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			gg.gormConfig.tableName = gormTableName
		}

		gg.gormConfig.stats = true

		// 创建模型
		gg.gormAutoMigrate(ctx)

	}

	return gg, nil
}

type ginGormBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w ginGormBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ginGormBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (gg *GinGorm) jsonUnmarshal(data string) (result any) {
	_ = gojson.Unmarshal([]byte(data), &result)
	return
}

// Middleware 中间件
func (gg *GinGorm) Middleware() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {

		// 开始时间
		startTime := gotime.Current().TimestampWithMillisecond()
		requestTime := gotime.Current().Time

		// 获取全部内容
		requestBody := gorequest.NewParams()
		queryParams := ginCtx.Request.URL.Query() // 请求URL参数
		for key, values := range queryParams {
			for _, value := range values {
				requestBody.Set(key, value)
			}
		}
		var dataMap map[string]any
		rawData, _ := ginCtx.GetRawData() // 请求内容参数
		if gojson.IsValidJSON(string(rawData)) {
			dataMap = gojson.JsonDecodeNoError(string(rawData))
		} else {
			dataMap = gojson.ParseQueryString(string(rawData))
		}
		for key, value := range dataMap {
			requestBody.Set(key, value)
		}

		// 重新赋值
		ginCtx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))

		blw := &ginGormBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
		ginCtx.Writer = blw

		// 处理请求
		ginCtx.Next()

		// 响应
		responseCode := ginCtx.Writer.Status()
		responseBody := blw.body.String()

		// 结束时间
		endTime := gotime.Current().TimestampWithMillisecond()
		responseTime := gotime.Current().Time

		go func() {

			// 记录
			gg.recordJson(ginCtx, requestTime, requestBody, responseTime, responseCode, responseBody, endTime-startTime, gorequest.ClientIp(ginCtx.Request))

		}()
	}
}
