// Package models  资产扩展
// date : 2023-07-23 09:06
// desc : 资产扩展
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取资产扩展
// @Description  根据主键获取资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.AssetExpansion
// @Router       /api/n9e/asset-expansion/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetExpansion, err := models.AssetExpansionGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if assetExpansion == nil {
		ginx.Bomb(404, "No such asset_expansion")
	}

	ginx.NewRender(c).Data(assetExpansion, nil)
}

// @Summary      查询资产扩展
// @Description  根据条件查询资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.AssetExpansion
// @Router       /api/n9e/asset-expansion/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.AssetExpansionCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.AssetExpansionGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建资产扩展
// @Description  创建资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetExpansion true "add assetExpansion"
// @Success      200
// @Router       /api/n9e/asset-expansion/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionAdd(c *gin.Context) {
	var f models.AssetExpansion
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      批量创建资产扩展
// @Description  批量创建资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        body  body   []models.AssetExpansion true "add assetExpansion"
// @Success      200
// @Router       /api/n9e/asset-expansion/batch [post]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionBatchAdd(c *gin.Context) {
	var f []models.AssetExpansion
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	for index := range f {
		f[index].CreatedBy = me.Username
	}

	// 更新模型
	err := models.BatchAdd(rt.Ctx, f)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新资产扩展
// @Description  更新资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetExpansion true "update assetExpansion"
// @Success      200
// @Router       /api/n9e/asset-expansion/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionPut(c *gin.Context) {
	var f models.AssetExpansion
	ginx.BindJSON(c, &f)

	old, err := models.AssetExpansionGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "asset_expansion not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除资产扩展
// @Description  根据主键删除资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/asset-expansion/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetExpansion, err := models.AssetExpansionGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if assetExpansion == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(assetExpansion.Del(rt.Ctx))
}
