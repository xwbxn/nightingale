// Package models  操作日志
// date : 2023-10-21 09:13
// desc : 操作日志
package router

import (
	"net/http"
	"time"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取操作日志
// @Description  根据主键获取操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.OperationLog
// @Router       /api/n9e/operation-log/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) operationLogGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	operationLog, err := models.OperationLogGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if operationLog == nil {
		ginx.Bomb(404, "No such operation_log")
	}

	ginx.NewRender(c).Data(operationLog, nil)
}

// @Summary      查询操作日志
// @Description  根据条件查询操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.OperationLog
// @Router       /api/n9e/operation-log/ [get]
// @Security     ApiKeyAuth
func (rt *Router) operationLogGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.OperationLogCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.OperationLogGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      过滤器查询操作日志
// @Description  过滤器条件查询操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Param        start query   int  false  "开始时间"
// @Param        end query   int  false  "结束时间"
// @Param        filterType query string  false  "类型"
// @Param        modelType query string  false  "模块类型"
// @Success      200  {array}  models.OperationLog
// @Router       /api/n9e/operation-log/filter [get]
// @Security     ApiKeyAuth
func (rt *Router) operationLogFilterGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")
	start := ginx.QueryInt64(c,"start", 0)
	end := ginx.QueryInt64(c,"end",time.Now().Unix())
	filterType := ginx.QueryStr(c,"filterType", "")
	modelTye := ginx.QueryStr(c, "modelType","")

	if end == -1{
		end = time.Now().Unix()
	}
	if end < start {
			ginx.Bomb(http.StatusOK, "Wrong date range selection")
		}
	
	total, err := models.FilterLogCount(rt.Ctx, query, start, end, filterType, modelTye)
	ginx.Dangerous(err)
	lst, err := models.FilterLogGets(rt.Ctx, query, (page-1)*limit, limit, start, end, filterType, modelTye)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建操作日志
// @Description  创建操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        body  body   models.OperationLog true "add operationLog"
// @Success      200
// @Router       /api/n9e/operation-log/ [post]
// @Security     ApiKeyAuth
func (rt *Router) operationLogAdd(c *gin.Context) {
	var f models.OperationLog
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新操作日志
// @Description  更新操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        body  body   models.OperationLog true "update operationLog"
// @Success      200
// @Router       /api/n9e/operation-log/ [put]
// @Security     ApiKeyAuth
func (rt *Router) operationLogPut(c *gin.Context) {
	var f models.OperationLog
	ginx.BindJSON(c, &f)

	old, err := models.OperationLogGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "operation_log not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除操作日志
// @Description  根据主键删除操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/operation-log/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) operationLogDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	operationLog, err := models.OperationLogGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if operationLog == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(operationLog.Del(rt.Ctx))
}
