package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"path"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
	"gorm.io/gorm"
)

type Asset struct {
	Id                 int64          `json:"id" gorm:"primaryKey"`
	Ident              string         `json:"ident"`
	GroupId            int64          `json:"group_id"`
	Name               string         `json:"name" cn:"名称" xml:"name"`
	Type               string         `json:"type" cn:"类型" xml:"type" source:"type=cache"`
	Ip                 string         `json:"ip" cn:"IP" xml:""`
	Manufacturers      string         `json:"manufacturers" cn:"厂商" xml:"manufacturers"`
	Position           string         `json:"position" cn:"资产位置" xml:"position"`
	Status             int64          `json:"status" xml:"status" cn:"状态" validate:"omitempty,oneof=0 1" source:"type=option,value=[下线;正常]"` //0: 未生效, 1: 已生效
	Label              string         `json:"label"`
	Tags               string         `json:"-"`
	TagsJSON           []string       `json:"tags" gorm:"-"`
	Memo               string         `json:"memo"`
	Configs            string         `json:"configs"`
	Params             string         `json:"params"`
	Plugin             string         `json:"plugin"`
	CreateAt           int64          `json:"create_at"`
	CreateBy           string         `json:"create_by"`
	UpdateAt           int64          `json:"update_at"`
	UpdateBy           string         `json:"update_by"`
	OrganizationId     int64          `json:"organization_id"`
	DirectoryId        int64          `json:"directory_id"`
	OptionalMetrics    string         `json:"-"`
	OptinalMetricsJson []*Metrics     `json:"optional_metrics" gorm:"-"` //巡检检查使用
	Dashboard          string         `json:"dashboard" gorm:"-"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at" swaggerignore:"true"`

	//下面的是健康检查使用，在memsto缓存中保存
	Health      int64                        `json:"health" gorm:"-"` //0: fail 1: ok
	HealthAt    int64                        `json:"-" gorm:"-"`
	Metrics     map[string]map[string]string `json:"-" gorm:"-"`
	Exps        []AssetsExpansion            `json:"exps" gorm:"-"`
	MetricsList []map[string]interface{}     `json:"metrics_list" gorm:"-"`
}

type AssetImport struct {
	Name          string `json:"name" cn:"名称" xml:"name"`
	Type          string `json:"type" cn:"类型" xml:"type"`
	Ip            string `json:"ip" cn:"IP" xml:"IP"`
	Manufacturers string `json:"manufacturers" cn:"厂商" xml:"manufacturers"`
	Position      string `json:"position" cn:"资产位置" xml:"position"`
	Status        string `json:"status" xml:"status" cn:"状态"` //0: 未生效, 1: 已生效

}

type Metrics struct {
	Name    string `json:"name"`
	Metrics string `json:"metrics" yaml:"metrics"`
}

type BaseProp struct {
	Name     string   `json:"name" yaml:"name"`
	Label    string   `json:"label"`
	Required string   `json:"required"`
	Type     string   `json:"type"`
	Options  []string `json:"options" yaml:"options"`
}

type ExtraProp struct {
	Cpu    *ExtraPropPart `json:"cpu" yaml:"cpu"`
	Memory *ExtraPropPart `json:"memory" yaml:"memory"`
	Disk   *ExtraPropPart `json:"disk" yaml:"disk"`
	Net    *ExtraPropPart `json:"net" yaml:"net"`
	Board  *ExtraPropPart `json:"board" yaml:"board"`
	Bios   *ExtraPropPart `json:"bios" yaml:"bios"`
	Bus    *ExtraPropPart `json:"bus" yaml:"bus"`
	Os     *ExtraPropPart `json:"os" yaml:"os"`
	Power  *ExtraPropPart `json:"power" yaml:"power"`
	Device *ExtraPropPart `json:"device" yaml:"device"`
}

type ExtraPropPart struct {
	Label string `json:"label"  yaml:"label"`
	Props []*struct {
		Name       string `json:"name" yaml:"name"`
		Label      string `json:"label" yaml:"label"`
		Type       string `json:"type" yaml:"type"`
		ItemsLimit int64  `json:"items_limit" yaml:"items_limit"`
		Items      []*struct {
			Name    string `json:"name" yaml:"name"`
			Label   string `json:"label" yaml:"label"`
			Type    string `json:"type" yaml:"type"`
			Options []*struct {
				Label string `json:"label" yaml:"label"`
				Value string `json:"value" yaml:"value"`
			} `json:"options" yaml:"options"`
		} `json:"items" yaml:"items"`
	} `json:"props" yaml:"props"`
}

type AssetType struct {
	Name            string                   `json:"name"`
	Plugin          string                   `json:"plugin"`
	Metrics         []*Metrics               `json:"metrics"`
	OptionalMetrics []string                 `json:"optional_metrics" yaml:"optional_metrics"`
	Category        string                   `json:"category"`
	Form            []map[string]interface{} `json:"form"`
	BaseProps       []*BaseProp              `json:"base_props" yaml:"base_props"`
	ExtraProps      ExtraProp                `json:"extra_props" yaml:"extra_props"`

	Dashboard string `json:"-"`
}

type AssetConfigs struct {
	Config []*AssetType `json:"config"`
}

func (ins *Asset) TableName() string {
	return "assets"
}

func (ins *Asset) Verify() error {
	return nil
}

func (ins *Asset) Add(ctx *ctx.Context) error {
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
	assetTypes, err := AssetTypeGetsAll()
	if err != nil {
		return err
	}

	for _, item := range assetTypes {
		if item.Name == ins.Type {
			ins.Plugin = item.Plugin
			break
		}
	}

	if err := Insert(ctx, ins); err != nil {
		return err
	}

	return nil
}

//西航
func (ins *Asset) AddXH(ctx *ctx.Context) (int64, error) {
	if err := ins.Verify(); err != nil {
		return 0, err
	}

	if ins.Type == "host" {
		if exists, err := Exists(DB(ctx).Where("ident = ? and plugin = 'host")); err != nil || exists {
			return 0, errors.New("duplicate host asset")
		}
	}

	now := time.Now().Unix()
	ins.CreateAt = now
	ins.UpdateAt = now
	ins.Status = 0
	assetTypes, err := AssetTypeGetsAll()
	if err != nil {
		return 0, err
	}

	for _, item := range assetTypes {
		if item.Name == ins.Type {
			ins.Plugin = item.Plugin
			break
		}
	}

	if err := Insert(ctx, ins); err != nil {
		return 0, err
	}

	return ins.Id, nil
}

func (ins *Asset) AddTx(tx *gorm.DB) (int64, error) {
	if err := ins.Verify(); err != nil {
		return 0, err
	}

	if ins.Type == "host" {
		if exists, err := Exists(tx.Where("ident = ? and plugin = 'host")); err != nil || exists {
			return 0, errors.New("duplicate host asset")
		}
	}

	now := time.Now().Unix()
	ins.CreateAt = now
	ins.UpdateAt = now
	ins.Status = 0
	assetTypes, err := AssetTypeGetsAll()
	if err != nil {
		return 0, err
	}

	for _, item := range assetTypes {
		if item.Name == ins.Type {
			ins.Plugin = item.Plugin
			break
		}
	}

	err = tx.Debug().Model(&Asset{}).Create(ins).Error
	if err != nil {
		tx.Rollback()
	}

	// if err := Insert(tx, ins); err != nil {
	// 	return err
	// }

	return ins.Id, nil
}

func (ins *Asset) Del(ctx *ctx.Context) error {
	if err := DB(ctx).Where("id=?", ins.Id).Delete(&Asset{}).Error; err != nil {
		return err
	}

	return nil
}

func (ins *Asset) Update(ctx *ctx.Context, selectField interface{}, selectFields ...interface{}) error {
	if err := ins.Verify(); err != nil {
		return err
	}

	if err := DB(ctx).Model(ins).Select(selectField, selectFields...).Updates(ins).Error; err != nil {
		return err
	}

	if err := DB(ctx).Model(ins).Select("status").Updates(map[string]interface{}{"status": 0}).Error; err != nil {
		return err
	}

	return nil
}

func (ins *Asset) UpdateTx(tx *gorm.DB, selectField interface{}, selectFields ...interface{}) error {
	if err := ins.Verify(); err != nil {
		return err
	}

	if err := tx.Model(ins).Select(selectField, selectFields...).Updates(ins).Error; err != nil {
		tx.Rollback()
	}

	if err := tx.Model(ins).Select("status").Updates(map[string]interface{}{"status": 0}).Error; err != nil {
		tx.Rollback()
	}

	return nil
}

//通过ids修改属性
func UpdateByIds(ctx *ctx.Context, ids []int64, name string, value interface{}) error {
	return DB(ctx).Model(&Asset{}).Where("id in ?", ids).Updates(map[string]interface{}{name: value}).Error
}

func AssetGetById(ctx *ctx.Context, id int64) (*Asset, error) {
	return AssetGet(ctx, "id = ?", id)
}

func AssetGet(ctx *ctx.Context, where string, args ...interface{}) (*Asset, error) {
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

func AssetGets(ctx *ctx.Context, bgid int64, query string, organizationId int64) ([]*Asset, error) {
	var lst []*Asset
	// session := DB(ctx).Where("1 = 1")
	session := DB(ctx).Find(&lst)
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

func AssetGetsAll(ctx *ctx.Context) ([]*Asset, error) {
	return AssetGets(ctx, -1, "", -1)
}

func AssetCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&Asset{}).Where(where, args...))
}

func AssetStatistics(ctx *ctx.Context) (*Statistics, error) {
	session := DB(ctx).Model(&Asset{}).Select("count(*) as total", "max(update_at) as last_updated")

	var stats []*Statistics
	err := session.Find(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats[0], nil
}

func AssetGenConfig(assetType string, f map[string]interface{}) (bytes.Buffer, error) {
	assetTypes, err := AssetTypeGetsAll()
	ginx.Dangerous(err)

	pluginName := ""
	for _, item := range assetTypes {
		if item.Name == assetType {
			pluginName = item.Plugin
			break
		}
	}
	filepath := path.Join("etc", "default", fmt.Sprintf("%s.toml", pluginName))

	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("unable to open config: %s", filepath)
		return bytes.Buffer{}, err
	}
	var content bytes.Buffer
	tpl.Execute(&content, f)
	return content, err
}

func AssetTypeGetsAll() ([]*AssetType, error) {
	fp := path.Join("etc", "assets.yaml")
	var assetConfig AssetConfigs
	err := file.ReadYaml(fp, &assetConfig)
	return assetConfig.Config, err
}

func AssetGetTags(ctx *ctx.Context, ids []string) ([]string, error) {
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

func (t *Asset) AddTags(ctx *ctx.Context, tags []string) error {
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

func (t *Asset) DelTags(ctx *ctx.Context, tags []string) error {
	for i := 0; i < len(tags); i++ {
		t.Tags = strings.ReplaceAll(t.Tags, tags[i]+" ", "")
	}

	return DB(ctx).Model(t).Updates(map[string]interface{}{
		"tags":      t.Tags,
		"update_at": time.Now().Unix(),
		"status":    0,
	}).Error
}

func AssetUpdateBgid(ctx *ctx.Context, ids []string, bgid int64, clearTags bool) error {
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

func AssetDel(ctx *ctx.Context, ids []string) error {
	if len(ids) == 0 {
		panic("ids empty")
	}
	return DB(ctx).Where("id in ?", ids).Delete(new(Asset)).Error
}

func AssetDelTx(tx *gorm.DB, ids []string) error {
	if len(ids) == 0 {
		panic("ids empty")
	}
	return tx.Where("id in ?", ids).Delete(new(Asset)).Error
}

func AssetUpdateNote(ctx *ctx.Context, ids []string, memo string) error {
	return DB(ctx).Model(&Asset{}).Where("id in ?", ids).Updates(map[string]interface{}{
		"memo":      memo,
		"update_at": time.Now().Unix(),
	}).Error
}

func AssetSetStatus(ctx *ctx.Context, ident string, status int64) error {
	return DB(ctx).Model(&Asset{}).Where("ident = ?", ident).Updates(map[string]interface{}{
		"status": status,
	}).Error
}

func AssetUpdateOrganization(ctx *ctx.Context, ids []string, organize_id int64) error {
	return DB(ctx).Model(&Asset{}).Where("id in ?", ids).Updates(map[string]interface{}{
		"organization_id": organize_id,
		"update_at":       time.Now().Unix(),
	}).Error
}

func (e *Asset) DB2FE() {
	json.Unmarshal([]byte(e.OptionalMetrics), &e.OptinalMetricsJson)
}

//获取某一类资产上月（自然月）环比值
func AssetMom(ctx *ctx.Context, aType string) (num int64, err error) {
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
	// err = DB(ctx).Debug().Unscoped().Find(&delLst).Error
	// err = DB(ctx).Debug().Unscoped().Where("IFNULL('deleted_at','kong')!='kong'").Find(&delLst).Error
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

//根据map查询
func AssetsGetsMap(ctx *ctx.Context, where map[string]interface{}) (lst []Asset, err error) {
	err = DB(ctx).Model(&Asset{}).Where(where).Find(&lst).Error
	return lst, err
}

//根据map统计数量
func AssetsCountMap(ctx *ctx.Context, where map[string]interface{}) (num int64, err error) {
	err = DB(ctx).Model(&Asset{}).Where(where).Count(&num).Error
	return num, err
}

//根据filter统计数量
func AssetsCountFilter(ctx *ctx.Context, aType string, query, queryType string) (num int64, err error) {
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

//根据filter查询
func AssetsGetsFilter(ctx *ctx.Context, aType string, query, queryType string, limit, offset int) (lst []Asset, err error) {
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
	return lst, err
}

//根据资产名称、类型、IP地址模糊匹配
func AssetIdByNameTypeIp(ctx *ctx.Context, query string) (ids []int64, err error) {
	query = "%" + query + "%"
	err = DB(ctx).Model(&Asset{}).Distinct().Where("name like ? or type like ? or ip like ?", query, query, query).Pluck("id", &ids).Error
	return ids, err
}

//获取ip
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

//根据ids获得资产
func AssetGetByIds(ctx *ctx.Context, ids []int64) ([]Asset, error) {
	var lst []Asset
	err := DB(ctx).Model(&Asset{}).Where("id in ?", ids).Find(&lst).Error
	return lst, err
}
