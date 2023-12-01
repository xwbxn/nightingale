// Package models  接口管理
// date : 2023-10-20 16:48
// desc : 接口管理
package router

import (
	"fmt"
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取接口管理
// @Description  根据主键获取接口管理
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.ApiService
// @Router       /api/n9e/api-service/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) apiServiceGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	apiService, err := models.ApiServiceGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if apiService == nil {
		ginx.Bomb(404, "No such api_service")
	}

	ginx.NewRender(c).Data(apiService, nil)
}

// @Summary      查询接口管理
// @Description  根据条件查询接口管理
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.ApiService
// @Router       /api/n9e/api-service/ [get]
// @Security     ApiKeyAuth
func (rt *Router) apiServiceGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.ApiServiceCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.ApiServiceGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建接口管理
// @Description  创建接口管理
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        body  body   models.ApiService true "add apiService"
// @Success      200
// @Router       /api/n9e/api-service/ [post]
// @Security     ApiKeyAuth
func (rt *Router) apiServiceAdd(c *gin.Context) {
	var f models.ApiService
	ginx.BindJSON(c, &f)

	if f.IsDangerous() {
		ginx.Bomb(400, "不合法的脚本")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新接口管理
// @Description  更新接口管理
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        body  body   models.ApiService true "update apiService"
// @Success      200
// @Router       /api/n9e/api-service/ [put]
// @Security     ApiKeyAuth
func (rt *Router) apiServicePut(c *gin.Context) {
	var f models.ApiService
	ginx.BindJSON(c, &f)

	if f.IsDangerous() {
		ginx.Bomb(400, "不合法的脚本")
	}

	old, err := models.ApiServiceGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "api_service not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除接口管理
// @Description  根据主键删除接口管理
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/api-service/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) apiServiceDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	apiService, err := models.ApiServiceGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if apiService == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(apiService.Del(rt.Ctx))
}

type apiServiceFE struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// @Summary      查询接口列表
// @Description  根据条件查询接口管理
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Success      200  {array}  apiServiceFE
// @Router       /api/n9e/api-service/options [get]
// @Security     ApiKeyAuth
func (rt *Router) apiServiceGetsOptions(c *gin.Context) {
	lst, err := models.ApiServiceGetsAll(rt.Ctx)
	ginx.Dangerous(err)

	data := make([]apiServiceFE, len(lst))
	for i, v := range lst {
		data[i].Code = fmt.Sprintf("/api/n9e/api-service/%d/execute", v.Id)
		data[i].Name = v.Name
	}

	ginx.NewRender(c).Data(data, nil)
}

// @Summary      执行接口
// @Description  执行接口
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "接口id"
// @Success      200
// @Router       /api/n9e/api-service/{id}/execute [get]
// @Security     ApiKeyAuth
func (rt *Router) apiServiceExecute(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	apiService, err := models.ApiServiceGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if apiService == nil {
		ginx.Bomb(404, "No such api_service")
	}

	api := rt.PromClients.GetCli(1)
	data, err := apiService.Execute(rt.Ctx, api)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(data, nil)
}
