// Package models  资产维保
// date : 2023-07-23 09:44
// desc : 资产维保
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取资产维保
// @Description  根据主键获取资产维保
// @Tags         资产维保
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.AssetMaintenance
// @Router       /api/n9e/asset-maintenance/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) assetMaintenanceGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetMaintenance, err := models.AssetMaintenanceGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if assetMaintenance == nil {
		ginx.Bomb(404, "No such asset_maintenance")
	}

	ginx.NewRender(c).Data(assetMaintenance, nil)
}

// @Summary      根据资产ID获取资产维保
// @Description  根据资产ID获取资产维保
// @Tags         资产维保
// @Accept       json
// @Produce      json
// @Param        asset    query    int  true  "资产ID"
// @Success      200  {object}  models.AssetMaintenanceVo
// @Router       /api/n9e/asset-maintenance/asset [get]
// @Security     ApiKeyAuth
func (rt *Router) assetMaintenanceGetAssetId(c *gin.Context) {
	assetId := ginx.QueryInt64(c, "asset", -1)
	assetMaintenanceVo, err := models.AssetMaintenanceVoGetByAssetId(rt.Ctx, assetId)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(assetMaintenanceVo, nil)
}

// @Summary      查询资产维保
// @Description  根据条件查询资产维保
// @Tags         资产维保
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.AssetMaintenance
// @Router       /api/n9e/asset-maintenance/ [get]
// @Security     ApiKeyAuth
func (rt *Router) assetMaintenanceGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.AssetMaintenanceCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.AssetMaintenanceGets(rt.Ctx, query, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建资产维保
// @Description  创建资产维保
// @Tags         资产维保
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetMaintenanceVo true "add AssetMaintenanceVo"
// @Success      200
// @Router       /api/n9e/asset-maintenance/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetMaintenanceAdd(c *gin.Context) {
	var fVo models.AssetMaintenanceVo
	ginx.BindJSON(c, &fVo)

	num, err := models.AssetMaintenanceCountMap(rt.Ctx, map[string]interface{}{"asset_id": fVo.AssetId})
	ginx.Dangerous(err)
	if num > 0 {
		ginx.Bomb(http.StatusOK, "该资产存在维保信息")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)

	//启动事务
	tx := models.DB(rt.Ctx).Begin()

	err = fVo.AddConfig(tx, me.Username)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()

	ginx.NewRender(c).Message(err)
}

// @Summary      更新资产维保
// @Description  更新资产维保
// @Tags         资产维保
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetMaintenanceVo true "update AssetMaintenanceVo"
// @Success      200
// @Router       /api/n9e/asset-maintenance/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetMaintenancePut(c *gin.Context) {
	var fVo models.AssetMaintenanceVo
	ginx.BindJSON(c, &fVo)

	old, err := models.AssetMaintenanceGetById(rt.Ctx, fVo.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "asset_maintenance not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	// fVo.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(models.MaintenanceUpdate(rt.Ctx, fVo.Id, me.Username, fVo))
}

// @Summary      删除资产维保
// @Description  根据主键删除资产维保
// @Tags         资产维保
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/asset-maintenance/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) assetMaintenanceDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetMaintenance, err := models.AssetMaintenanceGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if assetMaintenance == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(assetMaintenance.DelTx(rt.Ctx))
}
