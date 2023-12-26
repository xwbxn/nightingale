package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/poster"
	"gorm.io/gorm"

	"github.com/pkg/errors"
	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/str"
)

const (
	METRIC = "metric"
	HOST   = "host"

	PROMETHEUS = "prometheus"
)

type AlertRule struct {
	Id                    int64             `json:"id" gorm:"primaryKey"`
	GroupId               int64             `json:"group_id"` // busi group id
	AssetId               int64             `json:"asset_id"`
	AssetName             string            `json:"asset_name" gorm:"-"`
	AssetIp               string            `json:"asset_ip" gorm:"-"`
	MonitoringId          int64             `json:"monitoring_id"`
	Cate                  string            `json:"cate"`                                 // alert rule cate (prometheus|elasticsearch)
	DatasourceIds         string            `json:"-" gorm:"datasource_ids"`              // datasource ids
	DatasourceIdsJson     []int64           `json:"datasource_ids" gorm:"-"`              // for fe
	Cluster               string            `json:"cluster"`                              // take effect by clusters, seperated by space
	Name                  string            `json:"name"`                                 // rule name
	Note                  string            `json:"note"`                                 // will sent in notify
	Prod                  string            `json:"prod"`                                 // product empty means n9e
	Algorithm             string            `json:"algorithm"`                            // algorithm (''|holtwinters), empty means threshold
	AlgoParams            string            `json:"-" gorm:"algo_params"`                 // params algorithm need
	AlgoParamsJson        interface{}       `json:"algo_params" gorm:"-"`                 // for fe
	Delay                 int               `json:"delay"`                                // Time (in seconds) to delay evaluation
	Severity              int               `json:"severity"`                             // 1: Emergency 2: Warning 3: Notice
	Severities            []int             `json:"severities" gorm:"-"`                  // 1: Emergency 2: Warning 3: Notice
	Disabled              int               `json:"disabled"`                             // 0: enabled, 1: disabled
	PromForDuration       int               `json:"prom_for_duration"`                    // prometheus for, unit:s
	PromQl                string            `json:"prom_ql"`                              // just one ql
	RuleConfig            string            `json:"-" gorm:"rule_config"`                 // rule config
	RuleConfigCn          string            `json:"rule_config_cn" gorm:"rule_config_cn"` // rule config
	RuleConfigJson        interface{}       `json:"rule_config" gorm:"-"`                 // rule config for fe
	RuleConfigFe          string            `json:"rule_config_fe" gorm:"rule_config_fe"` // rule config for fe
	PromEvalInterval      int               `json:"prom_eval_interval"`                   // unit:s
	EnableStime           string            `json:"-"`                                    // split by space: "00:00 10:00 12:00"
	EnableStimeJSON       string            `json:"enable_stime" gorm:"-"`                // for fe
	EnableStimesJSON      []string          `json:"enable_stimes" gorm:"-"`               // for fe
	EnableEtime           string            `json:"-"`                                    // split by space: "00:00 10:00 12:00"
	EnableEtimeJSON       string            `json:"enable_etime" gorm:"-"`                // for fe
	EnableEtimesJSON      []string          `json:"enable_etimes" gorm:"-"`               // for fe
	EnableDaysOfWeek      string            `json:"-"`                                    // eg: "0 1 2 3 4 5 6 ; 0 1 2"
	EnableDaysOfWeekJSON  []string          `json:"enable_days_of_week" gorm:"-"`         // for fe
	EnableDaysOfWeeksJSON [][]string        `json:"enable_days_of_weeks" gorm:"-"`        // for fe
	EnableInBG            int               `json:"enable_in_bg"`                         // 0: global 1: enable one busi-group
	NotifyRecovered       int               `json:"notify_recovered"`                     // whether notify when recovery
	NotifyChannels        string            `json:"-"`                                    // split by space: sms voice email dingtalk wecom
	NotifyChannelsJSON    []string          `json:"notify_channels" gorm:"-"`             // for fe
	NotifyGroups          string            `json:"-"`                                    // split by space: 233 43
	NotifyGroupsObj       []UserGroup       `json:"notify_groups_obj" gorm:"-"`           // for fe
	NotifyGroupsJSON      []string          `json:"notify_groups" gorm:"-"`               // for fe
	NotifyRepeatStep      int               `json:"notify_repeat_step"`                   // notify repeat interval, unit: min
	NotifyMaxNumber       int               `json:"notify_max_number"`                    // notify: max number
	RecoverDuration       int64             `json:"recover_duration"`                     // unit: s
	Callbacks             string            `json:"-"`                                    // split by space: http://a.com/api/x http://a.com/api/y'
	CallbacksJSON         []string          `json:"callbacks" gorm:"-"`                   // for fe
	RunbookUrl            string            `json:"runbook_url"`                          // sop url
	AppendTags            string            `json:"-"`                                    // split by space: service=n9e mod=api
	AppendTagsJSON        []string          `json:"append_tags" gorm:"-"`                 // for fe
	Annotations           string            `json:"-"`                                    //
	AnnotationsJSON       map[string]string `json:"annotations" gorm:"-"`                 // for fe
	ExtraConfig           string            `json:"-" gorm:"extra_config"`                // extra config
	ExtraConfigJSON       interface{}       `json:"extra_config" gorm:"-"`                // for fe
	CreateAt              int64             `json:"create_at"`
	CreateBy              string            `json:"create_by"`
	UpdateAt              int64             `json:"update_at"`
	UpdateBy              string            `json:"update_by"`
	DeletedAt             gorm.DeletedAt    `gorm:"column:deleted_at" json:"deleted_at" swaggerignore:"true"`
}

type PromRuleConfig struct {
	Queries    []PromQuery `json:"queries"`
	Inhibit    bool        `json:"inhibit"`
	PromQl     string      `json:"prom_ql"`
	Severity   int         `json:"severity"`
	AlgoParams interface{} `json:"algo_params"`
}

type HostRuleConfig struct {
	Queries  []HostQuery   `json:"queries"`
	Triggers []HostTrigger `json:"triggers"`
	Inhibit  bool          `json:"inhibit"`
}

type PromQuery struct {
	PromQl   string `json:"prom_ql"`
	Severity int    `json:"severity"`
}

type HostTrigger struct {
	Type     string `json:"type"`
	Duration int    `json:"duration"`
	Percent  int    `json:"percent"`
	Severity int    `json:"severity"`
}

type RuleQuery struct {
	Queries  []interface{} `json:"queries"`
	Triggers []Trigger     `json:"triggers"`
}

type Trigger struct {
	Expressions interface{} `json:"expressions"`
	Mode        int         `json:"mode"`
	Exp         string      `json:"exp"`
	Severity    int         `json:"severity"`
}

type AlertRuleSimplify struct {
	Id           int64  `json:"id"`
	RuleConfigCn string `json:"rule_config_cn"`
	Relation     string `json:"relation"`
	Value        int64  `json:"value"`
	Severity     int    `json:"severity"`
}

func GetHostsQuery(queries []HostQuery) []map[string]interface{} {
	var query []map[string]interface{}
	for _, q := range queries {
		m := make(map[string]interface{})
		switch q.Key {
		case "group_ids":
			ids := ParseInt64(q.Values)
			if q.Op == "==" {
				m["group_id in (?)"] = ids
			} else {
				m["group_id not in (?)"] = ids
			}
		case "tags":
			lst := []string{}
			for _, v := range q.Values {
				if v == nil {
					continue
				}
				lst = append(lst, v.(string))
			}
			if q.Op == "==" {
				blank := " "
				for _, tag := range lst {
					m["tags like ?"+blank] = "%" + tag + "%"
					blank += " "
				}
			} else {
				blank := " "
				for _, tag := range lst {
					m["tags not like ?"+blank] = "%" + tag + "%"
					blank += " "
				}
			}
		case "hosts":
			lst := []string{}
			for _, v := range q.Values {
				if v == nil {
					continue
				}
				lst = append(lst, v.(string))
			}
			if q.Op == "==" {
				m["ident in (?)"] = lst
			} else {
				m["ident not in (?)"] = lst
			}
		}
		query = append(query, m)
	}
	return query
}

func ParseInt64(values []interface{}) []int64 {
	b, _ := json.Marshal(values)
	var ret []int64
	json.Unmarshal(b, &ret)
	return ret
}

type HostQuery struct {
	Key    string        `json:"key"`
	Op     string        `json:"op"`
	Values []interface{} `json:"values"`
}

func Str2Int(arr []string) []int64 {
	var ret []int64
	for _, v := range arr {
		i, _ := strconv.ParseInt(v, 10, 64)
		ret = append(ret, i)
	}
	return ret
}

func (ar *AlertRule) TableName() string {
	return "alert_rule"
}

func (ar *AlertRule) Verify() error {
	if ar.GroupId < 0 {
		return fmt.Errorf("GroupId(%d) invalid", ar.GroupId)
	}

	if IsAllDatasource(ar.DatasourceIdsJson) {
		ar.DatasourceIdsJson = []int64{0}
	}

	if str.Dangerous(ar.Name) {
		return errors.New("告警规则名称不符合命名规范，包含特殊字符")
	}

	if ar.Name == "" {
		return errors.New("名称为空")
	}

	if ar.Prod == "" {
		ar.Prod = METRIC
	}

	if ar.Cate == "" {
		ar.Cate = PROMETHEUS
	}

	if ar.RuleConfig == "" {
		return errors.New("告警规则配置为空")
	}

	if ar.PromEvalInterval <= 0 {
		ar.PromEvalInterval = 15
	}

	// check in front-end
	// if _, err := parser.ParseExpr(ar.PromQl); err != nil {
	// 	return errors.New("prom_ql parse error: %")
	// }

	ar.AppendTags = strings.TrimSpace(ar.AppendTags)
	arr := strings.Fields(ar.AppendTags)
	for i := 0; i < len(arr); i++ {
		if len(strings.Split(arr[i], "=")) != 2 {
			return fmt.Errorf("AppendTags(%s) invalid", arr[i])
		}
	}

	gids := strings.Fields(ar.NotifyGroups)
	for i := 0; i < len(gids); i++ {
		if _, err := strconv.ParseInt(gids[i], 10, 64); err != nil {
			return fmt.Errorf("NotifyGroups(%s) invalid", ar.NotifyGroups)
		}
	}

	return nil
}

func (ar *AlertRule) Add(ctx *ctx.Context) error {
	if err := ar.Verify(); err != nil {
		return err
	}

	exists, err := AlertRuleExists(ctx, 0, ar.GroupId, ar.DatasourceIdsJson, ar.Name)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("在同一资产下，告警名称不能重复")
	}

	now := time.Now().Unix()
	ar.CreateAt = now
	ar.UpdateAt = now

	return Insert(ctx, ar)
}

func (ar *AlertRule) AddTx(ctx *ctx.Context, tx *gorm.DB) error {
	if err := ar.Verify(); err != nil {
		return err
	}

	exists, err := AlertRuleExists(ctx, 0, ar.GroupId, ar.DatasourceIdsJson, ar.Name)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("在同一资产下，告警名称不能重复")
	}

	now := time.Now().Unix()
	ar.CreateAt = now
	ar.UpdateAt = now

	err = tx.Create(ar).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

func (ar *AlertRule) Update(ctx *ctx.Context, arf AlertRule) error {
	if ar.Name != arf.Name {
		exists, err := AlertRuleExists(ctx, ar.Id, ar.GroupId, ar.DatasourceIdsJson, arf.Name)
		if err != nil {
			return err
		}

		if exists {
			return errors.New("AlertRule already exists")
		}
	}

	err := arf.FE2DB(ctx)
	if err != nil {
		return err
	}

	arf.Id = ar.Id
	arf.GroupId = ar.GroupId
	arf.CreateAt = ar.CreateAt
	arf.CreateBy = ar.CreateBy
	arf.UpdateAt = time.Now().Unix()

	err = arf.Verify()
	if err != nil {
		return err
	}
	return DB(ctx).Model(ar).Select("*").Updates(arf).Error
}

func (ar *AlertRule) UpdateTx(ctx *ctx.Context, tx *gorm.DB, arf AlertRule) error {
	if ar.Name != arf.Name {
		exists, err := AlertRuleExists(ctx, ar.Id, ar.GroupId, ar.DatasourceIdsJson, arf.Name)
		if err != nil {
			return err
		}

		if exists {
			return errors.New("AlertRule already exists")
		}
	}

	err := arf.FE2DB(ctx)
	if err != nil {
		return err
	}

	arf.Id = ar.Id
	arf.GroupId = ar.GroupId
	arf.CreateAt = ar.CreateAt
	arf.CreateBy = ar.CreateBy
	arf.UpdateAt = time.Now().Unix()

	err = arf.Verify()
	if err != nil {
		return err
	}
	err = tx.Model(ar).Select("*").Updates(arf).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

func (ar *AlertRule) UpdateColumn(ctx *ctx.Context, column string, value interface{}) error {
	if value == nil {
		return nil
	}

	if column == "datasource_ids" {
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		return DB(ctx).Model(ar).UpdateColumn(column, string(b)).Error
	}

	if column == "severity" {
		severity := int(value.(float64))
		if ar.Cate == PROMETHEUS {
			var ruleConfig PromRuleConfig
			err := json.Unmarshal([]byte(ar.RuleConfig), &ruleConfig)
			if err != nil {
				return err
			}

			if len(ruleConfig.Queries) < 1 {
				ruleConfig.Severity = severity
				b, err := json.Marshal(ruleConfig)
				if err != nil {
					return err
				}
				return DB(ctx).Model(ar).UpdateColumn("rule_config", string(b)).Error
			}

			if len(ruleConfig.Queries) != 1 {
				return nil
			}

			ruleConfig.Queries[0].Severity = severity
			b, err := json.Marshal(ruleConfig)
			if err != nil {
				return err
			}
			return DB(ctx).Model(ar).UpdateColumn("rule_config", string(b)).Error
		} else if ar.Cate == HOST {
			var ruleConfig HostRuleConfig
			err := json.Unmarshal([]byte(ar.RuleConfig), &ruleConfig)
			if err != nil {
				return err
			}

			if len(ruleConfig.Triggers) != 1 {
				return nil
			}

			ruleConfig.Triggers[0].Severity = severity

			b, err := json.Marshal(ruleConfig)
			if err != nil {
				return err
			}
			return DB(ctx).Model(ar).UpdateColumn("rule_config", string(b)).Error
		} else {
			var ruleConfig RuleQuery
			err := json.Unmarshal([]byte(ar.RuleConfig), &ruleConfig)
			if err != nil {
				return err
			}

			if len(ruleConfig.Triggers) != 1 {
				return nil
			}

			ruleConfig.Triggers[0].Severity = severity
			b, err := json.Marshal(ruleConfig)
			if err != nil {
				return err
			}
			return DB(ctx).Model(ar).UpdateColumn("rule_config", string(b)).Error
		}
	}

	if column == "runbook_url" {
		url := value.(string)

		err := json.Unmarshal([]byte(ar.Annotations), &ar.AnnotationsJSON)
		if err != nil {
			return err
		}

		if ar.AnnotationsJSON == nil {
			ar.AnnotationsJSON = make(map[string]string)
		}

		ar.AnnotationsJSON["runbook_url"] = url

		b, err := json.Marshal(ar.AnnotationsJSON)
		if err != nil {
			return err
		}

		return DB(ctx).Model(ar).UpdateColumn("annotations", string(b)).Error
	}

	return DB(ctx).Model(ar).UpdateColumn(column, value).Error
}

func (ar *AlertRule) UpdateFieldsMap(ctx *ctx.Context, fields map[string]interface{}) error {
	return DB(ctx).Model(ar).Updates(fields).Error
}

// for v5 rule
func (ar *AlertRule) FillDatasourceIds() error {
	if ar.DatasourceIds != "" {
		json.Unmarshal([]byte(ar.DatasourceIds), &ar.DatasourceIdsJson)
		return nil
	}
	return nil
}

func (ar *AlertRule) FillSeverities() error {
	if ar.RuleConfig != "" {
		if ar.Cate == PROMETHEUS {
			var rule PromRuleConfig
			if err := json.Unmarshal([]byte(ar.RuleConfig), &rule); err != nil {
				return err
			}

			if len(rule.Queries) == 0 {
				ar.Severities = append(ar.Severities, rule.Severity)
				return nil
			}

			for i := range rule.Queries {
				ar.Severities = append(ar.Severities, rule.Queries[i].Severity)
			}
		} else {
			var rule HostRuleConfig
			if err := json.Unmarshal([]byte(ar.RuleConfig), &rule); err != nil {
				return err
			}
			for i := range rule.Triggers {
				ar.Severities = append(ar.Severities, rule.Triggers[i].Severity)
			}
		}
	}
	return nil
}

func (ar *AlertRule) FillNotifyGroups(ctx *ctx.Context, cache map[int64]*UserGroup) error {
	// some user-group already deleted ?
	count := len(ar.NotifyGroupsJSON)
	if count == 0 {
		ar.NotifyGroupsObj = []UserGroup{}
		return nil
	}

	exists := make([]string, 0, count)
	delete := false
	for i := range ar.NotifyGroupsJSON {
		id, _ := strconv.ParseInt(ar.NotifyGroupsJSON[i], 10, 64)

		ug, has := cache[id]
		if has {
			exists = append(exists, ar.NotifyGroupsJSON[i])
			ar.NotifyGroupsObj = append(ar.NotifyGroupsObj, *ug)
			continue
		}

		ug, err := UserGroupGetById(ctx, id)
		if err != nil {
			return err
		}

		if ug == nil {
			delete = true
		} else {
			exists = append(exists, ar.NotifyGroupsJSON[i])
			ar.NotifyGroupsObj = append(ar.NotifyGroupsObj, *ug)
			cache[id] = ug
		}
	}

	if delete {
		// some user-group already deleted
		ar.NotifyGroupsJSON = exists
		ar.NotifyGroups = strings.Join(exists, " ")
		DB(ctx).Model(ar).Update("notify_groups", ar.NotifyGroups)
	}

	return nil
}

func (ar *AlertRule) FE2DB(ctx *ctx.Context) error {

	if len(ar.EnableStimesJSON) > 0 {
		ar.EnableStime = strings.Join(ar.EnableStimesJSON, " ")
		ar.EnableEtime = strings.Join(ar.EnableEtimesJSON, " ")
	} else {
		ar.EnableStime = ar.EnableStimeJSON
		ar.EnableEtime = ar.EnableEtimeJSON
	}

	if len(ar.EnableDaysOfWeeksJSON) > 0 {
		for i := 0; i < len(ar.EnableDaysOfWeeksJSON); i++ {
			if len(ar.EnableDaysOfWeeksJSON) == 1 {
				ar.EnableDaysOfWeek = strings.Join(ar.EnableDaysOfWeeksJSON[i], " ")
			} else {
				if i == len(ar.EnableDaysOfWeeksJSON)-1 {
					ar.EnableDaysOfWeek += strings.Join(ar.EnableDaysOfWeeksJSON[i], " ")
				} else {
					ar.EnableDaysOfWeek += strings.Join(ar.EnableDaysOfWeeksJSON[i], " ") + ";"
				}
			}
		}
	} else {
		ar.EnableDaysOfWeek = strings.Join(ar.EnableDaysOfWeekJSON, " ")
	}

	ar.NotifyChannels = strings.Join(ar.NotifyChannelsJSON, " ")
	ar.NotifyGroups = strings.Join(ar.NotifyGroupsJSON, " ")
	ar.Callbacks = strings.Join(ar.CallbacksJSON, " ")
	ar.AppendTags = strings.Join(ar.AppendTagsJSON, " ")
	algoParamsByte, err := json.Marshal(ar.AlgoParamsJson)
	if err != nil {
		return fmt.Errorf("marshal algo_params err:%v", err)
	}
	ar.AlgoParams = string(algoParamsByte)

	if len(ar.DatasourceIdsJson) > 0 {
		idsByte, err := json.Marshal(ar.DatasourceIdsJson)
		if err != nil {
			return fmt.Errorf("marshal datasource_ids err:%v", err)
		}
		ar.DatasourceIds = string(idsByte)
	}

	if ar.RuleConfigJson == nil {
		query := PromQuery{
			PromQl:   ar.PromQl,
			Severity: ar.Severity,
		}
		ar.RuleConfigJson = PromRuleConfig{
			Queries: []PromQuery{query},
		}
	} else {
		// configOld, err := json.Marshal(ar.RuleConfigJson)
		// if err != nil {
		// 	return fmt.Errorf("格式化数据错误:%v", err)
		// }
		// ar.RuleConfigJsonFe = string(configOld)
		query := ar.RuleConfigJson.(map[string]interface{})
		// mon, monOk := query["queries"]["monitor_id"]

		queries := query["queries"].([]interface{})
		for index := range queries {
			config := queries[index].(map[string]interface{})
			newConfig := make(map[string]interface{}, 0)
			mon, monOk := config["monitor_id"]
			if monOk {
				monitoring, err := MonitoringGetById(ctx, int64(mon.(float64)))
				if err != nil || monitoring == nil {
					return err
				}
				ar.MonitoringId = monitoring.Id
				newConfig["prom_ql"] = monitoring.MonitoringSql + config["relation"].(string) + config["value"].(string)
				newConfig["severity"] = int64(config["severity"].(float64))
				queries[index] = newConfig
			}
		}

	}

	// json.Marshal  RuleConfigJson
	if ar.RuleConfigJson != nil {
		b, err := json.Marshal(ar.RuleConfigJson)
		if err != nil {
			return fmt.Errorf("marshal rule_config err:%v", err)
		}
		ar.RuleConfig = string(b)
		ar.PromQl = ""
	}

	if ar.AnnotationsJSON != nil {
		b, err := json.Marshal(ar.AnnotationsJSON)
		if err != nil {
			return fmt.Errorf("marshal annotations err:%v", err)
		}
		ar.Annotations = string(b)
	}

	if ar.ExtraConfigJSON != nil {
		b, err := json.Marshal(ar.ExtraConfigJSON)
		if err != nil {
			return fmt.Errorf("marshal extra_config err:%v", err)
		}
		ar.ExtraConfig = string(b)
	}

	return nil
}

func (ar *AlertRule) DB2FE() error {
	ar.EnableStimesJSON = strings.Fields(ar.EnableStime)
	ar.EnableEtimesJSON = strings.Fields(ar.EnableEtime)
	if len(ar.EnableEtimesJSON) > 0 {
		ar.EnableStimeJSON = ar.EnableStimesJSON[0]
		ar.EnableEtimeJSON = ar.EnableEtimesJSON[0]
	}

	cache := strings.Split(ar.EnableDaysOfWeek, ";")
	for i := 0; i < len(cache); i++ {
		ar.EnableDaysOfWeeksJSON = append(ar.EnableDaysOfWeeksJSON, strings.Fields(cache[i]))
	}
	if len(ar.EnableDaysOfWeeksJSON) > 0 {
		ar.EnableDaysOfWeekJSON = ar.EnableDaysOfWeeksJSON[0]
	}

	ar.NotifyChannelsJSON = strings.Fields(ar.NotifyChannels)
	ar.NotifyGroupsJSON = strings.Fields(ar.NotifyGroups)
	ar.CallbacksJSON = strings.Fields(ar.Callbacks)
	ar.AppendTagsJSON = strings.Fields(ar.AppendTags)
	json.Unmarshal([]byte(ar.AlgoParams), &ar.AlgoParamsJson)
	json.Unmarshal([]byte(ar.RuleConfig), &ar.RuleConfigJson)
	json.Unmarshal([]byte(ar.Annotations), &ar.AnnotationsJSON)
	json.Unmarshal([]byte(ar.ExtraConfig), &ar.ExtraConfigJSON)

	err := ar.FillDatasourceIds()
	return err
}

func AlertRuleDels(ctx *ctx.Context, ids []int64, bgid ...int64) error {
	for i := 0; i < len(ids); i++ {
		session := DB(ctx).Where("id = ?", ids[i])
		if len(bgid) > 0 && bgid[0] > 0 {
			session = session.Where("group_id = ?", bgid[0])
		}
		ret := session.Delete(&AlertRule{})
		if ret.Error != nil {
			return ret.Error
		}

		// 说明确实删掉了，把相关的活跃告警也删了，这些告警永远都不会恢复了，而且策略都没了，说明没人关心了
		if ret.RowsAffected > 0 {
			DB(ctx).Where("rule_id = ?", ids[i]).Delete(new(AlertCurEvent))
		}
	}

	return nil
}

func AlertRuleExists(ctx *ctx.Context, id, groupId int64, datasourceIds []int64, name string) (bool, error) {
	session := DB(ctx).Where("id <> ? and group_id = ? and name = ?", id, groupId, name)

	var lst []AlertRule
	err := session.Find(&lst).Error
	if err != nil {
		return false, err
	}
	if len(lst) == 0 {
		return false, nil
	}

	// match cluster
	for _, r := range lst {
		r.FillDatasourceIds()
		for _, id := range r.DatasourceIdsJson {
			if MatchDatasource(datasourceIds, id) {
				return true, nil
			}
		}
	}
	return false, nil
}

func AlertRuleGets(ctx *ctx.Context, groupId int64) ([]AlertRule, error) {
	session := DB(ctx).Where("group_id=?", groupId).Order("name")

	var lst []AlertRule
	err := session.Find(&lst).Error
	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].DB2FE()
		}
	}

	return lst, err
}

func AlertRuleGetsAll(ctx *ctx.Context) ([]*AlertRule, error) {
	if !ctx.IsCenter {
		lst, err := poster.GetByUrls[[]*AlertRule](ctx, "/v1/n9e/alert-rules?disabled=0")
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(lst); i++ {
			lst[i].FE2DB(ctx)
		}
		return lst, err
	}

	session := DB(ctx).Where("disabled = ?", 0)

	var lst []*AlertRule
	err := session.Find(&lst).Error
	if err != nil {
		return lst, err
	}

	if len(lst) == 0 {
		return lst, nil
	}

	for i := 0; i < len(lst); i++ {
		lst[i].DB2FE()
	}
	return lst, nil
}

func AlertRulesGetsBy(ctx *ctx.Context, prods []string, query, algorithm, cluster string, cates []string, disabled int) ([]*AlertRule, error) {
	session := DB(ctx)

	if len(prods) > 0 {
		session = session.Where("prod in (?)", prods)
	}

	if query != "" {
		arr := strings.Fields(query)
		for i := 0; i < len(arr); i++ {
			qarg := "%" + arr[i] + "%"
			session = session.Where("append_tags like ?", qarg)
		}
	}

	if algorithm != "" {
		session = session.Where("algorithm = ?", algorithm)
	}

	if cluster != "" {
		session = session.Where("cluster like ?", "%"+cluster+"%")
	}

	if len(cates) != 0 {
		session = session.Where("cate in (?)", cates)
	}

	if disabled != -1 {
		session = session.Where("disabled = ?", disabled)
	}

	var lst []*AlertRule
	err := session.Find(&lst).Error
	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].DB2FE()
		}
	}

	return lst, err
}

func AlertRuleGet(ctx *ctx.Context, where string, args ...interface{}) (*AlertRule, error) {
	var lst []*AlertRule
	err := DB(ctx).Where(where, args...).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}

	lst[0].DB2FE()

	return lst[0], nil
}

func AlertRuleGetById(ctx *ctx.Context, id int64) (*AlertRule, error) {
	return AlertRuleGet(ctx, "id=?", id)
}

func AlertRuleGetByIds(ctx *ctx.Context, ids []int64) ([]AlertRule, error) {
	var lst []AlertRule
	err := DB(ctx).Model(&AlertRule{}).Where("id in ?", ids).Find(&lst).Error
	return lst, err
}

func AlertRuleGetName(ctx *ctx.Context, id int64) (string, error) {
	var names []string
	err := DB(ctx).Model(new(AlertRule)).Where("id = ?", id).Pluck("name", &names).Error
	if err != nil {
		return "", err
	}

	if len(names) == 0 {
		return "", nil
	}

	return names[0], nil
}

func AlertRuleStatistics(ctx *ctx.Context) (*Statistics, error) {
	if !ctx.IsCenter {
		s, err := poster.GetByUrls[*Statistics](ctx, "/v1/n9e/statistic?name=alert_rule")
		return s, err
	}

	session := DB(ctx).Model(&AlertRule{}).Select("count(*) as total", "max(update_at) as last_updated").Where("disabled = ?", 0)

	var stats []*Statistics
	err := session.Find(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats[0], nil
}

func (ar *AlertRule) IsPrometheusRule() bool {
	return ar.Prod == METRIC && ar.Cate == PROMETHEUS
}

func (ar *AlertRule) IsHostRule() bool {
	return ar.Prod == HOST
}

func (ar *AlertRule) GetRuleType() string {
	if ar.Prod == METRIC {
		return ar.Cate
	}

	return ar.Prod
}

func (ar *AlertRule) GenerateNewEvent(ctx *ctx.Context) *AlertCurEvent {
	event := &AlertCurEvent{}
	ar.UpdateEvent(event)
	return event
}

func (ar *AlertRule) UpdateEvent(event *AlertCurEvent) {
	if event == nil {
		return
	}

	event.GroupId = ar.GroupId
	event.AssetId = ar.AssetId
	event.Cate = ar.Cate
	event.RuleId = ar.Id
	event.RuleName = ar.Name
	event.RuleNote = ar.Note
	event.RuleProd = ar.Prod
	event.RuleAlgo = ar.Algorithm
	event.PromForDuration = ar.PromForDuration
	event.PromQl = ar.PromQl
	event.RuleConfig = ar.RuleConfig
	event.RuleConfigJson = ar.RuleConfigJson
	event.PromEvalInterval = ar.PromEvalInterval
	event.Callbacks = ar.Callbacks
	event.CallbacksJSON = ar.CallbacksJSON
	event.RunbookUrl = ar.RunbookUrl
	event.NotifyRecovered = ar.NotifyRecovered
	event.NotifyChannels = ar.NotifyChannels
	event.NotifyChannelsJSON = ar.NotifyChannelsJSON
	event.NotifyGroups = ar.NotifyGroups
	event.NotifyGroupsJSON = ar.NotifyGroupsJSON
}

func AlertRuleUpgradeToV6(ctx *ctx.Context, dsm map[string]Datasource) error {
	var lst []*AlertRule
	err := DB(ctx).Find(&lst).Error
	if err != nil {
		return err
	}

	for i := 0; i < len(lst); i++ {
		var ids []int64
		if lst[i].Cluster == "$all" {
			ids = append(ids, 0)
		} else {
			clusters := strings.Fields(lst[i].Cluster)
			for j := 0; j < len(clusters); j++ {
				if ds, exists := dsm[clusters[j]]; exists {
					ids = append(ids, ds.Id)
				}
			}
		}

		b, err := json.Marshal(ids)
		if err != nil {
			continue
		}
		lst[i].DatasourceIds = string(b)

		if lst[i].PromQl == "" {
			continue
		}

		ruleConfig := PromRuleConfig{
			Queries: []PromQuery{
				{
					PromQl:   lst[i].PromQl,
					Severity: lst[i].Severity,
				},
			},
		}
		b, _ = json.Marshal(ruleConfig)
		lst[i].RuleConfig = string(b)

		m := make(map[string]string)
		if lst[i].RunbookUrl != "" {
			m["runbook_url"] = lst[i].RunbookUrl

			b, err = json.Marshal(m)
			if err != nil {
				continue
			}

			lst[i].Annotations = string(b)
		}

		if lst[i].Prod == "" {
			lst[i].Prod = METRIC
		}

		if lst[i].Cate == "" {
			lst[i].Cate = PROMETHEUS
		}

		err = lst[i].UpdateFieldsMap(ctx, map[string]interface{}{
			"datasource_ids": lst[i].DatasourceIds,
			"annotations":    lst[i].Annotations,
			"rule_config":    lst[i].RuleConfig,
			"prod":           lst[i].Prod,
			"cate":           lst[i].Cate,
		})
		if err != nil {
			logger.Errorf("update alert rule:%d datasource ids failed, %v", lst[i].Id, err)
		}

	}
	return nil
}

// TimeSpanMuteStrategy 根据规则配置的告警生效时间段过滤,如果产生的告警不在规则配置的告警生效时间段内,则不告警,即被mute
// 时间范围，左闭右开，默认范围：00:00-24:00
func (ar *AlertRule) TimeSpanMuteStrategy() bool {
	tm := time.Unix(time.Now().Unix(), 0)
	triggerTime := tm.Format("15:04")
	triggerWeek := strconv.Itoa(int(tm.Weekday()))

	enableStime := strings.Fields(ar.EnableStime)
	enableEtime := strings.Fields(ar.EnableEtime)
	enableDaysOfWeek := strings.Split(ar.EnableDaysOfWeek, ";")
	length := len(enableDaysOfWeek)
	// enableStime,enableEtime,enableDaysOfWeek三者长度肯定相同，这里循环一个即可
	for i := 0; i < length; i++ {
		enableDaysOfWeek[i] = strings.Replace(enableDaysOfWeek[i], "7", "0", 1)
		if !strings.Contains(enableDaysOfWeek[i], triggerWeek) {
			continue
		}

		if enableStime[i] < enableEtime[i] {
			if enableEtime[i] == "23:59" {
				// 02:00-23:59，这种情况做个特殊处理，相当于左闭右闭区间了
				if triggerTime < enableStime[i] {
					// mute, 即没生效
					continue
				}
			} else {
				// 02:00-04:00 或者 02:00-24:00
				if triggerTime < enableStime[i] || triggerTime >= enableEtime[i] {
					// mute, 即没生效
					continue
				}
			}
		} else if enableStime[i] > enableEtime[i] {
			// 21:00-09:00
			if triggerTime < enableStime[i] && triggerTime >= enableEtime[i] {
				// mute, 即没生效
				continue
			}
		}

		// 到这里说明当前时刻在告警规则的某组生效时间范围内，即没有 mute，直接返回 false
		return false
	}

	return true
}

// 根据asset_id批量删除告警规则（事务）
func AlertRuleDelTxByMonId(tx *gorm.DB, ids []string) error {
	for _, id := range ids {
		str := `"monitor_id":` + id
		err := tx.Where("rule_config_fe like ?", str).Delete(&AlertRule{}).Error
		if err != nil {
			tx.Rollback()
		}
	}

	return nil
}

func AlertRuleDelTx(tx *gorm.DB, alertRule AlertRule) error {
	err := tx.Delete(&alertRule).Error
	if err != nil {
		tx.Rollback()
	}

	return nil
}

func AlertRuleGetsFilter(ctx *ctx.Context, groupId int64, filter, query string, ids []int64, limit, offset int) ([]AlertRule, error) {
	session := DB(ctx).Where("group_id=?", groupId)
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("update_at DESC")
	}
	query = "%" + query + "%"
	if filter == "1" {
		session = session.Where("severity like ?", query)
	} else if filter == "2" {
		session = session.Where("asset_id in ?", ids)
	} else if filter == "3" {
		session = session.Where("asset_id like ?", query)
	} else if filter == "" {
		session = session.Where("id like ? or name like ? or asset_id like ? or severity like ? or asset_id in ?", query, query, query, query, ids)
	}

	var lst []AlertRule
	err := session.Model(&AlertRule{}).Find(&lst).Error
	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].DB2FE()
		}
	}

	return lst, err
}

func AlertRuleGetsTotal(ctx *ctx.Context, groupId int64, filter, query string, ids []int64) (int64, error) {
	session := DB(ctx).Where("group_id=?", groupId)
	query = "%" + query + "%"
	if filter == "1" {
		session = session.Where("severity like ?", query)
	} else if filter == "2" {
		session = session.Where("asset_id in ?", ids)
	} else if filter == "3" {
		session = session.Where("asset_id like ?", query)
	} else if filter == "" {
		session = session.Where("id like ? or name like ? or asset_id like ? or severity like ? or asset_id in ?", query, query, query, query, ids)
	}

	var num int64
	err := session.Model(&AlertRule{}).Count(&num).Error

	return num, err
}

func AlertRuleGetsFilterNew(ctx *ctx.Context, groupId int64, filter, query string, ids []int64, limit, offset int) ([]AlertRule, error) {
	session := DB(ctx)
	if groupId != -1 {
		session = session.Where("group_id=?", groupId)
	}
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("update_at DESC")
	}
	query = "%" + query + "%"
	if filter == "severity" {
		session = session.Where("severity like ?", query)
	} else if filter == "ip" || filter == "name" || filter == "type" {
		session = session.Where("asset_id in ?", ids)
	} else if filter == "rule_name" {
		session = session.Where("name like ?", query)
	} else if filter == "alert_rule" {
		session = session.Where("rule_config_cn like ?", query)
	} else if filter == "asset_id" {
		session = session.Where("asset_id like ?", query)
	}

	var lst []AlertRule
	err := session.Model(&AlertRule{}).Find(&lst).Error
	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].DB2FE()
		}
	}

	return lst, err
}

func AlertRuleGetsTotalNew(ctx *ctx.Context, groupId int64, filter, query string, ids []int64) (int64, error) {
	session := DB(ctx)
	if groupId != -1 {
		session = session.Where("group_id=?", groupId)
	}
	query = "%" + query + "%"
	if filter == "severity" {
		session = session.Where("severity like ?", query)
	} else if filter == "ip" || filter == "name" || filter == "type" {
		session = session.Where("asset_id in ?", ids)
	} else if filter == "rule_name" {
		session = session.Where("name like ?", query)
	} else if filter == "alert_rule" {
		session = session.Where("rule_config_cn like ?", query)
	} else if filter == "asset_id" {
		session = session.Where("asset_id like ?", query)
	}

	var num int64
	err := session.Model(&AlertRule{}).Count(&num).Error

	return num, err
}

// 根据规则名称模糊匹配
func AlertRuleIdByName(ctx *ctx.Context, query string) (ids []int64, err error) {
	query = "%" + query + "%"
	err = DB(ctx).Model(&AlertRule{}).Distinct().Where("name like ?", query).Pluck("id", &ids).Error
	return ids, err
}

// 根据告警规则模糊匹配
func AlertRuleIdByAlertRule(ctx *ctx.Context, query string) (ids []int64, err error) {
	query = "%" + query + "%"
	err = DB(ctx).Model(&AlertRule{}).Distinct().Where("rule_config_cn like ?", query).Pluck("id", &ids).Error
	return ids, err
}

func BuildAlertRule(ctx *ctx.Context, monitoring Monitoring, alertRule AlertRuleSimplify) (*AlertRule, error) {

	lst, err := AssetsGetsMap(ctx, map[string]interface{}{"id": monitoring.AssetId})
	if err != nil {
		return &AlertRule{}, err
	}
	if len(lst) == 0 {
		return &AlertRule{}, nil
	}
	rule := &AlertRule{
		AssetId:              monitoring.AssetId,
		GroupId:              lst[0].GroupId,
		Cate:                 "prometheus",
		DatasourceIds:        "[1]",
		Name:                 alertRule.RuleConfigCn,
		Prod:                 "metric",
		PromForDuration:      60,
		PromEvalInterval:     30,
		EnableStimeJSON:      "00:00",
		EnableEtimeJSON:      "00:00",
		EnableDaysOfWeekJSON: []string{"0", "1", "2", "3", "4", "5", "6"},
		NotifyRecovered:      1,
		RuleConfigCn:         alertRule.RuleConfigCn,
		NotifyRepeatStep:     60,
		Severity:             alertRule.Severity,
		MonitoringId:         monitoring.Id,
	}
	// config := map[string][]map[string]interface{}{}
	config := make(map[string]interface{})
	config["queries"] = make([]interface{}, 0)
	configData := make(map[string]interface{})
	configData["promql"] = monitoring.MonitoringSql
	severityN, _ := strconv.ParseFloat(strconv.Itoa(alertRule.Severity), 64)
	configData["severity"] = severityN
	monId, _ := strconv.ParseFloat(strconv.FormatInt(monitoring.Id, 10), 64)
	configData["monitor_id"] = monId
	configData["relation"] = alertRule.Relation
	configData["value"] = fmt.Sprintf("%d", alertRule.Value)
	config["queries"] = append(config["queries"].([]interface{}), configData)
	// configJson, err := json.Marshal(config){
	// 	{
	// 		"prom_ql":    monitoring.MonitoringSql,
	// 		"severity":   alertRule.Severity,
	// 		"monitor_id": monitoring.Id,
	// 		"relation":   alertRule.Relation,
	// 		"value":      fmt.Sprintf("%d", alertRule.Value),
	// 	},
	// }
	configJson, err := json.Marshal(config)
	if err != nil {
		return &AlertRule{}, nil
	}
	rule.RuleConfigJson = interface{}(config)
	rule.RuleConfigFe = string(configJson)
	// if err := rule.Add(ctx); err != nil {
	// 	return &AlertRule{}, nil
	// }
	return rule, nil
}

func AlertRuleGetByMap(ctx *ctx.Context, where map[string]interface{}) ([]AlertRule, error) {
	var lst []AlertRule
	err := DB(ctx).Where(where).Find(&lst).Error
	return lst, err
}
