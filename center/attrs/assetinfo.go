package attrs

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	prom_tool "github.com/ccfos/nightingale/v6/pkg/prom"
	"github.com/ccfos/nightingale/v6/prom"
	"github.com/google/uuid"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/toolkits/pkg/logger"
)

var DEFAULT_CLIENT int64 = 1
var attrSync Attr
var HOST string = "主机"

type Attr struct {
	Ctx        *ctx.Context
	Client     *prom.PromClientMap
	AssetCache *memsto.AssetCacheType
}

func StartAttrSync(ctx *ctx.Context, client *prom.PromClientMap, ac *memsto.AssetCacheType) {
	attrSync = Attr{
		Ctx:        ctx,
		Client:     client,
		AssetCache: ac,
	}
	attrSync.Init()
}

func (as *Attr) Init() {
	as.Client.GetCli(DEFAULT_CLIENT)
	go as.loopSync(context.Background())
}

func (as *Attr) loopSync(ctx context.Context) {
	duration := 20 * time.Second
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			as.syncAttrs()
		}
	}
}

func (as *Attr) syncAttrs() {
	logger.Debug("start sync asset attrs")
	assets := as.AssetCache.GetAll()
	for _, asset := range assets {
		if asset.Type == HOST {
			as.updateCpuAttrs(asset)
			as.updateMemAttrs(asset)
			as.updateBoardAttrs(asset)
			as.updateBIOSAttrs(asset)
		}
	}
}

func (as *Attr) updateCpuAttrs(asset *models.Asset) {
	cate := "cpu"
	client := as.Client.GetCli(DEFAULT_CLIENT)
	promql := prom_tool.InjectLabel("w_aviation_cpu_percent", "asset_id", strconv.Itoa(int(asset.Id)), labels.MatchEqual)
	logger.Debugf("update attr promql: %s", promql)
	value, _, _ := client.Query(context.Background(), promql, time.Now())
	items, ok := value.(model.Vector)
	if !ok || items.Len() == 0 {
		return
	}

	tx := as.Ctx.DB.Begin()
	err := models.AssetsExpansionDelMap(tx, map[string]interface{}{"assets_id": asset.Id, "config_category": cate})
	if err != nil {
		logger.Error(err)
		return
	}
	for _, v := range items {
		uid := uuid.New().String()
		var attrs []models.AssetsExpansion
		for labelName, labelValue := range v.Metric {
			if labelName == "cpu_model" {
				attr := createAttr(asset.Id, cate, "型号", "model", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if labelName == "cpu_arch" {
				attr := createAttr(asset.Id, cate, "架构", "arch", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if labelName == "cpu_cores" {
				attr := createAttr(asset.Id, cate, "物理核数量", "core_count", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if labelName == "cpu_threads" {
				attr := createAttr(asset.Id, cate, "逻辑核数量", "thread_count", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if labelName == "cpu_Mhz" {
				attr := createAttr(asset.Id, cate, "最大主频(Mhz)", "max_frequency", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
		}

		err = models.AssetsExpansionAddTx(tx, attrs)
		if err != nil {
			logger.Error(err)
			continue
		}

	}
	tx.Commit()
}

func (as *Attr) updateMemAttrs(asset *models.Asset) {
	cate := "memory"
	client := as.Client.GetCli(DEFAULT_CLIENT)
	promql := prom_tool.InjectLabel("w_aviation_mem_usedpercent", "asset_id", strconv.Itoa(int(asset.Id)), labels.MatchEqual)
	logger.Debugf("update attr promql: %s", promql)
	value, _, _ := client.Query(context.Background(), promql, time.Now())
	items, ok := value.(model.Vector)
	if !ok || items.Len() == 0 {
		return
	}

	tx := as.Ctx.DB.Begin()
	err := models.AssetsExpansionDelMap(tx, map[string]interface{}{"assets_id": asset.Id, "config_category": cate})
	if err != nil {
		logger.Error(err)
		return
	}
	for _, v := range items {
		uid := uuid.New().String()
		var attrs []models.AssetsExpansion
		for labelName, labelValue := range v.Metric {
			if labelName == "mem_mf" {
				attr := createAttr(asset.Id, cate, "品牌", "brand", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if labelName == "mem_speed" {
				attr := createAttr(asset.Id, cate, "主频", "frequency", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if strings.HasPrefix(string(labelName), "mem_total") {
				attr := createAttr(asset.Id, cate, "容量", "capacity", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if labelName == "mem_type" {
				attr := createAttr(asset.Id, cate, "类型", "type", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
		}

		err = models.AssetsExpansionAddTx(tx, attrs)
		if err != nil {
			logger.Error(err)
			continue
		}

	}
	tx.Commit()
}

func (as *Attr) updateBoardAttrs(asset *models.Asset) {
	cate := "board"
	client := as.Client.GetCli(DEFAULT_CLIENT)
	promql := prom_tool.InjectLabel("w_aviation_BaseBoard", "asset_id", strconv.Itoa(int(asset.Id)), labels.MatchEqual)
	logger.Debugf("update attr promql: %s", promql)
	value, _, _ := client.Query(context.Background(), promql, time.Now())
	items, ok := value.(model.Vector)
	if !ok || items.Len() == 0 {
		return
	}

	tx := as.Ctx.DB.Begin()
	err := models.AssetsExpansionDelMap(tx, map[string]interface{}{"assets_id": asset.Id, "config_category": cate})
	if err != nil {
		logger.Error(err)
		return
	}
	for _, v := range items {
		uid := uuid.New().String()
		var attrs []models.AssetsExpansion
		for labelName, labelValue := range v.Metric {
			if labelName == "Manufacturer" {
				attr := createAttr(asset.Id, cate, "厂商", "manufacturers", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if labelName == "SerialNumber" {
				attr := createAttr(asset.Id, cate, "序列号", "serial_num", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if labelName == "Product" {
				attr := createAttr(asset.Id, cate, "版本", "version", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
		}

		err = models.AssetsExpansionAddTx(tx, attrs)
		if err != nil {
			logger.Error(err)
			continue
		}

	}
	tx.Commit()
}

func (as *Attr) updateBIOSAttrs(asset *models.Asset) {
	cate := "bios"
	client := as.Client.GetCli(DEFAULT_CLIENT)
	promql := prom_tool.InjectLabel("w_aviation_BIOS", "asset_id", strconv.Itoa(int(asset.Id)), labels.MatchEqual)
	logger.Debugf("update attr promql: %s", promql)
	value, _, _ := client.Query(context.Background(), promql, time.Now())
	items, ok := value.(model.Vector)
	if !ok || items.Len() == 0 {
		return
	}

	tx := as.Ctx.DB.Begin()
	err := models.AssetsExpansionDelMap(tx, map[string]interface{}{"assets_id": asset.Id, "config_category": cate})
	if err != nil {
		logger.Error(err)
		return
	}
	for _, v := range items {
		uid := uuid.New().String()
		var attrs []models.AssetsExpansion
		for labelName, labelValue := range v.Metric {
			if labelName == "Manufacturer" {
				attr := createAttr(asset.Id, cate, "厂商", "manufacturers", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if labelName == "SMBIOSBIOSVersion" {
				attr := createAttr(asset.Id, cate, "版本", "version", uid, string(labelValue))
				attrs = append(attrs, attr)
			}
			if labelName == "ReleaseDate" {
				attr := createAttr(asset.Id, cate, "发行日期", "release_date", uid, string(labelValue)[0:8])
				attrs = append(attrs, attr)
			}
		}

		err = models.AssetsExpansionAddTx(tx, attrs)
		if err != nil {
			logger.Error(err)
			continue
		}

	}
	tx.Commit()
}

func createAttr(assetId int64, cc string, namecn string, name string, uid string, value string) models.AssetsExpansion {
	attr := models.AssetsExpansion{
		AssetsId:       assetId,
		ConfigCategory: cc,
		GroupId:        uid,
		NameCn:         namecn,
		Name:           name,
		Value:          value,
		CreatedBy:      "system",
	}
	return attr
}
