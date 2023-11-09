// Package models
// date : 2023-10-08 15:31
// desc :
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// Bigscreen  。
// 说明:
// 表名:bigscreen
// group: Bigscreen
// version:2023-10-08 15:31
type Bigscreen struct {
	Id        int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int         comment:主键        version:2023-10-08 15:31
	CreatedBy string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string       comment:创建人      version:2023-10-08 15:31
	CreatedAt int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int         comment:创建时间    version:2023-10-08 15:31
	UpdatedBy string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string       comment:更新人      version:2023-10-08 15:31
	UpdatedAt int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int         comment:更新时间    version:2023-10-08 15:31
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at" swaggerignore:"true"`
	Title     string         `gorm:"column:TITLE" json:"title" `   //type:string       comment:标题        version:2023-10-08 15:31
	Desc      string         `gorm:"column:DESC" json:"desc" `     //type:string       comment:简介        version:2023-10-08 15:31
	Config    string         `gorm:"column:CONFIG" json:"config" ` //type:string       comment:配置        version:2023-10-08 15:31
}

// TableName 表名:bigscreen，。
// 说明:
func (b *Bigscreen) TableName() string {
	return "bigscreen"
}

// 条件查询
func BigscreenGets(ctx *ctx.Context, query string, limit, offset int) ([]Bigscreen, error) {
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

	var lst []Bigscreen
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func BigscreenGetById(ctx *ctx.Context, id int64) (*Bigscreen, error) {
	var obj *Bigscreen
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func BigscreenGetsAll(ctx *ctx.Context) ([]Bigscreen, error) {
	var lst []Bigscreen
	err := DB(ctx).Find(&lst).Select("id", "title", "desc", "create_at", "create_by").Error

	return lst, err
}

// 增加
func (b *Bigscreen) Add(ctx *ctx.Context) error {
	// 这里写Bigscreen的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(b).Error
}

// 删除
func (b *Bigscreen) Del(ctx *ctx.Context) error {
	// 这里写Bigscreen的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(b).Error
}

// 更新
func (b *Bigscreen) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写Bigscreen的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(b).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func BigscreenCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&Bigscreen{}).Where(where, args...))
}
