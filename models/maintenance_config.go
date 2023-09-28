// Package models  维保服务项配置
// date : 2023-08-01 14:17
// desc : 维保服务项配置
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// MaintenanceServiceConfig  维保服务项配置。
// 说明:
// 表名:maintenance_service_config
// group: MaintenanceServiceConfig
// version:2023-08-01 14:17
type MaintenanceServiceConfig struct {
	Id                int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键            version:2023-08-01 14:17
	MaintenanceId     int64  `gorm:"column:MAINTENANCE_ID" json:"maintenance_id" `             //type:*int     comment:维保ID          version:2023-08-01 14:17
	ServiceOptionCode string `gorm:"column:SERVICE_OPTION_CODE" json:"service_option_code" `   //type:string   comment:服务选项编码    version:2023-08-01 14:17
	ServiceOptionKey  string `gorm:"column:SERVICE_OPTION_KEY" json:"service_option_key" `     //type:string   comment:服务选项标签    version:2023-08-01 14:17
	ServiceObjCode    string `gorm:"column:SERVICE_OBJ_CODE" json:"service_obj_code" `         //type:string   comment:服务对象编码    version:2023-08-01 14:17
	ServiceObjKey     string `gorm:"column:SERVICE_OBJ_KEY" json:"service_obj_key" `           //type:string   comment:服务对象标签    version:2023-08-01 14:17
	Deadline          int64  `gorm:"column:DEADLINE" json:"deadline" `                         //type:*int   comment:服务截止时间    version:2023-08-01 14:17
	CreatedBy         string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人          version:2023-08-01 14:17
	CreatedAt         int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间        version:2023-08-01 14:17
	UpdatedBy         string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人          version:2023-08-01 14:17
	UpdatedAt         int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间        version:2023-08-01 14:17
}

// TableName 表名:maintenance_service_config，维保服务项配置。
// 说明:
func (m *MaintenanceServiceConfig) TableName() string {
	return "maintenance_service_config"
}

// 条件查询
func MaintenanceServiceConfigGets(ctx *ctx.Context, query string, limit, offset int) ([]MaintenanceServiceConfig, error) {
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

	var lst []MaintenanceServiceConfig
	err := session.Find(&lst).Error

	return lst, err
}

// 按maintenance_id查询
func MaintenanceServiceConfigGetByMaintId(ctx *ctx.Context, maintId int64) ([]MaintenanceServiceConfig, error) {
	var lst []MaintenanceServiceConfig
	err := DB(ctx).Where("MAINTENANCE_ID = ?", maintId).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	return lst, nil
}

// 按id查询
func MaintenanceServiceConfigGetById(ctx *ctx.Context, id int64) (*MaintenanceServiceConfig, error) {
	var obj *MaintenanceServiceConfig
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func MaintenanceServiceConfigGetsAll(ctx *ctx.Context) ([]MaintenanceServiceConfig, error) {
	var lst []MaintenanceServiceConfig
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加维保服务项配置
func (m *MaintenanceServiceConfig) Add(ctx *ctx.Context) error {
	// 这里写MaintenanceServiceConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(m).Error
}

// 删除维保服务项配置
func (m *MaintenanceServiceConfig) Del(ctx *ctx.Context) error {
	// 这里写MaintenanceServiceConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(m).Error
}

// 更新维保服务项配置
func (m *MaintenanceServiceConfig) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写MaintenanceServiceConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(m).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func MaintenanceServiceConfigCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&MaintenanceServiceConfig{}).Where(where, args...))
}
