// Package models  设备厂商
// date : 2023-07-08 14:51
// desc : 设备厂商
package router

import (
	"net/http"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取设备厂商
// @Description  根据主键获取设备厂商
// @Tags         设备厂商
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DeviceProducer
// @Router       /api/n9e/device-producer/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceProducerGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceProducer, err := models.DeviceProducerGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if deviceProducer == nil {
		ginx.Bomb(404, "No such device_producer")
	}

	ginx.NewRender(c).Data(deviceProducer, nil)
}

// @Summary      获取设备厂商名称
// @Description  获取设备厂商名称
// @Tags         设备厂商
// @Accept       json
// @Produce      json
// @Param        type query   string  false  "厂商类型"
// @Success      200  {object}  []models.DeviceProducerNameVo
// @Router       /api/n9e/device-producer/getName [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceProducerGetNames(c *gin.Context) {

	prodType := ginx.QueryStr(c, "type", "")

	deviceProducerNames, err := models.FindNames(rt.Ctx, prodType)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(deviceProducerNames, nil)
}

// @Summary      查询设备厂商
// @Description  根据条件查询设备厂商
// @Tags         设备厂商
// @Accept       json
// @Produce      json
// @Param        type query   string  false  "厂商类型"
// @Param        query query   string  false  "简称/全称/中文名称/联系人"
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Success      200  {array}  models.DeviceProducer
// @Router       /api/n9e/device-producer/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceProducerGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	prodType := ginx.QueryStr(c, "type", "")
	query := ginx.QueryStr(c, "query", "")

	if prodType == "" {
		ginx.Bomb(http.StatusOK, "厂商类型为空,查询失败!")
		return
	}

	m := make(map[string]interface{})
	m["producer_type"] = prodType
	if query != "" {
		m["query"] = query
	}

	total, err := models.DeviceProducerCountMap(rt.Ctx, m)
	ginx.Dangerous(err)
	lst, err := models.DeviceProducerGetByPage[models.DeviceProducer](rt.Ctx, m, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建设备厂商
// @Description  创建设备厂商
// @Tags         设备厂商
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceProducer true "add deviceProducer"
// @Success      200
// @Router       /api/n9e/device-producer/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceProducerAdd(c *gin.Context) {
	var f models.DeviceProducer
	ginx.BindJSON(c, &f)

	// 生成一个validate实例
	validate := validator.New()
	errValidate := validate.Struct(f)
	if errValidate != nil {
		ginx.Bomb(http.StatusOK, errValidate.Error())
	}

	num, err := models.DeviceProducerCountMap(rt.Ctx, map[string]interface{}{"company_name": f.CompanyName})
	ginx.Dangerous(err)
	if num != 0 {
		ginx.Bomb(http.StatusOK, "厂商已存在")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err = models.AddProd[models.DeviceProducer](rt.Ctx, f)
	ginx.NewRender(c).Message(err)
}

// @Summary      EXCEL导入设备厂商
// @Description  EXCEL导入设备厂商
// @Tags         设备厂商
// @Accept       multipart/form-data
// @Produce      json
// @Param        type  formData   string true "add query"
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/device-producer/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importDeviceProducer(c *gin.Context) {
	prodType := c.Request.FormValue("type")
	logger.Debug(prodType)
	// prodTypeTemp, prodTypeOk := f["type"]
	// var prodType string
	if prodType == "" {
		ginx.Bomb(http.StatusOK, "厂商类型为空,数据导入失败!")
		return
	}

	file, _, err := c.Request.FormFile("file")
	c.Request.FormValue("")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "上传文件出错")
	}
	//读excel流
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "读取excel文件失败")
	}
	//解析excel的数据

	me := c.MustGet("user").(*models.User)
	var qty int = 0
	DeviceProducerLst := make([]models.DeviceProducer, 0)
	if prodType == "producer" {
		produceroVo, _, lxRrr := excels.ReadExce[models.ProduceroVo](xlsx, rt.Ctx)
		if lxRrr != nil {
			ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
			return
		}

		for _, entity := range produceroVo {

			num, err := models.DeviceProducerCountMap(rt.Ctx, map[string]interface{}{"company_name": entity.CompanyName})
			ginx.Dangerous(err)
			if num != 0 {
				ginx.Bomb(http.StatusOK, "厂商已存在")

			}
			// 循环体
			var f models.DeviceProducer
			f.ProducerType = prodType
			f.Alias = entity.Alias
			f.ChineseName = entity.ChineseName
			f.CompanyName = entity.CompanyName
			f.Official = entity.Official
			f.IsDomestic = entity.IsDomestic
			f.IsDisplayChinese = entity.IsDisplayChinese
			f.CreatedBy = me.Username
			DeviceProducerLst = append(DeviceProducerLst, f)
			qty++
		}
	} else if prodType == "third_party_maintenance" {
		maintenanceVo, _, lxRrr := excels.ReadExce[models.MaintenanceVo](xlsx, rt.Ctx)
		if lxRrr != nil {
			ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
			return
		}

		for _, entity := range maintenanceVo {
			num, err := models.DeviceProducerCountMap(rt.Ctx, map[string]interface{}{"company_name": entity.CompanyName})
			ginx.Dangerous(err)
			if num != 0 {
				ginx.Bomb(http.StatusOK, "厂商已存在")

			}
			// 循环体
			var f models.DeviceProducer
			f.ProducerType = prodType
			f.Alias = entity.Alias
			f.ChineseName = entity.ChineseName
			f.CompanyName = entity.CompanyName
			f.Official = entity.Official
			f.IsDomestic = entity.IsDomestic
			f.IsDisplayChinese = entity.IsDisplayChinese
			f.CreatedBy = me.Username
			DeviceProducerLst = append(DeviceProducerLst, f)
			qty++
		}
	} else if prodType == "supplier" {
		supplierVo, _, lxRrr := excels.ReadExce[models.SupplierVo](xlsx, rt.Ctx)
		if lxRrr != nil {
			ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
			return
		}

		for _, entity := range supplierVo {
			num, err := models.DeviceProducerCountMap(rt.Ctx, map[string]interface{}{"company_name": entity.CompanyName})
			ginx.Dangerous(err)
			if num != 0 {
				ginx.Bomb(http.StatusOK, "厂商已存在")

			}
			// 循环体
			var f models.DeviceProducer
			f.ProducerType = prodType
			f.Alias = entity.Alias
			f.ChineseName = entity.ChineseName
			f.CompanyName = entity.CompanyName
			f.Official = entity.Official
			f.IsDomestic = entity.IsDomestic
			f.IsDisplayChinese = entity.IsDisplayChinese
			f.CreatedBy = me.Username
			DeviceProducerLst = append(DeviceProducerLst, f)
			qty++
		}
	} else if prodType == "component_brand" {
		partVo, _, lxRrr := excels.ReadExce[models.PartVo](xlsx, rt.Ctx)
		if lxRrr != nil {
			ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
			return
		}

		for _, entity := range partVo {
			num, err := models.DeviceProducerCountMap(rt.Ctx, map[string]interface{}{"company_name": entity.CompanyName})
			ginx.Dangerous(err)
			if num != 0 {
				ginx.Bomb(http.StatusOK, "厂商已存在")

			}
			// 循环体
			var f models.DeviceProducer
			f.ProducerType = prodType
			f.Alias = entity.Alias
			f.ChineseName = entity.ChineseName
			f.CompanyName = entity.CompanyName
			f.Official = entity.Official
			f.IsDomestic = entity.IsDomestic
			f.IsDisplayChinese = entity.IsDisplayChinese
			f.CreatedBy = me.Username
			DeviceProducerLst = append(DeviceProducerLst, f)
			qty++
		}
	}
	err = models.BatchAddProd(rt.Ctx, DeviceProducerLst)
	ginx.NewRender(c).Data(qty, err)
}

// @Summary      EXCEL导出设备厂商
// @Description  EXCEL导出设备厂商
// @Tags         设备厂商
// @Accept       multipart/form-data
// @Produce      application/msexcel
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200
// @Router       /api/n9e/device-producer/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) downloadDeviceProducer(c *gin.Context) {

	f := make(map[string]interface{})
	ginx.BindJSON(c, &f)

	prodTypeTemp, prodTypeOk := f["type"]
	var prodType string
	if prodTypeOk {
		prodType = prodTypeTemp.(string)
		delete(f, "type")
		f["producer_type"] = prodType
	} else {
		ginx.Bomb(http.StatusOK, "厂商类型为空,数据导出失败!")
		return
	}

	idsTemp, idsOk := f["ids"]
	ids := make([]int64, 0)
	// var err error
	if idsOk {
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, int64(val.(float64)))
		}
		delete(f, "ids")
	}

	datas := make([]interface{}, 0)
	var err error

	if prodType == "producer" {
		var lst []models.ProduceroVo
		if idsOk {
			lst, err = models.DeviceProducerGetByIds[models.ProduceroVo](rt.Ctx, ids)
		} else {
			lst, err = models.DeviceProducerGetByPage[models.ProduceroVo](rt.Ctx, f, -1, -1)
		}
		ginx.Dangerous(err)

		if len(lst) > 0 {
			for _, v := range lst {
				datas = append(datas, v)
			}
			excels.NewMyExcel(prodType+"数据").ExportDataInfo(datas, "cn", rt.Ctx, c)
		} else {
			datas = append(datas, models.ProduceroVo{})
			excels.NewMyExcel(prodType+"数据").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
		}
	} else if prodType == "third_party_maintenance" {
		var lst []models.MaintenanceVo
		if idsOk {
			lst, err = models.DeviceProducerGetByIds[models.MaintenanceVo](rt.Ctx, ids)
		} else {
			lst, err = models.DeviceProducerGetByPage[models.MaintenanceVo](rt.Ctx, f, -1, -1)
		}
		ginx.Dangerous(err)

		if len(lst) > 0 {
			for _, v := range lst {
				datas = append(datas, v)
			}
			excels.NewMyExcel(prodType+"数据").ExportDataInfo(datas, "cn", rt.Ctx, c)
		} else {
			datas = append(datas, models.MaintenanceVo{})
			excels.NewMyExcel(prodType+"数据").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
		}
	} else if prodType == "supplier" {
		var lst []models.SupplierVo
		if idsOk {
			lst, err = models.DeviceProducerGetByIds[models.SupplierVo](rt.Ctx, ids)
		} else {
			lst, err = models.DeviceProducerGetByPage[models.SupplierVo](rt.Ctx, f, -1, -1)
		}
		ginx.Dangerous(err)

		if len(lst) > 0 {
			for _, v := range lst {
				datas = append(datas, v)
			}
			excels.NewMyExcel(prodType+"数据").ExportDataInfo(datas, "cn", rt.Ctx, c)
		} else {
			datas = append(datas, models.SupplierVo{})
			excels.NewMyExcel(prodType+"数据").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
		}
	} else if prodType == "component_brand" {
		var lst []models.PartVo
		if idsOk {
			lst, err = models.DeviceProducerGetByIds[models.PartVo](rt.Ctx, ids)
		} else {
			lst, err = models.DeviceProducerGetByPage[models.PartVo](rt.Ctx, f, -1, -1)
		}
		ginx.Dangerous(err)

		if len(lst) > 0 {
			for _, v := range lst {
				datas = append(datas, v)
			}
			excels.NewMyExcel(prodType+"数据").ExportDataInfo(datas, "cn", rt.Ctx, c)
		} else {
			datas = append(datas, models.PartVo{})
			excels.NewMyExcel(prodType+"数据").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
		}
	}
}

// @Summary      导出设备厂商模板
// @Description  导出设备厂商模板
// @Tags         设备厂商
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200
// @Router       /api/n9e/device-producer/templet/ [post]
// @Security     ApiKeyAuth
func (rt *Router) templetDeviceProducer(c *gin.Context) {

	f := make(map[string]interface{})
	ginx.BindJSON(c, &f)

	prodTypeTemp, prodTypeOk := f["type"]
	var prodType string
	if prodTypeOk {
		prodType = prodTypeTemp.(string)
	} else {
		ginx.Bomb(http.StatusOK, "厂商类型为空,模板下载失败!")
		return
	}

	datas := make([]interface{}, 0)

	if prodType == "producer" {
		datas = append(datas, models.ProduceroVo{})
	} else if prodType == "third_party_maintenance" {
		datas = append(datas, models.MaintenanceVo{})
	} else if prodType == "supplier" {
		datas = append(datas, models.SupplierVo{})
	} else if prodType == "component_brand" {
		datas = append(datas, models.PartVo{})
	}
	m := make(map[string]string)
	m["producer"] = "厂商"
	m["third_party_maintenance"] = "第三方维保服务商"
	m["supplier"] = "供应商"
	m["component_brand"] = "部件品牌"

	excels.NewMyExcel(m[prodType]+"模板").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
}

// @Summary      更新设备厂商
// @Description  更新设备厂商
// @Tags         设备厂商
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceProducer true "update deviceProducer"
// @Success      200
// @Router       /api/n9e/device-producer/ [put]
// @Security     ApiKeyAuth
func (rt *Router) deviceProducerPut(c *gin.Context) {
	var f models.DeviceProducer
	ginx.BindJSON(c, &f)

	// 生成一个validate实例
	validate := validator.New()
	errValidate := validate.Struct(f)
	if errValidate != nil {
		ginx.Bomb(http.StatusOK, errValidate.Error())
	}

	num, err := models.DeviceProducerCountMap(rt.Ctx, map[string]interface{}{"company_name": f.CompanyName})
	ginx.Dangerous(err)
	if num != 0 {
		ginx.Bomb(http.StatusOK, "厂商已存在")

	}

	old, err := models.DeviceProducerGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "device_producer not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除设备厂商
// @Description  根据主键删除设备厂商
// @Tags         设备厂商
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/device-producer/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) deviceProducerDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceProducer, err := models.DeviceProducerGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if deviceProducer == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(deviceProducer.Del(rt.Ctx))
}

// @Summary      批量删除设备厂商
// @Description  批量删除设备厂商
// @Tags         设备厂商
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "batch delete deviceProducer"
// @Success      200
// @Router       /api/n9e/device-producer/batch-del/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceProducerBatchDel(c *gin.Context) {

	var f []int64
	ginx.BindJSON(c, &f)

	deviceModels, err := models.DeviceModelGetByPros(rt.Ctx, f)
	ginx.Dangerous(err)
	if len(deviceModels) > 0 {
		var builder strings.Builder
		builder.WriteString("存在该厂商的设备，设备型号为")
		for index, val := range deviceModels {
			if index == len(deviceModels)-1 {
				builder.WriteString(val.Model)
				break
			}
			builder.WriteString(val.Model)
			builder.WriteString("、")
		}
		ginx.Bomb(http.StatusOK, builder.String())
		return
	}

	ginx.NewRender(c).Message(models.DeviceProducerBatchDel(rt.Ctx, f))
}
