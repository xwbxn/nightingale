// Package models  数据中心
// date : 2023-07-16 08:47
// desc : 数据中心
package models

import (
	"errors"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// Datacenter  数据中心。
// 说明:
// 表名:datacenter
// group: Datacenter
// version:2023-07-16 08:47
type Datacenter struct {
	Id                 int64   `gorm:"column:ID;primaryKey" json:"id" `                                            //type:*int       comment:主键            version:2023-07-11 15:49
	DatacenterName     string  `gorm:"column:DATACENTER_NAME" json:"datacenter_name" validate:"required"`          //type:string     comment:数据中心名称    version:2023-07-11 15:49
	DatacenterCode     string  `gorm:"column:DATACENTER_CODE" json:"datacenter_code" validate:"omitempty"`         //type:string     comment:数据中心编码    version:2023-07-11 15:49
	City               string  `gorm:"column:CITY" json:"city" validate:"required"`                                //type:string     comment:所在城市        version:2023-07-11 15:49
	Address            string  `gorm:"column:ADDRESS" json:"address" validate:"omitempty"`                         //type:string     comment:地址            version:2023-07-11 15:49
	DutyPersonOne      string  `gorm:"column:DUTY_PERSON_ONE" json:"duty_person_one" validate:"omitempty"`         //type:string     comment:责任人1         version:2023-07-11 15:49
	DutyPersonTwo      string  `gorm:"column:DUTY_PERSON_TWO" json:"duty_person_two" validate:"omitempty"`         //type:string     comment:责任人2         version:2023-07-11 15:49
	LoadBearing        float64 `gorm:"column:LOAD_BEARING" json:"load_bearing" validate:"omitempty,gte=0"`         //type:*float64   comment:承重            version:2023-07-11 15:49
	Area               float64 `gorm:"column:AREA" json:"area" validate:"omitempty,gte=0"`                         //type:*float64   comment:面积            version:2023-07-11 15:49
	BelongOrganization int64   `gorm:"column:BELONG_ORGANIZATION" json:"belong_organization" validate:"omitempty"` //type:string     comment:所属组织机构    version:2023-07-11 15:49
	Remark             string  `gorm:"column:REMARK" json:"remark" validate:"omitempty"`                           //type:string     comment:备注            version:2023-07-11 15:49
	CreatedBy          string  `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                   //type:string     comment:创建人          version:2023-07-11 15:49
	CreatedAt          int64   `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                   //type:*int       comment:创建时间        version:2023-07-11 15:49
	UpdatedBy          string  `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                   //type:string     comment:更新人          version:2023-07-11 15:49
	UpdatedAt          int64   `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                   //type:*int       comment:更新时间        version:2023-07-11 15:49
}

// TableName 表名:datacenter，数据中心。
// 说明:
func (d *Datacenter) TableName() string {
	return "datacenter"
}

// 条件查询
func DatacenterGets(ctx *ctx.Context, query string, limit, offset int) ([]Datacenter, error) {
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

	var lst []Datacenter
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DatacenterGetById(ctx *ctx.Context, id int64) (*Datacenter, error) {
	var obj *Datacenter
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DatacenterGetsAll(ctx *ctx.Context) ([]Datacenter, error) {
	var lst []Datacenter
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

//根据数据中心名称查询
func (d *Datacenter) DatacenterGetsByName(ctx *ctx.Context) (bool, error) {
	var lst []Datacenter
	err := DB(ctx).Where("DATACENTER_NAME", d.DatacenterName).Find(&lst).Error
	return len(lst) == 0, err
}

// 增加数据中心
func (d *Datacenter) Add(ctx *ctx.Context) error {
	// 这里写Datacenter的业务逻辑，通过error返回错误
	exist, err := d.DatacenterGetsByName(ctx)
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("数据中心名称已存在！")
	}

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除数据中心
func (d *Datacenter) Del(ctx *ctx.Context) error {
	// 这里写Datacenter的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新数据中心
func (d *Datacenter) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写Datacenter的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DatacenterCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&Datacenter{}).Where(where, args...))
}
