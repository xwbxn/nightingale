// Package models  数据中心
// date : 2023-07-11 15:49
// desc : 数据中心
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取数据中心
// @Description  根据主键获取数据中心
// @Tags         数据中心
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DataCenter
// @Router       /api/n9e/data-center/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) dataCenterGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	dataCenter, err := models.DataCenterGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if dataCenter == nil {
		ginx.Bomb(404, "No such data_center")
	}

	ginx.NewRender(c).Data(dataCenter, nil)
}

// @Summary      查询数据中心
// @Description  根据条件查询数据中心
// @Tags         数据中心
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.DataCenter
// @Router       /api/n9e/data-center/ [get]
// @Security     ApiKeyAuth
func (rt *Router) dataCenterGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DataCenterCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.DataCenterGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建数据中心
// @Description  创建数据中心
// @Tags         数据中心
// @Accept       json
// @Produce      json
// @Param        body  body   models.DataCenter true "add dataCenter"
// @Success      200
// @Router       /api/n9e/data-center/ [post]
// @Security     ApiKeyAuth
func (rt *Router) dataCenterAdd(c *gin.Context) {
	var f models.DataCenter
	ginx.BindJSON(c, &f)

	// 生成一个validate实例
	validate := validator.New()
	errValidate := validate.Struct(f)
	if errValidate != nil {
		ginx.Bomb(http.StatusOK, errValidate.Error())
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新数据中心
// @Description  更新数据中心
// @Tags         数据中心
// @Accept       json
// @Produce      json
// @Param        body  body   models.DataCenter true "update dataCenter"
// @Success      200
// @Router       /api/n9e/data-center/ [put]
// @Security     ApiKeyAuth
func (rt *Router) dataCenterPut(c *gin.Context) {
	var f models.DataCenter
	ginx.BindJSON(c, &f)

	old, err := models.DataCenterGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "data_center not found")
	}

	// 生成一个validate实例
	validate := validator.New()
	errValidate := validate.Struct(f)
	if errValidate != nil {
		ginx.Bomb(http.StatusOK, errValidate.Error())
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除数据中心
// @Description  根据主键删除数据中心
// @Tags         数据中心
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/data-center/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) dataCenterDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	dataCenter, err := models.DataCenterGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if dataCenter == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(dataCenter.Del(rt.Ctx))
}
