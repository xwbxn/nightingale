// Package models  机房信息
// date : 2023-07-16 09:01
// desc : 机房信息
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// ComputerRoom  机房信息。
// 说明:
// 表名:computer_room
// group: ComputerRoom
// version:2023-07-16 09:01
type ComputerRoom struct {
	Id                  int64   `gorm:"column:ID;primaryKey" json:"id" `                                                      //type:*int       comment:主键            version:2023-07-11 16:11
	RoomName            string  `gorm:"column:ROOM_NAME" json:"room_name" validate:"required"`                                //type:string     comment:名称            version:2023-07-11 16:11
	RoomCode            string  `gorm:"column:ROOM_CODE" json:"room_code" validate:"required"`                                //type:string     comment:编码            version:2023-07-11 16:11
	IdcLocation         string  `gorm:"column:IDC_LOCATION" json:"idc_location" validate:"omitempty"`                         //type:string     comment:所在IDC         version:2023-07-11 16:11
	Subgallery          string  `gorm:"column:SUBGALLERY" json:"subgallery" validate:"omitempty"`                             //type:string     comment:所属楼座        version:2023-07-11 16:11
	Floor               int64   `gorm:"column:FLOOR" json:"floor" validate:"omitempty"`                                       //type:*int       comment:所属楼层        version:2023-07-11 16:11
	Voltage             int64   `gorm:"column:VOLTAGE" json:"voltage" validate:"omitempty,min=0"`                             //type:*int       comment:电压            version:2023-07-11 16:11
	Electric            int64   `gorm:"column:ELECTRIC" json:"electric" validate:"omitempty,gte=0"`                           //type:*int       comment:电流            version:2023-07-11 16:11
	RowMax              int64   `gorm:"column:ROW_MAX" json:"row_max" validate:"required,gte=1,lte=9999"`                     //type:*int       comment:最大行数        version:2023-07-11 16:11
	ColumnMax           int64   `gorm:"column:COLUMN_MAX" json:"column_max" validate:"required,gte=1,lte=9999"`               //type:*int       comment:最大列数        version:2023-07-11 16:11
	CabinetNumber       int64   `gorm:"column:CABINET_NUMBER" json:"cabinet_number" validate:"required,gte=1,lte=999999"`     //type:*int       comment:可容纳机柜数    version:2023-07-11 16:11
	RoomBearingCapacity float64 `gorm:"column:ROOM_BEARING_CAPACITY" json:"room_bearing_capacity" validate:"omitempty,gte=0"` //type:*float64   comment:机房承重        version:2023-07-11 16:11
	RoomArea            float64 `gorm:"column:ROOM_AREA" json:"room_area" validate:"omitempty,gte=0"`                         //type:*float64   comment:机房面积        version:2023-07-11 16:11
	RatedPower          int64   `gorm:"column:RATED_POWER" json:"rated_power" validate:"omitempty,gte=0"`                     //type:*int       comment:额定功率        version:2023-07-11 16:11
	RoomPicture         string  `gorm:"column:ROOM_PICTURE" json:"room_picture" validate:"omitempty"`                         //type:string     comment:机房图片        version:2023-07-11 16:11
	DutyPersonOne       string  `gorm:"column:DUTY_PERSON_ONE" json:"duty_person_one" validate:"omitempty"`                   //type:string     comment:责任人1         version:2023-07-11 16:11
	DutyPersonTwo       string  `gorm:"column:DUTY_PERSON_two" json:"duty_person_two" validate:"omitempty"`                   //type:string     comment:责任人2         version:2023-07-11 16:11
	CreatedBy           string  `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                             //type:string     comment:创建人          version:2023-07-11 16:11
	CreatedAt           int64   `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                             //type:*int       comment:创建时间        version:2023-07-11 16:11
	UpdatedBy           string  `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                             //type:string     comment:更新人          version:2023-07-11 16:11
	UpdatedAt           int64   `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                             //type:*int       comment:更新时间        version:2023-07-11 16:11
}

// TableName 表名:computer_room，机房信息。
// 说明:
func (c *ComputerRoom) TableName() string {
	return "computer_room"
}

// 条件查询
func ComputerRoomGets(ctx *ctx.Context, query string, limit, offset int) ([]ComputerRoom, error) {
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

	var lst []ComputerRoom
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func ComputerRoomGetById(ctx *ctx.Context, id int64) (*ComputerRoom, error) {
	var obj *ComputerRoom
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func ComputerRoomGetsAll(ctx *ctx.Context) ([]ComputerRoom, error) {
	var lst []ComputerRoom
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加机房信息
func (c *ComputerRoom) Add(ctx *ctx.Context) error {
	// 这里写ComputerRoom的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(c).Error
}

// 删除机房信息
func (c *ComputerRoom) Del(ctx *ctx.Context) error {
	// 这里写ComputerRoom的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(c).Error
}

// 更新机房信息
func (c *ComputerRoom) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写ComputerRoom的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(c).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func ComputerRoomCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&ComputerRoom{}).Where(where, args...))
}
