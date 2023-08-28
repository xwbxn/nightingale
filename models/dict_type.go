// Package models  字典类别表
// date : 2023-07-21 08:48
// desc : 字典类别表
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/toolkits/pkg/logger"
)

// DictType  字典类别表。
// 说明:
// 表名:dict_type
// group: DictType
// version:2023-07-21 08:48
type DictType struct {
	Id        int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-07-21 08:48
	DictCode  string `gorm:"column:DICT_CODE" json:"dict_code" `                       //type:string   comment:字典编码    version:2023-07-21 08:48
	DictName  string `gorm:"column:DICT_NAME" json:"dict_name" `                       //type:string   comment:字典名称    version:2023-07-21 08:48
	Remark    string `gorm:"column:REMARK" json:"remark" `                             //type:string   comment:备注        version:2023-07-21 08:48
	CreatedBy string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-07-21 08:48
	CreatedAt int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-07-21 08:48
	UpdatedBy string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-07-21 08:48
	UpdatedAt int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-07-21 08:48
}

// TableName 表名:dict_type，字典类别表。
// 说明:
func (d *DictType) TableName() string {
	return "dict_type"
}

// 条件查询
func DictTypeGets(ctx *ctx.Context, query string, limit, offset int) ([]DictType, error) {
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

	var lst []DictType
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DictTypeGetByDictCode(ctx *ctx.Context, dictCode string) (bool, error) {
	var dictType *DictType
	err := DB(ctx).Where("DICT_CODE = ?", dictCode).Find(&dictType).Error
	logger.Debug(dictType)
	if err != nil || dictType.Id == 0 {
		return false, err
	}

	return true, nil
}

// 按id查询
func DictTypeGetById(ctx *ctx.Context, id int64) (*DictType, error) {
	var obj *DictType
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DictTypeGetsAll(ctx *ctx.Context) ([]DictType, error) {
	var lst []DictType
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加字典类别表
func (d *DictType) Add(ctx *ctx.Context) error {
	// 这里写DictType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除字典类别表
func (d *DictType) Del(ctx *ctx.Context) error {
	// 这里写DictType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新字典类别表
func (d *DictType) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DictType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DictTypeCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DictType{}).Where(where, args...))
}
