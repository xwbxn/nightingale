package health

import (
	"fmt"
	"time"

	"github.com/ccfos/nightingale/v6/models"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
)

const (
	LabelName     = "__name__"
	AssetId       = "asset_id"
	AssetInstance = "instance"
	HealthMetric  = "asset_up"
)

func ConvertMetricTimeSeries(value model.Value, mts *models.Metrics, asset *models.Asset) (lst []*prompb.TimeSeries) {
	// 用于健康检查使用的是query接口，返回一定是Vector类型
	result, ok := value.(model.Vector)
	if !ok {
		return
	}

	health := 0
	for _, resultValue := range result {
		health = 1
		asset.Metrics[mts.Name] = map[string]string{
			"label": string(resultValue.Metric[LabelName]),
			"value": fmt.Sprintf("%f", resultValue.Value),
			"name":  mts.Name,
		}
	}
	ts := convertHealthSeries(asset, health)
	lst = append(lst, ts)

	return
}

func convertHealthSeries(asset *models.Asset, health int) (ts *prompb.TimeSeries) {
	var healthLabels []prompb.Label
	healthLabels = append(healthLabels, prompb.Label{
		Name:  LabelName,
		Value: HealthMetric,
	})
	healthLabels = append(healthLabels, prompb.Label{
		Name:  AssetId,
		Value: fmt.Sprintf("%d", asset.Id),
	})
	healthLabels = append(healthLabels, prompb.Label{
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
