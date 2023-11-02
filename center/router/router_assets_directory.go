// Package models  资产目录
// date : 2023-9-19 11:35
// desc : 资产目录
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      生成资产目录树
// @Description  生成资产目录树
// @Tags         资产目录
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.AssetsDirTree
// @Router       /api/n9e/asset-directory/tree [get]
// @Security     ApiKeyAuth
func (rt *Router) assetsDirTreeGet(c *gin.Context) {
	var assetsDirTree *models.AssetsDirTree
	treeLst := make([]*models.AssetsDirTree, 0)
	lst, err := models.AssetsDirectoryGetsMap(rt.Ctx, map[string]interface{}{"sort": -1})
	ginx.Dangerous(err)
	if len(lst) == 0 {
		ginx.NewRender(c).Data(treeLst, nil)
	}
	assetsDirTree, err = models.BuildDirTree(rt.Ctx, lst[0].Id)
	ginx.Dangerous(err)

	assetsDirTree, err = models.AssetsDirCount(rt.Ctx, assetsDirTree)
	ginx.Dangerous(err)
	treeLst = assetsDirTree.Children
	logger.Debug(treeLst)

	ginx.NewRender(c).Data(treeLst, nil)
}

// @Summary      创建资产目录
// @Description  创建资产目录
// @Tags         资产目录
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetsDirectory true "add assetsDirectory"
// @Success      200
// @Router       /api/n9e/asset-directory [post]
// @Security     ApiKeyAuth
func (rt *Router) assetsDirectoryAdd(c *gin.Context) {
	var f models.AssetsDirectory
	ginx.BindJSON(c, &f)
	logger.Debug(f)

	assetsDirTree, err := models.BuildDirTree(rt.Ctx, f.ParentId)
	ginx.Dangerous(err)

	if len(assetsDirTree.Children) != 0 {
		f.Sort = assetsDirTree.Children[len(assetsDirTree.Children)-1].Id
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err = f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新资产目录
// @Description  更新资产目录
// @Tags         资产目录
// @Accept       json
// @Produce      json
// @Param        id    query    int64  true  "主键"
// @Param        name    query    string  true  "名称"
// @Success      200
// @Router       /api/n9e/asset-directory/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetsDirectoryPut(c *gin.Context) {
	id := ginx.QueryInt64(c, "id", -1)
	name := ginx.QueryStr(c, "name", "")
	if id == -1 || name == "" {
		ginx.Bomb(http.StatusOK, "参数错误")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	tx := models.DB(rt.Ctx).Begin()
	err := models.AssetsDirectoryUpdate(tx, map[string]interface{}{"id": id}, map[string]interface{}{"name": name, "updated_by": me.Username})
	tx.Commit()
	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(err)
}

// @Summary      删除资产目录
// @Description  根据主键删除资产目录
// @Tags         资产目录
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/asset-directory/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) assetsDirectoryDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")

	lst, err := models.AssetsGetsMap(rt.Ctx, map[string]interface{}{"directory_id": id})
	ginx.Dangerous(err)

	if len(lst) > 0 {
		ginx.Bomb(http.StatusOK, "该目录下存在资产，不可删除")
	}

	assetsDirectory, err := models.AssetsDirectoryGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)

	tx := models.DB(rt.Ctx).Begin()
	err = models.AssetsDirectoryUpdate(tx, map[string]interface{}{"sort": id}, map[string]interface{}{"sort": assetsDirectory.Sort, "updated_by": me.Username})
	ginx.Dangerous(err)
	err = assetsDirectory.Del(tx)
	tx.Commit()

	ginx.NewRender(c).Message(err)
}

// @Summary      移动资产目录
// @Description  移动资产目录
// @Tags         资产目录
// @Accept       json
// @Produce      json
// @Param        type    query    string  true  "类型(up/down)"
// @Param        id    query    int64  true  "主键"
// @Success      200
// @Router       /api/n9e/asset-directory/move [get]
// @Security     ApiKeyAuth
func (rt *Router) assetsDirectoryMove(c *gin.Context) {
	id := ginx.QueryInt64(c, "id", -1)
	moveType := ginx.QueryStr(c, "type", "")
	if id == -1 || moveType == "" {
		ginx.Bomb(http.StatusOK, "参数错误")
	}
	assetsDirectory, err := models.AssetsDirectoryGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	assetsDirTree := &models.AssetsDirTree{
		Id:       assetsDirectory.Id,
		Name:     assetsDirectory.Name,
		ParentId: assetsDirectory.ParentId,
		Sort:     assetsDirectory.Sort,
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	var up *models.AssetsDirectory
	var down, mid []*models.AssetsDirTree

	tx := models.DB(rt.Ctx).Begin()

	if moveType == "up" {
		if assetsDirectory.Sort == 0 {
			ginx.Bomb(http.StatusOK, "已是首位，不可上移")
		}
		up, err = models.AssetsDirectoryGetById(rt.Ctx, assetsDirectory.Sort)
		ginx.Dangerous(err)
		mid = append(mid, assetsDirTree)
		down, err = models.AssetsDirectoryGetsMap(rt.Ctx, map[string]interface{}{"sort": assetsDirectory.Id})
		ginx.Dangerous(err)

	} else if moveType == "down" {
		logger.Debug("down")
		mid, err = models.AssetsDirectoryGetsMap(rt.Ctx, map[string]interface{}{"sort": assetsDirectory.Id})
		ginx.Dangerous(err)
		if len(mid) == 0 {
			ginx.Bomb(http.StatusOK, "已是末位，不可下移")
		}
		logger.Debug(len(mid))
		up = assetsDirectory
		down, err = models.AssetsDirectoryGetsMap(rt.Ctx, map[string]interface{}{"sort": mid[0].Id})
		ginx.Dangerous(err)
	}
	err = models.AssetsDirectoryUpdate(tx, map[string]interface{}{"id": mid[0].Id}, map[string]interface{}{"sort": up.Sort, "updated_by": me.Username})
	ginx.Dangerous(err)
	err = models.AssetsDirectoryUpdate(tx, map[string]interface{}{"id": up.Id}, map[string]interface{}{"sort": mid[0].Id, "updated_by": me.Username})
	ginx.Dangerous(err)
	if len(down) != 0 {
		err = models.AssetsDirectoryUpdate(tx, map[string]interface{}{"id": down[0].Id}, map[string]interface{}{"sort": up.Id, "updated_by": me.Username})
		ginx.Dangerous(err)
	} else {
		down = nil
	}
	tx.Commit()
	ginx.NewRender(c).Message(err)
}
