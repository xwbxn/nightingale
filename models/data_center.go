// Package models  数据中心
// date : 2023-07-11 15:49
// desc : 数据中心
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// DataCenter  数据中心。
// 说明:
// 表名:data_center
// group: DataCenter
// version:2023-07-11 15:49
type DataCenter struct {
	Id                     int64   `gorm:"column:ID;primaryKey" json:"id" `                                                    //type:*int       comment:主键            version:2023-07-11 15:49
	DataCenterName         string  `gorm:"column:DATA_CENTER_NAME" json:"data_center_name" validate:"required"`                //type:string     comment:数据中心名称    version:2023-07-11 15:49
	DataCenterCode         string  `gorm:"column:DATA_CENTER_CODE" json:"data_center_code" validate:"omitempty"`               //type:string     comment:数据中心编码    version:2023-07-11 15:49
	City                   string  `gorm:"column:CITY" json:"city" validate:"required"`                                        //type:string     comment:所在城市        version:2023-07-11 15:49
	Address                string  `gorm:"column:ADDRESS" json:"address" validate:"omitempty"`                                 //type:string     comment:地址            version:2023-07-11 15:49
	DutyPersonOne          string  `gorm:"column:DUTY_PERSON_ONE" json:"duty_person_one" validate:"omitempty"`                 //type:string     comment:责任人1         version:2023-07-11 15:49
	DutyPersonTwo          string  `gorm:"column:DUTY_PERSON_TWO" json:"duty_person_two" validate:"omitempty"`                 //type:string     comment:责任人2         version:2023-07-11 15:49
	LoadBearing            float64 `gorm:"column:LOAD_BEARING" json:"load_bearing" validate:"omitempty,gte=0"`                 //type:*float64   comment:承重            version:2023-07-11 15:49
	Area                   float64 `gorm:"column:AREA" json:"area" validate:"omitempty,gte=0"`                                 //type:*float64   comment:面积            version:2023-07-11 15:49
	AffiliatedOrganization string  `gorm:"column:AFFILIATED_ORGANIZATION" json:"affiliated_organization" validate:"omitempty"` //type:string     comment:所属组织机构    version:2023-07-11 15:49
	Remark                 string  `gorm:"column:REMARK" json:"remark" validate:"omitempty"`                                   //type:string     comment:备注            version:2023-07-11 15:49
	CreatedBy              string  `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                           //type:string     comment:创建人          version:2023-07-11 15:49
	CreatedAt              int64   `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                           //type:*int       comment:创建时间        version:2023-07-11 15:49
	UpdatedBy              string  `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                           //type:string     comment:更新人          version:2023-07-11 15:49
	UpdatedAt              int64   `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                           //type:*int       comment:更新时间        version:2023-07-11 15:49
}

// TableName 表名:data_center，数据中心。
// 说明:
func (d *DataCenter) TableName() string {
	return "data_center"
}

// 条件查询
func DataCenterGets(ctx *ctx.Context, query string, limit, offset int) ([]DataCenter, error) {
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

	var lst []DataCenter
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DataCenterGetById(ctx *ctx.Context, id int64) (*DataCenter, error) {
	var obj *DataCenter
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DataCenterGetsAll(ctx *ctx.Context) ([]DataCenter, error) {
	var lst []DataCenter
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加数据中心
func (d *DataCenter) Add(ctx *ctx.Context) error {
	// 这里写DataCenter的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除数据中心
func (d *DataCenter) Del(ctx *ctx.Context) error {
	// 这里写DataCenter的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新数据中心
func (d *DataCenter) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DataCenter的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DataCenterCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DataCenter{}).Where(where, args...))
}
