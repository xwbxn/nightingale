// Package models  数据中心
// date : 2023-07-16 08:55
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
// @Success      200  {object}  models.Datacenter
// @Router       /api/n9e/datacenter/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) datacenterGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	datacenter, err := models.DatacenterGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if datacenter == nil {
		ginx.Bomb(404, "No such datacenter")
	}

	ginx.NewRender(c).Data(datacenter, nil)
}

// @Summary      查询数据中心
// @Description  根据条件查询数据中心
// @Tags         数据中心
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.Datacenter
// @Router       /api/n9e/datacenter/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) datacenterGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DatacenterCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.DatacenterGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
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
// @Param        body  body   models.Datacenter true "add datacenter"
// @Success      200
// @Router       /api/n9e/datacenter/ [post]
// @Security     ApiKeyAuth
func (rt *Router) datacenterAdd(c *gin.Context) {
	var f models.Datacenter
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
// @Param        body  body   models.Datacenter true "update datacenter"
// @Success      200
// @Router       /api/n9e/datacenter/ [put]
// @Security     ApiKeyAuth
func (rt *Router) datacenterPut(c *gin.Context) {
	var f models.Datacenter
	ginx.BindJSON(c, &f)

	// 生成一个validate实例
	validate := validator.New()
	errValidate := validate.Struct(f)
	if errValidate != nil {
		ginx.Bomb(http.StatusOK, errValidate.Error())
	}

	old, err := models.DatacenterGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "数据中心不存在!")
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
// @Router       /api/n9e/datacenter/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) datacenterDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	datacenter, err := models.DatacenterGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	//资产校验
	num, err := models.AssetCountByMap(rt.Ctx, map[string]interface{}{"datacenter_id": id})
	ginx.Dangerous(err)
	if num > 0 {
		ginx.Bomb(http.StatusOK, "该数据中心存在资产!")
	}
	//机房校验
	num, err = models.ComputerRoomByMap(rt.Ctx, map[string]interface{}{"idc_location": id})
	ginx.Dangerous(err)
	if num > 0 {
		ginx.Bomb(http.StatusOK, "该数据中心存在机房!")
	}
	//库房校验
	num, err = models.StoreroomManagementByMap(rt.Ctx, map[string]interface{}{"belong_idc": id})
	ginx.Dangerous(err)
	if num > 0 {
		ginx.Bomb(http.StatusOK, "该数据中心存在库房!")
	}
	//报废设备校验
	num, err = models.DeviceScrapByMap(rt.Ctx, map[string]interface{}{"old_datacenter": id})
	ginx.Dangerous(err)
	if num > 0 {
		ginx.Bomb(http.StatusOK, "该数据中心存在报废设备!")
	}

	if datacenter == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(datacenter.Del(rt.Ctx))
}

// @Summary      查询数据中心
// @Description  根据条件查询数据中心
// @Tags         数据中心
// @Accept       json
// @Produce      json
// @Param        id query   int     false  "数据中心ID"
// @Success      200  {array}  models.Datacenter
// @Router       /api/n9e/datacenter/statistics/ [get]
// @Security     ApiKeyAuth
func (rt *Router) datacenterIndexStatistics(c *gin.Context) {
	id := ginx.QueryInt64(c, "id", -1)
	if id == -1 {
		ginx.Bomb(http.StatusOK, "参数错误!")
	}
	ans := make(map[string]interface{}, 0)
	ans["room_num"] = 0
	ans["device_total"] = 0
	ans["nanotube_num"] = 0 //纳管开发未完成
	ans["online_num"] = 0
	ans["cabinet_num"] = 0

	//查询机房
	rooms, err := models.ComputerRoomGetByMap(rt.Ctx, map[string]interface{}{"idc_location": id})
	ginx.Dangerous(err)
	if len(rooms) == 0 {
		ginx.NewRender(c).Data(ans, nil)
	}
	ans["room_num"] = len(rooms)

	//查询机柜
	for _, room := range rooms {
		roomAns := make(map[string]interface{})
		//查询该机房全部设备数
		deviceNum, err := models.AssetCountByMap(rt.Ctx, map[string]interface{}{"equipment_room": room.Id})
		ginx.Dangerous(err)
		roomAns["device_num"] = deviceNum
		cabinetNum, err := models.DeviceCabinetCountByMap(rt.Ctx, map[string]interface{}{"belong_room": room.Id})
		ginx.Dangerous(err)
		roomAns["cabinet_num"] = cabinetNum

	}

	// total, err := models.DatacenterCount(rt.Ctx, query)
	// ginx.Dangerous(err)
	// lst, err := models.DatacenterGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	// ginx.Dangerous(err)

	// ginx.NewRender(c).Data(gin.H{
	// 	"list":  lst,
	// 	"total": total,
	// }, nil)
}
