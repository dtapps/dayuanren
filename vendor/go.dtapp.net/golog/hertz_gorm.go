package golog

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/requestid"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gotime"
	"gorm.io/gorm"
	"strings"
	"time"
)

// HertzGorm 框架日志
type HertzGorm struct {
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

// HertzGormFun *HertzGorm 框架日志驱动
type HertzGormFun func() *HertzGorm

// NewHertzGorm 创建框架实例化
func NewHertzGorm(ctx context.Context, gormClient *gorm.DB, gormTableName string) (*HertzGorm, error) {

	hg := &HertzGorm{}
	hg.setConfig(ctx)

	if gormClient == nil {
		hg.gormConfig.stats = false
	} else {

		hg.gormClient = gormClient

		if gormTableName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			hg.gormConfig.tableName = gormTableName
		}

		hg.gormConfig.stats = true

		// 创建模型
		hg.gormAutoMigrate(ctx)

	}

	return hg, nil
}

// Middleware 中间件
func (hg *HertzGorm) Middleware() app.HandlerFunc {
	return func(c context.Context, h *app.RequestContext) {

		var log = hertzGormLog{}
		log.RequestTime = gotime.Current().Time

		start := time.Now()
		h.Next(c)
		end := time.Now()
		latency := end.Sub(start).Microseconds()

		log.ResponseTime = gotime.Current().Time

		// 输出路由日志
		hlog.CtxTracef(c, "status=%d cost=%d method=%s full_path=%s client_ip=%s host=%s",
			h.Response.StatusCode(),
			latency,
			h.Request.Header.Method(),
			h.Request.URI().PathOriginal(),
			h.ClientIP(),
			h.Request.Host(),
		)

		log.RequestID = requestid.Get(h)
		log.RequestHost = string(h.Request.Host())
		log.RequestPath = string(h.Request.URI().Path())
		log.RequestQuery = gojson.JsonEncodeNoError(gojson.ParseQueryString(string(h.Request.QueryString())))
		log.RequestMethod = string(h.Request.Header.Method())
		log.RequestScheme = string(h.Request.Scheme())

		log.RequestContentType = string(h.ContentType())

		if strings.Contains(log.RequestContentType, consts.MIMEApplicationHTMLForm) {
			log.RequestBody = gojson.JsonEncodeNoError(gojson.ParseQueryString(string(h.Request.Body())))
		} else if strings.Contains(log.RequestContentType, consts.MIMEMultipartPOSTForm) {
			log.RequestBody = string(h.Request.Body())
		} else if strings.Contains(log.RequestContentType, consts.MIMEApplicationJSON) {
			log.RequestBody = gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(h.Request.Body())))
		} else {
			log.RequestBody = string(h.Request.Body())
		}

		log.RequestClientIP = h.ClientIP()
		log.RequestUserAgent = string(h.UserAgent())

		requestHeader := make(map[string][]string)
		h.Request.Header.VisitAll(func(k, v []byte) {
			requestHeader[string(k)] = append(requestHeader[string(k)], string(v))
		})
		log.RequestHeader = gojson.JsonEncodeNoError(requestHeader)

		log.RequestCostTime = latency

		responseHeader := make(map[string][]string)
		h.Response.Header.VisitAll(func(k, v []byte) {
			responseHeader[string(k)] = append(responseHeader[string(k)], string(v))
		})
		log.ResponseHeader = gojson.JsonEncodeNoError(responseHeader)

		log.ResponseStatusCode = h.Response.StatusCode()

		if gojson.IsValidJSON(string(h.Response.Body())) {
			log.ResponseBody = gojson.JsonEncodeNoError(gojson.JsonDecodeNoError(string(h.Response.Body())))
		} else {
			log.ResponseBody = string(h.Response.Body())
		}

		go func() {
			hg.gormRecord(c, log)
		}()

	}
}
