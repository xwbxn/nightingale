// Package models
// date : 2023-10-08 15:32
// desc :
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取
// @Description  根据主键获取
// @Tags
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.Bigscreen
// @Router       /api/n9e/bigscreen/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) bigscreenGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	bigscreen, err := models.BigscreenGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if bigscreen == nil {
		ginx.Bomb(404, "No such bigscreen")
	}

	ginx.NewRender(c).Data(bigscreen, nil)
}

// @Summary      查询
// @Description  根据条件查询
// @Tags
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.Bigscreen
// @Router       /api/n9e/bigscreen/ [get]
// @Security     ApiKeyAuth
func (rt *Router) bigscreenGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.BigscreenCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.BigscreenGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建
// @Description  创建
// @Tags
// @Accept       json
// @Produce      json
// @Param        body  body   models.Bigscreen true "add bigscreen"
// @Success      200
// @Router       /api/n9e/bigscreen/ [post]
// @Security     ApiKeyAuth
func (rt *Router) bigscreenAdd(c *gin.Context) {
	var f models.Bigscreen
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新
// @Description  更新
// @Tags
// @Accept       json
// @Produce      json
// @Param        body  body   models.Bigscreen true "update bigscreen"
// @Success      200
// @Router       /api/n9e/bigscreen/ [put]
// @Security     ApiKeyAuth
func (rt *Router) bigscreenPut(c *gin.Context) {
	var f models.Bigscreen
	ginx.BindJSON(c, &f)

	old, err := models.BigscreenGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "bigscreen not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除
// @Description  根据主键删除
// @Tags
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/bigscreen/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) bigscreenDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	bigscreen, err := models.BigscreenGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if bigscreen == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(bigscreen.Del(rt.Ctx))
}
