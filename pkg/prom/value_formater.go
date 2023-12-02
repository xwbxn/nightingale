package prom

import "github.com/prometheus/common/model"

// prom基础单位
// 时间：秒   流量：bit   存储：Byte
// 如果指标不属于基础单位，需要在指标配置的promql里转换为基础单位
var UNIT_LIST map[string]string = map[string]string{
	"b":  "b", // bit
	"Kb": "Kb",
	"Mb": "Mb",
	"Gb": "Gb",
	"B":  "B", // BYTE
	"KB": "KB",
	"MB": "MB",
	"GB": "GB",
	"MS": "毫秒",
	"S":  "秒", // second
	"M":  "分",
	"H":  "时",
	"D":  "天",
	"%":  "%",
}

// TODO: 实现单位的转换
func formatValue(unit string, value float64) float64 {
	switch unit { // key of UNIT_LIST
	case "b":
		return value
	default:
		return value
	}
}

func FormatPromValue(value model.Value, unit string) model.Value {
	switch value.Type() {
	case model.ValVector:
		items, ok := value.(model.Vector)
		if !ok {
			break
		}
		for _, item := range items {
			item.Metric["unit"] = model.LabelValue(unit)
			item.Value = model.SampleValue(formatValue(unit, float64(item.Value)))
		}

	case model.ValMatrix:
		items, ok := value.(model.Matrix)
		if !ok {
			break
		}

		for _, item := range items {
			if len(item.Values) == 0 {
				break
			}
			item.Metric["unit"] = model.LabelValue(unit)

			for _, v := range item.Values {
				v.Value = model.SampleValue(formatValue(unit, float64(v.Value)))
			}
		}
	default:
	}
	return value
}
