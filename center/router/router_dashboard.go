package router

import (
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

type ChartPure struct {
	Configs string `json:"configs"`
	Weight  int    `json:"weight"`
}

type ChartGroupPure struct {
	Name   string      `json:"name"`
	Weight int         `json:"weight"`
	Charts []ChartPure `json:"charts"`
}

type DashboardPure struct {
	Name        string           `json:"name"`
	Tags        string           `json:"tags"`
	Configs     string           `json:"configs"`
	ChartGroups []ChartGroupPure `json:"chart_groups"`
}

type MetricJson struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AssetJson struct {
	Id       int64        `json:"id"`
	Name     string       `json:"name"`
	Status   int64        `json:"status"`
	UpdateAt int64        `json:"update_at"`
	Category string       `json:"category"`
	Type     string       `json:"type"`
	Metrics  []MetricJson `json:"metrics"`
	GroupId  int64        `json:"group_id"`
	Tags     []string     `json:"tags"`
}

func (rt *Router) getDashboardAssets(c *gin.Context) {
	var data []*AssetJson
	lst := rt.assetCache.GetAll()
	for _, item := range lst {
		data = append(data, &AssetJson{
			Id:       item.Id,
			Name:     item.Name,
			Status:   item.Health,
			UpdateAt: item.HealthAt,
			Category: item.Type,
			Type:     item.Type,
			Metrics:  nil,
			GroupId:  item.GroupId,
			Tags:     item.TagsJSON,
		})
	}
	ginx.NewRender(c).Data(data, nil)
}
