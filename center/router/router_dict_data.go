// Package models  字典数据表
// date : 2023-07-21 08:51
// desc : 字典数据表
package router

import (
	"net/http"
	"strings"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
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
	typeCodeTemp := ginx.UrlParamStr(c, "code")
	typeCodeLst := strings.Split(typeCodeTemp, ",")
	var dictData []models.DictData
	var err error
	if len(typeCodeLst) == 1 {
		m := make(map[string]interface{})
		m["type_code"] = typeCodeLst[0]
		dictData, err = models.DictDataGetByMap(rt.Ctx, m)
		ginx.Dangerous(err)
	} else {
		dictData, err = models.DictDataGetByTypeCodes(rt.Ctx, typeCodeLst)
		ginx.Dangerous(err)
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

// @Summary      查询字典数据表
// @Description  根据条件查询字典数据表
// @Tags         字典数据表
// @Accept       json
// @Produce      json
// @Param        code query   string     true  "字典类型编码"
// @Param        query query   string  false  "查询条件"
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Success      200  {array}  models.DictData
// @Router       /api/n9e/dict-data/exp/ [get]
// @Security     ApiKeyAuth
func (rt *Router) dictDataGetExp(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	typeCode := ginx.QueryStr(c, "code", "")
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DictDataExpCount(rt.Ctx, query, typeCode)
	ginx.Dangerous(err)
	lst, err := models.DictDataGetExp(rt.Ctx, query, typeCode, limit, (page-1)*limit)
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
	err := models.DictDataBatchAdd(rt.Ctx, f)
	ginx.NewRender(c).Message(err)
}

// @Summary      创建字典数据表(单条)
// @Description  创建字典数据表(单条)
// @Tags         字典数据表
// @Accept       json
// @Produce      json
// @Param        body  body   models.DictData true "add dictData"
// @Success      200
// @Router       /api/n9e/dict-data/one/ [post]
// @Security     ApiKeyAuth
func (rt *Router) dictDataOneAdd(c *gin.Context) {
	var f models.DictData
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新字典数据表
// @Description  更新字典数据表
// @Tags         字典数据表
// @Accept       json
// @Produce      json
// @Param        typeCode    query    string  true  "编码"
// @Param        body  body   []models.DictData true "update dictData"
// @Success      200
// @Router       /api/n9e/dict-data/ [put]
// @Security     ApiKeyAuth
func (rt *Router) dictDataPut(c *gin.Context) {
	typeCode := ginx.QueryStr(c, "typeCode", "")

	m := make(map[string]interface{})
	m["type_code"] = typeCode

	err := models.DictDataDelByMap(rt.Ctx, m)
	ginx.Dangerous(err)

	var f []*models.DictData
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	for _, val := range f {
		val.CreatedBy = me.Username
	}

	// 更新模型
	err = models.DictDataBatchAdd(rt.Ctx, f)

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(err)
}

// @Summary      更新字典数据表
// @Description  更新字典数据表
// @Tags         字典数据表
// @Accept       json
// @Produce      json
// @Param        body  body   models.DictData true "update dictData"
// @Success      200
// @Router       /api/n9e/dict-data/one [put]
// @Security     ApiKeyAuth
func (rt *Router) dictDataOnePut(c *gin.Context) {

	var f models.DictData
	ginx.BindJSON(c, &f)

	old, err := models.DictDataGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if (&models.DictData{} == old) {
		ginx.Bomb(http.StatusOK, "数据不存在")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 更新模型
	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
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

// @Summary      批量删除资产扩展字段
// @Description  根据主键删除资产扩展字段
// @Tags         字典数据表
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "delete dictData"
// @Success      200
// @Router       /api/n9e/dict-data/asset-batch/ [post]
// @Security     ApiKeyAuth
func (rt *Router) dictDataBtachDel(c *gin.Context) {
	ids := make([]int64, 0)
	ginx.BindJSON(c, &ids)

	//查询字典数据
	dictDatas, err := models.DictDataGetBatchById(rt.Ctx, ids)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)
	if len(dictDatas) == 0 {
		ginx.Bomb(http.StatusOK, "dict_data not found")
	}

	tx := models.DB(rt.Ctx).Begin()
	for _, dictData := range dictDatas {
		m := make(map[string]interface{})
		m["property_category"] = dictData.TypeCode
		m["property_name"] = dictData.DictKey
		err = models.MapTxDel(tx, m)
		ginx.Dangerous(err)
	}
	err = models.DictDataDelByIds(tx, ids)
	ginx.Dangerous(err)
	tx.Commit()

	ginx.NewRender(c).Message(err)
}
