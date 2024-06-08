package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gorequest"
	"gorm.io/gorm"
)

// ApiGorm 接口日志
type ApiGorm struct {
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

// ApiGormFun 接口日志驱动
type ApiGormFun func() *ApiGorm

// NewApiGorm 创建接口实例化
func NewApiGorm(ctx context.Context, gormClient *gorm.DB, gormTableName string) (*ApiGorm, error) {

	gl := &ApiGorm{}

	gl.setConfig(ctx)

	if gormClient == nil {
		gl.gormConfig.stats = false
	} else {

		gl.gormClient = gormClient

		if gormTableName == "" {
			return nil, errors.New("没有设置表名")
		} else {
			gl.gormConfig.tableName = gormTableName
		}

		gl.gormConfig.stats = true

		// 创建模型
		gl.gormAutoMigrate(ctx)

	}

	return gl, nil
}

// Middleware 中间件
func (ag *ApiGorm) Middleware(ctx context.Context, request gorequest.Response) {
	if ag.gormConfig.stats {
		ag.gormMiddleware(ctx, request)
	}
}

// MiddlewareXml 中间件
func (ag *ApiGorm) MiddlewareXml(ctx context.Context, request gorequest.Response) {
	if ag.gormConfig.stats {
		ag.gormMiddlewareXml(ctx, request)
	}
}

// MiddlewareCustom 中间件
func (ag *ApiGorm) MiddlewareCustom(ctx context.Context, api string, request gorequest.Response) {
	if ag.gormConfig.stats {
		ag.gormMiddlewareCustom(ctx, api, request)
	}
}
