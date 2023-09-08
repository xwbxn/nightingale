// Package models  资产详情
// date : 2023-07-21 08:45
// desc : 资产详情
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取资产详情
// @Description  根据主键获取资产详情
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.AssetBasic
// @Router       /api/n9e/asset-basic/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetBasic, err := models.AssetBasicGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if assetBasic == nil {
		ginx.Bomb(404, "No such asset_basic")
	}

	ginx.NewRender(c).Data(assetBasic, nil)
}

// @Summary      查询资产详情
// @Description  根据条件查询资产详情
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.AssetBasic
// @Router       /api/n9e/asset-basic/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicGets(c *gin.Context) {

	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.AssetBasicCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.AssetBasicGets(rt.Ctx, query, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      查询资产详情（资产树）
// @Description  根据条件查询资产详情（资产树）
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        body  body   models.AssetBasicFindVo true "add AssetBasicFindVo"
// @Success      200  {array}  []models.AssetBasicDetailsVo
// @Router       /api/n9e/asset-basic/list/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicGetsByTree(c *gin.Context) {

	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)

	m := make(map[string]interface{})
	ginx.BindJSON(c, &m)

	total, err := models.AssetCountByMap(rt.Ctx, m)
	ginx.Dangerous(err)
	lst, err := models.AssetBasicGetsByMap(rt.Ctx, m, limit, (page-1)*limit)
	ginx.Dangerous(err)

	for index := range lst {
		//回填设备类型
		deviceType, _ := models.DeviceTypeGetById(rt.Ctx, lst[index].DeviceType)
		lst[index].DeviceTypeName = deviceType.Name
		//回填设备厂商
		deviceProducer, _ := models.DeviceProducerGetById(rt.Ctx, lst[index].DeviceProducer)
		lst[index].DeviceProducerName = deviceProducer.Alias
		//回填设备型号
		deviceModel, _ := models.DeviceModelGetById(rt.Ctx, lst[index].DeviceModel)
		lst[index].DeviceModelName = deviceModel.Name
		//TODO 回填所属组织(不清楚)

		//回填所在机房
		if lst[index].EquipmentRoom != 0 {
			computerRoom, _ := models.ComputerRoomGetById(rt.Ctx, lst[index].EquipmentRoom)
			lst[index].RoomName = computerRoom.RoomName
		}

		//回填所在机柜
		if lst[index].OwningCabinet != 0 {
			deviceCabinet, _ := models.DeviceCabinetGetById(rt.Ctx, lst[index].OwningCabinet)
			lst[index].CabinetName = deviceCabinet.CabinetName
		}
	}

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建资产详情
// @Description  创建资产详情
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetBasic true "add assetBasic"
// @Success      200
// @Router       /api/n9e/asset-basic/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicAdd(c *gin.Context) {
	var f models.AssetBasic
	ginx.BindJSON(c, &f)
	if f.ManagementIp == "" && f.SerialNumber == "" {
		ginx.Bomb(http.StatusOK, "管理IP和序列号不能同时为空!")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	id, err := f.Add(rt.Ctx)
	ginx.NewRender(c).Data(id, err)
}

// @Summary      更新资产详情
// @Description  更新资产详情
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetBasic true "update assetBasic"
// @Success      200
// @Router       /api/n9e/asset-basic/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicPut(c *gin.Context) {
	var f models.AssetBasic
	ginx.BindJSON(c, &f)

	old, err := models.AssetBasicGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "asset_basic not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除资产详情
// @Description  根据主键删除资产详情
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/asset-basic/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetBasic, err := models.AssetBasicGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if assetBasic == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(assetBasic.Del(rt.Ctx))
}

// @Summary      查询资产详情(表名、字段)
// @Description  根据表名查询资产详情
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        query    query   string     true  "表名"
// @Param        body  body   []string true "add array"
// @Success      200  {array} map[string]interface{}
// @Router       /api/n9e/asset-basic/table/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetFieldGetsByTableName(c *gin.Context) {

	var a []string
	ginx.BindJSON(c, &a)

	tableName := ginx.QueryStr(c, "query", "")
	//校验表名是否存在
	exist := models.HasTableByName(rt.Ctx, tableName)
	if !exist {
		ginx.Bomb(http.StatusOK, "table_name not found")
	}
	//校验该表中传入字段是否存在
	for _, val := range a {
		count, err := models.HasTableFieldByName(rt.Ctx, tableName, val)
		ginx.Dangerous(err)
		if count == 0 {
			ginx.Bomb(http.StatusOK, "field not found")
		}
	}

	//查询数据
	data, err := models.TableGetsByName(rt.Ctx, tableName, a)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  data,
		"total": len(data),
	}, nil)
}
