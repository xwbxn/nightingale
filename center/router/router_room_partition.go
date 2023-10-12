// Package models  机房分区表
// date : 2023-9-10 10:39
// desc : 机房分区表
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取机房分区表
// @Description  根据主键获取机房分区表
// @Tags         机房分区表
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.RoomPartition
// @Router       /api/n9e/room-partition/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) roomPartitionGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	roomPartition, err := models.RoomPartitionGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if roomPartition == nil {
		ginx.Bomb(404, "No such room_partition")
	}

	ginx.NewRender(c).Data(roomPartition, nil)
}

// @Summary      根据room_id查询机房分区表
// @Description  根据room_id查询机房分区表
// @Tags         机房分区表
// @Accept       json
// @Produce      json
// @Param        room_id query   int     true  "机房ID"
// @Success      200  {array}  models.RoomPartition
// @Router       /api/n9e/room-partition/room-id/ [get]
// @Security     ApiKeyAuth
func (rt *Router) roomPartitionGets(c *gin.Context) {
	roomId := ginx.QueryInt64(c, "room_id", -1)

	lst, err := models.RoomPartitionGetBymap(rt.Ctx, map[string]interface{}{"room_id": roomId})
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(lst, nil)
}

// @Summary      查询机房分区表
// @Description  根据条件查询机房分区表
// @Tags         机房分区表
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.RoomPartition
// @Router       /api/n9e/room-partition/ [get]
// @Security     ApiKeyAuth
// func (rt *Router) roomPartitionGets(c *gin.Context) {
// 	limit := ginx.QueryInt(c, "limit", 20)
// 	query := ginx.QueryStr(c, "query", "")

// 	total, err := models.RoomPartitionCount(rt.Ctx, query)
// 	ginx.Dangerous(err)
// 	lst, err := models.RoomPartitionGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
// 	ginx.Dangerous(err)

// 	ginx.NewRender(c).Data(gin.H{
// 		"list":  lst,
// 		"total": total,
// 	}, nil)
// }

// @Summary      创建机房分区表
// @Description  创建机房分区表
// @Tags         机房分区表
// @Accept       json
// @Produce      json
// @Param        body  body   models.RoomPartition true "add roomPartition"
// @Success      200
// @Router       /api/n9e/room-partition/ [post]
// @Security     ApiKeyAuth
func (rt *Router) roomPartitionAdd(c *gin.Context) {
	var f models.RoomPartition
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新机房分区表
// @Description  更新机房分区表
// @Tags         机房分区表
// @Accept       json
// @Produce      json
// @Param        body  body   models.RoomPartition true "update roomPartition"
// @Success      200
// @Router       /api/n9e/room-partition/ [put]
// @Security     ApiKeyAuth
func (rt *Router) roomPartitionPut(c *gin.Context) {
	var f models.RoomPartition
	ginx.BindJSON(c, &f)

	old, err := models.RoomPartitionGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "room_partition not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除机房分区表
// @Description  根据主键删除机房分区表
// @Tags         机房分区表
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/room-partition/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) roomPartitionDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	roomPartition, err := models.RoomPartitionGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if roomPartition == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(roomPartition.Del(rt.Ctx))
}
