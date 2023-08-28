// Package models  字典类别表
// date : 2023-07-21 08:48
// desc : 字典类别表
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取字典类别表
// @Description  根据主键获取字典类别表
// @Tags         字典类别表
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DictType
// @Router       /api/n9e/dict-type/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) dictTypeGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	dictType, err := models.DictTypeGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if dictType == nil {
		ginx.Bomb(404, "No such dict_type")
	}

	ginx.NewRender(c).Data(dictType, nil)
}

// @Summary      查询字典类别表
// @Description  根据条件查询字典类别表
// @Tags         字典类别表
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.DictType
// @Router       /api/n9e/dict-type/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) dictTypeGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DictTypeCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.DictTypeGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建字典类别表
// @Description  创建字典类别表
// @Tags         字典类别表
// @Accept       json
// @Produce      json
// @Param        body  body   models.DictType true "add dictType"
// @Success      200
// @Router       /api/n9e/dict-type/ [post]
// @Security     ApiKeyAuth
func (rt *Router) dictTypeAdd(c *gin.Context) {
	var f models.DictType
	ginx.BindJSON(c, &f)

	//校验缓存是否存在编码
	// typeVal := rt.Redis.Get(rt.Ctx.Ctx, f.DictCode)

	// rt.Redis.Set(rt.Ctx.Ctx, f.DictCode, f)

	//校验编码是否存在
	exist, err := models.DictTypeGetByDictCode(rt.Ctx, f.DictCode)
	logger.Debug("-========================")
	logger.Debug(exist)
	logger.Debug(f.DictCode)
	ginx.Dangerous(err)
	if exist {
		ginx.Bomb(http.StatusOK, "dictCode exist")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err = f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新字典类别表
// @Description  更新字典类别表
// @Tags         字典类别表
// @Accept       json
// @Produce      json
// @Param        body  body   models.DictType true "update dictType"
// @Success      200
// @Router       /api/n9e/dict-type/ [put]
// @Security     ApiKeyAuth
func (rt *Router) dictTypePut(c *gin.Context) {
	var f models.DictType
	ginx.BindJSON(c, &f)

	old, err := models.DictTypeGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "dict_type not found")
	}
	//校验编码是否存在
	exist, err := models.DictTypeGetByDictCode(rt.Ctx, f.DictCode)
	ginx.Dangerous(err)
	if exist {
		ginx.Bomb(http.StatusOK, "dictCode exist")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除字典类别表
// @Description  根据主键删除字典类别表
// @Tags         字典类别表
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/dict-type/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) dictTypeDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	dictType, err := models.DictTypeGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if dictType == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	//删除字典数据
	err = models.DictDataDelByDictCode(rt.Ctx, dictType.DictCode)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	ginx.NewRender(c).Message(dictType.Del(rt.Ctx))
}
