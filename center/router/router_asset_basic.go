// Package models  资产详情
// date : 2023-07-21 08:45
// desc : 资产详情
package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/gin-gonic/gin"
	"github.com/pochard/commons/randstr"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取资产详情
// @Description  根据主键获取资产详情
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.AssetBasicExpansionVo
// @Router       /api/n9e/asset-basic/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetBasicExpansionVo, err := models.AssetBasicGetById[models.AssetBasicExpansionVo](rt.Ctx, id)
	ginx.Dangerous(err)

	assetExpansions, err := models.AssetExpansionGetByMap(rt.Ctx, map[string]interface{}{"asset_id": id, "config_category": "form-basic"})
	ginx.Dangerous(err)

	if len(assetExpansions) > 0 {
		assetBasicExpansionVo.BasicExpansion = assetExpansions
	}

	if assetBasicExpansionVo == nil {
		ginx.Bomb(http.StatusOK, "数据不存在!")
	}

	ginx.NewRender(c).Data(assetBasicExpansionVo, nil)
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

// @Summary      查询资产详情(过滤器)
// @Description  根据条件查询资产详情(过滤器)
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200  {array}  models.AssetBasic
// @Router       /api/n9e/asset-basic/list/query/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicBatchGets(c *gin.Context) {

	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)

	query := make(map[string]interface{})
	ginx.BindJSON(c, &query)
	logger.Debug(query)

	manager, managerOk := query["manager"]
	if managerOk {
		delete(query, "manager")
		query["device_manager_one"] = manager
		query["device_manager_two"] = manager
		query["business_manager_one"] = manager
		query["business_manager_two"] = manager
	}
	queryCopy := make(map[string]interface{})
	for key := range query {
		queryCopy[key] = query[key]
	}
	total, err := models.AssetCountMap(rt.Ctx, query)
	ginx.Dangerous(err)

	lst, err := models.AssetByMap(rt.Ctx, queryCopy, limit, (page-1)*limit)
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
		deviceType, err := models.DeviceTypeGetById(rt.Ctx, lst[index].DeviceType)
		ginx.Dangerous(err)
		lst[index].DeviceTypeName = deviceType.Name
		//回填设备厂商
		deviceProducer, err := models.DeviceProducerGetById(rt.Ctx, lst[index].DeviceProducer)
		ginx.Dangerous(err)
		lst[index].DeviceProducerName = deviceProducer.Alias
		//回填设备型号
		deviceModel, err := models.DeviceModelGetById(rt.Ctx, lst[index].DeviceModel)
		ginx.Dangerous(err)
		lst[index].DeviceModelName = deviceModel.Name
		//TODO 回填所属组织(不清楚)

		//回填所在机房
		if lst[index].EquipmentRoom != 0 {
			computerRoom, err := models.ComputerRoomGetById(rt.Ctx, lst[index].EquipmentRoom)
			ginx.Dangerous(err)
			lst[index].RoomName = computerRoom.RoomName
		}

		//回填所在机柜
		if lst[index].OwningCabinet != 0 {
			deviceCabinet, err := models.DeviceCabinetGetById(rt.Ctx, lst[index].OwningCabinet)
			ginx.Dangerous(err)
			lst[index].CabinetName = deviceCabinet.CabinetName
		}
		//补充扩展字段
		assetExpansions, err := models.AssetExpansionGetByMap(rt.Ctx, map[string]interface{}{"asset_id": lst[index].Id, "config_category": "form-basic"})
		ginx.Dangerous(err)
		lst[index].BasicExpansion = assetExpansions
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
// @Param        body  body   models.AssetBasicExpansionVo true "add assetBasic"
// @Success      200
// @Router       /api/n9e/asset-basic/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicAdd(c *gin.Context) {

	var data map[string]interface{}
	var basic models.AssetBasic
	var expansions []models.AssetExpansion
	ginx.BindJSON(c, &data)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)

	_, expOk := data["basic_expansion"]
	if expOk {
		resByre, resByteErr := json.Marshal(data["basic_expansion"])
		if resByteErr != nil {
			ginx.Bomb(http.StatusOK, "参数解析失败!")
			return
		}
		jsonRes := json.Unmarshal(resByre, &expansions)
		if jsonRes != nil {
			ginx.Bomb(http.StatusOK, "参数解析失败!")
			return
		}
		delete(data, "basic_expansion")
		for index := range expansions {
			expansions[index].CreatedBy = me.Username
		}
	}

	resByre, resByteErr := json.Marshal(data)
	if resByteErr != nil {
		ginx.Bomb(http.StatusOK, "参数解析失败!")
		return
	}
	jsonRes := json.Unmarshal(resByre, &basic)
	if jsonRes != nil {
		ginx.Bomb(http.StatusOK, "参数解析失败!")
		return
	}
	basic.CreatedBy = me.Username

	if basic.ManagementIp == "" && basic.SerialNumber == "" {
		ginx.Bomb(http.StatusOK, "管理IP和序列号不能同时为空!")
	}

	// 更新模型
	//启动事务
	tx := models.DB(rt.Ctx).Begin()
	id, err := basic.Add(rt.Ctx, tx)
	ginx.Dangerous(err)
	for index := range expansions {
		expansions[index].AssetId = id
	}
	if expOk {
		err = models.AssetExpansionTxBatchAdd(tx, expansions)
		ginx.Dangerous(err)
	}
	tx.Commit()

	ginx.NewRender(c).Data(id, err)

}

// @Summary      更新资产详情
// @Description  更新资产详情
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetBasicExpansionVo true "update assetBasic"
// @Success      200
// @Router       /api/n9e/asset-basic/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicPut(c *gin.Context) {
	var f models.AssetBasicExpansionVo
	ginx.BindJSON(c, &f)

	old, err := models.AssetBasicGetById[models.AssetBasicExpansionVo](rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "asset_basic not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	oldAssetExpansions, err := models.AssetExpansionGetByMap(rt.Ctx, map[string]interface{}{"asset_id": old.Id, "config_category": "form-basic"})
	ginx.Dangerous(err)

	old.BasicExpansion = oldAssetExpansions

	tx := models.DB(rt.Ctx).Begin()
	// 可修改"*"为字段名称，实现更新部分字段功能
	err = old.Update(rt.Ctx, tx, me.Username, f, "*")
	ginx.Dangerous(err)

	for _, newAssetExpansion := range f.BasicExpansion {
		if newAssetExpansion.Id == 0 {
			err = newAssetExpansion.AddTx(tx)
			ginx.Dangerous(err)
		}
	}
	for _, oldAssetExpansion := range oldAssetExpansions {
		flog := false
		for _, newAssetExpansion := range f.BasicExpansion {
			if oldAssetExpansion.Id == newAssetExpansion.Id {
				flog = true
				err = oldAssetExpansion.Update(tx, me.Username, newAssetExpansion, "*")
				ginx.Dangerous(err)
				break
			}
		}
		if !flog {
			err = oldAssetExpansion.TxDel(tx)
			ginx.Dangerous(err)
		}
	}
	tx.Commit()
	ginx.NewRender(c).Message(err)
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
	assetBasic, err := models.AssetBasicGetById[models.AssetBasic](rt.Ctx, id)
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

	ginx.NewRender(c).Data(data, nil)
}

// @Summary      查询资产详情/维保/管理
// @Description  根据条件查询资产详情/维保/管理
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        asset query   int     true  "资产ID"
// @Success      200  {array}  models.AssetBasicMainMang
// @Router       /api/n9e/asset-basic/copy/ [get]
// @Security     ApiKeyAuth
func (rt *Router) assetCopyGetsByAssetId(c *gin.Context) {

	assetId := ginx.QueryInt64(c, "asset", -1)

	var assetBasicMainMang models.AssetBasicMainMang
	basic, err := models.AssetBasicGetById[models.AssetBasicExpansionVo](rt.Ctx, assetId)
	ginx.Dangerous(err)

	assetExpansions, err := models.AssetExpansionGetByMap(rt.Ctx, map[string]interface{}{"asset_id": basic.Id, "config_category": "form-basic"})
	ginx.Dangerous(err)
	basic.BasicExpansion = assetExpansions

	assetBasicMainMang.AssetBasicCopy = *basic

	maintVo, err := models.AssetMaintenanceVoGetByAssetId(rt.Ctx, assetId)
	ginx.Dangerous(err)
	assetBasicMainMang.AssetMaintenanceCopy = *maintVo

	manag, err := models.AssetManagementGetByAssetId(rt.Ctx, assetId)
	ginx.Dangerous(err)
	assetBasicMainMang.AssetManagementCopy = *manag

	ginx.NewRender(c).Data(assetBasicMainMang, nil)
}

// @Summary      复制资产详情/维保/管理
// @Description  复制资产详情/维保/管理
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetBasicMainMang true "add AssetBasicMainMang"
// @Success      200
// @Router       /api/n9e/asset-basic/copy/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetCopyAdd(c *gin.Context) {
	var f models.AssetBasicMainMang
	ginx.BindJSON(c, &f)

	//新增资产详情
	assetBasicExpansionVo := f.AssetBasicCopy
	if assetBasicExpansionVo.ManagementIp == "" && assetBasicExpansionVo.SerialNumber == "" {
		ginx.Bomb(http.StatusOK, "管理IP和序列号不能同时为空!")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	assetBasicExpansionVo.CreatedBy = me.Username

	assetBasic := models.AssetBasic{
		DeviceType:             assetBasicExpansionVo.DeviceType,
		ManagementIp:           assetBasicExpansionVo.ManagementIp,
		DeviceName:             assetBasicExpansionVo.DeviceName,
		SerialNumber:           assetBasicExpansionVo.SerialNumber,
		DeviceStatus:           assetBasicExpansionVo.DeviceStatus,
		ManagedState:           assetBasicExpansionVo.ManagedState,
		DeviceProducer:         assetBasicExpansionVo.DeviceProducer,
		DeviceModel:            assetBasicExpansionVo.DeviceModel,
		Subtype:                assetBasicExpansionVo.Subtype,
		OutlineStructure:       assetBasicExpansionVo.OutlineStructure,
		Specifications:         assetBasicExpansionVo.Specifications,
		UNumber:                assetBasicExpansionVo.UNumber,
		UseStorage:             assetBasicExpansionVo.UseStorage,
		DatacenterId:           assetBasicExpansionVo.DatacenterId,
		RelatedService:         assetBasicExpansionVo.RelatedService,
		ServicePath:            assetBasicExpansionVo.ServicePath,
		DeviceManagerOne:       assetBasicExpansionVo.DeviceManagerOne,
		DeviceManagerTwo:       assetBasicExpansionVo.DeviceManagerTwo,
		BusinessManagerOne:     assetBasicExpansionVo.BusinessManagerOne,
		BusinessManagerTwo:     assetBasicExpansionVo.BusinessManagerTwo,
		OperatingSystem:        assetBasicExpansionVo.OperatingSystem,
		Remark:                 assetBasicExpansionVo.Remark,
		AffiliatedOrganization: assetBasicExpansionVo.AffiliatedOrganization,
		EquipmentRoom:          assetBasicExpansionVo.EquipmentRoom,
		OwningCabinet:          assetBasicExpansionVo.OwningCabinet,
		Region:                 assetBasicExpansionVo.Region,
		CabinetLocation:        assetBasicExpansionVo.CabinetLocation,
		Abreast:                assetBasicExpansionVo.Abreast,
		LocationDescription:    assetBasicExpansionVo.LocationDescription,
		ExtensionTest:          assetBasicExpansionVo.ExtensionTest,
		CreatedBy:              me.Username,
	}
	expansions := assetBasicExpansionVo.BasicExpansion

	// 更新模型
	//启动事务
	tx := models.DB(rt.Ctx).Begin()
	assetId, err := assetBasic.Add(rt.Ctx, tx)
	ginx.Dangerous(err)
	for index := range expansions {
		expansions[index].AssetId = assetId
		expansions[index].CreatedBy = me.Username
	}
	err = models.AssetExpansionTxBatchAdd(tx, expansions)
	ginx.Dangerous(err)

	//新增资产维保
	assetMain := f.AssetMaintenanceCopy
	err = assetMain.AddConfig(tx, me.Username)
	ginx.Dangerous(err)

	//新增资产管理
	assetManag := f.AssetManagementCopy
	err = assetManag.AssetManagementAddTx(tx)
	ginx.Dangerous(err)
	tx.Commit()
	ginx.NewRender(c).Data(assetId, err)
}

// @Summary      批量删除资产
// @Description  根据主键批量删除资产
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "update assetBasic"
// @Success      200
// @Router       /api/n9e/asset-basic/del [post]
// @Security     ApiKeyAuth
func (rt *Router) assetBatchDel(c *gin.Context) {

	var assetId []int64
	ginx.BindJSON(c, &assetId)

	//开启事务
	tx := models.DB(rt.Ctx).Begin()

	//删除资产详情
	err := models.AssetBasicBatchDel(tx, assetId)
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

	tx.Commit()

	ginx.NewRender(c).Message(err)
}

// @Summary      修改资产状态
// @Description  根据资产ID修改资产状态
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        status query   int     true  "修改状态"
// @Param        body  body   []int true "add array"
// @Success      200
// @Router       /api/n9e/asset-basic/status/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicStatusUpdateByAssetId(c *gin.Context) {

	status := ginx.QueryInt64(c, "status", -1)
	var assetIds []int64
	ginx.BindJSON(c, &assetIds)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)

	ginx.NewRender(c).Message(models.UpdateStatus(rt.Ctx, assetIds, status, me.Username))
}

// @Summary      统计全部资产（总数、已上线/待上线/已下线）
// @Description  统计全部资产（总数、已上线/待上线/已下线）
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/n9e/asset-basic/statistics/ [get]
// @Security     ApiKeyAuth
func (rt *Router) assetStatisticsAll(c *gin.Context) {

	total, err := models.AssetCountByMap(rt.Ctx, map[string]interface{}{})
	ginx.Dangerous(err)

	m, err := models.AssetStatusCount(rt.Ctx)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  m,
		"total": total,
	}, nil)
}

// @Summary      导出资产模板
// @Description  导出资产模板
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.AssetBasicImport
// @Router       /api/n9e/asset-basic/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templeAsset(c *gin.Context) {

	datas := make([]interface{}, 0)

	datas = append(datas, models.AssetBasicImport{})
	typeM := make(map[string]interface{})
	typeM["type_code"] = "basic_expansion"

	expansionData, err := models.DictDataGetByMap(rt.Ctx, typeM)
	ginx.Dangerous(err)
	m := make([]map[string]string, 0)
	if len(expansionData) > 0 {
		for _, val := range expansionData {
			mExp := make(map[string]string)
			mExp["title"] = val.DictValue
			mExp["key"] = val.DictKey
			m = append(m, mExp)
		}
	}

	excels.NewMyExcel("设备模板").ExportTempletToWeb(datas, m, "cn", "source", rt.Ctx, c)
}

// @Summary      EXCEL导入资产
// @Description  EXCEL导入资产
// @Tags         资产详情
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/asset-basic/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importAsset(c *gin.Context) {

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "上传文件出错")
	}
	//读excel流
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "读取excel文件失败")
	}
	//解析excel的数据
	assetBasicImport, exps, lxRrr := excels.ReadExce[models.AssetBasicImport](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)

	tx := models.DB(rt.Ctx).Begin()

	for index, val := range assetBasicImport {
		var assetBasic models.AssetBasic
		var assetMaintenance models.AssetMaintenance
		var assetExpansion models.AssetExpansion
		var assetManagement models.AssetManagement
		assetBasic.DeviceType = val.DeviceType
		assetBasic.DeviceName = val.DeviceName
		assetBasic.SerialNumber = val.SerialNumber
		assetBasic.DeviceStatus = 1
		assetBasic.DeviceProducer = val.DeviceProducer
		//查询型号
		deviceModel, err := models.DeviceModelGetByNameOrModel(rt.Ctx, val.DeviceModel)
		ginx.Dangerous(err)
		if (models.DeviceModel{} == *deviceModel) {
			ginx.Bomb(http.StatusOK, fmt.Sprintf("第%d行数据,型号不存在!", index))
		}
		assetBasic.DeviceModel = deviceModel.Id
		assetBasic.Subtype = deviceModel.Subtype
		assetBasic.OutlineStructure = deviceModel.OutlineStructure
		assetBasic.Specifications = deviceModel.Specifications
		assetBasic.UNumber = deviceModel.UNumber
		assetBasic.OperatingSystem = val.OperatingSystem
		assetBasic.RelatedService = val.RelatedService
		assetBasic.ServicePath = val.ServicePath
		assetBasic.DatacenterId = val.DatacenterId
		//查询机房
		computerRoom, err := models.ComputerRoomGetByRoomName(rt.Ctx, val.EquipmentRoom)
		ginx.Dangerous(err)
		if (models.ComputerRoom{} == *computerRoom) {
			ginx.Bomb(http.StatusOK, fmt.Sprintf("第%d行数据,机房不存在!", index))
		}
		assetBasic.EquipmentRoom = computerRoom.Id
		//查询机柜
		deviceCabinet, err := models.DeviceCabinetGetByCabinetName(rt.Ctx, val.OwningCabinet)
		ginx.Dangerous(err)
		if (models.DeviceCabinet{} == deviceCabinet) {
			ginx.Bomb(http.StatusOK, fmt.Sprintf("第%d行数据,机柜不存在!", index))
		}
		assetBasic.OwningCabinet = deviceCabinet.Id
		assetBasic.CabinetLocation = val.CabinetLocation
		//TODO 机柜位置c：不清楚具体逻辑
		assetBasic.Abreast = val.Abreast
		assetBasic.Region = val.Region
		assetBasic.DeviceManagerOne = val.DeviceManagerOne
		assetBasic.DeviceManagerTwo = val.DeviceManagerTwo
		assetBasic.BusinessManagerOne = val.BusinessManagerOne
		assetBasic.BusinessManagerTwo = val.BusinessManagerTwo
		assetBasic.CreatedBy = me.CreateBy

		assetId, err := assetBasic.Add(rt.Ctx, tx)
		ginx.Dangerous(err)

		assetExpansion.AssetId = assetId
		assetExpansion.ConfigCategory = "form-netconfig"
		assetExpansion.PropertyCategory = "group_ip"
		assetExpansion.GroupId = randstr.RandomAlphanumeric(32)
		assetExpansion.PropertyNameCn = "生产IP"
		assetExpansion.PropertyName = "production_ip"
		assetExpansion.PropertyValue = val.ProductionIp
		assetExpansion.CreatedBy = me.CreateBy
		err = assetExpansion.AddTx(tx)
		ginx.Dangerous(err)

		groupId := randstr.RandomAlphanumeric(32)
		for key := range exps[index] {
			var assetExp models.AssetExpansion
			assetExp.AssetId = assetId
			assetExp.ConfigCategory = "form-basic"
			assetExp.PropertyCategory = "basic_expansion"
			assetExp.GroupId = groupId
			assetExp.PropertyNameCn = key
			dictData, err := models.DictDataGetByTypeCodeValue(rt.Ctx, "basic_expansion", key)
			ginx.Dangerous(err)
			assetExp.PropertyName = dictData.DictKey
			assetExp.PropertyValue = exps[index][key]
			assetExp.CreatedBy = me.Username
			err = assetExp.AddTx(tx)
			ginx.Dangerous(err)
		}

		assetMaintenance.AssetId = assetId
		assetMaintenance.MaintenanceType = val.MaintenanceType
		assetMaintenance.MaintenanceProvider = val.MaintenanceProvider
		assetMaintenance.StartAt = val.StartAt
		assetMaintenance.FinishAt = val.FinishAt
		assetMaintenance.CreatedBy = me.CreateBy

		err = assetMaintenance.AddTx(tx)
		ginx.Dangerous(err)

		assetManagement.AssetId = assetId
		assetManagement.AssetCode = val.AssetCode
		assetManagement.BelongDept = val.BelongDept
		assetManagement.EquipmentUse = val.EquipmentUse
		assetManagement.CreatedBy = me.CreateBy

		err = assetManagement.AssetManagementAddTx(tx)
		ginx.Dangerous(err)
	}
	tx.Commit()

	ginx.NewRender(c).Data(nil, err)
}

// @Summary      导出资产数据
// @Description  根据传入字段导出资产数据
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        body  body   map[string][]string true "update"
// @Success      200
// @Router       /api/n9e/asset-basic/outport [post]
// @Security     ApiKeyAuth
func (rt *Router) exportAsset(c *gin.Context) {

	var f map[string][]string
	ginx.BindJSON(c, &f)

	lst, err := models.AssetBasicGetsAll(rt.Ctx)
	ginx.Dangerous(err)

	datas := make([]interface{}, 0)
	if len(lst) > 0 {
		for _, v := range lst {
			datas = append(datas, v)

		}
	}

	selectFields := make([]string, 0)
	for Key := range f {
		for _, val := range f[Key] {
			selectFields = append(selectFields, val)
		}
	}

	// excels.NewMyExcel("资产数据").ExportDataSelect(datas, selectFields, "cn", rt.Ctx, c)

}

// @Summary      批量更新资产详情
// @Description  批量更新资产详情
// @Tags         资产详情
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]interface{} true "update"
// @Success      200
// @Router       /api/n9e/asset-basic/batch-update/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetBasicsUpdate(c *gin.Context) {
	var f map[string]interface{}
	ginx.BindJSON(c, &f)

	ids := f["assetIds"].([]interface{})
	delete(f, "assetIds")

	assetIds := []int64{}
	for _, id := range ids {
		assetIds = append(assetIds, int64(id.(float64)))
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f["updated_by"] = me.Username

	var err error

	_, depOk := f["user_department"]
	if depOk {
		err = models.UpdateManagMap(rt.Ctx, assetIds, f)
	} else {
		err = models.UpdateBasicMap(rt.Ctx, assetIds, f)
	}

	ginx.NewRender(c).Message(err)
}
