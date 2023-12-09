package health

import (
	"context"
	"fmt"
	"time"

	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/prom"
	"github.com/ccfos/nightingale/v6/pushgw/writer"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"

	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/str"
)

const (
	LabelName     = "__name__"
	AssetId       = "asset_id"
	AssetInstance = "instance"
	HealthMetric  = "asset_up"
)

type HealthRuleContext struct {
	datasourceId int64
	quit         chan struct{}

	asset       *models.Asset
	writers     *writer.WritersType
	promClients *prom.PromClientMap
	ac          *memsto.AssetCacheType
}

func NewHealthRuleContext(asset *models.Asset, datasourceId int64, promClients *prom.PromClientMap, writers *writer.WritersType, ac *memsto.AssetCacheType) *HealthRuleContext {
	return &HealthRuleContext{
		datasourceId: datasourceId,
		quit:         make(chan struct{}),
		asset:        asset,
		promClients:  promClients,
		writers:      writers,
		ac:           ac,
	}
}

func (hrc *HealthRuleContext) Key() string {
	return fmt.Sprintf("health-%d-%d-%s", hrc.datasourceId, hrc.asset.Id, hrc.asset.Name)
}

func (hrc *HealthRuleContext) Hash() string {
	h := []int64{}
	for _, m := range hrc.asset.Monitorings {
		h = append(h, m.Id)
		h = append(h, m.Status)
		h = append(h, m.UpdatedAt)
	}
	return str.MD5(fmt.Sprintf("%s_%d_%d_%v",
		hrc.asset.Name,
		hrc.asset.Id,
		hrc.datasourceId,
		h,
	))
}

func (hrc *HealthRuleContext) Prepare() {}

func (hrc *HealthRuleContext) Start() {
	logger.Infof("health_check:%s started", hrc.Key())
	// interval := hrc.rule.PromHealth_checkInterval
	interval := 15
	if interval <= 0 {
		interval = 15
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

	logger.Debugf("starting check health: %s", hrc.Key())
	metrics := []*models.AssetMetric{}
	for _, m := range hrc.asset.Monitorings {
		if m.Status == 0 {
			continue
		}
		promql := m.CompilePromQL()
		value, warnings, err := hrc.promClients.GetCli(hrc.datasourceId).Query(context.Background(), promql, time.Now())
		if err != nil {
			logger.Errorf("health:%s promql:%s, error:%v", hrc.Key(), promql, err)
			return
		}
		if len(warnings) > 0 {
			logger.Errorf("health:%s promql:%s, warnings:%v", hrc.Key(), promql, warnings)
			return
		}
		metrics = append(metrics, &models.AssetMetric{
			Name:      m.MonitoringName,
			Label:     promql,
			PromValue: value,
			Value:     "0",
		})
	}

	assetInCache, has := hrc.ac.Get(hrc.asset.Id)
	if !has {
		return
	}

	assetInCache.Metrics = metrics
	health := 0
	for _, m := range assetInCache.Metrics {
		vector, ok := m.PromValue.(model.Vector)
		if ok {
			for _, resultValue := range vector {
				health = 1
				m.Value = fmt.Sprintf("%f", resultValue.Value)
				break
			}
		}
		matrix, ok := m.PromValue.(model.Matrix)
		if ok {
			for _, resultValue := range matrix {
				if len(resultValue.Values) > 0 {
					health = 1
					m.Value = fmt.Sprintf("%f", resultValue.Values[0].Value)
				}
				break
			}
		}
	}

	assetInCache.Health = int64(health)
	assetInCache.HealthAt = time.Now().Unix()
	logger.Debugf("processing check health: %s, health: %d, %p", hrc.Key(), health, hrc.asset)

	ts := convertHealthTimeSeries(hrc, health)
	hrc.writers.PushSample("health_check", ts)
}

func convertHealthTimeSeries(hrc *HealthRuleContext, health int) *prompb.TimeSeries {
	var healthLabels []prompb.Label
	healthLabels = append(healthLabels, prompb.Label{
		Name:  LabelName,
		Value: HealthMetric,
	})
	healthLabels = append(healthLabels, prompb.Label{
		Name:  AssetId,
		Value: fmt.Sprintf("%d", hrc.asset.Id),
	})
	healthLabels = append(healthLabels, prompb.Label{
		Name:  AssetInstance,
		Value: hrc.asset.Label,
	})
	s := prompb.Sample{}
	s.Timestamp = time.Now().UnixNano() / 1e6
	s.Value = float64(health)
	ts := &prompb.TimeSeries{
		Labels:  healthLabels,
		Samples: []prompb.Sample{s},
	}
	return ts
}

func (hrc *HealthRuleContext) Stop() {
	logger.Infof("%s stopped", hrc.Key())
	close(hrc.quit)
}
