// Package models  资产管理
// date : 2023-08-04 09:45
// desc : 资产管理
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	optlock "gorm.io/plugin/optimisticlock"
)

// AssetManagement  资产管理。
// 说明:
// 表名:asset_management
// group: AssetManagement
// version:2023-08-04 09:45
type AssetManagement struct {
	Id             int64           `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-08-04 09:45
	AssetCode      string          `gorm:"column:ASSET_CODE" json:"asset_code" `                     //type:string   comment:资产编号    version:2023-08-04 09:45
	ShutdownLevel  int64           `gorm:"column:SHUTDOWN_LEVEL" json:"shutdown_level" `             //type:*int     comment:关机级别    version:2023-08-04 09:45
	ServiceLevel   int64           `gorm:"column:SERVICE_LEVEL" json:"service_level" `               //type:*int     comment:服务级别    version:2023-08-04 09:45
	BelongDept     int64           `gorm:"column:BELONG_DEPT" json:"belong_dept" `                   //type:*int     comment:所属部门    version:2023-08-04 09:45
	EquipmentUse   string          `gorm:"column:EQUIPMENT_USE" json:"equipment_use" `               //type:string   comment:设备用途    version:2023-08-04 09:45
	UserDepartment int64           `gorm:"column:USER_DEPARTMENT" json:"user_department" `           //type:*int     comment:使用部门    version:2023-08-04 09:45
	UsingSite      string          `gorm:"column:USING_SITE" json:"using_site" `                     //type:string   comment:使用地点    version:2023-08-04 09:45
	Version        optlock.Version `gorm:"column:VERSION" json:"version" swaggerignore:"true"`       //type:*int     comment:版本号      version:2023-08-04 09:45
	CreatedBy      string          `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-08-04 09:45
	CreatedAt      int64           `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-08-04 09:45
	UpdatedBy      string          `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-08-04 09:45
	UpdatedAt      int64           `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-08-04 09:45
}

// TableName 表名:asset_management，资产管理。
// 说明:
func (a *AssetManagement) TableName() string {
	return "asset_management"
}

// 条件查询
func AssetManagementGets(ctx *ctx.Context, query string, limit, offset int) ([]AssetManagement, error) {
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

	var lst []AssetManagement
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func AssetManagementGetById(ctx *ctx.Context, id int64) (*AssetManagement, error) {
	var obj *AssetManagement
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func AssetManagementGetsAll(ctx *ctx.Context) ([]AssetManagement, error) {
	var lst []AssetManagement
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加资产管理
func (a *AssetManagement) Add(ctx *ctx.Context) error {
	// 这里写AssetManagement的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 删除资产管理
func (a *AssetManagement) Del(ctx *ctx.Context) error {
	// 这里写AssetManagement的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(a).Error
}

// 更新资产管理
func (a *AssetManagement) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写AssetManagement的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(a).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func AssetManagementCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetManagement{}).Where(where, args...))
}
