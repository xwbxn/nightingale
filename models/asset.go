package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/ginx"
)

type Asset struct {
	Id                 int64      `json:"id" gorm:"primaryKey"`
	Ident              string     `json:"ident"`
	GroupId            int64      `json:"group_id"`
	Name               string     `json:"name"`
	Label              string     `json:"label"`
	Tags               string     `json:"-"`
	TagsJSON           []string   `json:"tags" gorm:"-"`
	Type               string     `json:"type"`
	Memo               string     `json:"memo"`
	Configs            string     `json:"configs"`
	Params             string     `json:"params"`
	Plugin             string     `json:"plugin"`
	Status             int64      `json:"status"` //0: 未生效, 1: 已生效
	CreateAt           int64      `json:"create_at"`
	CreateBy           string     `json:"create_by"`
	UpdateAt           int64      `json:"update_at"`
	UpdateBy           string     `json:"update_by"`
	OrganizeId         int64      `json:"organize_id"`
	OptionalMetrics    string     `json:"-"`
	OptinalMetricsJson []*Metrics `json:"optional_metrics" gorm:"-"` //巡检检查使用
	Dashboard          string     `json:"dashboard" gorm:"-"`

	//下面的是健康检查使用，在memsto缓存中保存
	Health   int64                        `json:"-" gorm:"-"` //0: fail 1: ok
	HealthAt int64                        `json:"-" gorm:"-"`
	Metrics  map[string]map[string]string `json:"-" gorm:"-"`
}

type Metrics struct {
	Name    string `json:"name"`
	Metrics string `json:"metrics"`
}

type AssetType struct {
	Name            string                   `json:"name"`
	Plugin          string                   `json:"plugin"`
	Metrics         []*Metrics               `json:"metrics"`
	OptionalMetrics []string                 `json:"optional_metrics" yaml:"optional_metrics"`
	Category        string                   `json:"category"`
	Form            []map[string]interface{} `json:"form"`

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
