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

// @Summary      更新资产扩展-西航
// @Description  更新资产扩展-西航
// @Tags         资产扩展-西航
// @Accept       json
// @Produce      json
// @Param        asset query   int     true  "资产ID"
// @Param        type query    string    true  "属性"
// @Param        body  body   []models.AssetsExpansion true "update assetsExpansion"
// @Success      200
// @Router       /api/n9e/xh/assets-expansion/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetsExpansionPut(c *gin.Context) {

	assetId := ginx.QueryInt64(c, "asset", -1)
	configCategory := ginx.QueryStr(c, "type", "")
	if configCategory == "" || assetId == -1 {
		ginx.Bomb(http.StatusOK, "参数错误")
	}
	var f []models.AssetsExpansion
	ginx.BindJSON(c, &f)

	tx := models.DB(rt.Ctx).Begin()

	err := models.AssetsExpansionDelMap(tx, map[string]interface{}{"assets_id": assetId, "config_category": configCategory})
	ginx.Dangerous(err)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	if len(f) > 0 {
		for index := range f {
			f[index].CreatedBy = me.Username
		}
		err = models.AssetsExpansionAddTx(tx, f)
		ginx.Dangerous(err)
	}
	tx.Commit()
	ginx.NewRender(c).Message(err)
}
