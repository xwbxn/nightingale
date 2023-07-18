package health

import (
	"fmt"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/models"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
)

const (
	LabelName     = "__name__"
	AssetId       = "asset_id"
	AssetInstance = "instance"
	HealthMetric  = "asset_health"
)

func ConvertMetricTimeSeries(value model.Value, atype *models.AssetType, assets []*models.Asset) (lst []*prompb.TimeSeries) {
	// 用于健康检查使用的是query接口，返回一定是Vector类型
	items, ok := value.(model.Vector)
	if !ok {
		return
	}

	labels := strings.Join(atype.Metrics, "|")

	for _, asset := range assets {
		health := 0
		var health_metrics []map[string]string
		for _, item := range items {
			if strings.Contains(labels, string(item.Metric[LabelName])) && strings.Contains(fmt.Sprintf("%d", asset.Id), string(item.Metric[AssetId])) {
				health_metrics = append(health_metrics, map[string]string{
					"name":  string(item.Metric[LabelName]),
					"value": fmt.Sprintf("%f", item.Value),
				})

				health = 1
			}
		}
		asset.Metrics = health_metrics
		ts := convertHealthSeries(asset, health)
		lst = append(lst, ts)
	}
	return
}

func convertHealthSeries(asset *models.Asset, health int) (ts *prompb.TimeSeries) {
	var healthLabels []*prompb.Label
	healthLabels = append(healthLabels, &prompb.Label{
		Name:  LabelName,
		Value: HealthMetric,
	})
	healthLabels = append(healthLabels, &prompb.Label{
		Name:  AssetId,
		Value: fmt.Sprintf("%d", asset.Id),
	})
	healthLabels = append(healthLabels, &prompb.Label{
		Name:  AssetInstance,
		Value: asset.Label,
	})

	s := prompb.Sample{}
	s.Timestamp = time.Now().UnixNano() / 1e6
	s.Value = float64(health)
	asset.Health = int64(s.Value)
	asset.HealthAt = time.Now().Unix()

	ts = &prompb.TimeSeries{
		Labels:  healthLabels,
		Samples: []prompb.Sample{s},
	}
	return
}
