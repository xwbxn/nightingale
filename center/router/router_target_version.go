// Package models  探针版本
// date : 2023-08-30 10:58
// desc : 探针版本
package router

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      导入最新探针
// @Description  导入最新探针
// @Tags         探针版本
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/target/version/ [post]
// @Security     ApiKeyAuth
func (rt *Router) importNewVersion(c *gin.Context) {

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "文件上传失败")
	}

	if fileHeader.Size > 1024*1024*50 {
		ginx.Bomb(http.StatusBadRequest, "文件超50MB")
	}
	fileName := strings.Split(fileHeader.Filename, "-")
	if len(fileName) != 4 {
		ginx.Bomb(http.StatusBadRequest, "文件命名错误")
	}
	if fileName[2] != "linux" && fileName[2] != "windows" && fileName[2] != "darwin" {
		ginx.Bomb(http.StatusBadRequest, "文件命名错误(操作系统)")
	}
	arch := strings.Split(fileName[3], ".")
	if arch[0] != "amd64" && arch[0] != "386" && arch[0] != "arm" && arch[0] != "arm64" {
		ginx.Bomb(http.StatusBadRequest, "文件命名错误(架构)")
	}
	if arch[len(arch)-1] != "gz" {
		ginx.Bomb(http.StatusBadRequest, "文件命名错误(后缀)")
	}

	content, err := io.ReadAll(file)
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, fmt.Sprintf("读取失败: %s", err.Error()))
		return
	}
	// 解析文件类型
	kind, err := filetype.Match(content)
	logger.Debug(kind.MIME.Subtype)
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, fmt.Sprintf("文件类型判断失败: %s", err.Error()))
		return
	}
	logger.Debug(kind.MIME.Value)
	if kind.MIME.Subtype != "gzip" {
		ginx.Bomb(http.StatusBadRequest, "文件格式错误")
	}
	// 设置路径,保存文件

	path := "etc/client/"

	_, err = PathExists(path)

	file_path := path + fileHeader.Filename
	c.SaveUploadedFile(fileHeader, file_path)

	ginx.NewRender(c).Message(err)
}

// @Summary      获取可用探针版本
// @Description  根据ident获取可用探针版本
// @Tags         探针版本
// @Accept       json
// @Produce      json
// @Param        ident    path    string  true  "ident"
// @Success      200  {array} array
// @Router       /api/n9e/target/{ident}/version [get]
// @Security     ApiKeyAuth
func (rt *Router) UsableVersionGet(c *gin.Context) {
	ident := ginx.UrlParamStr(c, "ident")
	target, has := rt.TargetCache.Get(ident)
	if !has {
		ginx.Bomb(404, "target not found")
	}
	os := target.OS
	arch := target.Arch
	FileInfo, err := ioutil.ReadDir("etc/client/")
	ginx.Dangerous(err)
	var array []string
	for index := range FileInfo {
		name := FileInfo[index].Name()
		partName := strings.Split(name, "-")
		archName := strings.Split(partName[len(partName)-1], ".")[0]
		if os == partName[2] && arch == archName {
			array = append(array, partName[1])
		}
	}

	ginx.NewRender(c).Data(array, nil)
}

// @Summary      更新所选探针
// @Description  更新所选探针
// @Tags         探针版本
// @Accept       json
// @Produce      json
// @Param        ident    path    string  true  "ident"
// @Param        body  body   map[string]interface{} true "update map"
// @Success      200
// @Router       /api/n9e/target/{ident}/version [put]
// @Security     ApiKeyAuth
func (rt *Router) targetVersionPut(c *gin.Context) {

	ident := ginx.UrlParamStr(c, "ident")

	var f map[string]interface{}
	ginx.BindJSON(c, &f)

	f["update_at"] = time.Now().Unix()
	err := models.UpdataVersion(rt.Ctx, f, ident)

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(err)
}

// @Summary      导入升级包
// @Description  导入升级包
// @Tags         探针版本
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/server/update/ [post]
// @Security     ApiKeyAuth
func (rt *Router) importUpgradePack(c *gin.Context) {

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "文件上传失败")
	}

	//将multipart.File转化为os.File
	tmpFile, err := ioutil.TempFile("", "file")
	ginx.Dangerous(err)
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)
	ginx.Dangerous(err)
	defer tmpFile.Close()

	//打开文件
	tarFile, err := os.Open(tmpFile.Name())
	ginx.Dangerous(err)
	defer tarFile.Close()

	//读取gzip压缩文件
	gr, err := gzip.NewReader(tarFile)
	ginx.Dangerous(err)
	defer gr.Close()

	fi, err := tarFile.Stat()
	ginx.Dangerous(err)

	path := "etc/update/"

	_, err = PathExists(path)

	//判断当前文件夹是否存在“n9e”
	now := time.Now()
	var newFileName string
	FileInfo, err := ioutil.ReadDir(path)
	ginx.Dangerous(err)
	for index := range FileInfo {
		name := FileInfo[index].Name()
		if name == "n9e" {
			timeUnix := strconv.FormatInt(now.Unix(), 10)
			newFileName = "/n9e-" + timeUnix + ".bak"
			newPath := path + newFileName
			err := os.Rename(path+"n9e", newPath)
			ginx.Dangerous(err)
		}
	}

	//创建“n9e”文件
	w, err := os.OpenFile(path+"n9e", os.O_CREATE|os.O_TRUNC|os.O_RDWR, fi.Mode())
	ginx.Dangerous(err)

	//写入“n9e”文件
	_, err = io.Copy(w, gr)
	ginx.Dangerous(err)
	defer w.Close()

	os.Chtimes(newFileName, now, gr.ModTime)

	ginx.NewRender(c).Message(err)

}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			ginx.Dangerous(err)
		} else {
			return true, nil
		}
	}
	return false, err
}
