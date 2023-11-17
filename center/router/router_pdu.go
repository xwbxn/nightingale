// Package models  PDU
// date : 2023-07-16 10:15
// desc : PDU
package router

import (
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取PDU
// @Description  根据主键获取PDU
// @Tags         PDU
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.Pdu
// @Router       /api/n9e/pdu/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) pduGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	pdu, err := models.PduGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if pdu == nil {
		ginx.Bomb(404, "No such pdu")
	}

	ginx.NewRender(c).Data(pdu, nil)
}

// @Summary      查询PDU
// @Description  根据条件查询PDU
// @Tags         PDU
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.Pdu
// @Router       /api/n9e/pdu/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) pduGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.PduCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.PduGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建PDU
// @Description  创建PDU
// @Tags         PDU
// @Accept       json
// @Produce      json
// @Param        body  body   models.Pdu true "add pdu"
// @Success      200
// @Router       /api/n9e/pdu/ [post]
// @Security     ApiKeyAuth
func (rt *Router) pduAdd(c *gin.Context) {
	var f models.Pdu
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新PDU
// @Description  更新PDU
// @Tags         PDU
// @Accept       json
// @Produce      json
// @Param        body  body   models.Pdu true "update pdu"
// @Success      200
// @Router       /api/n9e/pdu/ [put]
// @Security     ApiKeyAuth
func (rt *Router) pduPut(c *gin.Context) {
	var f models.Pdu
	ginx.BindJSON(c, &f)

	old, err := models.PduGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "pdu not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除PDU
// @Description  根据主键删除PDU
// @Tags         PDU
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/pdu/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) pduDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	pdu, err := models.PduGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if pdu == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(pdu.Del(rt.Ctx))
}

// @Summary      EXCEL导入PDU
// @Description  EXCEL导入PDU
// @Tags         PDU
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/pdu/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importpdu(c *gin.Context) {

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

	pdus, _, lxRrr := excels.ReadExce[models.Pdu](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	me := c.MustGet("user").(*models.User)
	var qty int = 0
	for _, entity := range pdus {
		// 循环体
		var f models.Pdu = entity
		f.CreatedBy = me.Username
		f.UpdatedAt = f.CreatedAt
		f.Add(rt.Ctx)
		qty++
	}
	ginx.NewRender(c).Data(qty, nil)
}

// @Summary      EXCEL导出PDU
// @Description  EXCEL导出PDU
// @Tags         PDU
// @Accept       multipart/form-data
// @Produce      application/msexcel
// @Param        query query   string  false  "导入查询条件"
// @Success      200
// @Router       /api/n9e/pdu/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) downloadpdu(c *gin.Context) {

	query := ginx.QueryStr(c, "query", "")
	list, err := models.PduGets(rt.Ctx, query, -1, ginx.Offset(c, -1)) //获取数据

	ginx.Dangerous(err)

	datas := make([]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			datas = append(datas, v)

		}
	}
	excels.NewMyExcel("PDU数据").ExportDataInfo(datas, "cn", rt.Ctx, c)

}

// @Summary      导出PDU模板
// @Description  导出PDU模板
// @Tags         PDU
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Pdu
// @Router       /api/n9e/pdu/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templetpdu(c *gin.Context) {

	datas := make([]interface{}, 0)

	datas = append(datas, models.Pdu{})

	excels.NewMyExcel("PDU模板").ExportTempletToWeb(datas, nil, "cn", "source", 1, rt.Ctx, c)
}
