// Package models  资产扩展
// date : 2023-07-23 09:04
// desc : 资产扩展
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// AssetExpansion  资产扩展。
// 说明:
// 表名:asset_expansion
// group: AssetExpansion
// version:2023-07-23 09:04
type AssetExpansion struct {
	Id               int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键                                          version:2023-07-23 09:04
	AssetId          int64  `gorm:"column:ASSET_ID" json:"asset_id" `                         //type:*int     comment:资产ID                                        version:2023-07-28 16:17
	ConfigCategory   int64  `gorm:"column:CONFIG_CATEGORY" json:"config_category" `           //type:*int     comment:配置类别(1:基本信息,2:硬件配置,3:网络配置)    version:2023-07-23 09:04
	PropertyCategory string `gorm:"column:PROPERTY_CATEGORY" json:"property_category" `       //type:string   comment:属性类别                                      version:2023-07-23 09:04
	GroupId          string `gorm:"column:GROUP_ID" json:"group_id" `                         //type:string   comment:分组ID                                        version:2023-07-23 09:04
	PropertyName     string `gorm:"column:PROPERTY_NAME" json:"property_name" `               //type:string   comment:属性名称                                      version:2023-07-23 09:04
	PropertyValue    string `gorm:"column:PROPERTY_VALUE" json:"property_value" `             //type:string   comment:属性值                                        version:2023-07-23 09:04
	AssociatedTable  string `gorm:"column:ASSOCIATED_TABLE" json:"associated_table" `         //type:string   comment:关联表名                                      version:2023-07-23 09:04
	CreatedBy        string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人                                        version:2023-07-23 09:04
	CreatedAt        int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间                                      version:2023-07-23 09:04
	UpdatedBy        string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人                                        version:2023-07-23 09:04
	UpdatedAt        int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间                                      version:2023-07-23 09:04
}

// TableName 表名:asset_expansion，资产扩展。
// 说明:
func (a *AssetExpansion) TableName() string {
	return "asset_expansion"
}

// 条件查询
func AssetExpansionGets(ctx *ctx.Context, query string, limit, offset int) ([]AssetExpansion, error) {
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

	var lst []AssetExpansion
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func AssetExpansionGetById(ctx *ctx.Context, id int64) (*AssetExpansion, error) {
	var obj *AssetExpansion
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func AssetExpansionGetsAll(ctx *ctx.Context) ([]AssetExpansion, error) {
	var lst []AssetExpansion
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加资产扩展
func (a *AssetExpansion) Add(ctx *ctx.Context) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 批量增加资产扩展
func BatchAdd(ctx *ctx.Context, a []AssetExpansion) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(&a).Error
}

// 删除资产扩展
func (a *AssetExpansion) Del(ctx *ctx.Context) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(a).Error
}

// 更新资产扩展
func (a *AssetExpansion) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(a).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func AssetExpansionCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetExpansion{}).Where(where, args...))
}
