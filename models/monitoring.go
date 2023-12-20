// Package models  监控
// date : 2023-10-08 16:45
// desc : 监控
package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	context "github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/prom"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/toolkits/pkg/logger"
	"gorm.io/gorm"
)

// Monitoring  监控。
// 说明:
// 表名:monitoring
// group: Monitoring
// version:2023-10-08 16:45
type Monitoring struct {
	Id             int64               `gorm:"column:ID;primaryKey" json:"id" `                          //type:BIGINT       comment:主键        version:2023-10-08 16:45
	AssetId        int64               `gorm:"column:ASSET_ID" json:"asset_id" `                         //type:BIGINT       comment:资产id      version:2023-10-08 16:45
	AssetIds       []int64             `gorm:"-" json:"asset_ids" `                                      //批量添加，对前端
	MonitoringName string              `gorm:"column:MONITORING_NAME" json:"monitoring_name" `           //type:string       comment:监控名称    version:2023-10-08 16:45
	DatasourceId   int64               `gorm:"column:DATASOURCE_ID" json:"datasource_id" `               //type:BIGINT       comment:数据源名称    version:2023-10-08 16:45
	MonitoringSql  string              `gorm:"column:MONITORING_SQL" json:"monitoring_sql" `             //type:string       comment:监控脚本    version:2023-10-08 16:45
	Status         int64               `gorm:"column:STATUS" json:"status" `                             //type:*int         comment:状态        version:2023-10-08 16:45
	IsAlarm        int64               `gorm:"column:IS_ALARM" json:"is_alarm" `                         //type:*int         comment:是否启用告警    version:2023-10-13 14:27
	TargetId       int64               `gorm:"column:TARGET_ID" json:"target_id" `                       //type:*int         comment:采集器      version:2023-10-08 16:45
	Config         string              `gorm:"column:CONFIG" json:"config" `                             //type:string       comment:配置信息    version:2023-10-13 14:20
	Remark         string              `gorm:"column:REMARK" json:"remark" `                             //type:string       comment:说明        version:2023-10-08 16:45
	Unit           string              `gorm:"column:UNIT" json:"unit"`                                  //type:string       comment:单位        version:2023-10-08 16:45
	Label          string              `gorm:"column:LABEL" json:"label"`                                //type:string       comment:标签        version:2023-10-08 16:45
	CreatedBy      string              `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string       comment:创建人      version:2023-10-08 16:45
	CreatedAt      int64               `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int         comment:创建时间    version:2023-10-08 16:45
	UpdatedBy      string              `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string       comment:更新人      version:2023-10-08 16:45
	UpdatedAt      int64               `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int         comment:更新时间    version:2023-10-08 16:45
	DeletedAt      gorm.DeletedAt      `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*time.Time   comment:删除时间    version:2023-10-08 16:45
	AlertRules     []AlertRuleSimplify `gorm:"-" json:"alert_rules" `
}

// MonitoringSimplify  监控简易信息。
// 说明:
// 表名:monitoring
// group: Monitoring
// version:2023-10-08 16:45
type MonitoringSimplify struct {
	Id int64 `json:"id"`
}

// TableName 表名:monitoring，监控。
// 说明:
func (m *Monitoring) TableName() string {
	return "monitoring"
}

// 查询所有
func MonitoringAllGets(ctx *context.Context, query string, limit, offset int) ([]Monitoring, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id DESC")
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

// 根据条件统计个数
func MonitoringMapCount(ctx *context.Context, where map[string]interface{}, query string,
	assetType string, datasource, assetId int64) (num int64, err error) {

	session := DB(ctx).Joins("LEFT JOIN assets ON monitoring.ASSET_ID = assets.id").
		Joins("LEFT JOIN datasource ON monitoring.datasource_id = datasource.id")

	if assetId != -1 {
		session = session.Where("monitoring.ASSET_ID = ? ", assetId)
	}
	if assetType != "" {
		session = session.Where("assets.type = ? ", assetType)
	}
	if datasource != -1 {
		session = session.Where("datasource.id = ? ", datasource)
	}

	var str strings.Builder
	vals := make([]interface{}, 0)
	if query != "" {
		query = "%" + query + "%"
		str.WriteString("( monitoring.MONITORING_NAME like ? or ")
		vals = append(vals, query)
		str.WriteString("assets.name like ? or ")
		vals = append(vals, query)
		str.WriteString("assets.type like ? or ")
		vals = append(vals, query)
		str.WriteString("assets.ip like ? )")
		vals = append(vals, query)

	}

	err = session.Model(&Monitoring{}).Where(str.String(), vals...).Count(&num).Error

	return num, err
}

// 条件查询
func MonitoringMapGets(ctx *context.Context, where map[string]interface{}, query string, limit, offset int,
	assetType string, datasource, assetId int64) (lst []Monitoring, err error) {
	session := DB(ctx).Joins("LEFT JOIN assets ON monitoring.ASSET_ID = assets.id").
		Joins("LEFT JOIN datasource ON monitoring.datasource_id = datasource.id")
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	if assetId != -1 {
		session = session.Where("monitoring.ASSET_ID = ? ", assetId)
	}
	if assetType != "" {
		session = session.Where("assets.type = ? ", assetType)
	}
	if datasource != -1 {
		session = session.Where("datasource.id = ? ", datasource)
	}

	var str strings.Builder
	vals := make([]interface{}, 0)
	// 这里使用列名的硬编码构造查询参数, 避免从前台传入造成注入风险
	if query != "" {
		query = "%" + query + "%"
		str.WriteString("( monitoring.MONITORING_NAME like ? or ")
		vals = append(vals, query)
		str.WriteString("assets.name like ? or ")
		vals = append(vals, query)
		str.WriteString("assets.type like ? or ")
		vals = append(vals, query)
		str.WriteString("assets.ip like ? )")
		vals = append(vals, query)
	}

	err = session.Model(&Monitoring{}).
		Select("monitoring.*").Where(str.String(), vals...).Find(&lst).Error

	return lst, err
}

// 根据条件统计个数(new)
func MonitoringMapCountNew(ctx *context.Context, filter, query, assetType string, assetId int64, assetIds []int64) (num int64, err error) {

	session := DB(ctx)
	if filter == "monitoring_name" {
		session = session.Where("MONITORING_NAME like ? ", "%"+query+"%")
		// } else if filter == "asset_name" {
		// 	ids, err := AssetIdByName(ctx, "%"+query+"%")
		// 	if err != nil {
		// 		return 0, err
		// 	}
		// 	session = session.Where("ASSET_ID in? ", ids)
	} else if filter == "status" {
		session = session.Where("STATUS like ? ", "%"+query+"%")
	} else if filter == "is_alarm" {
		session = session.Where("IS_ALARM like? ", "%"+query+"%")
	} else if filter == "asset_name" || filter == "asset_type" || filter == "asset_ip" {
		session = session.Where("ASSET_ID in ? ", assetIds)
	}
	if assetId != -1 {
		session = session.Where("ASSET_ID = ? ", assetId)
	}
	if assetType != "" {
		session = session.Where("ASSET_ID in ? ", assetIds)
		// 		return 0, err
	}
	// 	session = session.Where("ASSET_ID in? ", ids)
	// }
	// if len(assetIds) > 0 {
	// 	session = session.Where("ASSET_ID in ? ", assetIds)
	// }

	err = session.Model(&Monitoring{}).Count(&num).Error

	return num, err
}

// 条件查询(new)
func MonitoringMapGetsNew(ctx *context.Context, filter, query, assetType string, assetId int64, assetIds []int64, limit, offset int) (lst []Monitoring, err error) {
	session := DB(ctx)

	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("UPDATED_AT DESC")
	}

	if filter == "monitoring_name" {
		session = session.Where("MONITORING_NAME like ? ", "%"+query+"%")
		// } else if filter == "asset_name" {
		// 	ids, err := AssetIdByName(ctx, "%"+query+"%")
		// 	if err != nil {
		// 		return lst, err
		// 	}
		// 	session = session.Where("ASSET_ID in? ", ids)
	} else if filter == "status" {
		session = session.Where("STATUS like ? ", "%"+query+"%")
	} else if filter == "is_alarm" {
		session = session.Where("IS_ALARM like? ", "%"+query+"%")
	} else if filter == "asset_name" || filter == "asset_type" || filter == "asset_ip" {
		logger.Debug("42314444444444445433333333333")
		session = session.Where("ASSET_ID in ? ", assetIds)
	}

	if assetId != -1 {
		session = session.Where("ASSET_ID = ? ", assetId)
	}
	if assetType != "" {
		session = session.Where("ASSET_ID in ? ", assetIds)
	}
	// 	session = session.Where("ASSET_ID in? ", ids)
	// }
	// if len(assetIds) > 0 {
	// 	session = session.Where("ASSET_ID in ? ", assetIds)
	// }

	err = session.Model(&Monitoring{}).Find(&lst).Error

	return lst, err
}

// 按id查询
func MonitoringGetById(ctx *context.Context, id int64) (*Monitoring, error) {
	var obj *Monitoring
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按map查询
func MonitoringGetMap(ctx *context.Context, where map[string]interface{}) ([]Monitoring, error) {
	var lst []Monitoring
	err := DB(ctx).Where(where).Find(&lst).Error

	return lst, err
}

// 按ids查询
func MonitoringGetByBatchId(ctx *context.Context, ids []int64) ([]Monitoring, error) {
	var lst []Monitoring
	err := DB(ctx).Model(&Monitoring{}).Where("id in ?", ids).Find(&lst).Error
	return lst, err
}

// 查询所有
func MonitoringGetsAll(ctx *context.Context) ([]*Monitoring, error) {
	var lst []*Monitoring
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

func (m *Monitoring) Verify() error {
	return nil
}

// 增加监控
func (m *Monitoring) Add(ctx *context.Context) error {
	// 这里写Monitoring的业务逻辑，通过error返回错误
	if err := m.Verify(); err != nil {
		return err
	}

	return DB(ctx).Create(m).Error
}

// 增加监控
func (m *Monitoring) AddTx(tx *gorm.DB) error {
	// 这里写Monitoring的业务逻辑，通过error返回错误
	if err := m.Verify(); err != nil {
		tx.Rollback()
	}

	err := tx.Create(m).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 增加监控及告警规则
func (m *Monitoring) AddAndAlertRules(ctx *context.Context, me *User) error {
	err := ctx.Transaction(func(ctx *context.Context) error {
		if len(m.AssetIds) == 0 {

			if err := m.BuildMonitoringAndAlertRules(ctx, me, m.AssetId); err != nil {
				return err
			}

		} else {
			for _, id := range m.AssetIds {
				if err := m.BuildMonitoringAndAlertRules(ctx, me, id); err != nil {
					return err
				}
			}
		}

		return nil
	})
	return err
}

func (m *Monitoring) BuildMonitoringAndAlertRules(ctx *context.Context, me *User, assetId int64) error {
	now := time.Now().Unix()
	var monitor = &Monitoring{
		AssetId:        assetId,
		MonitoringName: m.MonitoringName,
		DatasourceId:   m.DatasourceId,
		MonitoringSql:  m.MonitoringSql,
		Status:         1, //默认启用
		TargetId:       m.TargetId,
		Remark:         m.Remark,
		Unit:           m.Unit,
		AlertRules:     m.AlertRules,
		CreatedBy:      me.Username,
		CreatedAt:      now,
		UpdatedBy:      me.Username,
		UpdatedAt:      now,
		Label:          m.Label,
	}
	// 这里写Monitoring的业务逻辑，通过error返回错误
	if err := monitor.Verify(); err != nil {
		return err
	}

	if err := DB(ctx).Create(monitor).Error; err != nil {
		return err
	}

	for _, alertRuleSimplify := range m.AlertRules {
		alertRule, err := BuildAlertRule(ctx, *monitor, alertRuleSimplify)
		if err != nil {
			return err
		}
		err = alertRule.FE2DB(ctx)
		if err != nil {
			return err
		}
		if err := alertRule.Add(ctx); err != nil {
			return err
		}
	}
	return nil
}

// 删除监控
func (m *Monitoring) Del(ctx *context.Context) error {
	return DB(ctx).Where("id=?", m.Id).Delete(&Monitoring{}).Error
}

// 删除监控
func BatchDelTx(tx *gorm.DB, ids []string) error {
	return tx.Where("id in ?", ids).Delete(&Monitoring{}).Error
}

// 删除监控并删除关联的告警规则
func BatchDel(ctx *context.Context, ids []string) error {
	if len(ids) == 0 {
		return errors.New("ids empty")
	}

	err := ctx.Transaction(func(ctx *context.Context) error {
		if err := DB(ctx).Where("id in ?", ids).Delete(&Monitoring{}).Error; err != nil {
			return err
		}
		if err := DB(ctx).Where("monitoring_id in ?", ids).Delete(&AlertRule{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// 更新监控
func (m *Monitoring) Update(ctx *context.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写Monitoring的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(m).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 更新监控
func (m *Monitoring) UpdateTx(tx *gorm.DB, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写Monitoring的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Model(m).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 根据条件统计个数
func MonitoringCount(ctx *context.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&Monitoring{}).Where(where, args...))
}

// 批量更改监控状态
func MonitoringUpdateStatus(ctx *context.Context, ids []int64, status, oType int64) error {
	if oType == 1 {
		return DB(ctx).Model(&Monitoring{}).Where("id in ?", ids).Updates(map[string]interface{}{"status": status, "updated_at": time.Now().Unix()}).Error
	} else if oType == 2 {
		return DB(ctx).Model(&Monitoring{}).Where("id in ?", ids).Updates(map[string]interface{}{"is_alarm": status, "updated_at": time.Now().Unix()}).Error
	}
	return nil
}

func (m *Monitoring) CompilePromQL() string {
	promql, err := prom.InjectLabel(m.MonitoringSql, "asset_id", strconv.Itoa(int(m.AssetId)), labels.MatchEqual)
	if err != nil {
		promql = ""
		logger.Error("compile promql error:", m.AssetId, m.MonitoringName, m.MonitoringSql)
	}
	return promql
}

func MonitoringStatistics(ctx *context.Context) (*Statistics, error) {
	session := DB(ctx).Model(&Monitoring{}).Select("count(*) as total", "max(updated_at) as last_updated")

	var stats []*Statistics
	err := session.Find(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats[0], nil
}
