// Package models  资产管理
// date : 2023-08-04 09:47
// desc : 资产管理
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取资产管理
// @Description  根据主键获取资产管理
// @Tags         资产管理
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.AssetManagement
// @Router       /api/n9e/asset-management/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) assetManagementGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetManagement, err := models.AssetManagementGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if assetManagement == nil {
		ginx.Bomb(404, "No such asset_management")
	}

	ginx.NewRender(c).Data(assetManagement, nil)
}

// @Summary      根据资产ID获取资产管理
// @Description  根据资产ID获取资产管理
// @Tags         资产管理
// @Accept       json
// @Produce      json
// @Param        asset    query    int  true  "资产ID"
// @Success      200  {object}  models.AssetManagement
// @Router       /api/n9e/asset-management/asset [get]
// @Security     ApiKeyAuth
func (rt *Router) assetManagementGetAI(c *gin.Context) {
	assetId := ginx.QueryInt64(c, "asset", -1)
	assetManagement, err := models.AssetManagementGetByAssetId(rt.Ctx, assetId)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(assetManagement, nil)
}

// @Summary      查询资产管理
// @Description  根据条件查询资产管理
// @Tags         资产管理
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.AssetManagement
// @Router       /api/n9e/asset-management/ [get]
// @Security     ApiKeyAuth
func (rt *Router) assetManagementGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.AssetManagementCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.AssetManagementGets(rt.Ctx, query, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建资产管理
// @Description  创建资产管理
// @Tags         资产管理
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetManagement true "add assetManagement"
// @Success      200
// @Router       /api/n9e/asset-management/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetManagementAdd(c *gin.Context) {
	var f models.AssetManagement
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新资产管理
// @Description  更新资产管理
// @Tags         资产管理
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetManagement true "update assetManagement"
// @Success      200
// @Router       /api/n9e/asset-management/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetManagementPut(c *gin.Context) {
	var f models.AssetManagement
	ginx.BindJSON(c, &f)

	old, err := models.AssetManagementGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "asset_management not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除资产管理
// @Description  根据主键删除资产管理
// @Tags         资产管理
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/asset-management/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) assetManagementDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetManagement, err := models.AssetManagementGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if assetManagement == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(assetManagement.Del(rt.Ctx))
}
