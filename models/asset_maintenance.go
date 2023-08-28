// Package models  资产维保
// date : 2023-07-23 09:44
// desc : 资产维保
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// AssetMaintenance  资产维保。
// 说明:
// 表名:asset_maintenance
// group: AssetMaintenance
// version:2023-07-23 09:44
type AssetMaintenance struct {
	Id                  int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-07-23 09:44
	MaintenanceType     string `gorm:"column:MAINTENANCE_TYPE" json:"maintenance_type" `         //type:string   comment:维保类型    version:2023-07-23 09:44
	MaintenanceProvider string `gorm:"column:MAINTENANCE_PROVIDER" json:"maintenance_provider" ` //type:string   comment:维保商      version:2023-07-23 09:44
	StartAt             int64  `gorm:"column:START_AT" json:"start_at" `                         //type:*int     comment:开始日期    version:2023-07-23 09:44
	FinishAt            int64  `gorm:"column:FINISH_AT" json:"finish_at" `                       //type:*int     comment:结束日期    version:2023-07-23 09:44
	MaintenancePeriod   string `gorm:"column:MAINTENANCE_PERIOD" json:"maintenance_period" `     //type:string   comment:维保期限    version:2023-07-23 09:44
	ProductionAt        int64  `gorm:"column:PRODUCTION_AT" json:"production_at" `               //type:*int     comment:出厂日期    version:2023-07-23 09:44
	Version             int64  `gorm:"column:VERSION" json:"version" `                           //type:*int     comment:版本号      version:2023-07-23 09:44
	CreatedBy           string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-07-23 09:44
	CreatedAt           int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-07-23 09:44
	UpdatedBy           string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-07-23 09:44
	UpdatedAt           int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-07-23 09:44
}

// TableName 表名:asset_maintenance，资产维保。
// 说明:
func (a *AssetMaintenance) TableName() string {
	return "asset_maintenance"
}

// 条件查询
func AssetMaintenanceGets(ctx *ctx.Context, query string, limit, offset int) ([]AssetMaintenance, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	// 这里使用列名的硬编码构造查询参数, 避免从前台传入造成注入风险
	if query != "" {
		q := "%" + query + "%"
		session = session.Where("id like ?", q)
	}

	var lst []AssetMaintenance
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func AssetMaintenanceGetById(ctx *ctx.Context, id int64) (*AssetMaintenance, error) {
	var obj *AssetMaintenance
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func AssetMaintenanceGetsAll(ctx *ctx.Context) ([]AssetMaintenance, error) {
	var lst []AssetMaintenance
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加资产维保
func (a *AssetMaintenance) Add(ctx *ctx.Context) error {
	// 这里写AssetMaintenance的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 删除资产维保
func (a *AssetMaintenance) Del(ctx *ctx.Context) error {
	// 这里写AssetMaintenance的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(a).Error
}

// 更新资产维保
func (a *AssetMaintenance) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写AssetMaintenance的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(a).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func AssetMaintenanceCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetMaintenance{}).Where(where, args...))
}
