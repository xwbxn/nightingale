package health

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/prom"
	"github.com/ccfos/nightingale/v6/pushgw/writer"

	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/str"
)

type HealthRuleContext struct {
	datasourceId int64
	quit         chan struct{}

	assetType *models.AssetType
	// writers     *writer.WritersType
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
		//writers:      writers,
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
	logger.Infof("eval:%s started", hrc.Key())
	// interval := hrc.rule.PromEvalInterval
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
				hrc.HealthCheck()
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}
	}()
}

func (hrc *HealthRuleContext) HealthCheck() {
	if hrc.assetType.Metrics == nil || len(hrc.assetType.Metrics) == 0 {
		logger.Errorf("eval:%s metrics is nil or empty", hrc.Key())
		return
	}

	promql := fmt.Sprintf("%s", strings.Join(hrc.assetType.Metrics, " or "))

	if hrc.promClients.IsNil(hrc.datasourceId) {
		logger.Errorf("eval:%s reader client is nil", hrc.Key())
		return
	}

	value, warnings, err := hrc.promClients.GetCli(hrc.datasourceId).Query(context.Background(), promql, time.Now())
	if err != nil {
		logger.Errorf("eval:%s promql:%s, error:%v", hrc.Key(), promql, err)
		return
	}

	if len(warnings) > 0 {
		logger.Errorf("eval:%s promql:%s, warnings:%v", hrc.Key(), promql, warnings)
		return
	}

	assetsOfType := hrc.assetCache.GetByType(hrc.assetType.Name)
	ts := ConvertMetricTimeSeries(value, hrc.assetType, assetsOfType)
	if len(ts) != 0 {
		hrc.promClients.GetWriterCli(hrc.datasourceId).Write(ts)
	}
}

func (hrc *HealthRuleContext) Stop() {
	logger.Infof("%s stopped", hrc.Key())
	close(hrc.quit)
}
