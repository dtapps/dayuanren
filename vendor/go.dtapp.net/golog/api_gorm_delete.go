package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gotime"
)

// GormDeleteData 删除N小时前数据
func (ag *ApiGorm) GormDeleteData(ctx context.Context, hour int64) error {
	return ag.GormDeleteDataCustom(ctx, ag.gormConfig.tableName, hour)
}

// GormDeleteDataCustom 删除N小时前数据
func (ag *ApiGorm) GormDeleteDataCustom(ctx context.Context, tableName string, hour int64) error {
	if ag.gormConfig.stats == false {
		return errors.New("没有驱动")
	}

	if tableName == "" {
		return errors.New("没有设置表名")
	}
	return ag.gormClient.WithContext(ctx).
		Table(tableName).
		Where("request_time < ?", gotime.Current().BeforeHour(hour).Format()).
		Delete(&apiGormLog{}).Error
}
