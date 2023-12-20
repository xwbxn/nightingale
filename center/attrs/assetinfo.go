package attrs

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	prom_tool "github.com/ccfos/nightingale/v6/pkg/prom"
	"github.com/ccfos/nightingale/v6/prom"
	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/toolkits/pkg/logger"
	"golang.org/x/exp/slices"
)

var DEFAULT_CLIENT int64 = 1
var attrSync Attr
var HOST string = "主机设备"

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
	logger.Debug("开始同步资产扩展属性")
	assets := as.AssetCache.GetAll()
	for _, asset := range assets {
		if asset.Type == HOST {
			as.updateExtraProps(asset)
		}
	}
}

func (as *Attr) updateExtraProps(asset *models.Asset) {
	atype, ok := as.AssetCache.GetType(asset.Type)
	if !ok {
		logger.Error("资产类型不存在: ", asset.Type)
		return
	}

	client := as.Client.GetCli(DEFAULT_CLIENT)
	for cate, prop := range atype.ExtraProps {
		if len(prop.Props) == 0 {
			logger.Warning("资产类型无扩展属性: ", asset.Type, cate)
			continue
		}

		promql := fmt.Sprintf("w_aviation_%s", cate) // 这里的基础属性使用了w_aviation插件，采集指标的后缀名称与assets.yaml中一致
		promql, _ = prom_tool.InjectLabel(promql, "asset_id", strconv.Itoa(int(asset.Id)), labels.MatchEqual)

		value, warning, err := client.Query(context.Background(), promql, time.Now())
		if len(warning) > 0 {
			logger.Errorf("查询资产错误: %s, %s", err.Error(), promql)
			continue
		}
		values, ok := value.(model.Vector)
		if !ok {
			logger.Errorf("查询资产错误: %s, %s", values, promql)
			continue
		}
		if len(values) == 0 {
			continue
		}

		var attrs []models.AssetsExpansion
		for _, item := range values {
			uid := uuid.New().String()
			for _, p := range prop.Props[0].Items {
				l := item.Metric
				v, exists := l[model.LabelName(p.Name)]

				var safeVal string = string(v)
				switch p.Unit {
				case models.Byte:
					if val, err := strconv.Atoi(string(v)); err == nil {
						safeVal = humanize.Bytes(uint64(val))
					}
				case models.Bit:
					if val, err := strconv.Atoi(string(v)); err == nil {
						safeVal = humanize.SI(float64(val), "b")
					}
				case models.Time:

				default:
					safeVal = strings.ReplaceAll(string(v), "\\", "/")
				}

				if exists {
					attr := createAttr(asset.Id, cate, p.Label, p.Name, uid, safeVal)
					attrs = append(attrs, attr)
				}
			}
		}

		query := map[string]interface{}{"assets_id": asset.Id, "config_category": cate}
		oldAttrs, err := models.AssetsExpansionGetsMap(as.Ctx, query)
		if err != nil {
			logger.Errorf("资产查询扩展属性错误: %s", err.Error())
		}

		notChange := slices.EqualFunc[models.AssetsExpansion](attrs, oldAttrs, func(s1, s2 models.AssetsExpansion) bool {
			return s1.Hash() == s2.Hash()
		})

		if !notChange {
			tx := as.Ctx.DB.Begin()
			err = models.AssetsExpansionDelMap(tx, map[string]interface{}{"assets_id": asset.Id, "config_category": cate})
			if err != nil {
				logger.Errorf("更新资产属性错误: %s", err.Error())
				tx.Rollback()
				continue
			}
			err = models.AssetsExpansionAddTx(tx, attrs)
			if err != nil {
				logger.Errorf("更新资产属性错误: %s", err.Error())
				tx.Rollback()
				continue
			}
			tx.Commit()
		}
	}
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
