// Package models  机柜组信息
// date : 2023-07-15 14:32
// desc : 机柜组信息
package models

import (
	"errors"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// CabinetGroup  机柜组信息。
// 说明:
// 表名:cabinet_group
// group: CabinetGroup
// version:2023-07-15 14:32
type CabinetGroup struct {
	Id               int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键          version:2023-07-15 14:32
	CabinetGroupCode string         `gorm:"column:CABINET_GROUP_CODE" json:"cabinet_group_code" `     //type:string   comment:机柜组编号    version:2023-07-15 14:32
	RoomId           int64          `gorm:"column:ROOM_ID" json:"room_id" `                           //type:string   comment:所属机房      version:2023-07-15 14:32
	CabinetGroupType string         `gorm:"column:CABINET_GROUP_TYPE" json:"cabinet_group_type" `     //type:string   comment:机柜组类型    version:2023-07-15 14:32
	Row              int64          `gorm:"column:ROW" json:"row" `                                   //type:*int     comment:行            version:2023-07-15 14:32
	StartColumn      int64          `gorm:"column:START_COLUMN" json:"start_column" `                 //type:*int     comment:开始列        version:2023-07-15 14:32
	Column           int64          `gorm:"column:COLUMN" json:"column" `                             //type:*int     comment:所在列        version:2023-07-15 14:32
	DutyPersonOne    string         `gorm:"column:DUTY_PERSON_ONE" json:"duty_person_one" `           //type:string   comment:责任人1       version:2023-07-15 14:32
	DutyPersonTwo    string         `gorm:"column:DUTY_PERSON_TWO" json:"duty_person_two" `           //type:string   comment:责任人2       version:2023-07-15 14:32
	UseNotes         string         `gorm:"column:USE_NOTES" json:"use_notes" `                       //type:string   comment:用途说明      version:2023-07-15 14:32
	CreatedBy        string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人        version:2023-07-15 14:32
	CreatedAt        int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间      version:2023-07-15 14:32
	UpdatedBy        string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人        version:2023-07-15 14:32
	UpdatedAt        int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间      version:2023-07-15 14:32
	DeletedAt        gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*int       comment:删除时间        version:2023-9-08 16:39
}

// TableName 表名:cabinet_group，机柜组信息。
// 说明:
func (c *CabinetGroup) TableName() string {
	return "cabinet_group"
}

// 条件查询
func CabinetGroupGets(ctx *ctx.Context, query string, limit, offset int) ([]CabinetGroup, error) {
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

	var lst []CabinetGroup
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func CabinetGroupGetById(ctx *ctx.Context, id int64) (*CabinetGroup, error) {
	var obj *CabinetGroup
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func CabinetGroupGetsAll(ctx *ctx.Context) ([]CabinetGroup, error) {
	var lst []CabinetGroup
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

//根据机房名称和机柜组编号查询
func (c *CabinetGroup) CabinetGroupGetsByCompNameAndCode(ctx *ctx.Context) (bool, error) {
	var lst []CabinetGroup
	err := DB(ctx).Where("ROOM_ID = ? AND CABINET_GROUP_CODE = ?", c.RoomId, c.CabinetGroupCode).Find(&lst).Error
	return len(lst) == 0, err
}

// 增加机柜组信息
func (c *CabinetGroup) Add(ctx *ctx.Context) error {
	// 这里写CabinetGroup的业务逻辑，通过error返回错误

	exist, err := c.CabinetGroupGetsByCompNameAndCode(ctx)
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("机房名称或编号已存在！")
	}

	// 实际向库中写入
	return DB(ctx).Create(c).Error
}

// 删除机柜组信息
func (c *CabinetGroup) Del(ctx *ctx.Context) error {
	// 这里写CabinetGroup的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(c).Error
}

// 更新机柜组信息
func (c *CabinetGroup) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写CabinetGroup的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(c).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func CabinetGroupCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&CabinetGroup{}).Where(where, args...))
}
