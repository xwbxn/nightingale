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

// @Summary      根据资产ID获取资产扩展
// @Description  根据资产ID获取资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        asset    query    int  true  "资产ID"
// @Success      200  {array}  []models.AssetExpansion
// @Router       /api/n9e/asset-expansion/asset [get]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionGetByAssetId(c *gin.Context) {
	assetId := ginx.QueryInt64(c, "asset", -1)
	m := make(map[string]interface{})
	m["asset_id"] = assetId

	lst, err := models.AssetExpansionGetByMap(rt.Ctx, m)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(lst, nil)
}

// @Summary      根据资产ID、类别或属性获取资产扩展
// @Description  根据资产ID、类别或属性获取资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetExpansion true "add AssetExpansion"
// @Success      200  {array}  []models.AssetExpansion
// @Router       /api/n9e/asset-expansion/map [post]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionGetByMap(c *gin.Context) {

	m := make(map[string]interface{})
	ginx.BindJSON(c, &m)

	lst, err := models.AssetExpansionGetByMap(rt.Ctx, m)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(lst, nil)
}

// @Summary      查询资产扩展
// @Description  根据条件查询资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.AssetExpansion
// @Router       /api/n9e/asset-expansion/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.AssetExpansionCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.AssetExpansionGets(rt.Ctx, query, limit, (page-1)*limit)
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
// @Router       /api/n9e/asset-expansion/batch/ [post]
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
	err := models.AssetExpansionBatchAdd(rt.Ctx, f)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新资产扩展
// @Description  更新资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        body  body   []models.AssetExpansion true "update assetExpansion"
// @Success      200
// @Router       /api/n9e/asset-expansion/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionPut(c *gin.Context) {

	var f []models.AssetExpansion
	ginx.BindJSON(c, &f)
	if len(f) == 0 {
		ginx.Bomb(http.StatusOK, "Not update asset_expansion")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	name := me.Username

	var err error
	m := make(map[string]interface{})

	//确定更新位置为硬件配置
	if f[0].ConfigCategory == "form-hardware-cfg" {

		//通过资产id、属性等查询
		m["ASSET_ID"] = f[0].AssetId
		m["CONFIG_CATEGORY"] = f[0].ConfigCategory
		m["PROPERTY_CATEGORY"] = f[0].PropertyCategory

		err = models.UpdateAssetExpansionGroup(rt.Ctx, m, f, name)
		ginx.Dangerous(err)

	} else if f[0].ConfigCategory == "form-netconfig" {

		m["ASSET_ID"] = f[0].AssetId
		m["CONFIG_CATEGORY"] = f[0].ConfigCategory

		err = models.UpdateAssetExpansionGroup(rt.Ctx, m, f, name)
		ginx.Dangerous(err)
	}

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(err)
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
