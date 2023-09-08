// Package models  设备型号
// date : 2023-07-08 14:57
// desc : 设备型号
package router

import (
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"

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
// @Success      200  {array}  models.DeviceModelDetailsVo
// @Router       /api/n9e/device-model/getmodel/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceModelGets(c *gin.Context) {
	deviceType := ginx.QueryInt(c, "deviceType", -1)

	lst, err := models.DeviceModelGetsByType(rt.Ctx, int64(deviceType))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": len(lst),
	}, nil)
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

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除设备型号
// @Description  根据主键删除设备型号
// @Tags         设备型号
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/device-model/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) deviceModelDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceModel, err := models.DeviceModelGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if deviceModel == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(deviceModel.Del(rt.Ctx))
}

// @Summary      导入设备型号数据
// @Description  导入型号
// @Tags         设备型号
// @Accept       multipart/form-data
// @Param        file formData file true "file"
// @Produce      json
// @Success      200
// @Router       /api/n9e/device-model/import [post]
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
	deviceModels, lxRrr := excels.ReadExce[models.DeviceModel](xlsx, rt.Ctx)
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
// @Param        query query   string  false  "查询条件"
// @Success      200  {object}  models.DeviceModel
// @Router       /api/n9e/device-model/outport [post]
// @Security     ApiKeyAuth
func (rt *Router) exportDeviceModels(c *gin.Context) {

	query := ginx.QueryStr(c, "query", "")
	list, err := models.DeviceModelGets(rt.Ctx, query, -1, ginx.Offset(c, -1)) //获取数据
	ginx.Dangerous(err)

	datas := make([]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			datas = append(datas, v)

		}
	}
	excels.NewMyExcel("设备型号数据").ExportDataInfo(datas, "cn", rt.Ctx, c)

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

	excels.NewMyExcel("设备型号模板").ExportTempletToWeb(datas, "cn", "source", rt.Ctx, c)
}
