// Package models  设备上线下线记录表
// date : 2023-08-27 16:34
// desc : 设备上线下线记录表
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// DeviceOnline  设备上线下线记录表。
// 说明:
// 表名:device_online
// group: DeviceOnline
// version:2023-08-27 16:34
type DeviceOnline struct {
	Id            int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键             version:2023-08-27 16:34
	DeviceStatus  int64  `gorm:"column:DEVICE_STATUS" json:"device_status" `               //type:*int     comment:类型             version:2023-08-27 17:00
	AssetId       int64  `gorm:"column:ASSET_ID" json:"asset_id" `                         //type:*int     comment:资产ID           version:2023-08-27 16:34
	Description   string `gorm:"column:DESCRIPTION" json:"description" `                   //type:string   comment:说明             version:2023-08-27 16:34
	LineAt        int64  `gorm:"column:LINE_AT" json:"line_at" `                           //type:*int     comment:上线/下线日期    version:2023-08-27 16:34
	LineDirectory int64  `gorm:"column:LINE_DIRECTORY" json:"line_directory" `             //type:*int     comment:上线/下线目录    version:2023-08-27 17:10
	CreatedBy     string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人           version:2023-08-27 16:34
	CreatedAt     int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间         version:2023-08-27 16:34
	UpdatedBy     string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人           version:2023-08-27 16:34
	UpdatedAt     int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间         version:2023-08-27 16:34
}

// TableName 表名:device_online，设备上线下线记录表。
// 说明:
func (d *DeviceOnline) TableName() string {
	return "device_online"
}

// 条件查询
func DeviceOnlineGets(ctx *ctx.Context, query string, limit, offset int) ([]DeviceOnline, error) {
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

	var lst []DeviceOnline
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DeviceOnlineGetById(ctx *ctx.Context, id int64) (*DeviceOnline, error) {
	var obj *DeviceOnline
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DeviceOnlineGetsAll(ctx *ctx.Context) ([]DeviceOnline, error) {
	var lst []DeviceOnline
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加设备上线下线记录表
func (d *DeviceOnline) Add(ctx *ctx.Context) error {
	// 这里写DeviceOnline的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 批量增加设备上线下线记录表
func DeviceOnlineTxBatchAdd(tx *gorm.DB, d []DeviceOnline) error {
	// 这里写DeviceOnline的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Debug().Create(d).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 删除设备上线下线记录表
func (d *DeviceOnline) Del(ctx *ctx.Context) error {
	// 这里写DeviceOnline的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新设备上线下线记录表
func (d *DeviceOnline) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DeviceOnline的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DeviceOnlineCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DeviceOnline{}).Where(where, args...))
}
