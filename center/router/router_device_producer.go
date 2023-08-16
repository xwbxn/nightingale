// Package models  设备厂商
// date : 2023-07-08 14:51
// desc : 设备厂商
package router

import (
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/toolkits/pkg/ginx"
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

// @Summary      查询设备厂商
// @Description  根据条件查询设备厂商
// @Tags         设备厂商
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.DeviceProducer
// @Router       /api/n9e/device-producer/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceProducerGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DeviceProducerCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.DeviceProducerGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
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

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      EXCEL导入设备厂商
// @Description  EXCEL导入设备厂商
// @Tags         设备厂商
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/device-producer/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importDeviceProducer(c *gin.Context) {

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

	deviceProducers, lxRrr := excels.ReadExce[models.DeviceProducer](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	me := c.MustGet("user").(*models.User)
	var qty int = 0
	for _, entity := range deviceProducers {
		// 循环体
		var f models.DeviceProducer = entity
		f.CreatedBy = me.Username
		f.UpdatedAt = f.CreatedAt
		f.Add(rt.Ctx)
		qty++
	}
	ginx.NewRender(c).Data(qty, nil)
}

// @Summary      EXCEL导出设备厂商
// @Description  EXCEL导出设备厂商
// @Tags         设备厂商
// @Accept       multipart/form-data
// @Produce      application/msexcel
// @Param        query query   string  false  "导入查询条件"
// @Success      200
// @Router       /api/n9e/device-producer/download-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) downloadDeviceProducer(c *gin.Context) {

	query := ginx.QueryStr(c, "query", "")
	list, err := models.DeviceProducerGets(rt.Ctx, query, 0, ginx.Offset(c, 0)) //获取数据

	ginx.Dangerous(err)

	datas := make([]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			datas = append(datas, v)

		}
	}
	excels.NewMyExcel("设备厂商数据").ExportDataInfo(datas, "cn", rt.Ctx, c)

}

// @Summary      导出设备厂商模板
// @Description  导出设备厂商模板
// @Tags         设备厂商
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.DeviceProducer
// @Router       /api/n9e/device-producer/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templetDeviceProducer(c *gin.Context) {

	datas := make([]interface{}, 0)

	datas = append(datas, models.DeviceProducer{})

	excels.NewMyExcel("设备厂商模板").ExportTempletToWeb(datas, "cn", "source", rt.Ctx, c)
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
