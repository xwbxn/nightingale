package prom

import (
	"fmt"
	"math"
	"strconv"

	"github.com/prometheus/common/model"
)

/* 声明全局变量 */
var dsMon map[int64][]string

func GetValue(value model.Value) float64 {

	var res float64 = 0

	if value.Type() == model.ValVector {
		items, ok := value.(model.Vector)
		if !ok {
			return res
		}

		for _, item := range items {
			if math.IsNaN(float64(item.Value)) {
				continue
			}
			floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(item.Value)), 64)
			res = floatVal
		}
	} else if value.Type() == model.ValMatrix {
		items, ok := value.(model.Matrix)
		if !ok {
			return res
		}

		for _, item := range items {
			if len(item.Values) == 0 {
				return res
			}
			if math.IsNaN(float64(item.Values[0].Value)) {
				continue
			}

			floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(item.Values[0].Value)), 64)
			res = floatVal
		}
	}
	return res
}
