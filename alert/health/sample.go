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
	HealthMetric  = "asset_health"
)

func ConvertToTimeSeries(value model.Value, atype *models.AssetType, assets []*models.Asset) (lst []*prompb.TimeSeries) {
	// 用于健康检查使用的是query接口，返回一定是Vector类型
	items, ok := value.(model.Vector)
	if !ok {
		return
	}

	for _, asset := range assets {
		var labels []*prompb.Label
		labels = append(labels, &prompb.Label{
			Name:  LabelName,
			Value: HealthMetric,
		})
		labels = append(labels, &prompb.Label{
			Name:  AssetId,
			Value: fmt.Sprintf("%d", asset.Id),
		})
		labels = append(labels, &prompb.Label{
			Name:  AssetInstance,
			Value: asset.Label,
		})

		s := prompb.Sample{}
		s.Timestamp = time.Now().UnixNano() / 1e6
		s.Value = 0

		for _, value := range items {
			asset_id, has := value.Metric[AssetId]
			if !has {
				continue
			}
			if asset_id == model.LabelValue(fmt.Sprintf("%d", asset.Id)) {
				s.Value = 1
				break
			}
		}

		asset.Health = int64(s.Value)
		asset.HealthAt = time.Now().Unix()
		lst = append(lst, &prompb.TimeSeries{
			Labels:  labels,
			Samples: []prompb.Sample{s},
		})
	}

	return
}
