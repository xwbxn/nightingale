// Package models  设备厂商
// date : 2023-07-08 14:43
// desc : 设备厂商
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

const DOWNLOADNUM = 6

var defaultHeight = 25.0

// DeviceProducer  设备厂商。
// 说明:
// 表名:device_producer
// group: DeviceProducer
// version:2023-07-08 14:43
type DeviceProducer struct {
	Id               int64  `gorm:"column:ID;primaryKey" json:"id" `                                                                                                 //type:*int     comment:主键            version:2023-07-08 14:43
	Alias            string `gorm:"column:ALIAS" json:"alias" cn:"厂商简称" validate:"required"`                                                                         //type:string   comment:厂商简称        version:2023-07-08 14:43
	ChineseName      string `gorm:"column:CHINESE_NAME" json:"chinese_name" cn:"中文名称" validate:"omitempty"`                                                          //type:string   comment:中文名称        version:2023-07-08 14:43
	CompanyName      string `gorm:"column:COMPANY_NAME" json:"company_name" cn:"公司全称" validate:"required"`                                                           //type:string   comment:公司全称        version:2023-07-08 14:43
	Official         string `gorm:"column:OFFICIAL" json:"official" cn:"官方站点" validate:"omitempty"`                                                                  //type:string   comment:官方站点        version:2023-07-08 14:43
	IsDomestic       int64  `gorm:"column:IS_DOMESTIC" json:"is_domestic" cn:"是否国产" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"`                 //type:*int     comment:是否国产        version:2023-07-08 14:43
	IsDisplayChinese int64  `gorm:"column:IS_DISPLAY_CHINESE" json:"is_display_chinese" cn:"是否显示中文" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"` //type:*int     comment:是否显示中文    version:2023-07-08 14:43
	CreatedBy        string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                        //type:string   comment:创建人          version:2023-07-08 14:43
	CreatedAt        int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                        //type:*int     comment:创建时间        version:2023-07-08 14:43
	UpdatedBy        string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                        //type:string   comment:更新人          version:2023-07-08 14:43
	UpdatedAt        int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                        //type:*int     comment:更新时间        version:2023-07-08 14:43
}

// TableName 表名:device_producer，设备厂商。
// 说明:
func (d *DeviceProducer) TableName() string {
	return "device_producer"
}

type DeviceProducerNameVo struct {
	Id    int64  `gorm:"column:ID;primaryKey" json:"id" `                         //type:*int     comment:主键            version:2023-07-08 14:43
	Alias string `gorm:"column:ALIAS" json:"alias" cn:"厂商简称" validate:"required"` //type:string   comment:厂商简称        version:2023-07-08 14:43
}

// 条件查询
func DeviceProducerGets(ctx *ctx.Context, query string, limit, offset int) ([]DeviceProducer, error) {
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

	var lst []DeviceProducer
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DeviceProducerGetById(ctx *ctx.Context, id int64) (*DeviceProducer, error) {
	var obj *DeviceProducer
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有厂商名称
func FindNames(ctx *ctx.Context) ([]DeviceProducerNameVo, error) {
	var obj []DeviceProducerNameVo
	err := DB(ctx).Model(&DeviceProducer{}).Select("ID", "ALIAS").Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DeviceProducerGetsAll(ctx *ctx.Context) ([]DeviceProducer, error) {
	var lst []DeviceProducer
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加设备厂商
func (d *DeviceProducer) Add(ctx *ctx.Context) error {
	// 这里写DeviceProducer的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除设备厂商
func (d *DeviceProducer) Del(ctx *ctx.Context) error {
	// 这里写DeviceProducer的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新设备厂商
func (d *DeviceProducer) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DeviceProducer的业务逻辑，通过error返回错误
	// s := GetChinese(*d)
	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DeviceProducerCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DeviceProducer{}).Where(where, args...))
}
