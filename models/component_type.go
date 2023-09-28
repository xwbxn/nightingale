// Package models  部件类型
// date : 2023-08-21 09:08
// desc : 部件类型
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// ComponentType  部件类型。
// 说明:
// 表名:component_type
// group: ComponentType
// version:2023-08-21 09:08
type ComponentType struct {
	Id               int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-08-21 09:08
	ComponentType    string `gorm:"column:COMPONENT_TYPE" json:"component_type" `             //type:string   comment:部件类型    version:2023-08-21 09:18
	Remark           string `gorm:"column:REMARK" json:"remark" `                             //type:string   comment:备注        version:2023-08-21 09:08
	ComponentPicture string `gorm:"column:COMPONENT_PICTURE" json:"component_picture" `       //type:string   comment:部件图      version:2023-08-21 09:08
	CreatedBy        string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-08-21 09:08
	CreatedAt        int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-08-21 09:08
	UpdatedBy        string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-08-21 09:08
	UpdatedAt        int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-08-21 09:08
}

// TableName 表名:component_type，部件类型。
// 说明:
func (c *ComponentType) TableName() string {
	return "component_type"
}

// 条件查询
func ComponentTypeGets(ctx *ctx.Context, query string, limit, offset int) ([]ComponentType, error) {
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

	var lst []ComponentType
	err := session.Find(&lst).Error

	return lst, err
}

// 根据map查询
func ComponentTypeGetMap(ctx *ctx.Context, m map[string]interface{}, limit, offset int) ([]ComponentType, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	var lst []ComponentType
	err := session.Debug().Where(m).Find(&lst).Error

	return lst, err
}

// 按id查询
func ComponentTypeGetById(ctx *ctx.Context, id int64) (*ComponentType, error) {
	var obj *ComponentType
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func ComponentTypeGetsAll(ctx *ctx.Context) ([]ComponentType, error) {
	var lst []ComponentType
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加部件类型
func (c *ComponentType) Add(ctx *ctx.Context) error {
	// 这里写ComponentType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(c).Error
}

// 删除部件类型
func (c *ComponentType) Del(ctx *ctx.Context) error {
	// 这里写ComponentType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(c).Error
}

// 批量删除部件类型
func ComponentTypeBatchDel(ctx *ctx.Context, c []ComponentType) error {
	// 这里写ComponentType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(c).Error
}

// 更新部件类型
func (c *ComponentType) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写ComponentType的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(c).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func ComponentTypeCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&ComponentType{}).Where(where, args...))
}

// 根据map统计个数
func ComponentTypeCountMap(ctx *ctx.Context, m map[string]interface{}) (num int64, err error) {
	return Count(DB(ctx).Debug().Model(&ComponentType{}).Where(m))
}
