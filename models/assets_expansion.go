// Package models  资产扩展-西航
// date : 2023-9-20 11:23
// desc : 资产扩展-西航
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// AssetsExpansion  资产扩展-西航。
// 说明:
// 表名:assets_expansion
// group: AssetsExpansion
// version:2023-9-20 11:23
type AssetsExpansion struct {
	Id             int64  `gorm:"column:id;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-9-20 11:23
	AssetsId       int64  `gorm:"column:assets_id" json:"assets_id" `                       //type:*int     comment:资产id      version:2023-9-20 14:33
	ConfigCategory string `gorm:"column:config_category" json:"config_category" `           //type:string       comment:配置类别(1:基本信息,2:硬件配置,3:网络配置)    version:2023-10-08 11:44
	GroupId        string `gorm:"column:group_id" json:"group_id" `                         //type:string       comment:分组ID                                        version:2023-10-08 11:44
	NameCn         string `gorm:"column:name_cn" json:"name_cn" `                           //type:string   comment:属性名称    version:2023-9-20 11:23
	Name           string `gorm:"column:name" json:"name" `                                 //type:string   comment:英文名称    version:2023-9-20 11:23
	Value          string `gorm:"column:value" json:"value" `                               //type:string   comment:属性值      version:2023-9-20 11:23
	CreatedBy      string `gorm:"column:created_by" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-9-20 11:23
	CreatedAt      int64  `gorm:"column:created_at" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-9-20 11:23
	UpdatedBy      string `gorm:"column:updated_by" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-9-20 11:23
	UpdatedAt      int64  `gorm:"column:updated_at" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-9-20 11:23
}

// TableName 表名:assets_expansion，资产扩展-西航。
// 说明:
func (a *AssetsExpansion) TableName() string {
	return "assets_expansion"
}

// 条件查询
func AssetsExpansionGets(ctx *ctx.Context, query string, limit, offset int) ([]AssetsExpansion, error) {
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

	var lst []AssetsExpansion
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func AssetsExpansionGetById(ctx *ctx.Context, id int64) (*AssetsExpansion, error) {
	var obj *AssetsExpansion
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func AssetsExpansionGetsAll(ctx *ctx.Context) ([]AssetsExpansion, error) {
	var lst []AssetsExpansion
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 根据map查询
func AssetsExpansionGetsMap(ctx *ctx.Context, where map[string]interface{}) ([]AssetsExpansion, error) {
	var lst []AssetsExpansion
	err := DB(ctx).Where(where).Find(&lst).Error

	return lst, err
}

// 根据map查询AssetIds
func AssetsExpansionGetAssetIdMap(ctx *ctx.Context, where map[string]interface{}) ([]int64, error) {
	var lst []int64
	err := DB(ctx).Where(where).Distinct().Pluck("assetI_id", &lst).Error

	return lst, err
}

// 增加资产扩展-西航
func (a *AssetsExpansion) Add(ctx *ctx.Context) error {
	// 这里写AssetsExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 增加资产扩展-西航
func AssetsExpansionAdd(ctx *ctx.Context, a []AssetsExpansion) error {
	// 这里写AssetsExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 增加资产扩展-西航
func AssetsExpansionAddTx(tx *gorm.DB, a []AssetsExpansion) error {
	// 这里写AssetsExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Create(a).Error
	if err != nil {
		tx.Rollback()
	}
	return nil
}

// 删除资产扩展-西航
func (a *AssetsExpansion) Del(ctx *ctx.Context) error {
	// 这里写AssetsExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(a).Error
}

// 通过AssetsIds删除资产扩展-西航
func AssetsExpansionDelAssetsIds(tx *gorm.DB, where []string) error {
	// 这里写AssetsExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Where("ASSETS_ID IN ?", where).Delete(&AssetsExpansion{}).Error
	if err != nil {
		tx.Rollback()
	}
	return nil
}

// 通过map删除资产扩展-西航
func AssetsExpansionDelMap(tx *gorm.DB, where map[string]interface{}) error {
	// 这里写AssetsExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Where(where).Delete(&AssetsExpansion{}).Error
	if err != nil {
		tx.Rollback()
	}
	return nil
}

// 更新资产扩展-西航
func (a *AssetsExpansion) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写AssetsExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(a).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 更新资产扩展-西航
func AssetsExpansionUpdateTx(tx *gorm.DB, where map[string]interface{}, updateFrom map[string]interface{}) error {
	// 这里写AssetsExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Model(&AssetsExpansion{}).Where(where).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
	if err != nil {
		tx.Rollback()
	}
	return nil
}

// 根据条件统计个数
func AssetsExpansionCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetsExpansion{}).Where(where, args...))
}
