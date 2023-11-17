package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/poster"
	"github.com/ccfos/nightingale/v6/pkg/tplx"
	"github.com/toolkits/pkg/logger"
	"gorm.io/gorm"
)

type AlertCurEvent struct {
	Id                       int64             `json:"id" gorm:"primaryKey"`
	AssetId                  int64             `json:"asset_id"`
	AssetName                string            `json:"asset_name" gorm:"-"`
	AssetIp                  string            `json:"asset_ip" gorm:"-"`
	Cate                     string            `json:"cate"`
	Cluster                  string            `json:"cluster"`
	DatasourceId             int64             `json:"datasource_id"`
	GroupId                  int64             `json:"group_id"`   // busi group id
	GroupName                string            `json:"group_name"` // busi group name
	Hash                     string            `json:"hash"`       // rule_id + vector_key
	RuleId                   int64             `json:"rule_id"`
	RuleName                 string            `json:"rule_name"`
	RuleNote                 string            `json:"rule_note"`
	RuleProd                 string            `json:"rule_prod"`
	RuleAlgo                 string            `json:"rule_algo"`
	Severity                 int               `json:"severity"`
	PromForDuration          int               `json:"prom_for_duration"`
	PromQl                   string            `json:"prom_ql"`
	RuleConfig               string            `json:"-" gorm:"rule_config"` // rule config
	RuleConfigJson           interface{}       `json:"rule_config" gorm:"-"` // rule config for fe
	PromEvalInterval         int               `json:"prom_eval_interval"`
	Callbacks                string            `json:"-"`                  // for db
	CallbacksJSON            []string          `json:"callbacks" gorm:"-"` // for fe
	RunbookUrl               string            `json:"runbook_url"`
	NotifyRecovered          int               `json:"notify_recovered"`
	NotifyChannels           string            `json:"-"`                          // for db
	NotifyChannelsJSON       []string          `json:"notify_channels" gorm:"-"`   // for fe
	NotifyGroups             string            `json:"-"`                          // for db
	NotifyGroupsJSON         []string          `json:"notify_groups" gorm:"-"`     // for fe
	NotifyGroupsObj          []*UserGroup      `json:"notify_groups_obj" gorm:"-"` // for fe
	TargetIdent              string            `json:"target_ident"`
	TargetNote               string            `json:"target_note"`
	TriggerTime              int64             `json:"trigger_time"`
	TriggerValue             string            `json:"trigger_value"`
	Tags                     string            `json:"-"`                         // for db
	TagsJSON                 []string          `json:"tags" gorm:"-"`             // for fe
	TagsMap                  map[string]string `json:"-" gorm:"-"`                // for internal usage
	Annotations              string            `json:"-"`                         //
	AnnotationsJSON          map[string]string `json:"annotations" gorm:"-"`      // for fe
	IsRecovered              bool              `json:"is_recovered" gorm:"-"`     // for notify.py
	NotifyUsersObj           []*User           `json:"notify_users_obj" gorm:"-"` // for notify.py
	LastEvalTime             int64             `json:"last_eval_time" gorm:"-"`   // for notify.py 上次计算的时间
	LastEscalationNotifyTime int64             `json:"last_escalation_notify_time" gorm:"-"`
	LastSentTime             int64             `json:"last_sent_time" gorm:"-"` // 上次发送时间
	NotifyCurNumber          int               `json:"notify_cur_number"`       // notify: current number
	FirstTriggerTime         int64             `json:"first_trigger_time"`      // 连续告警的首次告警时间
	ExtraConfig              interface{}       `json:"extra_config" gorm:"-"`
	Status                   int               `json:"status" gorm:"-"`
	Claimant                 string            `json:"claimant" gorm:"-"`
	SubRuleId                int64             `json:"sub_rule_id" gorm:"-"`
	DeletedAt                gorm.DeletedAt    `gorm:"column:deleted_at" json:"deleted_at" swaggerignore:"true"`
}

type AlterImport struct {
	RuleName         string `json:"rule_name" cn:"规则标题" xml:"rule_name"`
	AssetName        string `json:"asset_name" cn:"资产名称" xml:"asset_name"`
	AssetIp          string `json:"asset_ip" cn:"资产IP" xml:"asset_ip"`
	Severity         string `json:"severity" cn:"告警级别" xml:"severity" validate:"omitempty,oneof=0 1 2 3" source:"type=option,value=[请输入告警级别;紧急;一般;事件]"`
	PromSql          string `json:"prom_sql" cn:"指标" xml:"prom_sql"`
	TriggerTime      string `json:"trigger_time" cn:"触发时间" xml:"trigger_time" `
	TriggerValue     string `json:"trigger_value" cn:"触发时值" xml:"trigger_value"`
	RecoverTime      string `json:"recover_time" cn:"恢复时间" xml:"recover_time" `
	PromEvalInterval int    `json:"prom_eval_interval" cn:"执行频率/s" xml:"prom_eval_interval"`
	PromForDuration  int    `json:"prom_for_duration" cn:"持续时长/s" xml:"prom_for_duration"`
}

type AlterXml struct {
	AlterData []AlterImport `xml:"alter_data"`
}

func (e *AlertCurEvent) TableName() string {
	return "alert_cur_event"
}

func (e *AlertCurEvent) Add(ctx *ctx.Context) error {
	return Insert(ctx, e)
}

type AggrRule struct {
	Type  string
	Value string
}

func (e *AlertCurEvent) ParseRule(field string) error {
	f := e.GetField(field)
	f = strings.TrimSpace(f)

	if f == "" {
		return nil
	}

	var defs = []string{
		"{{$labels := .TagsMap}}",
		"{{$value := .TriggerValue}}",
	}

	text := strings.Join(append(defs, f), "")
	t, err := template.New(fmt.Sprint(e.RuleId)).Funcs(template.FuncMap(tplx.TemplateFuncMap)).Parse(text)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = t.Execute(&body, e)
	if err != nil {
		return err
	}

	if field == "rule_name" {
		e.RuleName = body.String()
	}

	if field == "rule_note" {
		e.RuleNote = body.String()
	}

	if field == "annotations" {
		e.Annotations = body.String()
		json.Unmarshal([]byte(e.Annotations), &e.AnnotationsJSON)
	}

	return nil
}

func (e *AlertCurEvent) GenCardTitle(rules []*AggrRule) string {
	arr := make([]string, len(rules))
	for i := 0; i < len(rules); i++ {
		rule := rules[i]

		if rule.Type == "field" {
			arr[i] = e.GetField(rule.Value)
		}

		if rule.Type == "tagkey" {
			arr[i] = e.GetTagValue(rule.Value)
		}

		if len(arr[i]) == 0 {
			arr[i] = "Null"
		}
	}
	return strings.Join(arr, "::")
}

func (e *AlertCurEvent) GetTagValue(tagkey string) string {
	for _, tag := range e.TagsJSON {
		i := strings.Index(tag, tagkey+"=")
		if i >= 0 {
			return tag[len(tagkey+"="):]
		}
	}
	return ""
}

func (e *AlertCurEvent) GetField(field string) string {
	switch field {
	case "cluster":
		return e.Cluster
	case "group_id":
		return fmt.Sprint(e.GroupId)
	case "group_name":
		return e.GroupName
	case "rule_id":
		return fmt.Sprint(e.RuleId)
	case "rule_name":
		return e.RuleName
	case "rule_note":
		return e.RuleNote
	case "severity":
		return fmt.Sprint(e.Severity)
	case "runbook_url":
		return e.RunbookUrl
	case "target_ident":
		return e.TargetIdent
	case "target_note":
		return e.TargetNote
	case "callbacks":
		return e.Callbacks
	case "annotations":
		return e.Annotations
	default:
		return ""
	}
}

func (e *AlertCurEvent) ToHis(ctx *ctx.Context) *AlertHisEvent {
	isRecovered := 0
	var recoverTime int64 = 0
	// var status int = 0 //添加Status参数
	if e.IsRecovered {
		isRecovered = 1
		recoverTime = e.LastEvalTime
	}

	return &AlertHisEvent{
		IsRecovered:      isRecovered,
		AssetId:          e.AssetId,
		Cate:             e.Cate,
		Cluster:          e.Cluster,
		DatasourceId:     e.DatasourceId,
		GroupId:          e.GroupId,
		GroupName:        e.GroupName,
		Hash:             e.Hash,
		RuleId:           e.RuleId,
		RuleName:         e.RuleName,
		RuleProd:         e.RuleProd,
		RuleAlgo:         e.RuleAlgo,
		RuleNote:         e.RuleNote,
		Severity:         e.Severity,
		PromForDuration:  e.PromForDuration,
		PromQl:           e.PromQl,
		PromEvalInterval: e.PromEvalInterval,
		RuleConfig:       e.RuleConfig,
		RuleConfigJson:   e.RuleConfigJson,
		Callbacks:        e.Callbacks,
		RunbookUrl:       e.RunbookUrl,
		NotifyRecovered:  e.NotifyRecovered,
		NotifyChannels:   e.NotifyChannels,
		NotifyGroups:     e.NotifyGroups,
		Annotations:      e.Annotations,
		AnnotationsJSON:  e.AnnotationsJSON,
		TargetIdent:      e.TargetIdent,
		TargetNote:       e.TargetNote,
		TriggerTime:      e.TriggerTime,
		TriggerValue:     e.TriggerValue,
		Tags:             e.Tags,
		RecoverTime:      recoverTime,
		LastEvalTime:     e.LastEvalTime,
		NotifyCurNumber:  e.NotifyCurNumber,
		FirstTriggerTime: e.FirstTriggerTime,
		Status:           e.Status, //添加状态参数
	}
}

func (e *AlertCurEvent) DB2FE() error {
	e.NotifyChannelsJSON = strings.Fields(e.NotifyChannels)
	e.NotifyGroupsJSON = strings.Fields(e.NotifyGroups)
	e.CallbacksJSON = strings.Fields(e.Callbacks)
	e.TagsJSON = strings.Split(e.Tags, ",,")
	json.Unmarshal([]byte(e.Annotations), &e.AnnotationsJSON)
	json.Unmarshal([]byte(e.RuleConfig), &e.RuleConfigJson)
	return nil
}

func (e *AlertCurEvent) FE2DB() {
	e.NotifyChannels = strings.Join(e.NotifyChannelsJSON, " ")
	e.NotifyGroups = strings.Join(e.NotifyGroupsJSON, " ")
	e.Callbacks = strings.Join(e.CallbacksJSON, " ")
	e.Tags = strings.Join(e.TagsJSON, ",,")
	b, _ := json.Marshal(e.AnnotationsJSON)
	e.Annotations = string(b)

	b, _ = json.Marshal(e.RuleConfigJson)
	e.RuleConfig = string(b)

}

func (e *AlertCurEvent) DB2Mem() {
	e.IsRecovered = false
	e.NotifyGroupsJSON = strings.Fields(e.NotifyGroups)
	e.CallbacksJSON = strings.Fields(e.Callbacks)
	e.NotifyChannelsJSON = strings.Fields(e.NotifyChannels)
	e.TagsJSON = strings.Split(e.Tags, ",,")
	e.TagsMap = make(map[string]string)
	for i := 0; i < len(e.TagsJSON); i++ {
		pair := strings.TrimSpace(e.TagsJSON[i])
		if pair == "" {
			continue
		}

		arr := strings.Split(pair, "=")
		if len(arr) != 2 {
			continue
		}

		e.TagsMap[arr[0]] = arr[1]
	}
}

// for webui
func (e *AlertCurEvent) FillNotifyGroups(ctx *ctx.Context, cache map[int64]*UserGroup) error {
	// some user-group already deleted ?
	count := len(e.NotifyGroupsJSON)
	if count == 0 {
		e.NotifyGroupsObj = []*UserGroup{}
		return nil
	}

	for i := range e.NotifyGroupsJSON {
		id, err := strconv.ParseInt(e.NotifyGroupsJSON[i], 10, 64)
		if err != nil {
			continue
		}

		ug, has := cache[id]
		if has {
			e.NotifyGroupsObj = append(e.NotifyGroupsObj, ug)
			continue
		}

		ug, err = UserGroupGetById(ctx, id)
		if err != nil {
			return err
		}

		if ug != nil {
			e.NotifyGroupsObj = append(e.NotifyGroupsObj, ug)
			cache[id] = ug
		}
	}

	return nil
}

func AlertCurEventTotal(ctx *ctx.Context, prods []string, bgid, stime, etime int64, severity int, dsIds []int64, cates []string, query string) (int64, error) {
	session := DB(ctx).Model(&AlertCurEvent{}).Where("trigger_time between ? and ?", stime, etime)

	if len(prods) != 0 {
		session = session.Where("rule_prod in ?", prods)
	}

	if bgid > 0 {
		session = session.Where("group_id = ?", bgid)
	}

	if severity >= 0 {
		session = session.Where("severity = ?", severity)
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

	return Count(session)
}

func AlertCurEventGets(ctx *ctx.Context, prods []string, bgid, stime, etime int64, severity int, dsIds []int64, cates []string, query string, limit, offset int) ([]AlertCurEvent, error) {
	session := DB(ctx)

	if stime != 0 {
		session = session.Where("trigger_time between ? and ?", stime, etime)
	}

	if len(prods) != 0 {
		session = session.Where("rule_prod in ?", prods)
	}

	if bgid > 0 {
		session = session.Where("group_id = ?", bgid)
	}

	if severity >= 0 {
		session = session.Where("severity = ?", severity)
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

	var lst []AlertCurEvent
	err := session.Order("id desc").Limit(limit).Offset(offset).Find(&lst).Error

	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].DB2FE()
		}
	}

	return lst, err
}

func AlertCurEventDel(ctx *ctx.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}

	return DB(ctx).Where("id in ?", ids).Delete(&AlertCurEvent{}).Error
}

func AlertCurEventDelByHash(ctx *ctx.Context, hash string) error {
	return DB(ctx).Where("hash = ?", hash).Delete(&AlertCurEvent{}).Error
}

func AlertCurEventExists(ctx *ctx.Context, where string, args ...interface{}) (bool, error) {
	return Exists(DB(ctx).Model(&AlertCurEvent{}).Where(where, args...))
}

func AlertCurEventGet(ctx *ctx.Context, where string, args ...interface{}) (*AlertCurEvent, error) {
	var lst []*AlertCurEvent
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

func AlertCurEventGetById(ctx *ctx.Context, id int64) (*AlertCurEvent, error) {
	return AlertCurEventGet(ctx, "id=?", id)
}

type AlertNumber struct {
	GroupId    int64
	GroupCount int64
}

// for busi_group list page
func AlertNumbers(ctx *ctx.Context, bgids []int64) (map[int64]int64, error) {
	ret := make(map[int64]int64)
	if len(bgids) == 0 {
		return ret, nil
	}

	var arr []AlertNumber
	err := DB(ctx).Model(&AlertCurEvent{}).Select("group_id", "count(*) as group_count").Where("group_id in ?", bgids).Group("group_id").Find(&arr).Error
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(arr); i++ {
		ret[arr[i].GroupId] = arr[i].GroupCount
	}

	return ret, nil
}

func AlertCurEventGetByIds(ctx *ctx.Context, ids []int64) ([]*AlertCurEvent, error) {
	var lst []*AlertCurEvent

	if len(ids) == 0 {
		return lst, nil
	}

	err := DB(ctx).Debug().Where("id in ?", ids).Order("id desc").Find(&lst).Error
	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].DB2FE()
		}
	}

	return lst, err
}

func AlertCurEventGetByRuleIdAndDsId(ctx *ctx.Context, ruleId int64, datasourceId int64) ([]*AlertCurEvent, error) {
	if !ctx.IsCenter {
		lst, err := poster.GetByUrls[[]*AlertCurEvent](ctx, "/v1/n9e/alert-cur-events-get-by-rid?rid="+strconv.FormatInt(ruleId, 10)+"&dsid="+strconv.FormatInt(datasourceId, 10))
		if err == nil {
			for i := 0; i < len(lst); i++ {
				lst[i].FE2DB()
			}
		}
		return lst, err
	}

	var lst []*AlertCurEvent
	err := DB(ctx).Where("rule_id=? and datasource_id = ?", ruleId, datasourceId).Find(&lst).Error
	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].DB2FE()
		}
	}
	return lst, err
}

func AlertCurEventGetMap(ctx *ctx.Context, cluster string) (map[int64]map[string]struct{}, error) {
	session := DB(ctx).Model(&AlertCurEvent{})
	if cluster != "" {
		session = session.Where("datasource_id = ?", cluster)
	}

	var lst []*AlertCurEvent
	err := session.Select("rule_id", "hash").Find(&lst).Error
	if err != nil {
		return nil, err
	}

	ret := make(map[int64]map[string]struct{})
	for i := 0; i < len(lst); i++ {
		rid := lst[i].RuleId
		hash := lst[i].Hash
		if _, has := ret[rid]; has {
			ret[rid][hash] = struct{}{}
		} else {
			ret[rid] = make(map[string]struct{})
			ret[rid][hash] = struct{}{}
		}
	}

	return ret, nil
}

// used by busi_group_overview
type AlertOverview struct {
	AlertNumber
	GroupName    string
	GroupTargets int64
	Emergency    int64
	Warning      int64
	Notice       int64
}

func (m *AlertCurEvent) UpdateFieldsMap(ctx *ctx.Context, fields map[string]interface{}) error {
	return DB(ctx).Model(m).Updates(fields).Error
}

func AlertCurEventUpgradeToV6(ctx *ctx.Context, dsm map[string]Datasource) error {
	var lst []*AlertCurEvent
	err := DB(ctx).Where("trigger_time > ?", time.Now().Unix()-3600*24*30).Find(&lst).Error
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

// 定义适用于前端返回的结构体
type FeAlert struct {
	Id             int64       `json:"id"`            //告警id
	Name           string      `json:"name"`          //告警规则名称 rulename
	Rule           string      `json:"rule"`          //告警规则 RuleConfigJson中prom_ql字段
	Severity       int         `json:"severity"`      //告警级别 0:紧急 1:警告 2:提醒 ，0为最高
	AssetId        int         `json:"asset_id"`      //资产id，对应TagsJSON中ident字段
	AssetName      string      `json:"asset_name"`    //资产名称
	TriggerTime    int64       `json:"trigger_time"`  //trigger_value
	TriggerValue   string      `json:"trigger_value"` //trigger_value
	OrganizeId     int         `json:"organize_id"`   //组织id
	OrganizeName   string      `json:"organize_name"` //组织name
	Url            string      `json:"url"`
	Tags           string      `json:"-"`
	TagsJSON       []string    `json:"-"`
	RuleConfig     string      `json:"-"`
	RuleConfigJson interface{} `json:"-"`
	Type           string      `json:"type"`
	Ip             string      `json:"ip"`
	// RuleConfigJson interface{} `json:"ruleconfigJson"`
}

type ruleConfigJson struct {
	Queries []map[string]interface{} `json:"queries"`
}

// 生成新的适用于前端页面的返回数据格式
func MakeFeAlert(dat []*AlertCurEvent) (Fe []*FeAlert) {
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

// 生成新的适用于前端页面的返回数据格式
// func AlertFeList(ctx *ctx.Context) ([]*FeAlert, error) {
// 	var dat []*AlertCurEvent
// 	var fedat []*FeAlert
// 	var assetID int
// 	err := DB(ctx).Find(&dat).Error
// 	for i := 0; i < len(dat); i++ {
// 		dat[i].DB2FE()

// 		for u := 0; u < len(dat[i].TagsJSON); u++ {
// 			s := strings.Split(dat[i].TagsJSON[u], "=")
// 			if s[0] == "asset_id" {
// 				assetID, _ = strconv.Atoi(s[1])
// 				break
// 			}
// 		}

// 		dic := &ruleConfigJson{}
// 		json.Unmarshal([]byte(dat[i].RuleConfig), dic)
// 		fedat = append(fedat, &FeAlert{
// 			Id:           dat[i].Id,
// 			Name:         dat[i].RuleName,
// 			Severity:     dat[i].Severity,
// 			TriggerTime:  dat[i].TriggerTime,
// 			TriggerValue: dat[i].TriggerValue,
// 			Rule:         dic.Queries[0]["prom_ql"].(string),
// 			AssetId:      assetID,
// 		})
// 	}
// 	return fedat, err
// }

func AlertFeList(ctx *ctx.Context) ([]*FeAlert, error) {
	var dat []*AlertCurEvent
	err := DB(ctx).Find(&dat).Error
	Fe := MakeFeAlert(dat)
	return Fe, err
}

//统计未处理告警
func AlertCurCount(ctx *ctx.Context) (num int64, err error) {
	err = DB(ctx).Debug().Model(&AlertCurEvent{}).Count(&num).Error
	return num, err
}

//通过资产id查询当前告警
func AlertFeListByAssetId(ctx *ctx.Context, where string) ([]*FeAlert, error) {
	where = "%asset_id=" + where + "%"
	var dat []*AlertCurEvent
	err := DB(ctx).Where("tags like ?", where).Find(&dat).Error
	Fe := MakeFeAlert(dat)
	return Fe, err
}

//西航
func AlertEventXHTotal(ctx *ctx.Context, alertType, severity, group, start, end int64, ids []int64, query string) (int64, error) {
	session := DB(ctx)
	if start != -1 {
		session = session.Where("trigger_time >= ?", start)
	}
	if end != -1 {
		session = session.Where("trigger_time <= ?", end)
	}
	if severity != -1 {
		session = session.Where("severity = ?", severity)
	}
	if group != -1 {
		session = session.Where("group_id = ?", group)
	}
	if query != "" {
		query = "%" + query + "%"
		// if fType == -1 {
		session = session.Where("asset_id in ? or id like ? or rule_name like ? or severity like ?", ids, query, query, query)
		// } else if fType == 1 {
		// 	session = session.Where("severity like ?", query)
		// } else if fType == 2 {
		// 	session = session.Where("asset_id in ?", ids)
		// }
	}
	var num int64
	var err error
	if alertType == 1 {
		err = session.Model(&AlertCurEvent{}).Count(&num).Error
	} else if alertType == 2 {
		err = session.Model(&AlertHisEvent{}).Count(&num).Error
	}

	return num, err
}

func AlertEventXHGets[T any](ctx *ctx.Context, alertType, severity, group, start, end int64, ids []int64, query string, limit, offset int) ([]T, error) {
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
	if severity != -1 {
		session = session.Where("severity = ?", severity)
	}
	if group != -1 {
		session = session.Where("group_id = ?", group)
	}
	if query != "" {
		query = "%" + query + "%"
		// if fType == -1 {
		session = session.Where("asset_id in ? or id like ? or rule_name like ? or severity like ?", ids, query, query, query)
		// } else if fType == 1 {
		// 	session = session.Where("severity like ?", query)
		// } else if fType == 2 {
		// 	session = session.Where("asset_id in ?", ids)
		// }
	}
	var lst []T
	var err error
	if alertType == 1 {
		err = session.Model(&AlertCurEvent{}).Find(&lst).Error
	} else if alertType == 2 {
		err = session.Where("is_recovered != 0").Model(&AlertHisEvent{}).Find(&lst).Error
	}
	// for index := range lst {
	// 	lst[index].DB2FE()
	// }
	return lst, err

}

func AlertCurEventDelByIds(ctx *ctx.Context, ids []int64) error {
	return DB(ctx).Where("id in ?", ids).Delete(&AlertCurEvent{}).Error
}

func AlertCurEventDelByIdsTx(tx *gorm.DB, ids []int64) error {
	err := tx.Where("id in ?", ids).Delete(&AlertCurEvent{}).Error
	if err != nil {
		tx.Rollback()
	}
	return nil
}
