// Package models  设备厂商
// date : 2023-07-08 14:43
// desc : 设备厂商
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// DeviceProducer  设备厂商。
// 说明:
// 表名:device_producer
// group: DeviceProducer
// version:2023-07-08 14:43
type DeviceProducer struct {
	Id               int64          `gorm:"column:ID;primaryKey" json:"id" `                                                                                                 //type:*int     comment:主键            version:2023-07-08 14:43
	ProducerType     string         `gorm:"column:PRODUCER_TYPE" json:"producer_type" `                                                                                      //type:string   comment:厂商类型        version:2023-08-20 09:09
	Alias            string         `gorm:"column:ALIAS" json:"alias" cn:"厂商简称" validate:"required"`                                                                         //type:string   comment:厂商简称        version:2023-07-08 14:43
	ChineseName      string         `gorm:"column:CHINESE_NAME" json:"chinese_name" cn:"中文名称" validate:"omitempty"`                                                          //type:string   comment:中文名称        version:2023-07-08 14:43
	CompanyName      string         `gorm:"column:COMPANY_NAME" json:"company_name" cn:"公司全称" validate:"required"`                                                           //type:string   comment:公司全称        version:2023-07-08 14:43
	ServiceTel       string         `gorm:"column:SERVICE_TEL" json:"service_tel" `                                                                                          //type:string   comment:服务电话        version:2023-08-20 09:09
	ServiceEmail     string         `gorm:"column:SERVICE_EMAIL" json:"service_email" `                                                                                      //type:string   comment:服务邮箱        version:2023-08-20 09:09
	Country          string         `gorm:"column:COUNTRY" json:"country" `                                                                                                  //type:string   comment:国家            version:2023-08-20 09:09
	City             string         `gorm:"column:CITY" json:"city" `                                                                                                        //type:string   comment:城市            version:2023-08-20 09:09
	Address          string         `gorm:"column:ADDRESS" json:"address" `                                                                                                  //type:string   comment:地址            version:2023-08-20 09:09
	Fax              string         `gorm:"column:FAX" json:"fax" `                                                                                                          //type:string   comment:传真            version:2023-08-20 09:09
	ContactPerson    string         `gorm:"column:CONTACT_PERSON" json:"contact_person" `                                                                                    //type:string   comment:联系人          version:2023-08-20 09:09
	ContactNumber    string         `gorm:"column:CONTACT_NUMBER" json:"contact_number" `                                                                                    //type:string   comment:联系人电话      version:2023-08-20 09:09
	ContactEmail     string         `gorm:"column:CONTACT_EMAIL" json:"contact_email" `                                                                                      //type:string   comment:联系人邮箱      version:2023-08-20 09:09
	Official         string         `gorm:"column:OFFICIAL" json:"official" cn:"官方站点" validate:"omitempty"`                                                                  //type:string   comment:官方站点        version:2023-07-08 14:43
	IsDomestic       int64          `gorm:"column:IS_DOMESTIC" json:"is_domestic" cn:"是否国产" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"`                 //type:*int     comment:是否国产        version:2023-07-08 14:43
	IsDisplayChinese int64          `gorm:"column:IS_DISPLAY_CHINESE" json:"is_display_chinese" cn:"是否显示中文" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"` //type:*int     comment:是否显示中文    version:2023-07-08 14:43
	CreatedBy        string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                        //type:string   comment:创建人          version:2023-07-08 14:43
	CreatedAt        int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                        //type:*int     comment:创建时间        version:2023-07-08 14:43
	UpdatedBy        string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                        //type:string   comment:更新人          version:2023-07-08 14:43
	UpdatedAt        int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                        //type:*int     comment:更新时间        version:2023-07-08 14:43
	DeletedAt        gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`                                                                        //type:*int       comment:删除时间        version:2023-9-08 16:39
}

type ProduceroVo struct {
	Id               int64          `gorm:"column:ID;primaryKey" json:"id" `                                                                                                 //type:*int     comment:主键            version:2023-07-08 14:43
	ProducerType     string         `gorm:"column:PRODUCER_TYPE" json:"producer_type" `                                                                                      //type:string   comment:厂商类型        version:2023-08-20 09:09
	Alias            string         `gorm:"column:ALIAS" json:"alias" cn:"厂商简称" validate:"required"`                                                                         //type:string   comment:厂商简称        version:2023-07-08 14:43
	ChineseName      string         `gorm:"column:CHINESE_NAME" json:"chinese_name" cn:"中文名称" validate:"omitempty"`                                                          //type:string   comment:中文名称        version:2023-07-08 14:43
	CompanyName      string         `gorm:"column:COMPANY_NAME" json:"company_name" cn:"公司全称" validate:"required"`                                                           //type:string   comment:公司全称        version:2023-07-08 14:43
	Official         string         `gorm:"column:OFFICIAL" json:"official" cn:"官方站点" validate:"omitempty"`                                                                  //type:string   comment:官方站点        version:2023-07-08 14:43
	IsDomestic       int64          `gorm:"column:IS_DOMESTIC" json:"is_domestic" cn:"是否国产" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"`                 //type:*int     comment:是否国产        version:2023-07-08 14:43
	IsDisplayChinese int64          `gorm:"column:IS_DISPLAY_CHINESE" json:"is_display_chinese" cn:"是否显示中文" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"` //type:*int     comment:是否显示中文    version:2023-07-08 14:43
	CreatedBy        string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                        //type:string   comment:创建人          version:2023-07-08 14:43
	CreatedAt        int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                        //type:*int     comment:创建时间        version:2023-07-08 14:43
	UpdatedBy        string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                        //type:string   comment:更新人          version:2023-07-08 14:43
	UpdatedAt        int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                        //type:*int     comment:更新时间        version:2023-07-08 14:43
	DeletedAt        gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`                                                                        //type:*int       comment:删除时间        version:2023-9-08 16:39
}

type MaintenanceVo struct {
	Id               int64          `gorm:"column:ID;primaryKey" json:"id" `                                                                                                 //type:*int     comment:主键            version:2023-07-08 14:43
	ProducerType     string         `gorm:"column:PRODUCER_TYPE" json:"producer_type" `                                                                                      //type:string   comment:厂商类型        version:2023-08-20 09:09
	Alias            string         `gorm:"column:ALIAS" json:"alias" cn:"维保服务商简称" validate:"required"`                                                                      //type:string   comment:厂商简称        version:2023-07-08 14:43
	ChineseName      string         `gorm:"column:CHINESE_NAME" json:"chinese_name" cn:"中文名称" validate:"omitempty"`                                                          //type:string   comment:中文名称        version:2023-07-08 14:43
	CompanyName      string         `gorm:"column:COMPANY_NAME" json:"company_name" cn:"公司全称" validate:"required"`                                                           //type:string   comment:公司全称        version:2023-07-08 14:43
	ServiceTel       string         `gorm:"column:SERVICE_TEL" json:"service_tel" cn:"服务电话"`                                                                                 //type:string   comment:服务电话        version:2023-08-20 09:09
	ServiceEmail     string         `gorm:"column:SERVICE_EMAIL" json:"service_email" cn:"服务邮箱"`                                                                             //type:string   comment:服务邮箱        version:2023-08-20 09:09
	Country          string         `gorm:"column:COUNTRY" json:"country" cn:"国家"`                                                                                           //type:string   comment:国家            version:2023-08-20 09:09
	City             string         `gorm:"column:CITY" json:"city" cn:"城市"`                                                                                                 //type:string   comment:城市            version:2023-08-20 09:09
	Address          string         `gorm:"column:ADDRESS" json:"address" cn:"地址"`                                                                                           //type:string   comment:地址            version:2023-08-20 09:09
	Fax              string         `gorm:"column:FAX" json:"fax" cn:"传真"`                                                                                                   //type:string   comment:传真            version:2023-08-20 09:09
	ContactPerson    string         `gorm:"column:CONTACT_PERSON" json:"contact_person" cn:"联系人"`                                                                            //type:string   comment:联系人          version:2023-08-20 09:09
	ContactNumber    string         `gorm:"column:CONTACT_NUMBER" json:"contact_number" cn:"联系人电话"`                                                                          //type:string   comment:联系人电话      version:2023-08-20 09:09
	ContactEmail     string         `gorm:"column:CONTACT_EMAIL" json:"contact_email" cn:"联系人邮箱"`                                                                            //type:string   comment:联系人邮箱      version:2023-08-20 09:09
	Official         string         `gorm:"column:OFFICIAL" json:"official" cn:"官方站点" validate:"omitempty"`                                                                  //type:string   comment:官方站点        version:2023-07-08 14:43
	IsDomestic       int64          `gorm:"column:IS_DOMESTIC" json:"is_domestic" cn:"是否国产" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"`                 //type:*int     comment:是否国产        version:2023-07-08 14:43
	IsDisplayChinese int64          `gorm:"column:IS_DISPLAY_CHINESE" json:"is_display_chinese" cn:"是否显示中文" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"` //type:*int     comment:是否显示中文    version:2023-07-08 14:43
	CreatedBy        string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                        //type:string   comment:创建人          version:2023-07-08 14:43
	CreatedAt        int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                        //type:*int     comment:创建时间        version:2023-07-08 14:43
	UpdatedBy        string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                        //type:string   comment:更新人          version:2023-07-08 14:43
	UpdatedAt        int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                        //type:*int     comment:更新时间        version:2023-07-08 14:43
	DeletedAt        gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`                                                                        //type:*int       comment:删除时间        version:2023-9-08 16:39
}

type SupplierVo struct {
	Id               int64          `gorm:"column:ID;primaryKey" json:"id" `                                                                                                 //type:*int     comment:主键            version:2023-07-08 14:43
	ProducerType     string         `gorm:"column:PRODUCER_TYPE" json:"producer_type" `                                                                                      //type:string   comment:厂商类型        version:2023-08-20 09:09
	Alias            string         `gorm:"column:ALIAS" json:"alias" cn:"供应商简称" validate:"required"`                                                                        //type:string   comment:厂商简称        version:2023-07-08 14:43
	ChineseName      string         `gorm:"column:CHINESE_NAME" json:"chinese_name" cn:"中文名称" validate:"omitempty"`                                                          //type:string   comment:中文名称        version:2023-07-08 14:43
	CompanyName      string         `gorm:"column:COMPANY_NAME" json:"company_name" cn:"公司全称" validate:"required"`                                                           //type:string   comment:公司全称        version:2023-07-08 14:43
	ServiceTel       string         `gorm:"column:SERVICE_TEL" json:"service_tel" cn:"服务电话"`                                                                                 //type:string   comment:服务电话        version:2023-08-20 09:09
	ServiceEmail     string         `gorm:"column:SERVICE_EMAIL" json:"service_email" cn:"服务邮箱"`                                                                             //type:string   comment:服务邮箱        version:2023-08-20 09:09
	Country          string         `gorm:"column:COUNTRY" json:"country" cn:"国家"`                                                                                           //type:string   comment:国家            version:2023-08-20 09:09
	City             string         `gorm:"column:CITY" json:"city" cn:"城市"`                                                                                                 //type:string   comment:城市            version:2023-08-20 09:09
	Address          string         `gorm:"column:ADDRESS" json:"address" cn:"地址"`                                                                                           //type:string   comment:地址            version:2023-08-20 09:09
	Fax              string         `gorm:"column:FAX" json:"fax" cn:"传真"`                                                                                                   //type:string   comment:传真            version:2023-08-20 09:09
	ContactPerson    string         `gorm:"column:CONTACT_PERSON" json:"contact_person" cn:"联系人"`                                                                            //type:string   comment:联系人          version:2023-08-20 09:09
	ContactNumber    string         `gorm:"column:CONTACT_NUMBER" json:"contact_number" cn:"联系人电话"`                                                                          //type:string   comment:联系人电话      version:2023-08-20 09:09
	ContactEmail     string         `gorm:"column:CONTACT_EMAIL" json:"contact_email" cn:"联系人邮箱"`                                                                            //type:string   comment:联系人邮箱      version:2023-08-20 09:09
	Official         string         `gorm:"column:OFFICIAL" json:"official" cn:"官方站点" validate:"omitempty"`                                                                  //type:string   comment:官方站点        version:2023-07-08 14:43
	IsDomestic       int64          `gorm:"column:IS_DOMESTIC" json:"is_domestic" cn:"是否国产" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"`                 //type:*int     comment:是否国产        version:2023-07-08 14:43
	IsDisplayChinese int64          `gorm:"column:IS_DISPLAY_CHINESE" json:"is_display_chinese" cn:"是否显示中文" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"` //type:*int     comment:是否显示中文    version:2023-07-08 14:43
	CreatedBy        string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                        //type:string   comment:创建人          version:2023-07-08 14:43
	CreatedAt        int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                        //type:*int     comment:创建时间        version:2023-07-08 14:43
	UpdatedBy        string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                        //type:string   comment:更新人          version:2023-07-08 14:43
	UpdatedAt        int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                        //type:*int     comment:更新时间        version:2023-07-08 14:43
	DeletedAt        gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`                                                                        //type:*int       comment:删除时间        version:2023-9-08 16:39
}

type PartVo struct {
	Id               int64          `gorm:"column:ID;primaryKey" json:"id" `                                                                                                 //type:*int     comment:主键            version:2023-07-08 14:43
	ProducerType     string         `gorm:"column:PRODUCER_TYPE" json:"producer_type" `                                                                                      //type:string   comment:厂商类型        version:2023-08-20 09:09
	Alias            string         `gorm:"column:ALIAS" json:"alias" cn:"部件品牌" validate:"required"`                                                                         //type:string   comment:厂商简称        version:2023-07-08 14:43
	ChineseName      string         `gorm:"column:CHINESE_NAME" json:"chinese_name" cn:"中文名称" validate:"omitempty"`                                                          //type:string   comment:中文名称        version:2023-07-08 14:43
	CompanyName      string         `gorm:"column:COMPANY_NAME" json:"company_name" cn:"公司全称" validate:"required"`                                                           //type:string   comment:公司全称        version:2023-07-08 14:43
	ServiceTel       string         `gorm:"column:SERVICE_TEL" json:"service_tel" cn:"服务电话"`                                                                                 //type:string   comment:服务电话        version:2023-08-20 09:09
	ServiceEmail     string         `gorm:"column:SERVICE_EMAIL" json:"service_email" cn:"服务邮箱"`                                                                             //type:string   comment:服务邮箱        version:2023-08-20 09:09
	Country          string         `gorm:"column:COUNTRY" json:"country" cn:"国家"`                                                                                           //type:string   comment:国家            version:2023-08-20 09:09
	City             string         `gorm:"column:CITY" json:"city" cn:"城市"`                                                                                                 //type:string   comment:城市            version:2023-08-20 09:09
	Address          string         `gorm:"column:ADDRESS" json:"address" cn:"地址"`                                                                                           //type:string   comment:地址            version:2023-08-20 09:09
	Fax              string         `gorm:"column:FAX" json:"fax" cn:"传真"`                                                                                                   //type:string   comment:传真            version:2023-08-20 09:09
	ContactPerson    string         `gorm:"column:CONTACT_PERSON" json:"contact_person" cn:"联系人"`                                                                            //type:string   comment:联系人          version:2023-08-20 09:09
	ContactNumber    string         `gorm:"column:CONTACT_NUMBER" json:"contact_number" cn:"联系人电话"`                                                                          //type:string   comment:联系人电话      version:2023-08-20 09:09
	ContactEmail     string         `gorm:"column:CONTACT_EMAIL" json:"contact_email" cn:"联系人邮箱"`                                                                            //type:string   comment:联系人邮箱      version:2023-08-20 09:09
	Official         string         `gorm:"column:OFFICIAL" json:"official" cn:"官方站点" validate:"omitempty"`                                                                  //type:string   comment:官方站点        version:2023-07-08 14:43
	IsDomestic       int64          `gorm:"column:IS_DOMESTIC" json:"is_domestic" cn:"是否国产" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"`                 //type:*int     comment:是否国产        version:2023-07-08 14:43
	IsDisplayChinese int64          `gorm:"column:IS_DISPLAY_CHINESE" json:"is_display_chinese" cn:"是否显示中文" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"` //type:*int     comment:是否显示中文    version:2023-07-08 14:43
	CreatedBy        string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                        //type:string   comment:创建人          version:2023-07-08 14:43
	CreatedAt        int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                        //type:*int     comment:创建时间        version:2023-07-08 14:43
	UpdatedBy        string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                        //type:string   comment:更新人          version:2023-07-08 14:43
	UpdatedAt        int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                        //type:*int     comment:更新时间        version:2023-07-08 14:43
	DeletedAt        gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`                                                                        //type:*int       comment:删除时间        version:2023-9-08 16:39
}

// TableName 表名:device_producer，设备厂商。
// 说明:
func (d *DeviceProducer) TableName() string {
	return "device_producer"
}

type DeviceProducerNameVo struct {
	Id           int64  `gorm:"column:ID;primaryKey" json:"id" `                         //type:*int     comment:主键            version:2023-07-08 14:43
	ProducerType string `gorm:"column:PRODUCER_TYPE" json:"producer_type" `              //type:string   comment:厂商类型        version:2023-08-20 09:09
	Alias        string `gorm:"column:ALIAS" json:"alias" cn:"厂商简称" validate:"required"` //type:string   comment:厂商简称        version:2023-07-08 14:43
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
	err := DB(ctx).Debug().Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按id批量查询
func DeviceProducerGetByIds[T any](ctx *ctx.Context, ids []int64) ([]T, error) {
	var obj []T
	err := DB(ctx).Debug().Model(&DeviceProducer{}).Where("ID IN ?", ids).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按map查询
func DeviceProducerGetByMap[T any](ctx *ctx.Context, m map[string]interface{}) ([]T, error) {
	var obj []T
	err := DB(ctx).Model(&DeviceProducer{}).Where(m).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按map分页查询
func DeviceProducerGetByPage[T any](ctx *ctx.Context, m map[string]interface{}, limit, offset int) ([]T, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	var obj []T
	var err error
	query, queryOk := m["query"]
	if queryOk {
		delete(m, "query")
		q := "%" + query.(string) + "%"
		err = session.Model(&DeviceProducer{}).Where(m).Where("ALIAS LIKE ? OR CHINESE_NAME LIKE ? OR COMPANY_NAME LIKE ? OR CONTACT_PERSON LIKE ?", q, q, q, q).Find(&obj).Error
	} else {
		err = session.Model(&DeviceProducer{}).Where(m).Find(&obj).Error
	}
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有厂商名称
func FindNames(ctx *ctx.Context, prodType string) ([]DeviceProducerNameVo, error) {
	var obj []DeviceProducerNameVo
	err := DB(ctx).Debug().Model(&DeviceProducer{}).Select("ID", "ALIAS").Where("PRODUCER_TYPE = ?", prodType).Find(&obj).Error
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
func AddProd[T any](ctx *ctx.Context, t T) error {
	// 这里写DeviceProducer的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Debug().Create(&t).Error
}

// 删除设备厂商
func BatchAddProd(ctx *ctx.Context, d []DeviceProducer) error {
	// 这里写DeviceProducer的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 删除设备厂商
func (d *DeviceProducer) Del(ctx *ctx.Context) error {
	// 这里写DeviceProducer的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 删除设备厂商
func DeviceProducerBatchDel(ctx *ctx.Context, d []int64) error {
	// 这里写DeviceProducer的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Where("ID IN ?", d).Delete(&DeviceProducer{}).Error
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

// 根据map统计个数
func DeviceProducerCountMap(ctx *ctx.Context, m map[string]interface{}) (num int64, err error) {
	query, queryOk := m["query"]
	if queryOk {
		delete(m, "query")
		q := "%" + query.(string) + "%"
		count, err := Count(DB(ctx).Model(&DeviceProducer{}).Where(m).Where("ALIAS LIKE ? OR CHINESE_NAME LIKE ? OR COMPANY_NAME LIKE ? OR CONTACT_PERSON LIKE ?", q, q, q, q))
		m["query"] = query
		return count, err
	} else {
		return Count(DB(ctx).Model(&DeviceProducer{}).Where(m))
	}
}
