// Package models  字典数据表
// date : 2023-07-21 08:50
// desc : 字典数据表
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// DictData  字典数据表。
// 说明:
// 表名:dict_data
// group: DictData
// version:2023-07-21 08:50
type DictData struct {
	Id        int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-07-21 08:50
	DictCode  string `gorm:"column:DICT_CODE" json:"dict_code" `                       //type:string   comment:字典编码    version:2023-07-21 08:50
	DictKey   string `gorm:"column:DICT_KEY" json:"dict_key" `                         //type:string   comment:字典标签    version:2023-07-21 08:50
	DictValue string `gorm:"column:DICT_VALUE" json:"dict_value" `                     //type:string   comment:字典键值    version:2023-07-21 08:50
	Remark    string `gorm:"column:REMARK" json:"remark" `                             //type:string   comment:备注        version:2023-07-21 08:50
	CreatedBy string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-07-21 08:50
	CreatedAt int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-07-21 08:50
	UpdatedBy string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-07-21 08:50
	UpdatedAt int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-07-21 08:50
}

// TableName 表名:dict_data，字典数据表。
// 说明:
func (d *DictData) TableName() string {
	return "dict_data"
}

// 条件查询
func DictDataGets(ctx *ctx.Context, query string, limit, offset int) ([]DictData, error) {
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

	var lst []DictData
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DictDataGetById(ctx *ctx.Context, id int64) (*DictData, error) {
	var obj *DictData
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按dict_code查询
func DictDataGetByDictCode(ctx *ctx.Context, dictCode string) ([]*DictData, error) {
	var obj []*DictData
	err := DB(ctx).Where("DICT_CODE = ?", dictCode).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按dict_code删除
func DictDataDelByDictCode(ctx *ctx.Context, dictCode string) error {
	err := DB(ctx).Where("DICT_CODE = ?", dictCode).Delete(&DictData{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 查询所有
func DictDataGetsAll(ctx *ctx.Context) ([]DictData, error) {
	var lst []DictData
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加字典数据表
func Add(ctx *ctx.Context, d []*DictData) error {
	// 这里写DictData的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除字典数据表
func (d *DictData) Del(ctx *ctx.Context) error {
	// 这里写DictData的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新字典数据表
func (d *DictData) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DictData的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DictDataCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DictData{}).Where(where, args...))
}
