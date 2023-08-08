// Package models  设备类型
// date : 2023-07-08 11:49
// desc : 设备类型
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取设备类型
// @Description  根据主键获取设备类型
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DeviceType
// @Router       /api/n9e/device-type/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceType, err := models.DeviceTypeGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if deviceType == nil {
		ginx.Bomb(404, "No such device_type")
	}

	ginx.NewRender(c).Data(deviceType, nil)
}

// @Summary      查询设备类型
// @Description  根据条件查询设备类型
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.DeviceType
// @Router       /api/n9e/device-type/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DeviceTypeCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.DeviceTypeGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建设备类型
// @Description  创建设备类型
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceType true "add deviceType"
// @Success      200
// @Router       /api/n9e/device-type/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeAdd(c *gin.Context) {
	var f models.DeviceType
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新设备类型
// @Description  更新设备类型
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceType true "update deviceType"
// @Success      200
// @Router       /api/n9e/device-type/ [put]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypePut(c *gin.Context) {
	var f models.DeviceType
	ginx.BindJSON(c, &f)

	old, err := models.DeviceTypeGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "device_type not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除设备类型
// @Description  根据主键删除设备类型
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/device-type/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceType, err := models.DeviceTypeGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if deviceType == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(deviceType.Del(rt.Ctx))
}
