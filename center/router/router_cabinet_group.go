// Package models  机柜组信息
// date : 2023-07-15 14:30
// desc : 机柜组信息
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取机柜组信息
// @Description  根据主键获取机柜组信息
// @Tags         机柜组信息
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.CabinetGroup
// @Router       /api/n9e/cabinet-group/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) cabinetGroupGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	cabinetGroup, err := models.CabinetGroupGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if cabinetGroup == nil {
		ginx.Bomb(404, "No such cabinet_group")
	}

	ginx.NewRender(c).Data(cabinetGroup, nil)
}

// @Summary      查询机柜组信息
// @Description  根据条件查询机柜组信息
// @Tags         机柜组信息
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.CabinetGroup
// @Router       /api/n9e/cabinet-group/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) cabinetGroupGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.CabinetGroupCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.CabinetGroupGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建机柜组信息
// @Description  创建机柜组信息
// @Tags         机柜组信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.CabinetGroup true "add cabinetGroup"
// @Success      200
// @Router       /api/n9e/cabinet-group/ [post]
// @Security     ApiKeyAuth
func (rt *Router) cabinetGroupAdd(c *gin.Context) {
	var f models.CabinetGroup
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新机柜组信息
// @Description  更新机柜组信息
// @Tags         机柜组信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.CabinetGroup true "update cabinetGroup"
// @Success      200
// @Router       /api/n9e/cabinet-group/ [put]
// @Security     ApiKeyAuth
func (rt *Router) cabinetGroupPut(c *gin.Context) {
	var f models.CabinetGroup
	ginx.BindJSON(c, &f)

	old, err := models.CabinetGroupGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "cabinet_group not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除机柜组信息
// @Description  根据主键删除机柜组信息
// @Tags         机柜组信息
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/cabinet-group/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) cabinetGroupDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	cabinetGroup, err := models.CabinetGroupGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if cabinetGroup == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(cabinetGroup.Del(rt.Ctx))
}
