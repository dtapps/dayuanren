package golog

import (
	"context"
	"errors"
	"go.dtapp.net/gotime"
)

// GormDeleteData 删除N小时前数据
func (hg *HertzGorm) GormDeleteData(ctx context.Context, hour int64) error {
	return hg.GormDeleteDataCustom(ctx, hg.gormConfig.tableName, hour)
}

// GormDeleteDataCustom 删除N小时前数据
func (hg *HertzGorm) GormDeleteDataCustom(ctx context.Context, tableName string, hour int64) error {
	if hg.gormConfig.stats == false {
		return errors.New("没有驱动")
	}

	if tableName == "" {
		return errors.New("没有设置表名")
	}
	return hg.gormClient.WithContext(ctx).
		Table(tableName).
		Where("request_time < ?", gotime.Current().BeforeHour(hour).Format()).
		Delete(&apiGormLog{}).Error
}
