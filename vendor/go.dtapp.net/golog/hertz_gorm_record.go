package golog

import (
	"context"
	"go.dtapp.net/gojson"
	"log"
)

// gormRecord 记录日志
func (hg *HertzGorm) gormRecord(ctx context.Context, data hertzGormLog) {
	if hg.gormConfig.stats == false {
		return
	}
	data.GoVersion = hg.config.GoVersion                         //【程序】GoVersion
	data.SdkVersion = hg.config.SdkVersion                       //【程序】SdkVersion
	data.SystemInfo = gojson.JsonEncodeNoError(hg.config.system) //【系统】SystemInfo

	err := hg.gormClient.WithContext(ctx).
		Table(hg.gormConfig.tableName).
		Create(&data).Error
	if err != nil {
		log.Printf("记录接口日志错误：%s", err)
		log.Printf("记录接口日志数据：%+v", data)
	}
}
