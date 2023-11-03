// Package models  设备类型表单配置表
// date : 2023-08-04 08:56
// desc : 设备类型表单配置表
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取设备类型表单配置表
// @Description  根据主键获取设备类型表单配置表
// @Tags         设备类型表单配置表
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DeviceTypeConfig
// @Router       /api/n9e/device-type_config/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeConfigGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceTypeConfig, err := models.DeviceTypeConfigGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if deviceTypeConfig == nil {
		ginx.Bomb(404, "No such device_type_config")
	}

	ginx.NewRender(c).Data(deviceTypeConfig, nil)
}

// @Summary      查询设备类型表单配置表
// @Description  根据条件查询设备类型表单配置表
// @Tags         设备类型表单配置表
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.DeviceTypeConfig
// @Router       /api/n9e/device-type_config/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeConfigGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DeviceTypeConfigCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.DeviceTypeConfigGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建设备类型表单配置表
// @Description  创建设备类型表单配置表
// @Tags         设备类型表单配置表
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceTypeConfig true "add deviceTypeConfig"
// @Success      200
// @Router       /api/n9e/device-type_config/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeConfigAdd(c *gin.Context) {
	var f models.DeviceTypeConfig
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新设备类型表单配置表
// @Description  更新设备类型表单配置表
// @Tags         设备类型表单配置表
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceTypeConfig true "update deviceTypeConfig"
// @Success      200
// @Router       /api/n9e/device-type_config/ [put]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeConfigPut(c *gin.Context) {
	var f models.DeviceTypeConfig
	ginx.BindJSON(c, &f)

	old, err := models.DeviceTypeConfigGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "device_type_config not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除设备类型表单配置表
// @Description  根据主键删除设备类型表单配置表
// @Tags         设备类型表单配置表
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/device-type_config/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeConfigDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceTypeConfig, err := models.DeviceTypeConfigGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if deviceTypeConfig == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(deviceTypeConfig.Del(rt.Ctx))
}
