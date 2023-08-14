// Package models  设备型号
// date : 2023-07-08 14:57
// desc : 设备型号
package router

import (
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

var (
	defaultSheetName = "Sheet1" //默认Sheet名称
	defaultHeight    = 25.0     //默认行高度
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
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.DeviceModel
// @Router       /api/n9e/device-model/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceModelGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DeviceModelCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.DeviceModelGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
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
	deviceModels, lxRrr := readExcel(xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	var qty int = 0
	for _, entity := range deviceModels {
		// 循环体
		var f models.DeviceModel = entity
		f.Add(rt.Ctx)
		qty++
	}
	ginx.NewRender(c).Data(qty, nil)
}

//ReadExcel .读取excel 转成切片
func readExcel(xlsx *excelize.File, ctx *ctx.Context) ([]models.DeviceModel, error) {
	//根据名字获取cells的内容，返回的是一个[][]string
	rows := xlsx.GetRows(xlsx.GetSheetName(xlsx.GetActiveSheetIndex()))
	//声明一个数组
	var deviceModels []models.DeviceModel
	fields := reflect.ValueOf(new(models.DeviceModel)).Elem()
	mapLit := make(map[int]string)
	for i, row := range rows {
		//去掉第一行是excel表头部分
		if i == 0 { //取得第一行的所有数据---execel表头
			for index, colCell := range row {
				mapLit[index] = colCell
			}
		} else {
			entity := &models.DeviceModel{}
			g := reflect.ValueOf(entity).Elem()
			for index, colCell := range row {
				title := mapLit[index]
				for i := 0; i < fields.NumField(); i++ {
					fieldInfo := fields.Type().Field(i)

					_, heardOk := fieldInfo.Tag.Lookup("cn")
					_, sourceOk := fieldInfo.Tag.Lookup("source")

					if heardOk && (fieldInfo.Tag.Get("cn") == title) {
						var results []int64
						var isDB = false
						if sourceOk && fieldInfo.Tag.Get("source") == "table" {
							isDB = true
							value := strings.Split(fieldInfo.Tag.Get("value"), ",")
							session := models.DB(ctx)
							session.Table(value[0]).Where(value[1]+" = ?", colCell).Pluck("id", &results)
						}

						switch fieldType := fieldInfo.Type.Kind(); fieldType {
						case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
							{
								if isDB {
									if len(results) > 0 {
										g.FieldByName(fieldInfo.Name).SetInt(results[0])
									}
								} else {
									s1, _ := strconv.Atoi(colCell)
									g.FieldByName(fieldInfo.Name).SetInt(int64(s1))
								}
							}
						case reflect.String:
							g.FieldByName(fieldInfo.Name).SetString(colCell)
						case reflect.Bool:
							g.FieldByName(fieldInfo.Name).SetBool(colCell == "true")
						default:
							log.Printf("field type %s not support yet", fieldType)
						}
					}
				}

			}
			deviceModels = append(deviceModels, *entity)
		}

	}
	return deviceModels, nil
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
	list, err := models.DeviceModelGets(rt.Ctx, query, 0, ginx.Offset(c, 0)) //获取数据
	ginx.Dangerous(err)

	datas := make([]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			datas = append(datas, v)

		}
	}
	excels.NewMyExcel("设备型号数据").ExportDataToWeb(datas, "cn", c)

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
