// Package models  资产扩展-西航
// date : 2023-9-20 11:24
// desc : 资产扩展-西航
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取资产扩展-西航
// @Description  根据主键获取资产扩展-西航
// @Tags         资产扩展-西航
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.AssetsExpansion
// @Router       /api/n9e/xh/assets-expansion/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) assetsExpansionGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetsExpansion, err := models.AssetsExpansionGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if assetsExpansion == nil {
		ginx.Bomb(404, "No such assets_expansion")
	}

	ginx.NewRender(c).Data(assetsExpansion, nil)
}

// @Summary      查询资产扩展-西航
// @Description  根据条件查询资产扩展-西航
// @Tags         资产扩展-西航
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.AssetsExpansion
// @Router       /api/n9e/xh/assets-expansion/ [get]
// @Security     ApiKeyAuth
func (rt *Router) assetsExpansionGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.AssetsExpansionCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.AssetsExpansionGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建资产扩展-西航
// @Description  创建资产扩展-西航
// @Tags         资产扩展-西航
// @Accept       json
// @Produce      json
// @Param        body  body   []models.AssetsExpansion true "add assetsExpansion"
// @Success      200
// @Router       /api/n9e/xh/assets-expansion/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetsExpansionAdd(c *gin.Context) {
	var f []models.AssetsExpansion
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	for index := range f {
		f[index].CreatedBy = me.Username
	}

	// 更新模型
	err := models.AssetsExpansionAdd(rt.Ctx, f)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新资产扩展-西航
// @Description  更新资产扩展-西航
// @Tags         资产扩展-西航
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetsExpansion true "update assetsExpansion"
// @Success      200
// @Router       /api/n9e/xh/assets-expansion/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetsExpansionPut(c *gin.Context) {
	var f models.AssetsExpansion
	ginx.BindJSON(c, &f)

	old, err := models.AssetsExpansionGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "assets_expansion not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除资产扩展-西航
// @Description  根据主键删除资产扩展-西航
// @Tags         资产扩展-西航
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/xh/assets-expansion/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) assetsExpansionDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetsExpansion, err := models.AssetsExpansionGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if assetsExpansion == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(assetsExpansion.Del(rt.Ctx))
}
