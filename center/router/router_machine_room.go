// Package models  机房信息
// date : 2023-07-11 16:11
// desc : 机房信息
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取机房信息
// @Description  根据主键获取机房信息
// @Tags         机房信息
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.MachineRoom
// @Router       /api/n9e/machine-room/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) machineRoomGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	machineRoom, err := models.MachineRoomGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if machineRoom == nil {
		ginx.Bomb(404, "No such machine_room")
	}

	ginx.NewRender(c).Data(machineRoom, nil)
}

// @Summary      查询机房信息
// @Description  根据条件查询机房信息
// @Tags         机房信息
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.MachineRoom
// @Router       /api/n9e/machine-room/ [get]
// @Security     ApiKeyAuth
func (rt *Router) machineRoomGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.MachineRoomCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.MachineRoomGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建机房信息
// @Description  创建机房信息
// @Tags         机房信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.MachineRoom true "add machineRoom"
// @Success      200
// @Router       /api/n9e/machine-room/ [post]
// @Security     ApiKeyAuth
func (rt *Router) machineRoomAdd(c *gin.Context) {
	var f models.MachineRoom
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

// @Summary      更新机房信息
// @Description  更新机房信息
// @Tags         机房信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.MachineRoom true "update machineRoom"
// @Success      200
// @Router       /api/n9e/machine-room/ [put]
// @Security     ApiKeyAuth
func (rt *Router) machineRoomPut(c *gin.Context) {
	var f models.MachineRoom
	ginx.BindJSON(c, &f)

	// 生成一个validate实例
	validate := validator.New()
	errValidate := validate.Struct(f)
	if errValidate != nil {
		ginx.Bomb(http.StatusOK, errValidate.Error())
	}

	old, err := models.MachineRoomGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "machine_room not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除机房信息
// @Description  根据主键删除机房信息
// @Tags         机房信息
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/machine-room/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) machineRoomDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	machineRoom, err := models.MachineRoomGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if machineRoom == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(machineRoom.Del(rt.Ctx))
}
