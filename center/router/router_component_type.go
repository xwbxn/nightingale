// Package models  部件类型
// date : 2023-08-21 09:08
// desc : 部件类型
package router

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取部件类型
// @Description  根据主键获取部件类型
// @Tags         部件类型
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.ComponentType
// @Router       /api/n9e/component-type/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) componentTypeGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	componentType, err := models.ComponentTypeGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if componentType == nil {
		ginx.Bomb(404, "No such component_type")
	}

	ginx.NewRender(c).Data(componentType, nil)
}

// @Summary      查询部件类型
// @Description  根据条件查询部件类型
// @Tags         部件类型
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "部件类型"
// @Success      200  {array}  models.ComponentType
// @Router       /api/n9e/component-type/ [get]
// @Security     ApiKeyAuth
func (rt *Router) componentTypeGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	comType := ginx.QueryStr(c, "query", "")

	m := make(map[string]interface{})
	if comType != "" {
		m["component_type"] = comType
	}

	total, err := models.ComponentTypeCountMap(rt.Ctx, m)
	ginx.Dangerous(err)
	lst, err := models.ComponentTypeGetMap(rt.Ctx, m, limit, (page-1)*limit)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建部件类型
// @Description  创建部件类型
// @Tags         部件类型
// @Accept       json
// @Produce      json
// @Param        body  body   models.ComponentType true "add componentType"
// @Success      200
// @Router       /api/n9e/component-type/ [post]
// @Security     ApiKeyAuth
func (rt *Router) componentTypeAdd(c *gin.Context) {
	var f models.ComponentType
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新部件类型
// @Description  更新部件类型
// @Tags         部件类型
// @Accept       json
// @Produce      json
// @Param        body  body   models.ComponentType true "update componentType"
// @Success      200
// @Router       /api/n9e/component-type/ [put]
// @Security     ApiKeyAuth
func (rt *Router) componentTypePut(c *gin.Context) {
	var f models.ComponentType
	ginx.BindJSON(c, &f)

	old, err := models.ComponentTypeGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "component_type not found")
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

// @Summary      批量删除部件类型
// @Description  批量删除部件类型
// @Tags         部件类型
// @Accept       json
// @Produce      json
// @Param        body  body   []models.ComponentType true "update componentType"
// @Success      200
// @Router       /api/n9e/component-type/batch-del/ [post]
// @Security     ApiKeyAuth
func (rt *Router) componentTypeBatchDel(c *gin.Context) {
	var f []models.ComponentType
	ginx.BindJSON(c, &f)

	err := models.ComponentTypeBatchDel(rt.Ctx, f)
	ginx.Dangerous(err)
	for _, val := range f {
		err = os.Remove(val.ComponentPicture)
		ginx.Dangerous(err)
	}

	ginx.NewRender(c).Message(err)
}

// @Summary      导入部件照片
// @Description  导入部件照片
// @Tags         部件类型
// @Accept       json
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/component-type/picture/ [post]
// @Security     ApiKeyAuth
func (rt *Router) componentTypePictureAdd(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "文件上传失败")
	}

	if fileHeader.Size > 1024*1024*5 {
		ginx.Bomb(http.StatusBadRequest, "文件超5MB")
	}
	fileName := strings.Split(fileHeader.Filename, ".")
	if fileName[len(fileName)-1] != "bmp" && fileName[len(fileName)-1] != "jpeg" && fileName[len(fileName)-1] != "jpg" && fileName[len(fileName)-1] != "png" {
		ginx.Bomb(http.StatusBadRequest, "文件格式错误")
	}
	// 设置路径,保存文件

	path := "etc/picture/"
	name := "component-type-" + strconv.FormatInt(time.Now().Unix(), 10) + "." + fileName[len(fileName)-1]

	_, err = PathExists(path)

	file_path := path + name
	c.SaveUploadedFile(fileHeader, file_path)

	ginx.NewRender(c).Data(file_path, err)
}
