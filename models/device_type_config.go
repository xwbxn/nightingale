// Package models  设备类型表单配置表
// date : 2023-08-04 08:47
// desc : 设备类型表单配置表
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// DeviceTypeConfig  设备类型表单配置表。
// 说明:
// 表名:device_type_config
// group: DeviceTypeConfig
// version:2023-08-04 08:47
type DeviceTypeConfig struct {
	Id        int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键            version:2023-08-04 08:47
	Name      string `gorm:"column:NAME" json:"name" `                                 //type:string   comment:名称            version:2023-08-04 08:47
	Type      int64  `gorm:"column:TYPE" json:"type" `                                 //type:*int     comment:设备类型        version:2023-08-04 08:47
	TypeName  string `gorm:"column:TYPE_NAME" json:"type_name" `                       //type:string   comment:设备类型名称    version:2023-08-04 08:47
	Config    string `gorm:"column:CONFIG" json:"config" `                             //type:string   comment:表单属性配置    version:2023-08-04 08:47
	CreatedBy string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人          version:2023-08-04 08:47
	CreatedAt int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间        version:2023-08-04 08:47
	UpdatedBy string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人          version:2023-08-04 08:47
	UpdatedAt int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间        version:2023-08-04 08:47
}

// TableName 表名:device_type_config，设备类型表单配置表。
// 说明:
func (d *DeviceTypeConfig) TableName() string {
	return "device_type_config"
}

// 条件查询
func DeviceTypeConfigGets(ctx *ctx.Context, query string, limit, offset int) ([]DeviceTypeConfig, error) {
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

	var lst []DeviceTypeConfig
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DeviceTypeConfigGetById(ctx *ctx.Context, id int64) (*DeviceTypeConfig, error) {
	var obj *DeviceTypeConfig
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DeviceTypeConfigGetsAll(ctx *ctx.Context) ([]DeviceTypeConfig, error) {
	var lst []DeviceTypeConfig
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加设备类型表单配置表
func (d *DeviceTypeConfig) Add(ctx *ctx.Context) error {
	// 这里写DeviceTypeConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除设备类型表单配置表
func (d *DeviceTypeConfig) Del(ctx *ctx.Context) error {
	// 这里写DeviceTypeConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新设备类型表单配置表
func (d *DeviceTypeConfig) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DeviceTypeConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DeviceTypeConfigCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DeviceTypeConfig{}).Where(where, args...))
}
