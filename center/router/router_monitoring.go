// Package models  监控
// date : 2023-10-08 16:45
// desc : 监控
package router

import (
	"net/http"

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
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.Monitoring
// @Router       /api/n9e/xh/monitoring/ [get]
// @Security     ApiKeyAuth
func (rt *Router) monitoringGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.MonitoringCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.MonitoringGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
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

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
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
