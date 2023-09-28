// Package models  资产树
// date : 2023-07-21 09:51
// desc : 资产树
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取资产树
// @Description  根据主键获取资产树
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.AssetTree
// @Router       /api/n9e/asset-tree/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) assetTreeGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetTree, err := models.AssetTreeGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if assetTree == nil {
		ginx.Bomb(404, "No such asset_tree")
	}

	ginx.NewRender(c).Data(assetTree, nil)
}

// @Summary      查询资产树
// @Description  根据条件查询资产树
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.AssetTree
// @Router       /api/n9e/asset-tree/ [get]
// @Security     ApiKeyAuth
func (rt *Router) assetTreeGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.AssetTreeCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.AssetTreeGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      根据资产ID查询资产树
// @Description  根据资产ID查询资产树
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "add ids"
// @Success      200  {array}  models.AssetTree
// @Router       /api/n9e/asset-tree/asset [post]
// @Security     ApiKeyAuth
func (rt *Router) assetTreeGetsByAssetId(c *gin.Context) {
	var f []int64
	ginx.BindJSON(c, &f)

	lst, err := models.AssetTreeBatchGetById(rt.Ctx, f)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(lst, err)
}

// @Summary      创建资产树
// @Description  创建资产树
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetTree true "add assetTree"
// @Success      200
// @Router       /api/n9e/asset-tree/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetTreeAdd(c *gin.Context) {
	var f models.AssetTree
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新资产树
// @Description  更新资产树
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetTree true "update assetTree"
// @Success      200
// @Router       /api/n9e/asset-tree/list/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetTreePut(c *gin.Context) {
	var f models.AssetTree
	ginx.BindJSON(c, &f)

	old, err := models.AssetTreeGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "asset_tree not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除资产树
// @Description  根据主键删除资产树
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/asset-tree/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) assetTreeDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")

	var assetTree models.AssetTree
	assetTree.Id = id

	assetNum, err := assetTree.AssetCountGet(rt.Ctx)
	logger.Debug("-------------")
	logger.Debug(assetNum)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if assetNum > 0 {
		ginx.Bomb(404, "This organization hava assets")
	}

	ginx.NewRender(c).Message(assetTree.Del(models.DB(rt.Ctx)))
}

// @Summary      获取资产树数据
// @Description  根据资产状态获取资产树数据
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        status query   int64     false  "设备状态"
// @Success      200  {object}  models.AssetTree
// @Router       /api/n9e/asset-tree/data [get]
// @Security     ApiKeyAuth
func (rt *Router) assetTreeGetALL(c *gin.Context) {

	status := ginx.QueryInt64(c, "status", 1)
	query := make(map[string]interface{})
	query["status"] = status
	assetTree, err := models.BuildAssetTree(rt.Ctx, query)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(assetTree, nil)
}

// @Summary      获取资产树数据（部分）
// @Description  根据资产状态获取资产树数据
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        status query   int64     true  "设备状态"
// @Param        type query   int64     true  "类型"
// @Success      200  {object}  models.AssetTree
// @Router       /api/n9e/asset-tree/part [get]
// @Security     ApiKeyAuth
func (rt *Router) assetTreeGetPart(c *gin.Context) {
	status := ginx.QueryInt64(c, "status", 1)
	dType := ginx.QueryInt64(c, "type", -1)

	deviceType, err := models.DeviceTypeGetById(rt.Ctx, dType)
	ginx.Dangerous(err)
	if (&models.DeviceType{} == deviceType) {
		ginx.Bomb(http.StatusOK, "type not found")
	}

	query := make(map[string]interface{})
	query["status"] = status
	query["name"] = deviceType.Name
	assetTree, err := models.BuildPartAssetTree(rt.Ctx, query)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(assetTree, nil)
}

// @Summary      获取资产数量
// @Description  根据主键获取资产数量
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        id query   int64     false  "主键"
// @Success      200  {object}  models.AssetTree
// @Router       /api/n9e/asset-tree/count [get]
// @Security     ApiKeyAuth
func (rt *Router) assetTreeGetCount(c *gin.Context) {
	id := ginx.QueryInt64(c, "id", 0)
	var assetTree models.AssetTree
	assetTree.Id = id
	asset, err := assetTree.AssetCountGet(rt.Ctx)
	ginx.Dangerous(err)

	// if assetTree == nil {
	// 	ginx.Bomb(404, "No such asset_tree")
	// }

	ginx.NewRender(c).Data(asset, nil)
}

// @Summary      更新资产树
// @Description  更新资产树
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "目标节点"
// @Param        body  body   []models.AssetTree true "update assetTree"
// @Success      200
// @Router       /api/n9e/asset-tree/parent/{id} [post]
// @Security     ApiKeyAuth
func (rt *Router) assetTreeUpdate(c *gin.Context) {
	var f []models.AssetTree
	ginx.BindJSON(c, &f)

	parentId := ginx.UrlParamInt64(c, "id")
	logger.Debug(parentId)

	ids := []int64{}
	for _, val := range f {
		if val.ParentId == parentId {
			ginx.Bomb(http.StatusOK, "原节点与转移节点一致！")
		}
		ids = append(ids, val.Id)
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)

	m := make(map[string]interface{})

	m["parent_id"] = parentId
	m["updated_by"] = me.Username
	logger.Debug(m)

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(models.UpdateTree(rt.Ctx, ids, m))
}

// @Summary      资产树转移
// @Description  资产树转移
// @Tags         资产树
// @Accept       json
// @Produce      json
// @Param        parent    query    int  true  "目标节点"
// @Param        body  body   []int true "update assetTree.ID"
// @Success      200
// @Router       /api/n9e/asset-tree/transfer [post]
// @Security     ApiKeyAuth
func (rt *Router) assetTreeTransfer(c *gin.Context) {
	var f []int64
	ginx.BindJSON(c, &f)
	// 添加审计信息
	me := c.MustGet("user").(*models.User)

	parentId := ginx.QueryInt64(c, "parent", -1)
	if parentId < 0 {
		ginx.Bomb(http.StatusOK, "目标节点错误!")
	}

	for _, val := range f {
		tx := models.DB(rt.Ctx).Begin()
		_, err := models.AssetTreeGetById(rt.Ctx, val)
		if err != nil {
			tx.Rollback()
			ginx.Dangerous(err)
		}

		err = models.UpdateTxTree(tx, map[string]interface{}{"id": val}, map[string]interface{}{"parent_id": parentId, "updated_by": me.Username})
		if err != nil {
			tx.Rollback()
			ginx.Dangerous(err)
		}
		tx.Commit()
	}
	// for _, val := range treeIds {
	// 	err := models.AssetTreeDelById(rt.Ctx, val)
	// 	ginx.Dangerous(err)
	// }

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(nil)
}
