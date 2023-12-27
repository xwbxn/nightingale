package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"
	"time"

	context "github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/prometheus/common/model"
	"github.com/toolkits/pkg/logger"
	"gorm.io/gorm"
)

type AssetMetric struct {
	Label     string      `json:"label"`
	Value     string      `json:"value"`
	Name      string      `json:"name"`
	PromValue model.Value `json:"-"`
}

type Asset struct {
	Id             int64                  `json:"id" gorm:"primaryKey"`
	Ident          string                 `json:"ident"`
	Name           string                 `json:"name" cn:"名称" xml:"name"`
	Type           string                 `json:"type" cn:"类型" xml:"type" source:"type=cache"`
	Ip             string                 `json:"ip" cn:"IP" xml:""`
	Manufacturers  string                 `json:"manufacturers" cn:"厂商" xml:"manufacturers"`
	Position       string                 `json:"position" cn:"资产位置" xml:"position"`
	Status         int64                  `json:"status" xml:"status" cn:"状态" validate:"omitempty,oneof=0 1" source:"type=option,value=[离线;正常]" ignore:"true"` // 纳管状态 0: 未生效, 1: 已生效
	StatusAt       int64                  `json:"status_at"`
	GroupId        int64                  `json:"group_id" cn:"业务组" source:"type=table,table=busi_group,property=id,field=name"`
	Label          string                 `json:"label"`
	Tags           string                 `json:"-"`
	TagsJSON       []string               `json:"tags" gorm:"-"`
	Memo           string                 `json:"memo"`
	Configs        string                 `json:"configs"`
	Params         string                 `json:"-"`
	ParamsJson     map[string]interface{} `json:"params" gorm:"-"`
	Plugin         string                 `json:"plugin"`
	CreateAt       int64                  `json:"create_at"`
	CreateBy       string                 `json:"create_by"`
	UpdateAt       int64                  `json:"update_at"`
	UpdateBy       string                 `json:"update_by"`
	OrganizationId int64                  `json:"organization_id"`
	DirectoryId    int64                  `json:"directory_id"`
	Dashboard      string                 `json:"dashboard" gorm:"-"`
	DeletedAt      gorm.DeletedAt         `gorm:"column:deleted_at" json:"deleted_at" swaggerignore:"true"`

	//下面的是健康检查使用，在memsto缓存中保存
	Health      int64             `json:"health" gorm:"-"` //0: fail 1: ok
	HealthAt    int64             `json:"-" gorm:"-"`
	Metrics     []*AssetMetric    `json:"metrics_list" gorm:"-"`
	Exps        []AssetsExpansion `json:"exps" gorm:"-"`
	Monitorings []*Monitoring     `json:"-" gorm:"-"`
}

type AssetImport struct {
	Name          string `json:"name" cn:"名称" xml:"name"`
	Type          string `json:"type" cn:"类型" xml:"type"`
	Ip            string `json:"ip" cn:"IP" xml:"IP"`
	Manufacturers string `json:"manufacturers" cn:"厂商" xml:"manufacturers"`
	Position      string `json:"position" cn:"资产位置" xml:"position"`
	Status        string `json:"status" xml:"status" cn:"状态"` //0: 未生效, 1: 已生效
	Group         string `json:"group" xml:"group" cn:"业务组"`
}

func (ins *Asset) TableName() string {
	return "assets"
}

func (ins *Asset) Verify() error {
	return nil
}

func (ins *Asset) Add(ctx *context.Context) error {
	if err := ins.Verify(); err != nil {
		return err
	}

	if ins.Type == "host" {
		if exists, err := Exists(DB(ctx).Where("ident = ? and plugin = 'host")); err != nil || exists {
			return errors.New("duplicate host asset")
		}
	}

	now := time.Now().Unix()
	ins.CreateAt = now
	ins.UpdateAt = now
	ins.Status = 0

	assetType, err := AssetTypeGet(ins.Type)
	if err != nil {
		return err
	}
	ins.Plugin = assetType.Plugin

	if err := Insert(ctx, ins); err != nil {
		return err
	}

	return nil
}

// 西航
func (ins *Asset) AddXH(ctx *context.Context) (int64, error) {
	if err := ins.Verify(); err != nil {
		return 0, err
	}

	now := time.Now().Unix()
	ins.CreateAt = now
	ins.UpdateAt = now
	ins.Status = 0

	// 回填plugin
	assetType, err := AssetTypeGet(ins.Type)
	if err != nil {
		return 0, err
	}
	ins.Plugin = assetType.Plugin
	// 生成默认采集配置
	err = ins.FillConfig(ctx)
	if err != nil {
		return 0, err
	}

	// 填充默认探针 TODO: 硬编码
	if assetType.Plugin != "host" {
		ins.Ident = "categraf-server"
	}

	err = ctx.Transaction(func(ctx *context.Context) error {
		if err := ins.Add(ctx); err != nil {
			return err
		}
		// 将asset.yaml中的默认指标入库
		for _, m := range assetType.Metrics {
			monitoring := &Monitoring{
				AssetId:        ins.Id,
				MonitoringName: m.Name,
				MonitoringSql:  m.Metrics,
				DatasourceId:   1, // TODO: 硬编码 暂定默认为1
				CreatedAt:      now,
				UpdatedAt:      now,
				Status:         1, // 默认启用
				Label:          m.Label,
				Unit:           m.Unit,
			}
			if err := monitoring.Add(ctx); err != nil {
				logger.Error("新增默认监控错误:", err)
				return err
			}

			for _, r := range m.Rules {
				rule := AlertRule{
					AssetId:          ins.Id,
					GroupId:          ins.GroupId,
					Cate:             "prometheus",
					DatasourceIds:    "[1]",
					Name:             r.Name,
					Prod:             "metric",
					PromForDuration:  60,
					PromEvalInterval: 30,
					EnableStime:      "00:00",
					EnableEtime:      "00:00",
					EnableDaysOfWeek: "0 1 2 3 4 5 6",
					NotifyRecovered:  1,
					RuleConfigCn:     m.Name,
					NotifyRepeatStep: 60,
					Severity:         int(r.Serverity),
					MonitoringId:     monitoring.Id,
				}
				config := map[string][]map[string]interface{}{}
				config["queries"] = []map[string]interface{}{
					{
						"prom_ql":    fmt.Sprintf("%s %s %d", m.Metrics, r.Relation, r.Value),
						"severity":   r.Serverity,
						"monitor_id": monitoring.Id,
						"relation":   r.Relation,
						"value":      fmt.Sprintf("%d", r.Value),
					},
				}
				configJson, err := json.Marshal(config)
				if err != nil {
					logger.Error("新增默认告警错误:", err)
					return err
				}
				rule.RuleConfig = string(configJson)
				config["queries"][0]["prom_ql"] = m.Metrics
				configJson, err = json.Marshal(config)
				if err != nil {
					logger.Error("新增默认告警错误:", err)
					return err
				}
				rule.RuleConfigFe = string(configJson)
				if err := rule.Add(ctx); err != nil {
					logger.Error("新增默认告警错误:", err)
					return err
				}
			}
		}

		// 为新增资产增加默认看板
		if err := ins.CreateBoard(ctx); err != nil {
			return err
		}
		return nil
	})

	return ins.Id, err
}

func (ins *Asset) CreateBoard(tx *context.Context) error {
	board := &Board{
		GroupId:  ins.GroupId,
		Name:     fmt.Sprintf("[%s]%s", ins.Ip, ins.Name),
		Ident:    "",
		Tags:     "",
		AssetId:  ins.Id,
		CreateBy: ins.CreateBy,
		UpdateBy: ins.UpdateBy,
	}

	err := board.Add(tx)
	if err != nil {
		return err
	}

	payload, err := BoardPayloadGetByAssetType(tx, ins.Type)
	if err == nil {
		BoardPayloadSave(tx, board.Id, payload)
	}
	return err
}

func (ins *Asset) Del(ctx *context.Context) error {

	err := ctx.Transaction(func(ctx *context.Context) error {
		if err := DB(ctx).Where("id=?", ins.Id).Delete(&Asset{}).Error; err != nil {
			return err
		}
		// 扩展属性
		if err := DB(ctx).Where("assets_id = ?", ins.Id).Delete(&AssetsExpansion{}).Error; err != nil {
			return err
		}
		// 监控
		if err := DB(ctx).Where("asset_id = ?", ins.Id).Delete(&Monitoring{}).Error; err != nil {
			return err
		}
		// 告警
		if err := DB(ctx).Where("asset_id = ?", ins.Id).Delete(&AlertRule{}).Error; err != nil {
			return err
		}
		// 看板
		var bids []int64
		DB(ctx).Model(&Board{}).Where("asset_id = ?", ins.Id).Pluck("id", &bids)
		if err := DB(ctx).Where("id in ?", bids).Delete(&Board{}).Error; err != nil {
			return err
		}
		if err := DB(ctx).Where("id in ?", bids).Delete(&BoardPayload{}).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func (ins *Asset) Update(ctx *context.Context, selectField interface{}, selectFields ...interface{}) error {
	if err := ins.Verify(); err != nil {
		return err
	}

	// 生成默认采集配置
	err := ins.FillConfig(ctx)
	if err != nil {
		return err
	}

	if err := DB(ctx).Model(ins).Select(selectField, selectFields...).Updates(ins).Error; err != nil {
		return err
	}

	return nil
}

func (ins *Asset) FillConfig(ctx *context.Context) error {
	assetType, err := AssetTypeGet(ins.Type)
	if err != nil {
		return err
	}
	ins.DB2FE()

	conf, err := assetType.GenConfig(ins.ParamsJson)
	if err != nil {
		return err
	}
	ins.Configs = conf
	return nil
}

// 通过ids修改属性
func UpdateByIds(ctx *context.Context, ids []int64, name string, value interface{}) error {
	return DB(ctx).Model(&Asset{}).Where("id in ?", ids).Updates(map[string]interface{}{name: value}).Error
}

func AssetGetById(ctx *context.Context, id int64) (*Asset, error) {
	return AssetGet(ctx, "id = ?", id)
}

func AssetGet(ctx *context.Context, where string, args ...interface{}) (*Asset, error) {
	var lst []*Asset
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

func AssetGets(ctx *context.Context, bgid int64, query string, organizationId int64) ([]*Asset, error) {
	var lst []*Asset
	session := DB(ctx).Where("1 = 1")
	if bgid >= 0 {
		session = session.Where("group_id = ?", bgid)
	}
	if organizationId >= 0 {
		session = session.Where("organization_id = ?", organizationId)
	}
	if query != "" {
		arr := strings.Fields(query)
		for i := 0; i < len(arr); i++ {
			qarg := "%" + arr[i] + "%"
			session = session.Where("name like ? or label like ? or tags like ?", qarg, qarg, qarg)
		}
	}

	err := session.Find(&lst).Error
	if err != nil {
		return nil, err
	}

	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].TagsJSON = strings.Fields(lst[i].Tags)
			lst[i].DB2FE()
		}
	}

	return lst, nil
}

func AssetGetsAll(ctx *context.Context) ([]*Asset, error) {
	return AssetGets(ctx, -1, "", -1)
}

func AssetCount(ctx *context.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&Asset{}).Where(where, args...))
}

func AssetStatistics(ctx *context.Context) (*Statistics, error) {
	session := DB(ctx).Model(&Asset{}).Select("count(*) as total", "max(update_at) as last_updated")

	var stats []*Statistics
	err := session.Find(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats[0], nil
}

func AssetGetTags(ctx *context.Context, ids []string) ([]string, error) {
	session := DB(ctx).Model(new(Asset))

	var arr []string
	if len(ids) > 0 {
		session = session.Where("id in ?", ids)
	}

	err := session.Select("distinct(tags) as tags").Pluck("tags", &arr).Error
	if err != nil {
		return nil, err
	}

	cnt := len(arr)
	if cnt == 0 {
		return []string{}, nil
	}

	set := make(map[string]struct{})
	for i := 0; i < cnt; i++ {
		tags := strings.Fields(arr[i])
		for j := 0; j < len(tags); j++ {
			set[tags[j]] = struct{}{}
		}
	}

	cnt = len(set)
	ret := make([]string, 0, cnt)
	for key := range set {
		ret = append(ret, key)
	}

	sort.Strings(ret)

	return ret, err
}

func (t *Asset) AddTags(ctx *context.Context, tags []string) error {
	for i := 0; i < len(tags); i++ {
		if !strings.Contains(t.Tags, tags[i]+" ") {
			t.Tags += tags[i] + " "
		}
	}

	arr := strings.Fields(t.Tags)
	sort.Strings(arr)

	return DB(ctx).Model(t).Updates(map[string]interface{}{
		"tags":      strings.Join(arr, " ") + " ",
		"update_at": time.Now().Unix(),
		"status":    0,
	}).Error
}

func (t *Asset) DelTags(ctx *context.Context, tags []string) error {
	for i := 0; i < len(tags); i++ {
		t.Tags = strings.ReplaceAll(t.Tags, tags[i]+" ", "")
	}

	return DB(ctx).Model(t).Updates(map[string]interface{}{
		"tags":      t.Tags,
		"update_at": time.Now().Unix(),
		"status":    0,
	}).Error
}

func AssetUpdateBgid(ctx *context.Context, ids []string, bgid int64, clearTags bool) error {
	fields := map[string]interface{}{
		"group_id":  bgid,
		"update_at": time.Now().Unix(),
		"status":    0,
	}

	if clearTags {
		fields["tags"] = ""
	}

	return DB(ctx).Model(&Asset{}).Where("id in ?", ids).Updates(fields).Error
}

func AssetBatchDel(ctx *context.Context, ids []string) error {
	if len(ids) == 0 {
		return errors.New("ids empty")
	}

	err := ctx.Transaction(func(ctx *context.Context) error {
		if err := DB(ctx).Where("id in ?", ids).Delete(&Asset{}).Error; err != nil {
			return err
		}
		if err := DB(ctx).Where("assets_id in ?", ids).Delete(&AssetsExpansion{}).Error; err != nil {
			return err
		}
		if err := DB(ctx).Where("asset_id in ?", ids).Delete(&Monitoring{}).Error; err != nil {
			return err
		}
		if err := DB(ctx).Where("asset_id in ?", ids).Delete(&AlertRule{}).Error; err != nil {
			return err
		}

		var bids []int64
		DB(ctx).Model(&Board{}).Where("asset_id in ?", ids).Pluck("id", &bids)
		if err := DB(ctx).Where("id in ?", bids).Delete(&Board{}).Error; err != nil {
			return err
		}
		if err := DB(ctx).Where("id in ?", bids).Delete(&BoardPayload{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

func AssetUpdateNote(ctx *context.Context, ids []string, memo string) error {
	return DB(ctx).Model(&Asset{}).Where("id in ?", ids).Updates(map[string]interface{}{
		"memo":      memo,
		"update_at": time.Now().Unix(),
	}).Error
}

func AssetSetStatus(ctx *context.Context, ident string, status int64) error {
	return DB(ctx).Model(&Asset{}).Where("ident = ?", ident).Updates(map[string]interface{}{
		"status":    status,
		"status_at": time.Now().Unix(),
	}).Error
}

func AssetUpdateOrganization(ctx *context.Context, ids []string, organize_id int64) error {
	return DB(ctx).Model(&Asset{}).Where("id in ?", ids).Updates(map[string]interface{}{
		"organization_id": organize_id,
		"update_at":       time.Now().Unix(),
	}).Error
}

func (e *Asset) DB2FE() {
	if e.StatusAt > 0 && time.Now().Unix()-e.StatusAt > 60*2 { // 超过2分钟设置为失联
		e.Status = 0
	}
	json.Unmarshal([]byte(e.Params), &e.ParamsJson)
}

func (e *Asset) FE2DB() error {
	paramsBytes, err := json.Marshal(e.ParamsJson)
	if err != nil {
		return err
	}
	e.Params = string(paramsBytes)
	return nil
}

// 获取某一类资产上月（自然月）环比值
func AssetMom(ctx *context.Context, aType string) (num int64, err error) {
	//统计上月新增资产
	var insertNum int64
	now := time.Now()
	start := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location()).Unix()
	end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Unix()

	err = DB(ctx).Model(&Asset{}).Where("create_at >= ? AND create_at < ?", start, end).Where("type = ?", aType).Count(&insertNum).Error
	if err != nil {
		return 0, errors.New("查询环比数据出错")
	}
	//统计上月删除资产
	var delLst []Asset
	// err = DB(ctx).Unscoped().Find(&delLst).Error
	// err = DB(ctx).Unscoped().Where("IFNULL('deleted_at','kong')!='kong'").Find(&delLst).Error
	err = DB(ctx).Unscoped().Where("type = ?", aType).Where("`assets`.`deleted_at` IS NOT NULL").Find(&delLst).Error
	if err != nil {
		return 0, errors.New("查询环比数据出错")
	}
	for _, val := range delLst {
		delUnix := val.DeletedAt.Time.Unix()
		if delUnix >= start && delUnix < end {
			insertNum--
		}
	}
	return insertNum, err
}

// 根据map查询
func AssetsGetsMap(ctx *context.Context, where map[string]interface{}) (lst []*Asset, err error) {
	err = DB(ctx).Model(&Asset{}).Where(where).Find(&lst).Error
	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].TagsJSON = strings.Fields(lst[i].Tags)
			lst[i].DB2FE()
		}
	}
	return lst, err
}

// 根据map统计数量
func AssetsCountMap(ctx *context.Context, where map[string]interface{}) (num int64, err error) {
	err = DB(ctx).Model(&Asset{}).Where(where).Count(&num).Error
	return num, err
}

// 统计ip是否存在
func IpCount(ctx *context.Context, ip string, aTypes []string) (num int64, err error) {
	err = DB(ctx).Model(&Asset{}).Where("ip = ? and type in ?", ip, aTypes).Count(&num).Error
	return num, err
}

// 根据filter统计数量
func AssetsCountFilter(ctx *context.Context, aType string, query, queryType string) (num int64, err error) {
	session := DB(ctx)

	if queryType == "" {
		if query != "" {
			session = session.Where("name like ? or type like ? or ip like ? or asset_status like ? or os like ? or location like ?", query, query, query, query, query, query)
		}
	} else {
		session = session.Where(queryType+" like ?", query)
	}

	// if len(dirs) != 0 {
	// 	logger.Debug(dirs)
	// 	session = session.Where("directory_id in ?", dirs)
	// }
	if aType != "" {
		session = session.Where("type = ?", aType)
	}

	err = session.Model(&Asset{}).Count(&num).Error
	return num, err
}

// 根据filter查询
func AssetsGetsFilter(ctx *context.Context, aType string, query, queryType string, limit, offset int) (lst []*Asset, err error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id DESC")
	}
	if queryType == "" {
		if query != "" {
			logger.Debug(query)
			session = session.Where("name like ? or type like ? or ip like ? or asset_status like ? or os like ? or location like ?", query, query, query, query, query, query)
		}
	} else {
		session = session.Where(queryType+" like ?", query)
	}
	// logger.Debug(dirs)
	// if len(dirs) != 0 {
	// 	session = session.Where("directory_id in ?", dirs)
	// }
	if aType != "" {
		session = session.Where("type = ?", aType)
	}

	err = session.Model(&Asset{}).Find(&lst).Error

	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].TagsJSON = strings.Fields(lst[i].Tags)
			lst[i].DB2FE()
		}
	}

	return lst, err
}

// 根据filter统计数量(new)
func AssetsCountFilterNew(ctx *context.Context, filter, query, aType string) (num int64, err error) {
	session := DB(ctx)

	if aType != "" {
		session = session.Where("type like ?", "%"+aType+"%")
	}

	if filter == "ip" {
		session = session.Where("ip like ?", "%"+query+"%")
	} else if filter == "name" {
		session = session.Where("name like?", "%"+query+"%")
	} else if filter == "status" {
		session = session.Where("status like ?", "%"+query+"%")
	} else if filter == "group_id" {
		session = session.Where("group_id like ?", "%"+query+"%")
	} else if filter == "position" {
		session = session.Where("position like ?", "%"+query+"%")
	} else if filter == "manufacturers" {
		session = session.Where("manufacturers like?", "%"+query+"%")
		// } else if filter == "type" {
		// 	session = session.Where("type like?", "%"+query+"%")
	} else if filter == "os" {
		ids, err := AssetsExpansionGetAssetIdMap(ctx, map[string]interface{}{"config_category": "os", "name": "name"}, "%"+query+"%")
		if err != nil || len(ids) == 0 {
			return 0, err
		}
		session = session.Where("id in ?", ids)
	}
	err = session.Model(&Asset{}).Count(&num).Error
	return num, err
}

// 根据filter查询(new)
func AssetsGetsFilterNew(ctx *context.Context, filter, query, aType string, limit, offset int) (lst []*Asset, err error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("update_at DESC")
	}

	if aType != "" {
		session = session.Where("type like ?", "%"+aType+"%")
	}

	if filter == "ip" {
		session = session.Where("ip like ?", "%"+query+"%")
	} else if filter == "name" {
		session = session.Where("name like?", "%"+query+"%")
	} else if filter == "status" {
		session = session.Where("status like ?", "%"+query+"%")
	} else if filter == "group_id" {
		session = session.Where("group_id like ?", "%"+query+"%")
	} else if filter == "position" {
		session = session.Where("position like ?", "%"+query+"%")
	} else if filter == "manufacturers" {
		session = session.Where("manufacturers like?", "%"+query+"%")
		// } else if filter == "type" {
		// 	session = session.Where("type like?", "%"+query+"%")
	} else if filter == "os" {
		ids, err := AssetsExpansionGetAssetIdMap(ctx, map[string]interface{}{"config_category": "os", "name": "name"}, "%"+query+"%")
		if err != nil {
			return lst, err
		}
		if len(ids) == 0 {
			return make([]*Asset, 0), nil
		}
		session = session.Where("id in ?", ids)
	}

	err = session.Model(&Asset{}).Find(&lst).Error

	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].TagsJSON = strings.Fields(lst[i].Tags)
			lst[i].DB2FE()
		}
	}

	return lst, err
}

// 根据资产名称、类型、IP地址模糊匹配
func AssetIdByNameTypeIp(ctx *context.Context, query string) (ids []int64, err error) {
	query = "%" + query + "%"
	err = DB(ctx).Model(&Asset{}).Distinct().Where("name like ? or type like ? or ip like ?", query, query, query).Pluck("id", &ids).Error
	return ids, err
}

// 根据资产名称模糊匹配
func AssetIdByName(ctx *context.Context, query string) (ids []int64, err error) {
	query = "%" + query + "%"
	err = DB(ctx).Model(&Asset{}).Distinct().Where("name like ?", query).Pluck("id", &ids).Error
	return ids, err
}

// 根据资产类型模糊匹配
func AssetIdByType(ctx *context.Context, query string) (ids []int64, err error) {
	query = "%" + query + "%"
	err = DB(ctx).Model(&Asset{}).Distinct().Where("type like ?", query).Pluck("id", &ids).Error
	return ids, err
}

// 获取ip
func GetIP(asset *Asset) (ip, name string) {
	ip = ""
	name = asset.Name
	address := net.ParseIP(asset.Label)
	if address != nil {
		ip = asset.Label
	} else {
		//将名称中含有的IP拿出来

		compileRegex := regexp.MustCompile("（(.*?)）") // 中文括号，例如：华南地区（广州） -> 广州
		matchArr := compileRegex.FindStringSubmatch(asset.Name)

		if len(matchArr) == 0 {
			compileRegex = regexp.MustCompile("\\((.*?)\\)") // 兼容英文括号并取消括号的转义，例如：华南地区 (广州) -> 广州。
			matchArr = compileRegex.FindStringSubmatch(asset.Name)
		}
		if len(matchArr) != 0 {
			ipTemp := matchArr[len(matchArr)-1]
			if len(strings.Split(ipTemp, ".")) == 4 {
				ip = ipTemp
				nameTemp := strings.Split(asset.Name, matchArr[0])
				name = nameTemp[0]
			}
		}
	}
	return ip, name
}

// 根据ids获得资产
func AssetGetByIds(ctx *context.Context, ids []int64) ([]*Asset, error) {
	var lst []*Asset
	err := DB(ctx).Model(&Asset{}).Where("id in ?", ids).Find(&lst).Error

	if err == nil {
		for i := 0; i < len(lst); i++ {
			lst[i].TagsJSON = strings.Fields(lst[i].Tags)
			lst[i].DB2FE()
		}
	}

	return lst, err
}

// 数据校验
func Verify(asset Asset) error {
	//非空字段校验
	if asset.Name == "" {
		return errors.New("名称不能为空")
	} else if asset.Ip == "" {
		return errors.New("IP不能为空")
	} else if asset.Type == "" {
		return errors.New("类型不能为空")
	} else if asset.GroupId == 0 {
		return errors.New("业务组不能为空")
	}
	//IP格式校验
	err := net.ParseIP(asset.Ip)
	if err == nil {
		return errors.New("IP格式不正确")
	}
	return nil
}
