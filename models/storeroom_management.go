// Package models  库房信息
// date : 2023-08-21 14:09
// desc : 库房信息
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// StoreroomManagement  库房信息。
// 说明:
// 表名:storeroom_management
// group: StoreroomManagement
// version:2023-08-21 14:09
type StoreroomManagement struct {
	Id               int64          `gorm:"column:ID;primaryKey" json:"id" `                                                                                       //type:*int     comment:主键        version:2023-08-21 14:09
	RoomNumber       string         `gorm:"column:ROOM_NUMBER" json:"room_number" cn:"房间号"`                                                                        //type:string   comment:房间号      version:2023-08-21 14:09
	BelongIdc        int64          `gorm:"column:BELONG_IDC" json:"belong_idc" cn:"所属IDC" source:"type=table,table=datacenter,property=id,field=datacenter_name"` //type:*int     comment:所属IDC     version:2023-08-21 14:09
	RoomAddress      string         `gorm:"column:ROOM_ADDRESS" json:"room_address" cn:"房间地址"`                                                                     //type:string   comment:房间地址    version:2023-08-21 14:09
	DutyBy           string         `gorm:"column:DUTY_BY" json:"duty_by" cn:"责任人"`                                                                                //type:string   comment:责任人      version:2023-08-21 14:09
	ShelfInformation string         `gorm:"column:SHELF_INFORMATION" json:"shelf_information" cn:"货架信息"`                                                           //type:string   comment:货架信息    version:2023-08-21 14:09
	ContactNumber    string         `gorm:"column:CONTACT_NUMBER" json:"contact_number" cn:"联系电话"`                                                                 //type:string   comment:联系电话    version:2023-08-21 14:09
	CreatedBy        string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                              //type:string   comment:创建人      version:2023-08-21 14:09
	CreatedAt        int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                              //type:*int     comment:创建时间    version:2023-08-21 14:09
	UpdatedBy        string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                              //type:string   comment:更新人      version:2023-08-21 14:09
	UpdatedAt        int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                              //type:*int     comment:更新时间    version:2023-08-21 14:09
	DeletedAt        gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`                                                              //type:*int       comment:删除时间        version:2023-9-08 16:39
}

// TableName 表名:storeroom_management，库房信息。
// 说明:
func (s *StoreroomManagement) TableName() string {
	return "storeroom_management"
}

// 条件查询
func StoreroomManagementGets(ctx *ctx.Context, query string, limit, offset int) ([]StoreroomManagement, error) {
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

	var lst []StoreroomManagement
	err := session.Find(&lst).Error

	return lst, err
}

// 根据房价号或房间地址模糊查询
func StoreroomNumAddressGets(ctx *ctx.Context, query string, limit, offset int) ([]StoreroomManagement, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	where := "%" + query + "%"

	var lst []StoreroomManagement
	err := session.Where("ROOM_NUMBER LIKE ? OR ROOM_ADDRESS LIKE ?", where, where).Find(&lst).Error

	return lst, err
}

// 按id查询
func StoreroomManagementGetById(ctx *ctx.Context, id int64) (*StoreroomManagement, error) {
	var obj *StoreroomManagement
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func StoreroomManagementGetsAll(ctx *ctx.Context) ([]StoreroomManagement, error) {
	var lst []StoreroomManagement
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加库房信息
func (s *StoreroomManagement) Add(ctx *ctx.Context) error {
	// 这里写StoreroomManagement的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(s).Error
}

// 删除库房信息
func (s *StoreroomManagement) Del(ctx *ctx.Context) error {
	// 这里写StoreroomManagement的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(s).Error
}

// 批量删除库房信息
func StoreroomManagementBatchDel(ctx *ctx.Context, s []int64) error {
	// 这里写StoreroomManagement的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Where("ID IN ?", s).Delete(&StoreroomManagement{}).Error
}

// 更新库房信息
func (s *StoreroomManagement) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写StoreroomManagement的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(s).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func StoreroomManagementCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&StoreroomManagement{}).Where(where, args...))
}

// 根据传入条件统计个数
func StoreroomManagementByMap(ctx *ctx.Context, where map[string]interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&StoreroomManagement{}).Where(where))
}

// 根据房价号或房间地址统计个数
func StoreroomNumAddressCount(ctx *ctx.Context, where string) (num int64, err error) {
	where = "%" + where + "%"
	return Count(DB(ctx).Model(&StoreroomManagement{}).Where("ROOM_NUMBER LIKE ? OR ROOM_ADDRESS LIKE ?", where, where))
}
