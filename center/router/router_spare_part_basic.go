// Package models  备件基础数据
// date : 2023-08-20 15:43
// desc : 备件基础数据
package router

import (
	"net/http"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	picture "github.com/ccfos/nightingale/v6/pkg/picture"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取备件基础数据
// @Description  根据主键获取备件基础数据
// @Tags         备件基础数据
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.SparePartBasic
// @Router       /api/n9e/spare-part_basic/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) sparePartBasicGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	sparePartBasic, err := models.SparePartBasicGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if sparePartBasic == nil {
		ginx.Bomb(404, "No such spare_part_basic")
	}

	ginx.NewRender(c).Data(sparePartBasic, nil)
}

// @Summary      查询备件基础数据
// @Description  根据条件查询备件基础数据
// @Tags         备件基础数据
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.SparePartBasic
// @Router       /api/n9e/spare-part_basic/ [get]
// @Security     ApiKeyAuth
func (rt *Router) sparePartBasicGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.SparePartBasicCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.SparePartBasicGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建备件基础数据
// @Description  创建备件基础数据
// @Tags         备件基础数据
// @Accept       json
// @Produce      json
// @Param        body  body   models.SparePartBasic true "add sparePartBasic"
// @Success      200
// @Router       /api/n9e/spare-part_basic/ [post]
// @Security     ApiKeyAuth
func (rt *Router) sparePartBasicAdd(c *gin.Context) {
	var f models.SparePartBasic
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新备件基础数据
// @Description  更新备件基础数据
// @Tags         备件基础数据
// @Accept       json
// @Produce      json
// @Param        body  body   models.SparePartBasic true "update sparePartBasic"
// @Success      200
// @Router       /api/n9e/spare-part_basic/ [put]
// @Security     ApiKeyAuth
func (rt *Router) sparePartBasicPut(c *gin.Context) {
	var f models.SparePartBasic
	ginx.BindJSON(c, &f)

	old, err := models.SparePartBasicGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "spare_part_basic not found")
	}
	oldPicture := old.ComponentPicture

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	err = old.Update(rt.Ctx, f, "*")
	ginx.Dangerous(err)
	if oldPicture != f.ComponentPicture {
		err = os.Remove(oldPicture)
	}
	ginx.NewRender(c).Message(err)
}

// @Summary      批量删除备件基础数据
// @Description  批量删除备件基础数据
// @Tags         备件基础数据
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "batch delete sparePartBasic"
// @Success      200
// @Router       /api/n9e/spare-part_basic/batch-del/ [post]
// @Security     ApiKeyAuth
func (rt *Router) sparePartBasicBatchDel(c *gin.Context) {
	var f []int64
	ginx.BindJSON(c, &f)

	sparePartBasicpictures, err := models.ComponentPictureGetById(rt.Ctx, f)
	ginx.Dangerous(err)

	err = models.SparePartBasicBatchDel(rt.Ctx, f)
	ginx.Dangerous(err)
	for _, sparePartBasicpicture := range sparePartBasicpictures {
		if sparePartBasicpicture.ComponentPicture == "" {
			continue
		}
		err = os.Remove(sparePartBasicpicture.ComponentPicture)
		ginx.Dangerous(err)
	}

	ginx.NewRender(c).Message(err)
}

// @Summary      导入部件照片
// @Description  导入部件照片
// @Tags         备件基础数据
// @Accept       json
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/spare-part_basic/picture/ [post]
// @Security     ApiKeyAuth
func (rt *Router) sparePartBasicPictureAdd(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "文件上传失败")
	}

	suffix, err := picture.VerifyPicture(fileHeader)
	ginx.Dangerous(err)

	filePath, err := picture.GeneratePictureName("spare-part", suffix)
	ginx.Dangerous(err)

	c.SaveUploadedFile(fileHeader, filePath)

	ginx.NewRender(c).Data(filePath, err)
}

// @Summary      导入备件基础信息数据
// @Description  导入备件基础信息数据
// @Tags         备件基础数据
// @Accept       multipart/form-data
// @Param        file formData file true "file"
// @Produce      json
// @Success      200
// @Router       /api/n9e/spare-part_basic/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importsSparePartBasic(c *gin.Context) {
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
	sparePartBasics, _, lxRrr := excels.ReadExce[models.SparePartBasic](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	me := c.MustGet("user").(*models.User)
	var qty int = 0
	for _, entity := range sparePartBasics {
		// 循环体
		logger.Debug("--------------------------")
		logger.Debug(entity)
		var f models.SparePartBasic = entity
		f.CreatedBy = me.Username
		f.UpdatedAt = f.CreatedAt
		f.Add(rt.Ctx)
		qty++
	}
	ginx.NewRender(c).Data(qty, nil)
}

// @Summary      导出备件基础信息数据
// @Description  导出备件基础信息数据
// @Tags         备件基础数据
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.SparePartBasic
// @Router       /api/n9e/spare-part_basic/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) exportSparePartBasic(c *gin.Context) {

	list, err := models.SparePartBasicGetsAll(rt.Ctx) //获取数据
	ginx.Dangerous(err)

	datas := make([]interface{}, 0)
	if len(list) > 0 {
		for _, v := range list {
			datas = append(datas, v)

		}
	}
	excels.NewMyExcel("备件基础信息数据导出").ExportDataInfo(datas, "cn", rt.Ctx, c)

}

// @Summary      导出备件基础信息模板
// @Description  导出备件基础信息模板
// @Tags         备件基础数据
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.SparePartBasic
// @Router       /api/n9e/spare-part_basic/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templetSparePartBasic(c *gin.Context) {

	datas := make([]interface{}, 0)

	datas = append(datas, models.SparePartBasic{})

	excels.NewMyExcel("备件基础信息导入模板").ExportTempletToWeb(datas, nil, "cn", "source", 1, rt.Ctx, c)
}
