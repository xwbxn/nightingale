// Package models  配线架信息
// date : 2023-07-16 10:16
// desc : 配线架信息
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

// @Summary      获取配线架信息
// @Description  根据主键获取配线架信息
// @Tags         配线架信息
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DistributionFrame
// @Router       /api/n9e/distribution-frame/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) distributionFrameGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	distributionFrame, err := models.DistributionFrameGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if distributionFrame == nil {
		ginx.Bomb(404, "No such distribution_frame")
	}

	ginx.NewRender(c).Data(distributionFrame, nil)
}

// @Summary      查询配线架信息
// @Description  根据条件查询配线架信息
// @Tags         配线架信息
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.DistributionFrame
// @Router       /api/n9e/distribution-frame/list/ [get]
// @Security     ApiKeyAuth
func (rt *Router) distributionFrameGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DistributionFrameCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.DistributionFrameGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建配线架信息
// @Description  创建配线架信息
// @Tags         配线架信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.DistributionFrame true "add distributionFrame"
// @Success      200
// @Router       /api/n9e/distribution-frame/ [post]
// @Security     ApiKeyAuth
func (rt *Router) distributionFrameAdd(c *gin.Context) {
	var f models.DistributionFrame
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

// @Summary      更新配线架信息
// @Description  更新配线架信息
// @Tags         配线架信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.DistributionFrame true "update distributionFrame"
// @Success      200
// @Router       /api/n9e/distribution-frame/ [put]
// @Security     ApiKeyAuth
func (rt *Router) distributionFramePut(c *gin.Context) {
	var f models.DistributionFrame
	ginx.BindJSON(c, &f)

	// 生成一个validate实例
	validate := validator.New()
	errValidate := validate.Struct(f)
	if errValidate != nil {
		ginx.Bomb(http.StatusOK, errValidate.Error())
	}

	old, err := models.DistributionFrameGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "distribution_frame not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除配线架信息
// @Description  根据主键删除配线架信息
// @Tags         配线架信息
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/distribution-frame/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) distributionFrameDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	distributionFrame, err := models.DistributionFrameGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if distributionFrame == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(distributionFrame.Del(rt.Ctx))
}

// @Summary      EXCEL导入配线架信息
// @Description  EXCEL导入配线架信息
// @Tags         配线架信息
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/distribution-frame/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importDistributionFrame(c *gin.Context) {

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
	distributionFrames, lxRrr := excels.ReadExce[models.DistributionFrame](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	me := c.MustGet("user").(*models.User)
	var qty int = 0
	for _, entity := range distributionFrames {
		// 循环体
		var f models.DistributionFrame = entity
		f.Unumber++
		f.CreatedBy = me.Username
		f.UpdatedAt = f.CreatedAt
		f.Add(rt.Ctx)
		qty++
	}
	ginx.NewRender(c).Data(qty, nil)
}

// @Summary      EXCEL导出配线架信息
// @Description  EXCEL导出配线架信息
// @Tags         配线架信息
// @Accept       multipart/form-data
// @Produce      application/msexcel
// @Param        query query   string  false  "导入查询条件"
// @Success      200
// @Router       /api/n9e/distribution-frame/download-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) downloadDistributionFrame(c *gin.Context) {

	query := ginx.QueryStr(c, "query", "")
	list, err := models.DistributionFrameGets(rt.Ctx, query, 0, ginx.Offset(c, 0)) //获取数据

	ginx.Dangerous(err)

	datas := make([]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			datas = append(datas, v)

		}
	}
	excels.NewMyExcel("配线架数据").ExportDataInfo(datas, "cn", rt.Ctx, c)

}

// @Summary      导出配线架信息模板
// @Description  导出配线架信息模板
// @Tags         配线架信息
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.DistributionFrame
// @Router       /api/n9e/distribution-frame/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templetDistributionFrame(c *gin.Context) {

	datas := make([]interface{}, 0)

	datas = append(datas, models.DistributionFrame{})

	excels.NewMyExcel("设备厂商模板").ExportTempletToWeb(datas, "cn", "source", rt.Ctx, c)
}
