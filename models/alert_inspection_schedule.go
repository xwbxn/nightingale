package models

import (
	"encoding/json"
	"strings"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

type AlertInspectionSchedule struct {
	Id          int64  `json:"id" gorm:"primaryKey"`
	PlanName    string `json:"paln_name"`   // 巡检计划名称
	Supervisor  string `json:"supervisor"`  // 巡检负责人
	Description string `json:"description"` // 巡检计划描述
	Area        string `json:"area"`        // 巡检区域
	Scope       string `json:"scope"`       // 巡检范围
	Report      string `json:"report"`      // 巡检报备
	Time        int64  `json:"time"`        // 巡检时间
	Receiver    string `json:"receiver"`    // 报告接收人
	State       string `json:"state"`       // 状态
	HandleAt    int64  `json:"handle_at"`   // 创建时间
	HandleBy    string `json:"handle_by"`   // 创建人
	UpdateAt    int64  `json:"update_at"`   // 修改人
	UpdateBy    string `json:"update_by"`   // 修改时间
	Reset       string `json:"reset"`       // 停用
}

func (e *AlertInspectionSchedule) TableName() string {
	return "alert_inspection_schedule"
}

/**
添加数据
**/
func (e *AlertInspectionSchedule) Add(ctx *ctx.Context) error {
	return Insert(ctx, e)
}

/**
更新所有字段
**/
func (e *AlertInspectionSchedule) UpdateAllFields(ctx *ctx.Context) error {
	return DB(ctx).Model(e).Select("*").Updates(e).Error
}

/**
更新部分字段
**/
func (e *AlertInspectionSchedule) Update(ctx *ctx.Context, selectField interface{}, selectFields ...interface{}) error {
	return DB(ctx).Model(e).Select(selectField, selectFields...).Updates(e).Error
}

/**
删除字段
**/
func (e *AlertInspectionSchedule) Del(ctx *ctx.Context) error {
	return DB(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Where("id=?", e.Id).Delete(&AlertInspectionSchedule{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (e *AlertInspectionSchedule) DB2FE() {
	e.NotifyChannelsJSON = strings.Fields(e.NotifyChannels)
	e.NotifyGroupsJSON = strings.Fields(e.NotifyGroups)
	e.CallbacksJSON = strings.Fields(e.Callbacks)
	e.TagsJSON = strings.Split(e.Tags, ",,")
	json.Unmarshal([]byte(e.Annotations), &e.AnnotationsJSON)
	json.Unmarshal([]byte(e.RuleConfig), &e.RuleConfigJson)
}

func AlertInspectionScheduleGets(ctx *ctx.Context, prods []string, bgid, stime, etime int64, severity int, recovered int, dsIds []int64, cates []string, query string, limit, offset int, status int) ([]AlertHisEvent, error) {
	session := DB(ctx).Where("HandleAt between ? and ?", stime, etime)

	if len(prods) != 0 {
		session = session.Where("rule_prod in ?", prods)
	}

	if bgid > 0 {
		session = session.Where("group_id = ?", bgid)
	}

	if severity >= 0 {
		session = session.Where("severity = ?", severity)
	}

	if recovered >= 0 {
		session = session.Where("is_recovered = ?", recovered)
	}
	if status >= 0 {
		session = session.Where("status = ?", status)
	}

	if len(dsIds) > 0 {
		session = session.Where("datasource_id in ?", dsIds)
	}

	if len(cates) > 0 {
		session = session.Where("cate in ?", cates)
	}

	if query != "" {
		arr := strings.Fields(query)
		for i := 0; i < len(arr); i++ {
			qarg := "%" + arr[i] + "%"
			session = session.Where("rule_name like ? or tags like ?", qarg, qarg)
		}
	}

	var lst []AlertInspectionSchedule
	err := session.Order("id desc").Limit(limit).Offset(offset).Find(&lst).Error

	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].DB2FE()
		}
	}

	return lst, err
}
