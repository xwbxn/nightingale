package router

import (
	"github.com/ccfos/nightingale/v6/models"
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

/***
  对外提供的资产信息列表
  author:guoxp
*/

type AssetJson struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	Status   int64               `json:"status"`
	UpdateAt int64               `json:"update_at"`
	Category string              `json:"category"`
	Type     string              `json:"type"`
	Metrics  []map[string]string `json:"metrics"`
	GroupId  int64               `json:"group_id"`
	Tags     []string            `json:"tags"`
}

func (rt *Router) getDashboardAssetsByFE(c *gin.Context) {

	category := ginx.QueryStr(c, "category", "")
	atype := ginx.QueryStr(c, "atype", "")
	groupId := ginx.QueryInt64(c, "group_id", -1)

	var data []*AssetJson
	lst := rt.assetCache.GetAll()
	for _, item := range lst {
		ar, _ := rt.assetCache.GetType(item.Type)

		if category != "" && ar.Category != category {
			continue
		}
		if atype != "" && item.Type != atype {
			continue
		}
		if groupId > -1 && item.OrganizeId != groupId { //这里使用orgid作为group返回查询条件
			continue
		}

		data = append(data, &AssetJson{
			Id:       item.Id,
			Name:     item.Name,
			Status:   item.Health,
			UpdateAt: item.HealthAt,
			Category: ar.Category,
			Type:     item.Type,
			Metrics:  item.Metrics,
			GroupId:  item.OrganizeId, //这里使用orgid作为group返回查询条件
			Tags:     item.TagsJSON,
		})
	}
	//ws.SetMessage(1, data) //socket推送内容

	ginx.NewRender(c).Data(data, nil)
}

func (rt *Router) getOrganizeTreeByFE(c *gin.Context) {
	list, err := models.OrganizeTreeGetsFE(rt.Ctx)
	ginx.Dangerous(err)
	ginx.NewRender(c).Data(list, nil)
}

func (rt *Router) getAlertListByFE(c *gin.Context) {
	list, err := models.AlertFeList(rt.Ctx)

	//从资产缓存中更新orgid
	for _, item := range list {
		asset, has := rt.assetCache.Get(item.Id)
		if has {
			item.OrganizeId = int(asset.OrganizeId)
		}
	}

	ginx.NewRender(c).Data(list, err)
}
