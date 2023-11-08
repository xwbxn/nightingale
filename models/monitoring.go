// Package models  监控
// date : 2023-10-08 16:45
// desc : 监控
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// Monitoring  监控。
// 说明:
// 表名:monitoring
// group: Monitoring
// version:2023-10-08 16:45
type Monitoring struct {
	Id             int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:BIGINT       comment:主键        version:2023-10-08 16:45
	AssetId        int64          `gorm:"column:ASSET_ID" json:"asset_id" `                         //type:BIGINT       comment:资产id      version:2023-10-08 16:45
	MonitoringName string         `gorm:"column:MONITORING_NAME" json:"monitoring_name" `           //type:string       comment:监控名称    version:2023-10-08 16:45
	MonitoringSql  string         `gorm:"column:MONITORING_SQL" json:"monitoring_sql" `             //type:string       comment:监控脚本    version:2023-10-08 16:45
	Status         int64          `gorm:"column:STATUS" json:"status" `                             //type:*int         comment:状态        version:2023-10-08 16:45
	TargetId       int64          `gorm:"column:TARGET_ID" json:"target_id" `                       //type:*int         comment:采集器      version:2023-10-08 16:45
	Remark         string         `gorm:"column:REMARK" json:"remark" `                             //type:string       comment:说明        version:2023-10-08 16:45
	CreatedBy      string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string       comment:创建人      version:2023-10-08 16:45
	CreatedAt      int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int         comment:创建时间    version:2023-10-08 16:45
	UpdatedBy      string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string       comment:更新人      version:2023-10-08 16:45
	UpdatedAt      int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int         comment:更新时间    version:2023-10-08 16:45
	DeletedAt      gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*time.Time   comment:删除时间    version:2023-10-08 16:45
}

// TableName 表名:monitoring，监控。
// 说明:
func (m *Monitoring) TableName() string {
	return "monitoring"
}

// 条件查询
func MonitoringGets(ctx *ctx.Context, query string, limit, offset int) ([]Monitoring, error) {
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

	var lst []Monitoring
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func MonitoringGetById(ctx *ctx.Context, id int64) (*Monitoring, error) {
	var obj *Monitoring
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func MonitoringGetsAll(ctx *ctx.Context) ([]Monitoring, error) {
	var lst []Monitoring
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加监控
func (m *Monitoring) Add(ctx *ctx.Context) error {
	// 这里写Monitoring的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(m).Error
}

// 删除监控
func (m *Monitoring) Del(ctx *ctx.Context) error {
	// 这里写Monitoring的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(m).Error
}

// 更新监控
func (m *Monitoring) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写Monitoring的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(m).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func MonitoringCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&Monitoring{}).Where(where, args...))
}
