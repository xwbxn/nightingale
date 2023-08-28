// Package models  字典数据表
// date : 2023-07-21 08:51
// desc : 字典数据表
package router

import (
	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取字典数据表
// @Description  根据主键获取字典数据表
// @Tags         字典数据表
// @Accept       json
// @Produce      json
// @Param        code    path    string  false  "编码"
// @Success      200  {object}  models.DictData
// @Router       /api/n9e/dict-data/{code} [get]
// @Security     ApiKeyAuth
func (rt *Router) dictDataGet(c *gin.Context) {
	dictCode := ginx.UrlParamStr(c, "code")
	dictData, err := models.DictDataGetByDictCode(rt.Ctx, dictCode)
	ginx.Dangerous(err)
	logger.Debug(dictData)
	if dictData == nil {
		ginx.Bomb(404, "No such dict_data")
	}

	ginx.NewRender(c).Data(dictData, nil)
}

// @Summary      查询字典数据表
// @Description  根据条件查询字典数据表
// @Tags         字典数据表
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.DictData
// @Router       /api/n9e/dict-data/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) dictDataGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DictDataCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.DictDataGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建字典数据表
// @Description  创建字典数据表
// @Tags         字典数据表
// @Accept       json
// @Produce      json
// @Param        body  body   []models.DictData true "add dictData"
// @Success      200
// @Router       /api/n9e/dict-data/ [post]
// @Security     ApiKeyAuth
func (rt *Router) dictDataAdd(c *gin.Context) {
	var f []*models.DictData
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	for _, val := range f {
		val.CreatedBy = me.Username
	}

	// 更新模型
	err := models.Add(rt.Ctx, f)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新字典数据表
// @Description  更新字典数据表
// @Tags         字典数据表
// @Accept       json
// @Produce      json
// @Param        dictCode    query    string  true  "编码"
// @Param        body  body   []models.DictData true "update dictData"
// @Success      200
// @Router       /api/n9e/dict-data/ [put]
// @Security     ApiKeyAuth
func (rt *Router) dictDataPut(c *gin.Context) {
	dictCode := ginx.QueryStr(c, "dictCode", "")

	logger.Debug("---------------")
	logger.Debug(dictCode)
	err := models.DictDataDelByDictCode(rt.Ctx, dictCode)
	ginx.Dangerous(err)

	var f []*models.DictData
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	for _, val := range f {
		val.CreatedBy = me.Username
	}

	// 更新模型
	err = models.Add(rt.Ctx, f)

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(err)
}

// @Summary      删除字典数据表
// @Description  根据主键删除字典数据表
// @Tags         字典数据表
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/dict-data/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) dictDataDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	dictData, err := models.DictDataGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if dictData == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(dictData.Del(rt.Ctx))
}
