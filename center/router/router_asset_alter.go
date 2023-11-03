// Package models  资产变更
// date : 2023-08-04 14:50
// desc : 资产变更
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取资产变更
// @Description  根据主键获取资产变更
// @Tags         资产变更
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.AssetAlter
// @Router       /api/n9e/asset-alter/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) assetAlterGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetAlter, err := models.AssetAlterGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if assetAlter == nil {
		ginx.Bomb(404, "No such asset_alter")
	}

	ginx.NewRender(c).Data(assetAlter, nil)
}

// @Summary      获取资产变更列表
// @Description  根据资产ID获取资产变更列表
// @Tags         资产变更
// @Accept       json
// @Produce      json
// @Param        asset    query    int  true  "资产ID"
// @Param        start    query    int  false  "变更开始日期"
// @Param        end    query    int  false  "变更结束日期"
// @Param        event    query    string  false  "变更事项"
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Success      200  {array}  models.AssetAlterVo
// @Router       /api/n9e/asset-alter/asset [get]
// @Security     ApiKeyAuth
func (rt *Router) assetAlterGetByAssetId(c *gin.Context) {
	assetId := ginx.QueryInt64(c, "asset", -1)
	limit := ginx.QueryInt(c, "limit", 20)
	page := ginx.QueryInt(c, "page", 1)
	start := ginx.QueryInt64(c, "start", -1)
	end := ginx.QueryInt64(c, "end", -1)
	event := ginx.QueryStr(c, "event", "")

	m := make(map[string]interface{})
	if start != -1 {
		m["start"] = start

		//开始时间不大于结束时间
		if end < start && end != -1 {
			ginx.Bomb(http.StatusOK, "Wrong date range selection")
		}
	}
	if end != -1 {
		m["end"] = end
	}
	m["asset_id"] = assetId
	if event != "" {
		m["alter_event_key"] = event
	}

	//查询资产详情
	assetBasic, err := models.AssetBasicGetById[models.AssetBasic](rt.Ctx, assetId)
	ginx.Dangerous(err)
	if assetBasic == nil {
		ginx.Bomb(404, "No such asset")
	}

	total, err := models.AssetAlterCountByMap(rt.Ctx, m)
	ginx.Dangerous(err)

	assetAlterVo, err := models.AssetAlterGetByMap(rt.Ctx, m, limit, (page-1)*limit)
	ginx.Dangerous(err)

	for index := range assetAlterVo {
		assetAlterVo[index].ManagementIp = assetBasic.ManagementIp
		assetAlterVo[index].SerialNumber = assetBasic.SerialNumber
		//回填设备类型
		deviceType, _ := models.DeviceTypeGetById(rt.Ctx, assetBasic.DeviceType)
		assetAlterVo[index].DeviceTypeName = deviceType.Name
		//回填设备厂商
		deviceProducer, _ := models.DeviceProducerGetById(rt.Ctx, assetBasic.DeviceProducer)
		assetAlterVo[index].DeviceProducerName = deviceProducer.Alias
		//回填设备型号
		deviceModel, _ := models.DeviceModelGetById(rt.Ctx, assetBasic.DeviceModel)
		assetAlterVo[index].DeviceModelName = deviceModel.Name
		//回填所在机房
		if assetBasic.EquipmentRoom != 0 {
			computerRoom, _ := models.ComputerRoomGetById(rt.Ctx, assetBasic.EquipmentRoom)
			assetAlterVo[index].RoomName = computerRoom.RoomName
		}

		//回填所在机柜
		if assetBasic.OwningCabinet != 0 {
			deviceCabinet, _ := models.DeviceCabinetGetById(rt.Ctx, assetBasic.OwningCabinet)
			assetAlterVo[index].CabinetName = deviceCabinet.CabinetName
		}
		assetAlterVo[index].UNumber = assetBasic.UNumber
	}
	ginx.NewRender(c).Data(gin.H{
		"list":  assetAlterVo,
		"total": total,
	}, nil)
}

// @Summary      查询资产变更
// @Description  根据条件查询资产变更
// @Tags         资产变更
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.AssetAlter
// @Router       /api/n9e/asset-alter/list [get]
// @Security     ApiKeyAuth
func (rt *Router) assetAlterGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.AssetAlterCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.AssetAlterGets(rt.Ctx, query, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建资产变更
// @Description  创建资产变更
// @Tags         资产变更
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetAlter true "add assetAlter"
// @Success      200
// @Router       /api/n9e/asset-alter/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetAlterAdd(c *gin.Context) {
	var f models.AssetAlter
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新资产变更
// @Description  更新资产变更
// @Tags         资产变更
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetAlter true "update assetAlter"
// @Success      200
// @Router       /api/n9e/asset-alter/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetAlterPut(c *gin.Context) {
	var f models.AssetAlter
	ginx.BindJSON(c, &f)

	old, err := models.AssetAlterGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "asset_alter not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除资产变更
// @Description  根据主键删除资产变更
// @Tags         资产变更
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/asset-alter/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) assetAlterDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetAlter, err := models.AssetAlterGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if assetAlter == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(assetAlter.Del(rt.Ctx))
}

// @Summary      EXCEL导出资产变更记录
// @Description  EXCEL导出资产变更记录
// @Tags         资产变更
// @Accept       multipart/form-data
// @Produce      application/msexcel
// @Param        query query   int  true  "资产ID"
// @Success      200
// @Router       /api/n9e/asset-alter/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) downloadAssetAlter(c *gin.Context) {

	assetId := ginx.QueryInt64(c, "query", -1)
	//获取数据
	where := make(map[string]interface{})
	where["ASSET_ID"] = assetId
	lst, err := models.AssetAlterGetByMap(rt.Ctx, where, -1, -1)
	ginx.Dangerous(err)

	datas := make([]interface{}, 0)
	if len(lst) > 0 {
		for _, v := range lst {
			datas = append(datas, v)

		}
	}
	excels.NewMyExcel("资产变更数据").ExportDataInfo(datas, "cn", rt.Ctx, c)

}
