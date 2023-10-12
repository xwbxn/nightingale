// Package models  设备类型
// date : 2023-07-08 11:44
// desc : 设备类型
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// DeviceType  设备类型。
// 说明:
// 表名:device_type
// group: DeviceType
// version:2023-07-08 11:44
type DeviceType struct {
	Id        int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-07-08 11:44
	Name      string         `gorm:"column:NAME" json:"name" `                                 //type:string   comment:名称        version:2023-07-08 11:44
	Types     int64          `gorm:"column:TYPES" json:"types" `                               //type:*int     comment:类别(1:设备类型;2:备件设备类型)    version:2023-08-21 11:31
	CreatedBy string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-07-08 11:44
	CreatedAt int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-07-08 11:44
	UpdatedBy string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-07-08 11:44
	UpdatedAt int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-07-08 11:44
	DeletedAt gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*int       comment:删除时间        version:2023-9-08 16:39
}

// TableName 表名:device_type，设备类型。
// 说明:
func (d *DeviceType) TableName() string {
	return "device_type"
}

// 条件查询
func DeviceTypeGets(ctx *ctx.Context, query string, limit, offset int) ([]DeviceType, error) {
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

	var lst []DeviceType
	err := session.Find(&lst).Error

	return lst, err
}

// 根据map查询
func DeviceTypeGetMap(ctx *ctx.Context, m map[string]interface{}, limit, offset int) ([]DeviceType, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	var lst []DeviceType
	err := session.Where(m).Find(&lst).Error

	return lst, err
}

// 按id查询
func DeviceTypeGetById(ctx *ctx.Context, id int64) (*DeviceType, error) {
	var obj *DeviceType
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按id查询
// func DeviceTypeIdGetByName(ctx *ctx.Context, name string) (int, error) {
// 	var obj *DeviceType
// 	err := DB(ctx).Take(&obj, id).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return obj, nil
// }

// 查询所有
func DeviceTypeGetsAll(ctx *ctx.Context) ([]DeviceType, error) {
	var lst []DeviceType
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加设备类型
func (d *DeviceType) Add(ctx *ctx.Context) error {
	// 这里写DeviceType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除设备类型
func (d *DeviceType) Del(ctx *ctx.Context) error {
	// 这里写DeviceType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 批量删除设备类型
func DeviceTypeBatchDel(ctx *ctx.Context, d []DeviceType) error {
	// 这里写DeviceType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新设备类型
func (d *DeviceType) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DeviceType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DeviceTypeCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DeviceType{}).Where(where, args...))
}

// 根据map统计个数
func DeviceTypeCountMap(ctx *ctx.Context, m map[string]interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DeviceType{}).Where(m))
}
