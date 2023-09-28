// Package models  库房信息
// date : 2023-08-21 14:10
// desc : 库房信息
package router

import (
	"net/http"

	"github.com/360EntSecGroup-Skylar/excelize"
	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取库房信息
// @Description  根据主键获取库房信息
// @Tags         库房信息
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.StoreroomManagement
// @Router       /api/n9e/storeroom-management/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) storeroomManagementGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	storeroomManagement, err := models.StoreroomManagementGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if storeroomManagement == nil {
		ginx.Bomb(404, "No such storeroom_management")
	}

	ginx.NewRender(c).Data(storeroomManagement, nil)
}

// @Summary      查询库房信息
// @Description  根据条件查询库房信息
// @Tags         库房信息
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.StoreroomManagement
// @Router       /api/n9e/storeroom-management/ [get]
// @Security     ApiKeyAuth
func (rt *Router) storeroomManagementGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.StoreroomNumAddressCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.StoreroomNumAddressGets(rt.Ctx, query, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建库房信息
// @Description  创建库房信息
// @Tags         库房信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.StoreroomManagement true "add storeroomManagement"
// @Success      200
// @Router       /api/n9e/storeroom-management/ [post]
// @Security     ApiKeyAuth
func (rt *Router) storeroomManagementAdd(c *gin.Context) {
	var f models.StoreroomManagement
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新库房信息
// @Description  更新库房信息
// @Tags         库房信息
// @Accept       json
// @Produce      json
// @Param        body  body   models.StoreroomManagement true "update storeroomManagement"
// @Success      200
// @Router       /api/n9e/storeroom-management/ [put]
// @Security     ApiKeyAuth
func (rt *Router) storeroomManagementPut(c *gin.Context) {
	var f models.StoreroomManagement
	ginx.BindJSON(c, &f)

	old, err := models.StoreroomManagementGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "storeroom_management not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      批量删除库房信息
// @Description  批量删除库房信息
// @Tags         库房信息
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "batch delete storeroomManagement"
// @Success      200
// @Router       /api/n9e/storeroom-management/batch-del/ [post]
// @Security     ApiKeyAuth
func (rt *Router) storeroomManagementBatchDel(c *gin.Context) {
	var f []int64
	ginx.BindJSON(c, &f)

	ginx.NewRender(c).Message(models.StoreroomManagementBatchDel(rt.Ctx, f))
}

// @Summary      导入库房信息
// @Description  导入库房信息
// @Tags         库房信息
// @Accept       multipart/form-data
// @Param        file formData file true "file"
// @Produce      json
// @Success      200
// @Router       /api/n9e/storeroom-management/import [post]
// @Security     ApiKeyAuth
func (rt *Router) importsStoreroomManagement(c *gin.Context) {
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
	storeroomManagements, _, lxRrr := excels.ReadExce[models.StoreroomManagement](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	me := c.MustGet("user").(*models.User)
	var qty int = 0
	for _, entity := range storeroomManagements {
		// 循环体
		var f models.StoreroomManagement = entity
		f.CreatedBy = me.Username
		f.UpdatedAt = f.CreatedAt
		f.Add(rt.Ctx)
		qty++
	}
	ginx.NewRender(c).Data(qty, nil)
}

// @Summary      导出库房信息模板
// @Description  导出库房信息模板
// @Tags         库房信息
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.StoreroomManagement
// @Router       /api/n9e/storeroom-management/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templetStoreroomManagement(c *gin.Context) {

	datas := make([]interface{}, 0)

	datas = append(datas, models.StoreroomManagement{})

	excels.NewMyExcel("库房信息导入模板").ExportTempletToWeb(datas, nil, "cn", "source", rt.Ctx, c)
}
