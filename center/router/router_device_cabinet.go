// Package models  机柜信息
// date : 2023-07-11 13:56
// desc : 机柜信息
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

// @Summary      获取机柜信息
// @Description  根据主键获取机柜信息
// @Tags         机柜信息
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DeviceCabinet
// @Router       /api/n9e/device-cabinet/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceCabinetGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceCabinet, err := models.DeviceCabinetGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if deviceCabinet == nil {
		ginx.Bomb(404, "No such device_cabinet")
	}

	ginx.NewRender(c).Data(deviceCabinet, nil)
}

// @Summary      获取机柜名称列表
// @Description  根据机房ID获取机柜名称列表
// @Tags         机柜信息
// @Accept       json
// @Produce      json
// @Param        roomId    query    string  true  "机房ID"
// @Success      200  {object}  models.DeviceCabinetNameVo
// @Router       /api/n9e/device-cabinet/getNames/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceCabinetNameGet(c *gin.Context) {
	roomId := ginx.QueryInt64(c, "roomId", -1)
	deviceCabinet, err := models.DeviceCabinetGetByRoomId(rt.Ctx, roomId)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(deviceCabinet, nil)
}

// @Summary      查询机柜信息
// @Description  根据条件查询机柜信息
// @Tags         机柜信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceCabinet true "add deviceCabinet"
// @Success      200  {array}  models.DeviceCabinet
// @Router       /api/n9e/device-cabinet/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceCabinetGets(c *gin.Context) {
	// limit := ginx.QueryInt(c, "limit", 20)
	// query := ginx.QueryStr(c, "query", "")
	f := make(map[string]interface{})
	ginx.BindJSON(c, &f)
	var limit, page int
	limitTemp, limitOk := f["limit"]
	if !limitOk {
		limit = -1
	} else {
		limit = int(limitTemp.(float64))
	}
	pageTemp, pageOk := f["page"]
	if !pageOk {
		page = -1
	} else {
		page = int(pageTemp.(float64))
	}

	//TODO 满载机柜未判定
	// var spaceStatus string
	// status, statusOk := f["space_status"]
	// if statusOk {
	// 	spaceStatus = status.(string)
	// 	delete(f, "space_status")
	// }

	total, err := models.DeviceCabinetCountByMap(rt.Ctx, f)
	ginx.Dangerous(err)
	lst, err := models.DeviceCabinetGetByMap(rt.Ctx, f, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建机柜信息
// @Description  创建机柜信息
// @Tags         机柜信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceCabinet true "add deviceCabinet"
// @Success      200
// @Router       /api/n9e/device-cabinet/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceCabinetAdd(c *gin.Context) {
	var f models.DeviceCabinet
	ginx.BindJSON(c, &f)

	// 生成一个validate实例
	validate := validator.New()
	errValidate := validate.Struct(f)
	if errValidate != nil {
		ginx.Bomb(http.StatusOK, errValidate.Error())
	}

	//添加机柜编号默认值
	// if f.CabinetCode == "" {
	// 	//生成随机数
	// 	min := 100
	// 	max := 999
	// 	rand.Seed(time.Now().UnixNano())
	// 	num := rand.Intn(max-min-1) + min + 1

	// 	var build strings.Builder
	// 	build.WriteString(f.EquipmentRoom)
	// 	build.WriteString("_")
	// 	build.WriteString(strconv.Itoa(num))
	// 	f.CabinetCode = build.String()
	// }

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新机柜信息
// @Description  更新机柜信息
// @Tags         机柜信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceCabinet true "update deviceCabinet"
// @Success      200
// @Router       /api/n9e/device-cabinet/ [put]
// @Security     ApiKeyAuth
func (rt *Router) deviceCabinetPut(c *gin.Context) {
	var f models.DeviceCabinet
	ginx.BindJSON(c, &f)

	// 生成一个validate实例
	validate := validator.New()
	errValidate := validate.Struct(f)
	if errValidate != nil {
		ginx.Bomb(http.StatusOK, errValidate.Error())
	}

	old, err := models.DeviceCabinetGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "device_cabinet not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除机柜信息
// @Description  根据主键删除机柜信息
// @Tags         机柜信息
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/device-cabinet/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) deviceCabinetDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceCabinet, err := models.DeviceCabinetGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if deviceCabinet == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(deviceCabinet.Del(rt.Ctx))
}

// @Summary      EXCEL导入机柜信息
// @Description  EXCEL导入机柜信息
// @Tags         机柜信息
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/device-cabinet/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importDeviceCabinet(c *gin.Context) {

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

	deviceCabinets, _, lxRrr := excels.ReadExce[models.DeviceCabinet](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	me := c.MustGet("user").(*models.User)
	var qty int = 0
	for _, entity := range deviceCabinets {
		// 循环体
		var f models.DeviceCabinet = entity
		f.CreatedBy = me.Username
		f.UpdatedAt = f.CreatedAt
		f.Add(rt.Ctx)
		qty++
	}
	ginx.NewRender(c).Data(qty, nil)
}

// @Summary      EXCEL导出机柜信息
// @Description  EXCEL导出机柜信息
// @Tags         机柜信息
// @Accept       multipart/form-data
// @Produce      application/msexcel
// @Param        query query   string  false  "导入查询条件"
// @Success      200
// @Router       /api/n9e/device-cabinet/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) downloadDeviceCabinet(c *gin.Context) {

	query := ginx.QueryStr(c, "query", "")
	list, err := models.DeviceCabinetGets(rt.Ctx, query, -1, ginx.Offset(c, -1)) //获取数据
	ginx.Dangerous(err)

	datas := make([]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			datas = append(datas, v)

		}
	}
	excels.NewMyExcel("机柜数据").ExportDataInfo(datas, "cn", rt.Ctx, c)

}

// @Summary      导出机柜信息模板
// @Description  导出机柜信息模板
// @Tags         机柜信息
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.DeviceCabinet
// @Router       /api/n9e/device-cabinet/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templetDeviceCabinet(c *gin.Context) {

	datas := make([]interface{}, 0)

	datas = append(datas, models.DeviceCabinet{})

	excels.NewMyExcel("设备厂商模板").ExportTempletToWeb(datas, nil, "cn", "source", 1, rt.Ctx, c)
}
