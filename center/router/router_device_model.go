// Package models  设备型号
// date : 2023-07-08 14:57
// desc : 设备型号
package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	picture "github.com/ccfos/nightingale/v6/pkg/picture"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取设备型号
// @Description  根据主键获取设备型号
// @Tags         设备型号
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DeviceModel
// @Router       /api/n9e/device-model/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceModelGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceModel, err := models.DeviceModelGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if deviceModel == nil {
		ginx.Bomb(404, "No such device_model")
	}

	ginx.NewRender(c).Data(deviceModel, nil)
}

// @Summary      查询设备型号
// @Description  根据条件查询设备型号
// @Tags         设备型号
// @Accept       json
// @Produce      json
// @Param        deviceType query   int     false  "类型"
// @Param        producer query   int     false  "厂商"
// @Param        query query   string     false  "型号/带外版本/描述"
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Success      200  {array}  models.DeviceModelDetailsVo
// @Router       /api/n9e/device-model/getmodel/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceModelGets(c *gin.Context) {
	deviceType := ginx.QueryInt64(c, "deviceType", -1)
	producer := ginx.QueryInt64(c, "producer", -1)
	query := ginx.QueryStr(c, "query", "")
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)

	m := make(map[string]interface{})
	if deviceType != -1 {
		m["device_type"] = deviceType
	}
	if producer != -1 {
		m["PRODUCER_ID"] = producer
	}
	if query != "" {
		m["query"] = query
	}

	total, err := models.DeviceModelCountMap(rt.Ctx, m)
	ginx.Dangerous(err)

	lst, err := models.DeviceModelGetsByType(rt.Ctx, m, limit, (page-1)*limit)
	ginx.Dangerous(err)

	// for index := range lst {
	// 	deviceProducer, err := models.DeviceProducerGetById(rt.Ctx, lst[index].ProducerId)
	// 	ginx.Dangerous(err)
	// 	lst[index].Alias = deviceProducer.Alias
	// }

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      查询设备型号
// @Description  根据条件查询设备型号
// @Tags         设备型号
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "add ids"
// @Success      200  {array}  models.DeviceModelDetailsVo
// @Router       /api/n9e/device-model/batch/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceModelBatchGets(c *gin.Context) {
	ids := make([]int64, 0)
	ginx.BindJSON(c, &ids)

	if len(ids) == 0 {
		ginx.NewRender(c).Data(nil, nil)
	}

	lst, err := models.DeviceModelGetsByIds(rt.Ctx, ids)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(lst, nil)
}

// @Summary      创建设备型号
// @Description  创建设备型号
// @Tags         设备型号
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceModel true "add deviceModel"
// @Success      200
// @Router       /api/n9e/device-model/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceModelAdd(c *gin.Context) {
	var f models.DeviceModel
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      导入设备照片
// @Description  导入设备照片
// @Tags         设备型号
// @Accept       json
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/device-model/picture/ [post]
// @Security     ApiKeyAuth
func (rt *Router) DeviceModelpictureAdd(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "文件上传失败")
	}

	suffix, err := picture.VerifyPicture(fileHeader)
	ginx.Dangerous(err)

	// 设置路径,保存文件
	filePath, err := picture.GeneratePictureName("device-model", suffix)
	ginx.Dangerous(err)

	c.SaveUploadedFile(fileHeader, filePath)

	ginx.NewRender(c).Data(filePath, err)
}

// @Summary      更新设备型号
// @Description  更新设备型号
// @Tags         设备型号
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceModel true "update deviceModel"
// @Success      200
// @Router       /api/n9e/device-model/ [put]
// @Security     ApiKeyAuth
func (rt *Router) deviceModelPut(c *gin.Context) {
	var f models.DeviceModel
	ginx.BindJSON(c, &f)

	old, err := models.DeviceModelGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "device_model not found")
	}
	oldPicture := old.Picture

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	err = old.Update(rt.Ctx, f, "*")
	ginx.Dangerous(err)
	if oldPicture != f.Picture {
		err = os.Remove(oldPicture)
	}
	ginx.NewRender(c).Message(err)
}

// @Summary      批量删除设备型号
// @Description  批量删除设备型号
// @Tags         设备型号
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "delete models"
// @Success      200
// @Router       /api/n9e/device-model/batch-del/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceModelBatchDel(c *gin.Context) {
	var f []int64
	ginx.BindJSON(c, &f)

	assetBasics, err := models.AssetBasicCountByModels(rt.Ctx, f)
	ginx.Dangerous(err)
	if len(assetBasics) > 0 {
		var builder strings.Builder
		builder.WriteString("存在设备型号为")
		for index, val := range assetBasics {
			model, err := models.DeviceModelGetById(rt.Ctx, val.DeviceModel)
			ginx.Dangerous(err)
			if index == len(assetBasics)-1 {
				builder.WriteString(model.Name)
				break
			}
			builder.WriteString(model.Name)
			builder.WriteString("、")
		}
		builder.WriteString("的资产,不可删除")
		ginx.Bomb(http.StatusOK, builder.String())
		return
	}
	//查询照片
	modelPictures, err := models.PicturesGetByIds(rt.Ctx, f)
	ginx.Dangerous(err)

	err = models.DeviceModelBatchDel(rt.Ctx, f)
	ginx.Dangerous(err)
	for _, modelPicture := range modelPictures {
		if modelPicture.Picture == "" {
			continue
		}
		err = os.Remove(modelPicture.Picture)
		ginx.Dangerous(err)
	}

	ginx.NewRender(c).Message(err)
}

// @Summary      导入设备型号数据
// @Description  导入型号
// @Tags         设备型号
// @Accept       multipart/form-data
// @Param        file formData file true "file"
// @Produce      json
// @Success      200
// @Router       /api/n9e/device-model/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importDeviceModels(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "上传文件出错")
		return
	}
	//读excel流
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "读取excel文件失败")
		return
	}
	//解析excel的数据
	deviceModels, _, lxRrr := excels.ReadExce[models.DeviceModel](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	me := c.MustGet("user").(*models.User)
	var qty int = 0
	for _, entity := range deviceModels {
		// 循环体
		logger.Debug("--------------------------")
		logger.Debug(entity)
		var f models.DeviceModel = entity
		f.CreatedBy = me.Username
		f.UpdatedAt = f.CreatedAt
		f.Add(rt.Ctx)
		qty++
	}
	ginx.NewRender(c).Data(qty, nil)
}

// @Summary      导出设备型号数据
// @Description  根据条件导出数据
// @Tags         设备型号
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200  {object}  models.DeviceModel
// @Router       /api/n9e/device-model/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) exportDeviceModels(c *gin.Context) {

	// deviceType := ginx.QueryInt64(c, "deviceType", -1)
	// producer := ginx.QueryInt64(c, "producer", -1)
	// query := ginx.QueryStr(c, "query", "")

	f := make(map[string]interface{})
	ginx.BindJSON(c, &f)
	idsTemp, idsOk := f["ids"]
	ids := make([]int64, 0)
	var lst []models.DeviceModel
	var err error
	if idsOk {
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, int64(val.(float64)))
		}
		lst, err = models.DeviceModelGetsByIds(rt.Ctx, ids)
		ginx.Dangerous(err)
	} else {
		// m := make(map[string]interface{})
		// deviceType, deviceTypeOk := f["deviceType"]
		// if deviceTypeOk {
		// 	m["device_type"] = int64(deviceType.(float64))
		// }
		// producer, producerOk := f["producer"]
		// if producerOk {
		// 	m["producer_id"] = int64(producer.(float64))
		// }
		// query, queryOk := f["query"]
		// if queryOk {
		// 	m["query"] = query
		// }

		lst, err = models.DeviceModelGetsByType(rt.Ctx, f, -1, -1)
		ginx.Dangerous(err)
	}

	datas := make([]interface{}, 0)
	if len(lst) > 0 {
		for _, v := range lst {
			datas = append(datas, v)
		}
		excels.NewMyExcel("设备型号数据").ExportDataInfo(datas, "cn", rt.Ctx, c)
	} else {
		datas = append(datas, models.DeviceModel{})
		excels.NewMyExcel("设备型号数据").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
	}

}

// @Summary      导出设备型号模板
// @Description  导出设备型号模板
// @Tags         设备型号
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.DeviceModel
// @Router       /api/n9e/device-model/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templetDeviceModels(c *gin.Context) {

	datas := make([]interface{}, 0)

	datas = append(datas, models.DeviceModel{})

	excels.NewMyExcel("设备型号模板").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
}
