// Package models  监控
// date : 2023-10-08 16:45
// desc : 监控
package router

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/prom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取监控
// @Description  根据主键获取监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.Monitoring
// @Router       /api/n9e/xh/monitoring/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) monitoringGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	monitoring, err := models.MonitoringGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if monitoring == nil {
		ginx.Bomb(404, "No such monitoring")
	}

	ginx.NewRender(c).Data(monitoring, nil)
}

// @Summary      根据资产id获取监控
// @Description  根据主键获取监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        id    query    int  true  "主键"
// @Success      200  {object}  models.Monitoring
// @Router       /api/n9e/xh/monitoring/asset [get]
// @Security     ApiKeyAuth
func (rt *Router) monitoringGetByAssetId(c *gin.Context) {
	id := ginx.QueryInt64(c, "id", -1)
	monitorings, err := models.MonitoringGetMap(rt.Ctx, map[string]interface{}{"asset_id": id})
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(monitorings, err)
}

// @Summary      查询监控
// @Description  根据条件查询监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        assetId query   int     false  "资产id"
// @Param        assetType query   int     false  "资产类型"
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "搜索栏"
// @Param        filter    query    string  false  "筛选框(“monitoring_name”：监控名称;“asset_name”：资产名称;“asset_ip”：资产IP;“asset_type”：资产类型;“status”：监控状态;“is_alarm”：是否启用告警)"
// @Success      200  {array}  models.Monitoring
// @Router       /api/n9e/xh/monitoring/filter [get]
// @Security     ApiKeyAuth
func (rt *Router) monitoringGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	assetType := ginx.QueryStr(c, "assetType", "")
	// query := ginx.QueryStr(c, "query", "")
	// datasource := ginx.QueryInt64(c, "datasource", -1)
	assetId := ginx.QueryInt64(c, "assetId", -1)

	// m := make(map[string]interface{})
	// if dataSource != -1 {
	// 	m["datasource.id"] = dataSource
	// }
	// if assetType != ""{
	// 	m["assets.type"] = assetType
	// }
	// total, err := models.MonitoringMapCount(rt.Ctx, m, query, assetType, datasource, assetId)
	// ginx.Dangerous(err)
	// lst, err := models.MonitoringMapGets(rt.Ctx, m, query, limit, (page-1)*limit, assetType, datasource, assetId)
	// ginx.Dangerous(err)

	// ginx.NewRender(c).Data(gin.H{
	// 	"list":  lst,
	// 	"total": total,
	// }, nil)

	filter := ginx.QueryStr(c, "filter", "")
	query := ginx.QueryStr(c, "query", "")

	assetIds := make([]int64, 0)
	assets := rt.assetCache.GetAll()
	logger.Debug(len(assets))
	for _, asset := range assets {
		if assetType != "" && strings.Contains(asset.Type, assetType) {
			if filter == "asset_ip" && strings.Contains(asset.Ip, query) {
				assetIds = append(assetIds, asset.Id)
			} else if filter == "asset_type" && strings.Contains(asset.Type, query) {
				assetIds = append(assetIds, asset.Id)
			} else if filter == "asset_name" && strings.Contains(asset.Name, query) {
				assetIds = append(assetIds, asset.Id)
			} else if filter == "" {
				assetIds = append(assetIds, asset.Id)
			}
		} else if assetType == "" {
			if filter == "asset_ip" && strings.Contains(asset.Ip, query) {
				logger.Debug(asset.Ip)
				assetIds = append(assetIds, asset.Id)
			} else if filter == "asset_type" && strings.Contains(asset.Type, query) {
				assetIds = append(assetIds, asset.Id)
			} else if filter == "asset_name" && strings.Contains(asset.Name, query) {
				assetIds = append(assetIds, asset.Id)
			}
		}
	}

	total, err := models.MonitoringMapCountNew(rt.Ctx, filter, query, assetId, assetIds)
	ginx.Dangerous(err)
	lst, err := models.MonitoringMapGetsNew(rt.Ctx, filter, query, assetId, assetIds, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      查询所有监控
// @Description  查询所有监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.Monitoring
// @Router       /api/n9e/xh/monitoring/ [get]
// @Security     ApiKeyAuth
func (rt *Router) monitoringAllGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.MonitoringCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.MonitoringAllGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建监控
// @Description  创建监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        body  body   models.Monitoring true "add monitoring"
// @Success      200
// @Router       /api/n9e/xh/monitoring/ [post]
// @Security     ApiKeyAuth
func (rt *Router) monitoringAdd(c *gin.Context) {
	var f models.Monitoring
	ginx.BindJSON(c, &f)
	me := c.MustGet("user").(*models.User)
	now := time.Now().Unix()
	var err error

	if len(f.AssetIds) == 0 {
		var monitor = models.Monitoring{
			AssetId:        f.AssetId,
			MonitoringName: f.MonitoringName,
			DatasourceId:   f.DatasourceId,
			MonitoringSql:  f.MonitoringSql,
			Status:         f.Status,
			TargetId:       f.TargetId,
			Remark:         f.Remark,
			Unit:           f.Unit,
			CreatedBy:      me.Username,
			CreatedAt:      now,
			UpdatedBy:      me.Username,
			UpdatedAt:      now,
		}

		// 更新模型
		err = monitor.Add(rt.Ctx)
	} else {
		tx := models.DB(rt.Ctx).Begin()
		for _, id := range f.AssetIds {
			var monitor = models.Monitoring{
				AssetId:        id,
				MonitoringName: f.MonitoringName,
				DatasourceId:   f.DatasourceId,
				MonitoringSql:  f.MonitoringSql,
				Status:         f.Status,
				TargetId:       f.TargetId,
				Remark:         f.Remark,
				Unit:           f.Unit,
				CreatedBy:      me.Username,
				CreatedAt:      now,
				UpdatedBy:      me.Username,
				UpdatedAt:      now,
			}
			err = monitor.AddTx(tx)
		}
		tx.Commit()
	}

	ginx.NewRender(c).Message(err)
}

// @Summary      更新监控
// @Description  更新监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        body  body   models.Monitoring true "update monitoring"
// @Success      200
// @Router       /api/n9e/xh/monitoring/ [put]
// @Security     ApiKeyAuth
func (rt *Router) monitoringPut(c *gin.Context) {
	var f models.Monitoring
	ginx.BindJSON(c, &f)

	old, err := models.MonitoringGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "monitoring not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除监控
// @Description  根据主键删除监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/xh/monitoring/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) monitoringDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	monitoring, err := models.MonitoringGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if monitoring == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	tx := models.DB(rt.Ctx).Begin()
	ids := make([]string, 0)
	ids = append(ids, strconv.FormatInt(monitoring.Id, 10))
	//删除监控
	err = models.BatchDelTx(tx, ids)
	ginx.Dangerous(err)
	//删除资产告警规则
	err = models.AlertRuleDelTxByMonId(tx, ids)
	tx.Commit()
	ginx.NewRender(c).Message(err)
}

// @Summary      批量删除监控-西航
// @Description  批量删除监控-西航
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        body  body  map[string]string   true "add ids"
// @Success      200
// @Router       /api/n9e/xh/monitoring/batch-del/ [post]
// @Security     ApiKeyAuth
func (rt *Router) monitoringDelXH(c *gin.Context) {
	var f map[string]interface{}
	ginx.BindJSON(c, &f)

	idsTemp, idsOk := f["ids"]
	ids := make([]string, 0)
	// var err error
	if idsOk {
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, val.(string))
		}
	} else {
		ginx.Bomb(http.StatusOK, "参数错误")
	}
	tx := models.DB(rt.Ctx).Begin()
	//删除监控
	err := models.BatchDelTx(tx, ids)
	ginx.Dangerous(err)
	//删除资产告警规则
	err = models.AlertRuleDelTxByMonId(tx, ids)
	tx.Commit()

	ginx.NewRender(c).Message(err)
}

// @Summary      监控开关
// @Description  根据主键更新监控状态
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param		 type	query	int		true "类型"
// @Param		 status	query	int		true "新状态"
// @Param        body	body	[]int64 true "add ids"
// @Success      200
// @Router       /api/n9e/xh/monitoring/status [post]
// @Security     ApiKeyAuth
func (rt *Router) monitoringStatusUpdate(c *gin.Context) {
	status := ginx.QueryInt64(c, "status", -1)
	oType := ginx.QueryInt64(c, "type", -1)
	if status == -1 || oType == -1 {
		ginx.Bomb(http.StatusOK, "参数错误")
	}

	var ids []int64
	ginx.BindJSON(c, &ids)

	ginx.NewRender(c).Message(models.MonitoringUpdateStatus(rt.Ctx, ids, status, oType))
}

// @Summary      获取监控数据
// @Description  获取监控数据
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        start query   int     false  "开始时间"
// @Param        end query   int     false  "结束时间"
// @Param        body  body   map[string]interface{} true "update query"
// @Success      200
// @Router       /api/n9e/xh/monitoring/data [post]
// @Security     ApiKeyAuth
func (rt *Router) monitoringData(c *gin.Context) {

	startT := ginx.QueryInt64(c, "start", -1)
	endT := ginx.QueryInt64(c, "end", -1)
	if startT == -1 || endT == -1 {
		ginx.Bomb(http.StatusOK, "参数为空")
	}
	if startT >= endT {
		ginx.Bomb(http.StatusOK, "时间区间错误")
	}

	var f map[string]interface{}
	ginx.BindJSON(c, &f)

	idsTemp, idsOk := f["ids"]
	ids := make([]int64, 0)
	// var err error
	if idsOk {
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, int64(val.(float64)))
		}
	}
	lst, err := models.MonitoringGetByBatchId(rt.Ctx, ids)
	ginx.Dangerous(err)

	end := time.Unix(endT, 0)
	start := time.Unix(startT, 0)
	var r prom.Range
	if endT-startT <= 60*60*24 {
		r = prom.Range{Start: start, End: end, Step: prom.DefaultStep * 2}
	} else if endT-startT > 60*60*24 && endT-startT <= 60*60*24*7 {
		r = prom.Range{Start: start, End: end, Step: prom.DefaultStep * 2 * 60}
	} else if endT-startT > 60*60*24*7 {
		r = prom.Range{Start: start, End: end, Step: prom.DefaultStep * 2 * 60 * 60}
	}

	values := make([]model.Value, 0)
	for _, monitoring := range lst {
		sql, err := prom.InjectLabel(monitoring.MonitoringSql, "asset_id", fmt.Sprintf("%d", monitoring.AssetId), labels.MatchEqual)
		ginx.Dangerous(err)

		value, warnings, err := rt.PromClients.GetCli(monitoring.DatasourceId).QueryRange(context.Background(), sql, r)

		if len(warnings) > 0 {
			logger.Error(err)
			return
		}

		// values = append(values, prom.FormatPromValue(value, monitoring.Unit))
		values = append(values, value)
	}
	ginx.NewRender(c).Data(values, err)
}

type monitoringOption struct {
	Label    string             `json:"label"`
	Value    string             `json:"value"`
	PromQL   string             `json:"promql"`
	Children []monitoringOption `json:"children"`
}

// 获取监控指标前端选择列表
func (rt *Router) monitoringGetOptions(c *gin.Context) {
	lst, err := models.MonitoringGetsAll(rt.Ctx)
	ginx.Dangerous(err)

	assets, err := models.AssetGetsAll(rt.Ctx)
	ginx.Dangerous(err)

	data := make([]monitoringOption, len(assets))
	for i, asset := range assets {
		data[i].Label = asset.Name
		data[i].Value = strconv.Itoa(int(asset.Id))
		data[i].Children = make([]monitoringOption, 0)
		for _, m := range lst {
			if m.AssetId == asset.Id {
				data[i].Children = append(data[i].Children, monitoringOption{
					Label:  m.MonitoringName,
					Value:  m.CompilePromQL(),
					PromQL: m.CompilePromQL(),
				})
			}
		}
	}

	ginx.NewRender(c).Data(data, nil)
}

// @Summary      获取监控指标单位
// @Description  获取监控指标单位
// @Tags         监控
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/n9e/xh/monitoring/unit [get]
// @Security     ApiKeyAuth
func (rt *Router) monitoringUnit(c *gin.Context) {
	ginx.NewRender(c).Data(prom.UNIT_LIST, nil)
}
