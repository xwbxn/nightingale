// Package models  设备类型
// date : 2023-07-08 11:49
// desc : 设备类型
package router

import (
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取设备类型
// @Description  根据主键获取设备类型
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DeviceType
// @Router       /api/n9e/device-type/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceType, err := models.DeviceTypeGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if deviceType == nil {
		ginx.Bomb(404, "No such device_type")
	}

	ginx.NewRender(c).Data(deviceType, nil)
}

// @Summary      查询设备类型
// @Description  根据条件查询设备类型
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "查询条件"
// @Param        types query   string  true  "类别"
// @Success      200  {array}  models.DeviceType
// @Router       /api/n9e/device-type/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	deviceType := ginx.QueryStr(c, "query", "")
	types := ginx.QueryStr(c, "types", "")

	m := make(map[string]interface{})
	m["types"] = types
	if deviceType != "" {
		m["name"] = deviceType
	}

	total, err := models.DeviceTypeCountMap(rt.Ctx, m)
	ginx.Dangerous(err)
	lst, err := models.DeviceTypeGetMap(rt.Ctx, m, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建设备类型
// @Description  创建设备类型
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceType true "add deviceType"
// @Success      200
// @Router       /api/n9e/device-type/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeAdd(c *gin.Context) {
	var f models.DeviceType
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新设备类型
// @Description  更新设备类型
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceType true "update deviceType"
// @Success      200
// @Router       /api/n9e/device-type/ [put]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypePut(c *gin.Context) {
	var f models.DeviceType
	ginx.BindJSON(c, &f)

	old, err := models.DeviceTypeGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "device_type not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      批量删除设备类型
// @Description  批量删除设备类型
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceType true "update deviceType"
// @Success      200
// @Router       /api/n9e/device-type/batch-del/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceTypeBatchDel(c *gin.Context) {
	var f []models.DeviceType
	ginx.BindJSON(c, &f)

	ginx.NewRender(c).Message(models.DeviceTypeBatchDel(rt.Ctx, f))
}

// @Summary      导入设备类型
// @Description  导入设备类型
// @Tags         设备类型
// @Accept       multipart/form-data
// @Param        file formData file true "file"
// @Produce      json
// @Success      200
// @Router       /api/n9e/device-type/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importsDeviceType(c *gin.Context) {
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
	deviceTypes, _, lxRrr := excels.ReadExce[models.DeviceType](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	me := c.MustGet("user").(*models.User)
	var qty int = 0
	for _, entity := range deviceTypes {
		// 循环体
		var f models.DeviceType = entity
		f.Types = 2
		f.CreatedBy = me.Username
		f.UpdatedAt = f.CreatedAt
		f.Add(rt.Ctx)
		qty++
	}
	ginx.NewRender(c).Data(qty, nil)
}

// @Summary      导出设备类型模板
// @Description  导出设备类型模板
// @Tags         设备类型
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.DeviceType
// @Router       /api/n9e/device-type/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templetDeviceType(c *gin.Context) {

	datas := make([]interface{}, 0)

	datas = append(datas, models.DeviceType{})

	excels.NewMyExcel("设备类型导入模板").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
}
