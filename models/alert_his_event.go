package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/toolkits/pkg/logger"
	"gorm.io/gorm"
)

type AlertHisEvent struct {
	Id                 int64             `json:"id" gorm:"primaryKey"`
	AssetId            int64             `json:"asset_id"`
	AssetName          string            `json:"asset_name" gorm:"-"`
	AssetIp            string            `json:"asset_ip" gorm:"-"`
	Cate               string            `json:"cate"`
	IsRecovered        int               `json:"is_recovered"`
	DatasourceId       int64             `json:"datasource_id"`
	Cluster            string            `json:"cluster"`
	GroupId            int64             `json:"group_id"`
	GroupName          string            `json:"group_name"` // busi group name
	Hash               string            `json:"hash"`
	RuleId             int64             `json:"rule_id"`
	RuleName           string            `json:"rule_name"`
	RuleNote           string            `json:"rule_note"`
	RuleProd           string            `json:"rule_prod"`
	RuleAlgo           string            `json:"rule_algo"`
	Severity           int               `json:"severity"`
	PromForDuration    int               `json:"prom_for_duration"`
	PromQl             string            `json:"prom_ql"`
	RuleConfig         string            `json:"-" gorm:"rule_config"` // rule config
	RuleConfigJson     interface{}       `json:"rule_config" gorm:"-"` // rule config for fe
	PromEvalInterval   int               `json:"prom_eval_interval"`
	Callbacks          string            `json:"-"`
	CallbacksJSON      []string          `json:"callbacks" gorm:"-"`
	RunbookUrl         string            `json:"runbook_url"`
	NotifyRecovered    int               `json:"notify_recovered"`
	NotifyChannels     string            `json:"-"`
	NotifyChannelsJSON []string          `json:"notify_channels" gorm:"-"`
	NotifyGroups       string            `json:"-"`
	NotifyGroupsJSON   []string          `json:"notify_groups" gorm:"-"`
	NotifyGroupsObj    []UserGroup       `json:"notify_groups_obj" gorm:"-"`
	TargetIdent        string            `json:"target_ident"`
	TargetNote         string            `json:"target_note"`
	TriggerTime        int64             `json:"trigger_time"`
	TriggerValue       string            `json:"trigger_value"`
	RecoverTime        int64             `json:"recover_time"`
	LastEvalTime       int64             `json:"last_eval_time"`
	Tags               string            `json:"-"`
	TagsJSON           []string          `json:"tags" gorm:"-"`
	Annotations        string            `json:"-"`
	AnnotationsJSON    map[string]string `json:"annotations" gorm:"-"` // for fe
	NotifyCurNumber    int               `json:"notify_cur_number"`    // notify: current number
	FirstTriggerTime   int64             `json:"first_trigger_time"`   // 连续告警的首次告警时间
	Status             int               `json:"status"`               // 状态
	Remark             string            `json:"remark"`               // 备注
	HandleAt           int64             `json:"handle_at"`            // 创建时间
	HandleBy           string            `json:"handle_by"`            // 创建人
	ExtraConfig        interface{}       `json:"extra_config" gorm:"-"`
	DeletedAt          gorm.DeletedAt    `gorm:"column:deleted_at" json:"deleted_at" swaggerignore:"true"`
}

func (e *AlertHisEvent) TableName() string {
	return "alert_his_event"
}

func (e *AlertHisEvent) Add(ctx *ctx.Context) error {
	return Insert(ctx, e)
}

func (e *AlertHisEvent) DB2FE() {
	e.NotifyChannelsJSON = strings.Fields(e.NotifyChannels)
	e.NotifyGroupsJSON = strings.Fields(e.NotifyGroups)
	e.CallbacksJSON = strings.Fields(e.Callbacks)
	e.TagsJSON = strings.Split(e.Tags, ",,")
	json.Unmarshal([]byte(e.Annotations), &e.AnnotationsJSON)
	json.Unmarshal([]byte(e.RuleConfig), &e.RuleConfigJson)
}

func (e *AlertHisEvent) FillNotifyGroups(ctx *ctx.Context, cache map[int64]*UserGroup) error {
	// some user-group already deleted ?
	count := len(e.NotifyGroupsJSON)
	if count == 0 {
		e.NotifyGroupsObj = []UserGroup{}
		return nil
	}

	for i := range e.NotifyGroupsJSON {
		id, err := strconv.ParseInt(e.NotifyGroupsJSON[i], 10, 64)
		if err != nil {
			continue
		}

		ug, has := cache[id]
		if has {
			e.NotifyGroupsObj = append(e.NotifyGroupsObj, *ug)
			continue
		}

		ug, err = UserGroupGetById(ctx, id)
		if err != nil {
			return err
		}

		if ug != nil {
			e.NotifyGroupsObj = append(e.NotifyGroupsObj, *ug)
			cache[id] = ug
		}
	}

	return nil
}

func AlertHisEventTotal(ctx *ctx.Context, prods []string, bgid, stime, etime int64, severity int, recovered int, dsIds []int64, cates []string, query string, status int) (int64, error) {
	session := DB(ctx).Model(&AlertHisEvent{}).Where("last_eval_time between ? and ?", stime, etime)

	if len(prods) > 0 {
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

	if len(dsIds) > 0 {
		session = session.Where("datasource_id in ?", dsIds)
	}

	if len(cates) > 0 {
		session = session.Where("cate in ?", cates)
	}
	if status >= 0 {
		session = session.Where("status = ?", status)
	}

	if query != "" {
		arr := strings.Fields(query)
		for i := 0; i < len(arr); i++ {
			qarg := "%" + arr[i] + "%"
			session = session.Where("rule_name like ? or tags like ?", qarg, qarg)
		}
	}

	return Count(session)
}

func AlertHisEventGets(ctx *ctx.Context, prods []string, bgid, stime, etime int64, severity int, recovered int, dsIds []int64, cates []string, query string, limit, offset int, status int) ([]AlertHisEvent, error) {
	session := DB(ctx).Where("last_eval_time between ? and ?", stime, etime)

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

	var lst []AlertHisEvent
	err := session.Order("id desc").Limit(limit).Offset(offset).Find(&lst).Error

	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].DB2FE()
		}
	}

	return lst, err
}

//西航
func AlertHisEventXHTotal(ctx *ctx.Context, fType, start, end, group int64, ids []int64, query string) (int64, error) {
	session := DB(ctx)
	if start != -1 {
		session = session.Where("trigger_time >= ?", start)
	}
	if end != -1 {
		session = session.Where("trigger_time <= ?", end)
	}
	if group == -1 {
		session = session.Where("group_id = ?", group)
	}
	if query != "" {
		query = "%" + query + "%"
		if fType == -1 {
			session = session.Where("asset_id in ? or id like ? or rule_name like ? or severity like ?", ids, query, query, query)
		} else if fType == 1 {
			session = session.Where("severity like ?", query)
		} else if fType == 2 {
			session = session.Where("asset_id in ?", ids)
		}
	}
	var num int64
	err := session.Model(&AlertHisEvent{}).Count(&num).Error
	return num, err
}

func AlertHisEventXHGets(ctx *ctx.Context, fType, start, end, group int64, ids []int64, query string, limit, offset int) ([]AlertHisEvent, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id DESC")
	}
	if start != -1 {
		session = session.Where("trigger_time >= ?", start)
	}
	if end != -1 {
		session = session.Where("trigger_time <= ?", end)
	}
	if group == -1 {
		session = session.Where("group_id = ?", group)
	}
	if query != "" {
		query = "%" + query + "%"
		if fType == -1 {
			session = session.Where("asset_id in ? or id like ? or rule_name like ? or severity like ?", ids, query, query, query)
		} else if fType == 1 {
			session = session.Where("severity like ?", query)
		} else if fType == 2 {
			session = session.Where("asset_id in ?", ids)
		}
	}
	var lst []AlertHisEvent
	err := session.Model(&AlertHisEvent{}).Find(&lst).Error

	for index := range lst {
		lst[index].DB2FE()
	}

	return lst, err

}

func AlertHisEventDelByIds(ctx *ctx.Context, ids []int64) error {
	return DB(ctx).Debug().Where("id in ?", ids).Delete(&AlertHisEvent{}).Error
}

func AlertHisEventDelByIdsTx(tx *gorm.DB, ids []int64) error {
	err := tx.Debug().Where("id in ?", ids).Delete(&AlertHisEvent{}).Error
	if err != nil {
		tx.Rollback()
	}
	return nil
}

func AlertHisEventGet(ctx *ctx.Context, where string, args ...interface{}) (*AlertHisEvent, error) {
	var lst []*AlertHisEvent
	err := DB(ctx).Where(where, args...).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}

	lst[0].DB2FE()
	lst[0].FillNotifyGroups(ctx, make(map[int64]*UserGroup))

	return lst[0], nil
}

func AlertHisEventGetById(ctx *ctx.Context, id int64) (*AlertHisEvent, error) {
	return AlertHisEventGet(ctx, "id=?", id)
}

func AlertHisEventGetByIds(ctx *ctx.Context, ids []int64) ([]AlertHisEvent, error) {
	var lst []AlertHisEvent
	err := DB(ctx).Model(&AlertHisEvent{}).Where("id in ?", ids).Find(&lst).Error
	return lst, err
}

func (m *AlertHisEvent) UpdateFieldsMap(ctx *ctx.Context, fields map[string]interface{}) error {
	return DB(ctx).Model(m).Updates(fields).Error
}

/*
*
通过IsRecovered判断指标是否恢复
*
*/
func (m *AlertHisEvent) UpdateStatus(ctx *ctx.Context, id int64, status int, remark string, handleBy string) error {

	modes, err := AlertHisEventGetById(ctx, id)
	if modes == nil {
		logger.Errorf("Sorry,Event isn't recovered, %s", err)
	}
	modes.Status = status
	modes.Remark = remark
	modes.HandleAt = time.Now().Unix()
	modes.HandleBy = handleBy

	return DB(ctx).Model(m).Updates(modes).Error
}

func AlertHisEventUpgradeToV6(ctx *ctx.Context, dsm map[string]Datasource) error {
	var lst []*AlertHisEvent
	err := DB(ctx).Where("trigger_time > ?", time.Now().Unix()-3600*24*30).Limit(10000).Order("id desc").Find(&lst).Error
	if err != nil {
		return err
	}

	for i := 0; i < len(lst); i++ {
		ds, exists := dsm[lst[i].Cluster]
		if !exists {
			continue
		}
		lst[i].DatasourceId = ds.Id

		ruleConfig := PromRuleConfig{
			Queries: []PromQuery{
				{
					PromQl:   lst[i].PromQl,
					Severity: lst[i].Severity,
				},
			},
		}
		b, _ := json.Marshal(ruleConfig)
		lst[i].RuleConfig = string(b)

		if lst[i].RuleProd == "" {
			lst[i].RuleProd = METRIC
		}

		if lst[i].Cate == "" {
			lst[i].Cate = PROMETHEUS
		}

		err = lst[i].UpdateFieldsMap(ctx, map[string]interface{}{
			"datasource_id": lst[i].DatasourceId,
			"rule_config":   lst[i].RuleConfig,
			"rule_prod":     lst[i].RuleProd,
			"cate":          lst[i].Cate,
		})
		if err != nil {
			logger.Errorf("update alert rule:%d datasource ids failed, %v", lst[i].Id, err)
		}
	}
	return nil
}

func EventPersist(ctx *ctx.Context, event *AlertCurEvent) error {
	has, err := AlertCurEventExists(ctx, "hash=?", event.Hash)
	if err != nil {
		return fmt.Errorf("event_persist_check_exists_fail: %v rule_id=%d hash=%s", err, event.RuleId, event.Hash)
	}

	alertRule, err := AlertRuleGetById(ctx, event.RuleId)
	if err != nil {
		return err
	}
	event.AssetId = alertRule.AssetId

	his := event.ToHis(ctx)

	// 不管是告警还是恢复，全量告警里都要记录
	if err := his.Add(ctx); err != nil {
		return fmt.Errorf("add his event error:%v", err)
	}

	if has {
		// 活跃告警表中有记录，删之
		err = AlertCurEventDelByHash(ctx, event.Hash)
		if err != nil {
			return fmt.Errorf("event_del_cur_fail: %v hash=%s", err, event.Hash)
		}

		if !event.IsRecovered {
			// 恢复事件，从活跃告警列表彻底删掉，告警事件，要重新加进来新的event
			// use his id as cur id
			event.Id = his.Id
			if event.Id > 0 {
				if err := event.Add(ctx); err != nil {
					return fmt.Errorf("add cur event err:%v", err)
				}
			}
		}

		return nil
	}

	if event.IsRecovered {
		// alert_cur_event表里没有数据，表示之前没告警，结果现在报了恢复，神奇....理论上不应该出现的
		return nil
	}

	// use his id as cur id
	event.Id = his.Id
	if event.Id > 0 {
		if err := event.Add(ctx); err != nil {
			return fmt.Errorf("add cur event error:%v", err)
		}
	}

	return nil
}

//统计未处理告警
func AlertHisCount(ctx *ctx.Context) (num int64, err error) {
	err = DB(ctx).Debug().Model(&AlertHisEvent{}).Where("status != 0").Count(&num).Error
	return num, err
}

//map统计历史告警
func AlertHisCountMap(ctx *ctx.Context, where map[string]interface{}) (num int64, err error) {
	err = DB(ctx).Debug().Model(&AlertHisEvent{}).Where(where).Count(&num).Error
	return num, err
}

//统计未处理告警
func TodayAlertCount(ctx *ctx.Context) (num int64, err error) {
	// 获取当前日期
	now := time.Now()
	// 根据当前日期设置零点时间
	zeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 获取零点时间戳
	zeroTimestamp := zeroTime.Unix()
	err = DB(ctx).Debug().Model(&AlertHisEvent{}).Where("trigger_time >= ?", zeroTimestamp).Count(&num).Error
	return num, err
}

//过滤器统计历史告警个数
func AlertHisCountFilter(ctx *ctx.Context, where map[string]interface{}, dateRange int64, query string, ids []int64) (num int64, err error) {
	session := DB(ctx)
	if dateRange != -1 {
		session = session.Where("trigger_time >= ?", time.Now().Unix()-dateRange)
	}
	if query != "" {
		querystr := "%" + query + "%"

		idsStr := make([]interface{}, 0)
		sql := "id like ? or rule_name like ? or rule_note like ? or severity like ? or tags like ? "
		idsStr = append(idsStr, querystr)
		idsStr = append(idsStr, querystr)
		idsStr = append(idsStr, querystr)
		idsStr = append(idsStr, querystr)
		idsStr = append(idsStr, "%asset_id="+query+"%")
		if len(ids) > 0 {
			sql += "or tags REGEXP "
			for index, val := range ids {
				if index == 0 {
					sql += "?"
				} else {
					sql += "|?"
				}
				str := "asset_id=" + strconv.FormatInt(val, 10)
				idsStr = append(idsStr, str)
			}
			logger.Debug(ids)
			logger.Debug(idsStr)
			// session = session.Or(sql, idsStr)
		}

		session = session.Where(sql, idsStr...)

	}
	err = session.Debug().Model(&AlertHisEvent{}).Where(where).Count(&num).Error
	return num, err
}

//过滤器统计历史告警
func AlertHisFilter(ctx *ctx.Context, where map[string]interface{}, dateRange int64, query string, ids []int64, limit, offset int) ([]*FeAlert, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("ID")
	}

	if dateRange != -1 {
		session = session.Where("trigger_time >= ?", time.Now().Unix()-dateRange)
	}
	if query != "" {
		querystr := "%" + query + "%"

		idsStr := make([]interface{}, 0)
		sql := "id like ? or rule_name like ? or rule_note like ? or severity like ? or tags like ? "
		idsStr = append(idsStr, querystr)
		idsStr = append(idsStr, querystr)
		idsStr = append(idsStr, querystr)
		idsStr = append(idsStr, querystr)
		idsStr = append(idsStr, "%asset_id="+query+"%")
		if len(ids) > 0 {
			sql += "or tags REGEXP "
			for index, val := range ids {
				if index == 0 {
					sql += "?"
				} else {
					sql += "|?"
				}
				str := "asset_id=" + strconv.FormatInt(val, 10)
				idsStr = append(idsStr, str)
			}
			logger.Debug(ids)
			logger.Debug(idsStr)
		}
		session = session.Where(sql, idsStr...)

	}
	var dat []*AlertHisEvent
	err := session.Debug().Model(&AlertHisEvent{}).Where(where).Find(&dat).Error
	Fe := MakeHisFeAlert(dat)
	return Fe, err
}

type RuleIdAndName struct {
	RuleId   int64  `json:"rule_id"`
	RuleName string `json:"rule_name"`
}

//统计所有的ruleName
func RuleNameGet(ctx *ctx.Context) (result []RuleIdAndName, err error) {
	err = DB(ctx).Model(&AlertHisEvent{}).Distinct("rule_id", "rule_name").Find(&result).Error
	return result, err
}

type GroupIdAndName struct {
	GroupId   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
}

//统计所有的groupName
func GroupNameGet(ctx *ctx.Context) (result []GroupIdAndName, err error) {
	err = DB(ctx).Model(&AlertHisEvent{}).Distinct("group_id", "group_name").Find(&result).Error
	return result, err
}

// 生成新的适用于前端页面的返回数据格式
func MakeHisFeAlert(dat []*AlertHisEvent) (Fe []*FeAlert) {
	Fe = make([]*FeAlert, 0) //避免返回null
	var assetID int
	for i := 0; i < len(dat); i++ {
		dat[i].DB2FE()
		for u := 0; u < len(dat[i].TagsJSON); u++ {
			s := strings.Split(dat[i].TagsJSON[u], "=")
			if s[0] == "asset_id" {
				assetID, _ = strconv.Atoi(s[1])
				break
			}
		}
		dic := &ruleConfigJson{}
		json.Unmarshal([]byte(dat[i].RuleConfig), dic)
		Fe = append(Fe, &FeAlert{
			Id:           dat[i].Id,
			Name:         dat[i].RuleName,
			Severity:     dat[i].Severity,
			TriggerTime:  dat[i].TriggerTime,
			TriggerValue: dat[i].TriggerValue,
			Rule:         dic.Queries[0]["prom_ql"].(string),
			AssetId:      assetID,
		})
	}
	return Fe
}
