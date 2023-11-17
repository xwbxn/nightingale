// Package models  设备报废
// date : 2023-9-08 14:43
// desc : 设备报废
package router

import (
	"net/http"
	"time"

	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取设备报废
// @Description  根据主键获取设备报废
// @Tags         设备报废
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DeviceScrap
// @Router       /api/n9e/device-scrap/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceScrapGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceScrap, err := models.DeviceScrapGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if deviceScrap == nil {
		ginx.Bomb(404, "No such device_scrap")
	}

	ginx.NewRender(c).Data(deviceScrap, nil)
}

// @Summary      查询设备报废
// @Description  根据条件查询设备报废
// @Tags         设备报废
// @Accept       json
// @Produce      json
// @Param        filter query   string     false  "查询范围"
// @Param        datacenter_id query   int  false  "数据中心"
// @Param        equipment_room query   int  false  "机房"
// @Param        query query   string  false  "查询条件"
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Success      200  {array}  models.DeviceScrap
// @Router       /api/n9e/device-scrap/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceScrapGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	dataRange := ginx.QueryStr(c, "range", "")
	room := ginx.QueryInt64(c, "equipment_room", -1)
	datacenter := ginx.QueryInt64(c, "datacenter_id", -1)
	query := ginx.QueryStr(c, "query", "")

	var dataRan int64
	dataRan = -1
	if dataRange == "in_3_months_scrap" {
		dataRan = time.Now().AddDate(0, -3, 0).Unix()
	}
	if dataRange == "in_6_months_scrap" {
		dataRan = time.Now().AddDate(0, -6, 0).Unix()
	}
	if dataRange == "in_this_year_scrap" {
		dataRan = time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	}
	if dataRange == "in_this_year_scrap" {
		dataRan = time.Now().AddDate(0, -12, 0).Unix()
	}

	total, err := models.DeviceScrapFindCount(rt.Ctx, query, dataRan, datacenter, room)
	ginx.Dangerous(err)
	lst, err := models.DeviceScrapGets(rt.Ctx, query, dataRan, datacenter, room, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建设备报废
// @Description  创建设备报废
// @Tags         设备报废
// @Accept       json
// @Produce      json
// @Param        tree query   int  false  "查询条件"
// @Param        body  body   models.DeviceScrap true "add deviceScrap"
// @Success      200
// @Router       /api/n9e/device-scrap/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceScrapAdd(c *gin.Context) {
	var f models.DeviceScrap
	ginx.BindJSON(c, &f)
	logger.Debug(f)

	assetTotal, err := models.DeviceScrapByMap(rt.Ctx, map[string]interface{}{"asset_id": f.AssetId})
	ginx.Dangerous(err)
	if assetTotal > 0 {
		ginx.Bomb(http.StatusOK, "该资产已报废!")
	}

	query := ginx.QueryInt64(c, "query", 0)
	m := make(map[string]interface{})
	m["serial_number"] = f.SerialNumber
	m["device_producer"] = f.DeviceProducer
	m["device_model"] = f.DeviceModel
	m["device_type"] = f.DeviceType
	if f.DeviceName != "" {
		m["device_name"] = f.DeviceName
	}
	if f.OldManagementIp != "" {
		m["management_ip"] = f.OldManagementIp
	}

	assetBasic, err := models.AssetBasicGetsByMap(rt.Ctx, m, -1, -1)
	ginx.Dangerous(err)
	if len(assetBasic) == 0 {
		ginx.Bomb(http.StatusOK, "资产数据错误")
	}
	f.AssetId = assetBasic[0].Id

	datacenter, err := models.ComputerRoomGetById(rt.Ctx, assetBasic[0].EquipmentRoom)
	ginx.Dangerous(err)
	f.OldDatacenter = datacenter.Id

	assetManagement, err := models.AssetManagementGetByAssetId(rt.Ctx, assetBasic[0].Id)
	ginx.Dangerous(err)
	f.AssetCode = assetManagement.AssetCode

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	tx := models.DB(rt.Ctx).Begin()

	err = models.UpdateTxTree(tx, map[string]interface{}{"property_id": assetBasic[0].Id, "type": "asset", "updated_by": me.Username}, map[string]interface{}{"parent_id": query})
	ginx.Dangerous(err)

	models.UpdateTxStatus(tx, []int64{assetBasic[0].Id}, 4, me.Username)
	// 更新模型
	err = f.AddTx(tx)
	ginx.NewRender(c).Message(err)

	tx.Commit()
}

// @Summary      更新设备报废
// @Description  更新设备报废
// @Tags         设备报废
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceScrap true "update deviceScrap"
// @Success      200
// @Router       /api/n9e/device-scrap/ [put]
// @Security     ApiKeyAuth
// func (rt *Router) deviceScrapPut(c *gin.Context) {
// 	var f models.DeviceScrap
// 	ginx.BindJSON(c, &f)

// 	old, err := models.DeviceScrapGetById(rt.Ctx, f.Id)
// 	ginx.Dangerous(err)
// 	if old == nil {
// 		ginx.Bomb(http.StatusOK, "device_scrap not found")
// 	}

// 	// 添加审计信息
// 	me := c.MustGet("user").(*models.User)
// 	f.UpdatedBy = me.Username

// 	// 可修改"*"为字段名称，实现更新部分字段功能
// 	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
// }

// @Summary      批量删除设备报废
// @Description  根据主键批量删除设备报废
// @Tags         设备报废
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "delete deviceScrap"
// @Success      200
// @Router       /api/n9e/device-scrap/batch [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceScrapBatchDel(c *gin.Context) {
	var f []int64
	ginx.BindJSON(c, &f)
	deviceScraps, err := models.DeviceScrapBatchGetById(rt.Ctx, f)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)
	assetId := make([]int64, 0)
	for _, val := range deviceScraps {
		assetId = append(assetId, val.AssetId)
	}

	//开启事务
	tx := models.DB(rt.Ctx).Begin()

	//删除资产详情
	err = models.AssetBasicBatchDel(tx, assetId)
	ginx.Dangerous(err)

	//删除资产树数据
	err = models.AssetTreeBatchDel(tx, assetId)
	ginx.Dangerous(err)

	//删除资产扩展
	err = models.AssetExpansionBatchDel(tx, assetId)
	ginx.Dangerous(err)
	//删除资产维保
	err = models.AssetMaintenanceBatchDel(rt.Ctx, tx, assetId)
	ginx.Dangerous(err)
	//删除资产管理
	err = models.AssetManagementBatchDel(tx, assetId)
	ginx.Dangerous(err)
	//删除资产变更
	err = models.AssetAlterBatchDel(tx, assetId)
	ginx.Dangerous(err)
	//删除报废数据
	err = models.DeviceScrapBatchDel(tx, f)
	ginx.Dangerous(err)

	tx.Commit()

	ginx.NewRender(c).Message(err)
}

// @Summary      导出设备报废数据
// @Description  根据条件导出设备报废数据
// @Tags         设备报废
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200  {object}  models.DeviceModel
// @Router       /api/n9e/device-scrap/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) exportDeviceScraps(c *gin.Context) {

	f := make(map[string]interface{})
	ginx.BindJSON(c, &f)

	idsTemp, idsOk := f["ids"]
	var deviceScraps []models.DeviceScrap
	var err error
	if idsOk {
		ids := make([]int64, 0)
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, int64(val.(float64)))
		}
		logger.Debug(ids)
		// ids := idsTemp.([]int64)
		deviceScraps, err = models.DeviceScrapBatchGetById(rt.Ctx, ids)
		ginx.Dangerous(err)
	} else {
		var dataRan int64
		dataRan = -1
		dataRanTemp, dataRanTempOk := f["data_range"]
		if dataRanTempOk {
			if dataRanTemp.(string) == "in_3_months_scrap" {
				dataRan = time.Now().AddDate(0, -3, 0).Unix()
			}
			if dataRanTemp.(string) == "in_6_months_scrap" {
				dataRan = time.Now().AddDate(0, -6, 0).Unix()
			}
			if dataRanTemp.(string) == "in_this_year_scrap" {
				dataRan = time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
			}
			if dataRanTemp.(string) == "in_1_year_scrap" {
				dataRan = time.Now().AddDate(0, -12, 0).Unix()
			}
		}

		var datacenter int64
		datacenter = -1
		dc, dcOk := f["datacenter"]
		if dcOk {
			datacenter = int64(dc.(float64))
		}

		var room int64
		room = -1
		r, rOk := f["room"]
		if rOk {
			room = int64(r.(float64))
		}

		query := ""
		q, qOk := f["query"]
		if qOk {
			query = q.(string)
		}

		deviceScraps, err = models.DeviceScrapGets(rt.Ctx, query, dataRan, datacenter, room, -1, -1)
		ginx.Dangerous(err)

	}

	datas := make([]interface{}, 0)
	if len(deviceScraps) == 0 {
		datas = append(datas, models.DeviceScrap{})
		excels.NewMyExcel("资产报废数据").ExportTempletToWeb(datas, nil, "cn", "source", 0, rt.Ctx, c)
	} else {
		for _, v := range deviceScraps {
			datas = append(datas, v)
		}
		excels.NewMyExcel("资产报废数据").ExportDataInfo(datas, "cn", rt.Ctx, c)
	}
}
