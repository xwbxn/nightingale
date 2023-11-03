// Package models  资产扩展
// date : 2023-07-23 09:06
// desc : 资产扩展
package router

import (
	"net/http"
	"reflect"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/gin-gonic/gin"
	"github.com/pochard/commons/randstr"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取资产扩展
// @Description  根据主键获取资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.AssetExpansion
// @Router       /api/n9e/asset-expansion/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetExpansion, err := models.AssetExpansionGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if assetExpansion == nil {
		ginx.Bomb(404, "No such asset_expansion")
	}

	ginx.NewRender(c).Data(assetExpansion, nil)
}

// @Summary      根据资产ID获取资产扩展
// @Description  根据资产ID获取资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        asset    query    int  true  "资产ID"
// @Success      200  {array}  []models.AssetExpansion
// @Router       /api/n9e/asset-expansion/asset [get]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionGetByAssetId(c *gin.Context) {
	assetId := ginx.QueryInt64(c, "asset", -1)
	m := make(map[string]interface{})
	m["asset_id"] = assetId

	lst, err := models.AssetExpansionGetByMap(rt.Ctx, m)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(lst, nil)
}

// @Summary      根据资产ID、类别或属性获取资产扩展
// @Description  根据资产ID、类别或属性获取资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetExpansion true "add AssetExpansion"
// @Success      200  {array}  []models.AssetExpansion
// @Router       /api/n9e/asset-expansion/map [post]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionGetByMap(c *gin.Context) {

	m := make(map[string]interface{})
	ginx.BindJSON(c, &m)

	lst, err := models.AssetExpansionGetByMap(rt.Ctx, m)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(lst, nil)
}

// @Summary      查询资产扩展
// @Description  根据条件查询资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.AssetExpansion
// @Router       /api/n9e/asset-expansion/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.AssetExpansionCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.AssetExpansionGets(rt.Ctx, query, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建资产扩展
// @Description  创建资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        body  body   models.AssetExpansion true "add assetExpansion"
// @Success      200
// @Router       /api/n9e/asset-expansion/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionAdd(c *gin.Context) {
	var f models.AssetExpansion
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      批量创建资产扩展
// @Description  批量创建资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        body  body   []models.AssetExpansion true "add assetExpansion"
// @Success      200
// @Router       /api/n9e/asset-expansion/batch/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionBatchAdd(c *gin.Context) {
	var f []models.AssetExpansion
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	for index := range f {
		f[index].CreatedBy = me.Username
	}

	// 更新模型
	err := models.AssetExpansionBatchAdd(rt.Ctx, f)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新资产扩展
// @Description  更新资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        body  body   []models.AssetExpansion true "update assetExpansion"
// @Success      200
// @Router       /api/n9e/asset-expansion/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionPut(c *gin.Context) {

	var f []models.AssetExpansion
	ginx.BindJSON(c, &f)
	if len(f) == 0 {
		ginx.Bomb(http.StatusOK, "未传入资产扩展信息")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	name := me.Username

	var err error
	m := make(map[string]interface{})

	//确定更新位置为硬件配置
	if f[0].ConfigCategory == "form-hardware-cfg" {

		//通过资产id、属性等查询
		m["ASSET_ID"] = f[0].AssetId
		m["CONFIG_CATEGORY"] = f[0].ConfigCategory
		m["PROPERTY_CATEGORY"] = f[0].PropertyCategory

		err = models.UpdateAssetExpansionGroup(rt.Ctx, m, f, name)
		ginx.Dangerous(err)

	} else if f[0].ConfigCategory == "form-netconfig" {

		m["ASSET_ID"] = f[0].AssetId
		m["CONFIG_CATEGORY"] = f[0].ConfigCategory

		err = models.UpdateAssetExpansionGroup(rt.Ctx, m, f, name)
		ginx.Dangerous(err)
	}

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(err)
}

// @Summary      删除资产扩展
// @Description  根据主键删除资产扩展
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/asset-expansion/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) assetExpansionDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assetExpansion, err := models.AssetExpansionGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if assetExpansion == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(assetExpansion.Del(rt.Ctx))
}

// @Summary      导出网络配置模板
// @Description  导出网络配置模板
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.AssetNetWork
// @Router       /api/n9e/asset-expansion/netconfig/templet/ [post]
// @Security     ApiKeyAuth
func (rt *Router) templeAssetNetWork(c *gin.Context) {

	datas := make([]interface{}, 0)

	datas = append(datas, models.AssetNetWork{})

	excels.NewMyExcel("网络配置").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
}

// @Summary      EXCEL导入网络配置
// @Description  EXCEL导入网络配置
// @Tags         资产扩展
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/asset-expansion/netconfig/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importAssetNetWork(c *gin.Context) {

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
	assetNetWorks, _, lxRrr := excels.ReadExce[models.AssetNetWork](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	me := c.MustGet("user").(*models.User)
	var qty int = 0
	tx := models.DB(rt.Ctx).Begin()
	for _, entity := range assetNetWorks {
		// 循环体
		var f models.AssetExpansion
		assetBasic, err := models.AssetBasicGetsByMap(rt.Ctx, map[string]interface{}{"device_name": entity.DeviceName, "serial_number": entity.SerialNumber, "management_ip": entity.ManagementIp}, -1, -1)
		ginx.Dangerous(err)
		if len(assetBasic) == 0 {
			ginx.Bomb(http.StatusBadRequest, "参数错误")
		}
		f.AssetId = assetBasic[0].Id
		f.ConfigCategory = "form-netconfig"
		f.GroupId = randstr.RandomAlphanumeric(32)
		f.CreatedBy = me.Username

		t := reflect.TypeOf(entity) // 注意，obj不能为指针类型，否则会报：panic recovered: reflect: NumField of non-struct type
		v := reflect.ValueOf(entity)
		//获取到该结构体有几个字段
		num := v.NumField()

		if entity.Type == 0 {
			f.PropertyCategory = "group_management"
		} else {
			f.PropertyCategory = "group_ip"
		}
		//遍历结构体的所有字段
		for i := 0; i < num; i++ {
			tagCn := t.Field(i).Tag.Get("cn")
			tagVal := t.Field(i).Tag.Get("prop")
			val := v.Field(i).String()
			if tagCn == "设备IP" || tagCn == "设备序列号" || tagCn == "设备名称" || tagCn == "类型" {
				continue
			}
			if val == "" {
				continue
			}

			if tagCn == "IP" {
				if entity.Type == 0 {
					f.PropertyNameCn = "带外IP"
					f.PropertyName = "ext_out_band_ip"
				} else {
					f.PropertyNameCn = "生产IP"
					f.PropertyName = "production_ip"
				}
				f.PropertyValue = val
			} else if tagCn == "交换机端口" {
				f.PropertyNameCn = "对应端口"
				f.PropertyName = tagVal
				f.PropertyValue = val
			} else if tagCn == "配线架名称" {
				f.PropertyNameCn = "对接配线架"
				f.PropertyName = tagVal
				f.PropertyValue = val
			} else if tagCn == "配线架端口" {
				f.PropertyNameCn = "对接端口"
				f.PropertyName = tagVal
				f.PropertyValue = val
			} else if tagCn == "连接用户名" {
				f.PropertyNameCn = "用户名"
				f.PropertyName = tagVal
				f.PropertyValue = val
			} else if tagCn == "连接密码" {
				f.PropertyNameCn = "密码"
				f.PropertyName = tagVal
				f.PropertyValue = val
			} else if tagCn == "远程端口" {
				f.PropertyNameCn = "端口"
				f.PropertyName = tagVal
				f.PropertyValue = val
			} else {
				f.PropertyNameCn = tagCn
				f.PropertyName = tagVal
				f.PropertyValue = val
			}
			f.Id = 0
			err = f.AddTx(tx)
			ginx.Dangerous(err)
			qty++
		}

	}
	tx.Commit()
	ginx.NewRender(c).Data(qty, nil)
}

// @Summary      导出网络配置
// @Description  根据条件导出网络配置
// @Tags         资产扩展
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200  {object}  models.DeviceModel
// @Router       /api/n9e/asset-expansion/netconfig/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) exportAssetNetWork(c *gin.Context) {

	// status := ginx.QueryInt64(c, "status", -1)
	f := make(map[string]interface{})
	ginx.BindJSON(c, &f)
	idsTemp, idsOk := f["ids"]
	ids := make([]int64, 0)
	var assetBasics []models.AssetBasicDetailsVo
	var err error
	if idsOk {
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, int64(val.(float64)))
		}
		assetBasics, err = models.AssetBasicGetsByIds(rt.Ctx, ids, -1, -1)
		ginx.Dangerous(err)
	} else {
		status, statusOk := f["status"]
		if statusOk {
			delete(f, "status")
			f["device_status"] = status
		}

		//根据状态查询资产
		assetBasics, err = models.AssetBasicGetsByMap(rt.Ctx, f, -1, -1)
		ginx.Dangerous(err)
		for _, val := range assetBasics {
			ids = append(ids, val.Id)
		}
	}

	//根据asset_id查询网络配置
	lst, err := models.AssetNetConfigGetByAssetId(rt.Ctx, ids)
	ginx.Dangerous(err)

	assetNetWorkLst := make([]models.AssetNetWork, 0)
	mTemp := make(map[string]map[string]interface{})
	for _, val := range lst {
		exsit, ok := mTemp[val.GroupId]
		if ok {
			exsit[val.PropertyName] = val.PropertyValue
		} else {
			group := make(map[string]interface{}, 0)
			group[val.PropertyName] = val.PropertyValue
			group["asset_id"] = val.AssetId
			mTemp[val.GroupId] = group
		}
	}

	for key := range mTemp {
		var assetNetWork models.AssetNetWork
		t := reflect.TypeOf(assetNetWork) // 注意，obj不能为指针类型，否则会报：panic recovered: reflect: NumField of non-struct type
		v := reflect.ValueOf(&assetNetWork).Elem()
		//获取到该结构体有几个字段
		num := v.NumField()

		logger.Debug(mTemp[key])
		//遍历结构体的所有字段
		for i := 0; i < num; i++ {
			tagCn := t.Field(i).Tag.Get("cn")
			tag := t.Field(i).Tag.Get("prop")
			logger.Debug(tagCn)
			logger.Debug(tag)

			//补充资产基本信息
			for _, assetBasic := range assetBasics {

				if assetBasic.Id == mTemp[key]["asset_id"] {
					if tagCn == "设备IP" {
						sliceValue := reflect.ValueOf(assetBasic.ManagementIp)
						v.Field(i).Set(sliceValue)
					} else if tagCn == "设备序列号" {
						sliceValue := reflect.ValueOf(assetBasic.SerialNumber)
						v.Field(i).Set(sliceValue)
					} else if tagCn == "设备名称" {
						sliceValue := reflect.ValueOf(assetBasic.DeviceName)
						v.Field(i).Set(sliceValue)
					}

					break
				}
			}

			bandIp, bandOk := mTemp[key]["ext_out_band_ip"]
			if bandOk {
				if tagCn == "类型" {
					v.Field(i).SetInt(0)
				}
				if tagCn == "IP" {
					sliceValue := reflect.ValueOf(bandIp)
					v.Field(i).Set(sliceValue)
				}
			}

			prodIp, prodOk := mTemp[key]["production_ip"]
			if prodOk {
				if tagCn == "类型" {
					v.Field(i).SetInt(1)
				}
				if tagCn == "IP" {
					sliceValue := reflect.ValueOf(prodIp)
					v.Field(i).Set(sliceValue)
				}
			}

			for key, val := range mTemp[key] {

				if tag == key {
					sliceValue := reflect.ValueOf(val)
					v.Field(i).Set(sliceValue)
				}
			}
		}
		assetNetWorkLst = append(assetNetWorkLst, assetNetWork)
	}

	datas := make([]interface{}, 0)
	if len(assetNetWorkLst) <= 0 {
		datas = append(datas, models.AssetNetWork{})
		excels.NewMyExcel("网络配置").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
	} else {
		for _, v := range assetNetWorkLst {
			datas = append(datas, v)
		}
		excels.NewMyExcel("网络配置").ExportDataInfo(datas, "cn", rt.Ctx, c)
	}
}
