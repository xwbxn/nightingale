// Package models  字典类别表
// date : 2023-07-21 08:48
// desc : 字典类别表
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/toolkits/pkg/logger"
	"gorm.io/gorm"
)

// DictType  字典类别表。
// 说明:
// 表名:dict_type
// group: DictType
// version:2023-07-21 08:48
type DictType struct {
	Id        int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-07-21 08:48
	TypeCode  string         `gorm:"column:TYPE_CODE" json:"type_code" `                       //type:string   comment:字典编码    version:2023-08-01 14:10
	DictName  string         `gorm:"column:DICT_NAME" json:"dict_name" `                       //type:string   comment:字典名称    version:2023-07-21 08:48
	IsVisible string         `gorm:"column:IS_VISIBLE" json:"is_visible" `                     //type:string   comment:是否可见    version:2023-08-22 08:58
	Remark    string         `gorm:"column:REMARK" json:"remark" `                             //type:string   comment:备注        version:2023-07-21 08:48
	CreatedBy string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-07-21 08:48
	CreatedAt int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-07-21 08:48
	UpdatedBy string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-07-21 08:48
	UpdatedAt int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-07-21 08:48
	DeletedAt gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*int       comment:删除时间        version:2023-9-08 16:39
}

type DictTypeVo struct {
	Id        int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-07-21 08:48
	TypeCode  string         `gorm:"column:TYPE_CODE" json:"type_code" `                       //type:string   comment:字典编码    version:2023-08-01 14:10
	DictName  string         `gorm:"column:DICT_NAME" json:"dict_name" `                       //type:string   comment:字典名称    version:2023-07-21 08:48
	Remark    string         `gorm:"column:REMARK" json:"remark" `                             //type:string   comment:备注        version:2023-07-21 08:48
	CreatedBy string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-07-21 08:48
	CreatedAt int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-07-21 08:48
	UpdatedBy string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-07-21 08:48
	UpdatedAt int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-07-21 08:48
	DeletedAt gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*int       comment:删除时间        version:2023-9-08 16:39
}

// TableName 表名:dict_type，字典类别表。
// 说明:
func (d *DictType) TableName() string {
	return "dict_type"
}

// 条件查询
func DictTypeGets(ctx *ctx.Context, query string, limit, offset int) ([]DictTypeVo, error) {
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

	var lst []DictTypeVo
	err := session.Model(&DictType{}).Where("IS_VISIBLE = 'YES'").Find(&lst).Error

	return lst, err
}

// 按id查询
func DictTypeGetByTypeCodeCode(ctx *ctx.Context, typeCode string) (bool, error) {
	var dictType *DictType
	err := DB(ctx).Where("TYPE_CODE = ?", typeCode).Where("IS_VISIBLE = 'YES'").Find(&dictType).Error
	logger.Debug(dictType)
	if err != nil || dictType.Id == 0 {
		return false, err
	}

	return true, nil
}

// 按id查询
func DictTypeGetById(ctx *ctx.Context, id int64) (*DictTypeVo, error) {
	var obj *DictTypeVo
	err := DB(ctx).Model(&DictType{}).Where("IS_VISIBLE = 'YES'").Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DictTypeGetsAll(ctx *ctx.Context) ([]DictTypeVo, error) {
	var lst []DictTypeVo
	err := DB(ctx).Model(&DictType{}).Where("IS_VISIBLE = 'YES'").Find(&lst).Error

	return lst, err
}

// 增加字典类别表
func (d *DictType) Add(ctx *ctx.Context) error {
	// 这里写DictType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除字典类别表
func (d *DictTypeVo) Del(ctx *ctx.Context) error {
	// 这里写DictType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(&DictType{}).Where("IS_VISIBLE = 'YES'").Delete(d).Error
}

// 更新字典类别表
func (d *DictType) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DictType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Where("IS_VISIBLE = 'YES'").Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DictTypeCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DictType{}).Where("IS_VISIBLE = 'YES'").Where(where, args...))
}
