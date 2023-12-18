package router

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取证书列表
// @Description  获取证书列表
// @Tags         许可管理
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/n9e/xh/license/list [get]
// @Security     ApiKeyAuth
func (rt *Router) licenseGetsXH(c *gin.Context) {
	ginx.NewRender(c).Data(rt.licenseCache.GetLicenseAll(), nil)
}

// @Summary      上传证书
// @Description  上传证书
// @Tags         许可管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        crt formData file true "file"
// @Param        key formData file true "file"
// @Success      200
// @Router       /api/n9e/xh/license/add-file [post]
// @Security     ApiKeyAuth
func (rt *Router) licenseAddXH(c *gin.Context) {
	crt, crtHeader, err := c.Request.FormFile("crt")
	ginx.Dangerous(err)
	key, keyHeader, err := c.Request.FormFile("key")
	ginx.Dangerous(err)
	now := strconv.FormatInt(time.Now().Unix(), 10)
	crtFileName := strings.Split(crtHeader.Filename, ".")[0] + "-" + now + ".crt"

	_, err = PathExists("etc/license/crt/")
	ginx.Dangerous(err)
	crtPath := "etc/license/crt/" + crtFileName
	// 创建目标文件
	dst, err := os.Create(crtPath)
	ginx.Dangerous(err)
	defer dst.Close()

	defer crt.Close()
	// 将文件内容复制到目标文件中
	if _, err := io.Copy(dst, crt); err != nil {
		ginx.Bomb(http.StatusOK, "保存文件失败")
	}

	keyFileName := strings.Split(keyHeader.Filename, ".")[0] + "-" + now + ".key"
	_, err = PathExists("etc/license/key/")
	ginx.Dangerous(err)
	keyPath := "etc/license/key/" + keyFileName
	// 创建目标文件
	dst, err = os.Create(keyPath)
	ginx.Dangerous(err)
	defer dst.Close()

	defer key.Close()
	// 将文件内容复制到目标文件中
	if _, err := io.Copy(dst, key); err != nil {
		ginx.Bomb(http.StatusOK, "保存文件失败")
	}
	ginx.NewRender(c).Message(nil)
}

// @Summary      更新证书
// @Description  更新证书
// @Tags         许可管理
// @Accept       multipart/form-data
// @Produce      json
// @Param        id formData string true "证书id"
// @Param        crt formData file false "file"
// @Param        key formData file false "file"
// @Success      200
// @Router       /api/n9e/xh/license/update [put]
// @Security     ApiKeyAuth
func (rt *Router) licenseUpdateXH(c *gin.Context) {
	id := c.Request.FormValue("id")
	logger.Debug(id)
	crt, _, err := c.Request.FormFile("crt")
	ginx.Dangerous(err)
	key, _, err := c.Request.FormFile("key")
	ginx.Dangerous(err)

	if crt == nil || key == nil {
		ginx.Bomb(http.StatusOK, "请上传证书或密钥")
	}
	// if strings.Contains(crtHeader.Filename, id) || strings.Contains(keyHeader.Filename, id) {
	// 	ginx.Bomb(http.StatusOK, "证书或密钥文件已存在")
	// }

	crtFilePath := "etc/license/crt/xihang-" + id + ".crt"
	keyFilePath := "etc/license/key/xihang-" + id + ".key"

	err = os.Rename(crtFilePath, "etc/license/crt/xihang-"+id+"-"+strconv.FormatInt(time.Now().Unix(), 10)+".crt"+".bak")
	ginx.Dangerous(err)
	err = os.Rename(keyFilePath, "etc/license/key/xihang-"+id+"-"+strconv.FormatInt(time.Now().Unix(), 10)+".key"+".bak")
	ginx.Dangerous(err)

	//保存新文件
	// 创建目标文件
	dst, err := os.Create(crtFilePath)
	ginx.Dangerous(err)
	defer dst.Close()

	defer crt.Close()
	// 将文件内容复制到目标文件中
	if _, err := io.Copy(dst, crt); err != nil {
		ginx.Bomb(http.StatusOK, "更新文件失败")
	}

	// 创建目标文件
	dst, err = os.Create(keyFilePath)
	ginx.Dangerous(err)
	defer dst.Close()

	defer key.Close()
	// 将文件内容复制到目标文件中
	if _, err := io.Copy(dst, key); err != nil {
		ginx.Bomb(http.StatusOK, "更新文件失败")
	}
	ginx.NewRender(c).Message(nil)
}

// @Summary      创建许可配置
// @Description  创建许可配置
// @Tags         许可管理
// @Accept       json
// @Produce      json
// @Param        body  body   models.LicenseConfig true "add licenseConfig"
// @Success      200
// @Router       /api/n9e/xh/license-config/ [post]
// @Security     ApiKeyAuth
func (rt *Router) licenseConfigAdd(c *gin.Context) {
	var f models.LicenseConfig
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新许可配置
// @Description  更新许可配置
// @Tags         许可管理
// @Accept       json
// @Produce      json
// @Param        body  body   models.LicenseConfig true "update licenseConfig"
// @Success      200
// @Router       /api/n9e/xh/license-config/ [put]
// @Security     ApiKeyAuth
func (rt *Router) licenseConfigPut(c *gin.Context) {
	var f models.LicenseConfig
	ginx.BindJSON(c, &f)

	old, err := models.LicenseConfigGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "许可配置未发现")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f))
}

// @Summary      获取许可配置
// @Description  获取许可配置
// @Tags         许可管理
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/n9e/xh/license-config [get]
// @Security     ApiKeyAuth
func (rt *Router) licenseConfigGetsXH(c *gin.Context) {
	licenseConfig, err := models.LicenseConfigGets(rt.Ctx)
	ginx.NewRender(c).Data(licenseConfig, err)
}

// @Summary      批量导出证书
// @Description  批量导出证书
// @Tags         许可管理
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]interface{} false "{“ids”:[111, 222]}"
// @Success      200
// @Router       /api/n9e/xh/license/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) exportLicenseXH(c *gin.Context) {
	var f map[string]interface{}
	ginx.BindJSON(c, &f)

	idsTemp, idsOk := f["ids"]
	ids := make([]int64, 0)
	if idsOk {
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, int64(val.(float64)))
		}
	}

	// 创建临时文件来存储压缩包

	zipFile, err := os.CreateTemp("", "download*.zip")
	ginx.Dangerous(err)
	defer os.Remove(zipFile.Name())

	// 创建 Zip Writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// licenses := make([]*License, 0)
	fileInfos, err := ioutil.ReadDir("etc/license/crt/")
	if err != nil {
		ginx.Bomb(http.StatusOK, "文件读取失败")
	}
	for _, fileInfo := range fileInfos {
		if strings.Contains(fileInfo.Name(), ".bak") || strings.Contains(fileInfo.Name(), "readme.txt") {
			continue
		}
		if len(ids) > 0 {
			for _, id := range ids {

				if strings.Contains(fileInfo.Name(), strconv.FormatInt(id, 10)) {
					file, err := os.Open("etc/license/crt/" + fileInfo.Name())
					ginx.Dangerous(err)
					defer file.Close()
					// 创建 Zip 文件
					zipEntry, err := zipWriter.Create(filepath.Base("etc/license/crt/" + fileInfo.Name()))
					ginx.Dangerous(err)
					// 将文件内容复制到 Zip 文件中
					_, err = io.Copy(zipEntry, file)
					ginx.Dangerous(err)

					//写入key
					file, err = os.Open("etc/license/key/" + strings.Split(fileInfo.Name(), ".")[0] + ".key")
					ginx.Dangerous(err)
					defer file.Close()
					// 创建 Zip 文件
					zipEntry, err = zipWriter.Create(filepath.Base("etc/license/key/" + strings.Split(fileInfo.Name(), ".")[0] + ".key"))
					ginx.Dangerous(err)
					// 将文件内容复制到 Zip 文件中
					_, err = io.Copy(zipEntry, file)
					ginx.Dangerous(err)
				}
			}
		} else {
			file, err := os.Open("etc/license/crt/" + fileInfo.Name())
			ginx.Dangerous(err)
			defer file.Close()
			// 创建 Zip 文件
			zipEntry, err := zipWriter.Create(filepath.Base("etc/license/crt/" + fileInfo.Name()))
			ginx.Dangerous(err)
			// 将文件内容复制到 Zip 文件中
			_, err = io.Copy(zipEntry, file)
			ginx.Dangerous(err)

			//写入key
			file, err = os.Open("etc/license/key/" + strings.Split(fileInfo.Name(), ".")[0] + ".key")
			ginx.Dangerous(err)
			defer file.Close()
			// 创建 Zip 文件
			zipEntry, err = zipWriter.Create(filepath.Base("etc/license/key/" + strings.Split(fileInfo.Name(), ".")[0] + ".key"))
			ginx.Dangerous(err)
			// 将文件内容复制到 Zip 文件中
			_, err = io.Copy(zipEntry, file)
			ginx.Dangerous(err)
		}

	}

	// 关闭zip文件写入器
	err = zipWriter.Close()
	ginx.Dangerous(err)
	// 设置响应头，告诉浏览器返回的是一个压缩包文件
	c.Header("Content-Type", "application/octet-stream")
	// 设置文件下载的名称
	c.Header("Content-Disposition", "attachment; filename=licenses.zip")
	// 将临时文件内容写入响应体
	c.File(zipFile.Name())
}
