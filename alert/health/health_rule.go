package health

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/models"
	pm "github.com/ccfos/nightingale/v6/pkg/prom"
	"github.com/ccfos/nightingale/v6/prom"
	"github.com/ccfos/nightingale/v6/pushgw/writer"
	"github.com/prometheus/prometheus/model/labels"

	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/str"
)

type HealthRuleContext struct {
	datasourceId int64
	quit         chan struct{}

	assetType   *models.AssetType
	writers     *writer.WritersType
	promClients *prom.PromClientMap
	assetCache  *memsto.AssetCacheType
}

func NewHealthRuleContext(atype *models.AssetType, datasourceId int64, promClients *prom.PromClientMap, writers *writer.WritersType, cache *memsto.AssetCacheType) *HealthRuleContext {
	return &HealthRuleContext{
		datasourceId: datasourceId,
		quit:         make(chan struct{}),
		assetType:    atype,
		promClients:  promClients,
		assetCache:   cache,
		writers:      writers,
	}
}

func (hrc *HealthRuleContext) Key() string {
	return fmt.Sprintf("health-%d-%s", hrc.datasourceId, hrc.assetType.Name)
}

func (hrc *HealthRuleContext) Hash() string {
	return str.MD5(fmt.Sprintf("%s_%s_%d",
		hrc.assetType.Name,
		hrc.assetType.Category,
		hrc.datasourceId,
	))
}

func (hrc *HealthRuleContext) Prepare() {}

func (hrc *HealthRuleContext) Start() {
	logger.Infof("health_check:%s started", hrc.Key())
	// interval := hrc.rule.PromHealth_checkInterval
	interval := 15
	if interval <= 0 {
		interval = 10
	}
	go func() {
		for {
			select {
			case <-hrc.quit:
				return
			default:
				hrc.AssetMetricsCheck()
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}
	}()
}

func (hrc *HealthRuleContext) AssetMetricsCheck() {
	assetsOfType := hrc.assetCache.GetByType(hrc.assetType.Name)
	for _, asset := range assetsOfType {
		metrics := []*models.Metrics{}
		if hrc.assetType.Metrics != nil {
			metrics = append(metrics, hrc.assetType.Metrics...)
		}
		if len(metrics) == 0 {
			logger.Warningf("health: asset %d metrics is nil or empty", asset.Id)
			continue
		}

		for _, m := range metrics {
			pm.InjectLabel(m.Metrics, "asset_id", strconv.Itoa(int(asset.Id)), labels.MatchEqual)
			value, warnings, err := hrc.promClients.GetCli(hrc.datasourceId).Query(context.Background(), m.Metrics, time.Now())
			if err != nil {
				logger.Errorf("health:%s promql:%s, error:%v", hrc.Key(), m.Metrics, err)
				return
			}

			if len(warnings) > 0 {
				logger.Errorf("health:%s promql:%s, warnings:%v", hrc.Key(), m.Metrics, warnings)
				return
			}
			ts := ConvertMetricTimeSeries(value, m, asset)
			for _, series := range ts {
				hrc.writers.PushSample("health_check", series)
			}
		}
	}
}

func (hrc *HealthRuleContext) Stop() {
	logger.Infof("%s stopped", hrc.Key())
	close(hrc.quit)
}
