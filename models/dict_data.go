// Package models  字典数据表
// date : 2023-07-21 08:50
// desc : 字典数据表
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// DictData  字典数据表。
// 说明:
// 表名:dict_data
// group: DictData
// version:2023-07-21 08:50
type DictData struct {
	Id        int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-07-21 08:50
	TypeCode  string `gorm:"column:TYPE_CODE" json:"type_code" `                       //type:string   comment:字典编码    version:2023-08-01 14:10
	DictKey   string `gorm:"column:DICT_KEY" json:"dict_key" `                         //type:string   comment:字典标签    version:2023-07-21 08:50
	DictValue string `gorm:"column:DICT_VALUE" json:"dict_value" `                     //type:string   comment:字典键值    version:2023-07-21 08:50
	Sn        int64  `gorm:"column:SN" json:"sn" `                                     //type:*int     comment:序号        version:2023-08-25 10:06
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

// 条件查询扩展
func DictDataGetExp(ctx *ctx.Context, where string, typeCode string, limit, offset int) ([]DictData, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	where = "%" + where + "%"

	var lst []DictData
	err := session.Debug().Where("TYPE_CODE = ?", typeCode).Where("DICT_VALUE LIKE ? OR REMARK LIKE ?", where, where).Find(&lst).Error

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

// 按id批量查询
func DictDataGetBatchById(ctx *ctx.Context, ids []int64) ([]DictData, error) {
	var obj []DictData
	err := DB(ctx).Debug().Where("ID IN ?", ids).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按map查询
func DictDataGetByMap(ctx *ctx.Context, m map[string]interface{}) ([]DictData, error) {
	var obj []DictData
	err := DB(ctx).Where(m).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按type_code查询
func DictDataGetByTypeCode(ctx *ctx.Context, typeCode string) ([]DictData, error) {
	var obj []DictData
	err := DB(ctx).Where("TYPE_CODE = ?", typeCode).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按type_code和dict_value查询
func DictDataGetByTypeCodeValue(ctx *ctx.Context, typeCode, value string) (DictData, error) {
	var obj DictData
	err := DB(ctx).Where("TYPE_CODE = ? AND DICT_VALUE = ?", typeCode, value).Find(&obj).Error
	if err != nil {
		return obj, err
	}

	return obj, nil
}

// 按map删除
func DictDataDelByMap(ctx *ctx.Context, m map[string]interface{}) error {
	err := DB(ctx).Where(m).Delete(&DictData{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 按map删除(事务)
func DictDataDelTxByMap(tx *gorm.DB, m map[string]interface{}) error {
	err := tx.Debug().Where(m).Delete(&DictData{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 按ids删除(事务)
func DictDataDelByIds(tx *gorm.DB, ids []int64) error {
	err := tx.Debug().Where("ID IN ?", ids).Delete(&DictData{}).Error
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

// 批量增加字典数据表
func DictDataBatchAdd(ctx *ctx.Context, d []*DictData) error {
	// 这里写DictData的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 增加字典数据表
func (d *DictData) Add(ctx *ctx.Context) error {
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

// 根据条件统计扩展个数
func DictDataExpCount(ctx *ctx.Context, where string, typeCode string) (num int64, err error) {
	where = "%" + where + "%"
	return Count(DB(ctx).Model(&DictData{}).Where("TYPE_CODE = ?", typeCode).Where("DICT_VALUE LIKE ? OR REMARK LIKE ?", where, where))
}
