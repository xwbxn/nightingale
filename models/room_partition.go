// Package models  机房分区表
// date : 2023-9-10 10:37
// desc : 机房分区表
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// RoomPartition  机房分区表。
// 说明:
// 表名:room_partition
// group: RoomPartition
// version:2023-9-10 10:37
type RoomPartition struct {
	Id            int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-9-10 10:37
	RoomId        int64          `gorm:"column:ROOM_ID" json:"room_id" `                           //type:*int     comment:机房ID      version:2023-9-10 11:14
	Name          string         `gorm:"column:NAME" json:"name" `                                 //type:string   comment:分区名称    version:2023-9-10 10:37
	StartRow      int64          `gorm:"column:START_ROW" json:"start_row" `                       //type:*int     comment:起始行      version:2023-9-10 10:37
	StartColumn   int64          `gorm:"column:START_COLUMN" json:"start_column" `                 //type:*int     comment:起始列      version:2023-9-10 10:37
	Height        int64          `gorm:"column:HEIGHT" json:"height" `                             //type:*int     comment:高度        version:2023-9-10 10:37
	Width         int64          `gorm:"column:WIDTH" json:"width" `                               //type:*int     comment:宽度        version:2023-9-10 10:37
	SpaceType     string         `gorm:"column:SPACE_TYPE" json:"space_type" `                     //type:string   comment:空间类型    version:2023-9-10 11:09
	Description   string         `gorm:"column:DESCRIPTION" json:"description" `                   //type:string   comment:位置描述    version:2023-9-10 10:37
	DutyPersonOne string         `gorm:"column:DUTY_PERSON_ONE" json:"duty_person_one" `           //type:string   comment:责任人1     version:2023-9-10 10:37
	DutyPersonTwo string         `gorm:"column:DUTY_PERSON_two" json:"duty_person_two" `           //type:string   comment:责任人2     version:2023-9-10 10:37
	Remark        string         `gorm:"column:REMARK" json:"remark" `                             //type:string   comment:备注        version:2023-9-10 10:37
	CreatedBy     string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-9-10 10:37
	CreatedAt     int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-9-10 10:37
	UpdatedBy     string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-9-10 10:37
	UpdatedAt     int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-9-10 10:37
	DeletedAt     gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:string   comment:删除时间    version:2023-9-10 10:37
}

// TableName 表名:room_partition，机房分区表。
// 说明:
func (r *RoomPartition) TableName() string {
	return "room_partition"
}

// 条件查询
func RoomPartitionGets(ctx *ctx.Context, query string, limit, offset int) ([]RoomPartition, error) {
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

	var lst []RoomPartition
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func RoomPartitionGetById(ctx *ctx.Context, id int64) (*RoomPartition, error) {
	var obj *RoomPartition
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按map查询
func RoomPartitionGetBymap(ctx *ctx.Context, where map[string]interface{}) ([]RoomPartition, error) {
	var obj []RoomPartition
	err := DB(ctx).Model(&RoomPartition{}).Where(where).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func RoomPartitionGetsAll(ctx *ctx.Context) ([]RoomPartition, error) {
	var lst []RoomPartition
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加机房分区表
func (r *RoomPartition) Add(ctx *ctx.Context) error {
	// 这里写RoomPartition的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(r).Error
}

// 删除机房分区表
func (r *RoomPartition) Del(ctx *ctx.Context) error {
	// 这里写RoomPartition的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(r).Error
}

// 更新机房分区表
func (r *RoomPartition) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写RoomPartition的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(r).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func RoomPartitionCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&RoomPartition{}).Where(where, args...))
}
