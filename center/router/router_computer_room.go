// Package models  机房信息
// date : 2023-07-16 09:04
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
// @Success      200  {object}  models.ComputerRoom
// @Router       /api/n9e/computer-room/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) computerRoomGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	computerRoom, err := models.ComputerRoomGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if computerRoom == nil {
		ginx.Bomb(404, "No such computer_room")
	}

	ginx.NewRender(c).Data(computerRoom, nil)
}

// @Summary      获取机房信息
// @Description  根据主键获取机房信息
// @Tags         机房信息
// @Accept       json
// @Produce      json
// @Param        idcLocation    query    int  true  "数据中心Id"
// @Success      200  {object}  []models.ComputerRoomNameVo
// @Router       /api/n9e/computer-room/datacenterId/ [get]
// @Security     ApiKeyAuth
func (rt *Router) computerRoomNameGet(c *gin.Context) {
	IdcLocation := ginx.QueryInt(c, "idcLocation", -1)
	computerRoom, err := models.ComputerRoomNameGetByIdc(rt.Ctx, int64(IdcLocation))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(computerRoom, nil)
}

// @Summary      查询机房信息
// @Description  根据条件查询机房信息
// @Tags         机房信息
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.ComputerRoom
// @Router       /api/n9e/computer-room/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) computerRoomGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.ComputerRoomCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.ComputerRoomGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
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
// @Param        body  body   models.ComputerRoom true "add computerRoom"
// @Success      200
// @Router       /api/n9e/computer-room/ [post]
// @Security     ApiKeyAuth
func (rt *Router) computerRoomAdd(c *gin.Context) {
	var f models.ComputerRoom
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
// @Param        body  body   models.ComputerRoom true "update computerRoom"
// @Success      200
// @Router       /api/n9e/computer-room/ [put]
// @Security     ApiKeyAuth
func (rt *Router) computerRoomPut(c *gin.Context) {
	var f models.ComputerRoom
	ginx.BindJSON(c, &f)

	// 生成一个validate实例
	validate := validator.New()
	errValidate := validate.Struct(f)
	if errValidate != nil {
		ginx.Bomb(http.StatusOK, errValidate.Error())
	}

	old, err := models.ComputerRoomGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "computer_room not found")
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
// @Router       /api/n9e/computer-room/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) computerRoomDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	computerRoom, err := models.ComputerRoomGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if computerRoom == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(computerRoom.Del(rt.Ctx))
}
