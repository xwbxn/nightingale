// Package models  监控
// date : 2023-10-08 16:45
// desc : 监控
package router

import (
	"net/http"
	"time"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
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

// @Summary      查询监控
// @Description  根据条件查询监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "搜索栏"
// @Param        assetType query   string  false  "资产类型"
// @Param        datasource query   int  false  "数据源类型"
// @Success      200  {array}  models.Monitoring
// @Router       /api/n9e/xh/monitoring/filter [get]
// @Security     ApiKeyAuth
func (rt *Router) monitoringGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	assetType := ginx.QueryStr(c, "assetType","")
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")
	datasource := ginx.QueryInt(c, "datasource", -1)

	m := make(map[string]interface{})
	// if dataSource != -1 {
	// 	m["datasource.id"] = dataSource
	// }
	// if assetType != ""{
	// 	m["assets.type"] = assetType
	// }
	total, err := models.MonitoringMapCount(rt.Ctx, m, query, assetType, datasource)
	ginx.Dangerous(err)
	lst, err := models.MonitoringMapGets(rt.Ctx, m, query, limit, (page-1)*limit, assetType, datasource)
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

type monitoring struct {
	Id             int64  `json:"id"`
	AssetId        int64  `json:"asset_id"`
	MonitoringName string `json:"monitoring_name"`
	DatasourceId   int64  `json:"datasource_id"`
	MonitoringSql  string `json:"monitoring_sql"`
	Status         int64  `json:"status"`
	TargetId       int64  `json:"target_id"`
	Remark         string `json:"remark"`
	CreatedBy      string `json:"created_by"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedBy      string `json:"updated_by"`
	UpdatedAt      int64  `json:"updated_at"`
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
	var f monitoring
	ginx.BindJSON(c, &f)
	me := c.MustGet("user").(*models.User)

	var monitor = models.Monitoring{
		AssetId:        f.AssetId,
		MonitoringName: f.MonitoringName,
		DatasourceId:   f.DatasourceId,
		MonitoringSql:  f.MonitoringSql,
		Status:         f.Status,
		TargetId:       f.TargetId,
		Remark:         f.Remark,
		CreatedBy:      me.Username,
		CreatedAt:      time.Now().Unix(),
		UpdatedBy:      me.Username,
		UpdatedAt:      time.Now().Unix(),
	}

	// 更新模型
	err := monitor.Add(rt.Ctx)
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
	ginx.NewRender(c).Message(monitoring.Del(rt.Ctx))
}

// @Summary      监控开关
// @Description  根据主键更新监控状态
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param		 status	query	int		true "新状态"
// @Param        body	body	[]int64 true "add ids"
// @Success      200
// @Router       /api/n9e/xh/monitoring/status/{id} [post]
// @Security     ApiKeyAuth
func (rt *Router) monitoringStatusUpdate(c *gin.Context) {
	status := ginx.QueryInt64(c, "status", 1)
	var ids []int64
	ginx.BindJSON(c, &ids)

	ginx.NewRender(c).Message(models.MonitoringUpdateStatus(rt.Ctx, ids, status))
}
