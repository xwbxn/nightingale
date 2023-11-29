// Package models  许可配置
// date : 2023-10-29 15:12
// desc : 许可配置
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// LicenseConfig  许可配置。
// 说明:
// 表名:license_config
// group: LicenseConfig
// version:2023-10-29 15:12
type LicenseConfig struct {
	Id        int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int         comment:主键          version:2023-10-29 15:12
	Days      int64          `gorm:"column:DAYS" json:"days" `                                 //type:*int         comment:剩余天数      version:2023-10-29 15:12
	Nodes     int64          `gorm:"column:NODES" json:"nodes" `                               //type:*int         comment:剩余节点数    version:2023-10-29 15:12
	Frequency string         `gorm:"column:FREQUENCY" json:"frequency" `                       //type:string       comment:提醒频率      version:2023-10-29 15:12
	Email     string         `gorm:"column:EMAIL" json:"email" `                               //type:string       comment:邮箱          version:2023-10-29 15:12
	CreatedBy string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string       comment:创建人        version:2023-10-29 15:12
	CreatedAt int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int         comment:创建时间      version:2023-10-29 15:12
	UpdatedBy string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string       comment:更新人        version:2023-10-29 15:12
	UpdatedAt int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int         comment:更新时间      version:2023-10-29 15:12
	DeletedAt gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*time.Time   comment:删除时间      version:2023-10-29 15:12
}

// TableName 表名:license_config，许可配置。
// 说明:
func (l *LicenseConfig) TableName() string {
	return "license_config"
}

// 按id查询
func LicenseConfigGetById(ctx *ctx.Context, id int64) (*LicenseConfig, error) {
	var obj *LicenseConfig
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 增加许可配置
func (l *LicenseConfig) Add(ctx *ctx.Context) error {
	// 这里写LicenseConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(l).Error
}

// 更新许可配置
func (l *LicenseConfig) Update(ctx *ctx.Context, updateFrom LicenseConfig) error {
	// 这里写LicenseConfig的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(l).Updates(updateFrom).Error
}

func LicenseStatistics(ctx *ctx.Context) (*Statistics, error) {
	session := DB(ctx).Model(&Asset{}).Select("count(*) as total", "max(UPDATED_AT) as last_updated")

	var stats []*Statistics
	err := session.Find(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats[0], nil
}

// 查询所有
func LicenseConfigGetsAll(ctx *ctx.Context) ([]*LicenseConfig, error) {
	var lst []*LicenseConfig
	err := DB(ctx).Find(&lst).Error

	return lst, err
}
