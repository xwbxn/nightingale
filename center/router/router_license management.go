package router

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

var PASSWORD = []byte("mypassword")

type License struct {
	Id             int64  `json:"id"`
	SerialNumber   string `json:"serial_number"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	PermissionNode int64  `json:"permission_node"`
	UsedNode       int64  `json:"used_node"`
}

// @Summary      获取证书列表
// @Description  获取证书列表
// @Tags         许可管理
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/n9e/xh/license/list [get]
// @Security     ApiKeyAuth
func (rt *Router) licenseGetsXH(c *gin.Context) {
	fileInfos, err := ioutil.ReadDir("etc/license/crt/")
	if err != nil {
		ginx.Bomb(http.StatusOK, "文件读取失败")
	}
	licenses := make([]License, 0)
	for _, fileInfo := range fileInfos {
		if strings.Contains(fileInfo.Name(), ".bak") {
			continue
		}
		var license License
		crtData, err := ioutil.ReadFile("etc/license/crt/" + fileInfo.Name())
		if err != nil {
			ginx.Bomb(http.StatusOK, "文件读取失败")
		}
		keyData, err := ioutil.ReadFile("etc/license/key/" + strings.Split(fileInfo.Name(), ".")[0] + ".key")
		if err != nil {
			ginx.Bomb(http.StatusOK, "文件读取失败")
		}
		// 解密证书并获取证书
		certBlock, _ := pem.Decode(crtData)
		// if certBlock == nil || certBlock.Type != "CERTIFICAET" {
		// 	ginx.Bomb(http.StatusOK, "证书文件格式错误")
		// 	return
		// }

		certificate, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			ginx.Bomb(http.StatusOK, "证书格式化失败")
		}
		id, _ := strconv.ParseInt(strings.Split(strings.Split(fileInfo.Name(), ".")[0], "-")[1], 10, 64)
		license.Id = id
		license.StartTime = certificate.NotBefore.Unix()
		license.EndTime = certificate.NotAfter.Unix()

		var value []byte
		for _, val := range certificate.Extensions {
			if val.Id.String() == "1.2.3.4" {

				_, err := asn1.Unmarshal(val.Value, &value)
				if err != nil {
					ginx.Bomb(http.StatusOK, "额外属性解析失败")
				}
				break
			}
		}

		// 解密私钥PEM并获取私钥
		keytBlock, _ := pem.Decode(keyData)
		// if keytBlock == nil || keytBlock.Type != "RSA PRIVATE KEY" {
		// 	ginx.Bomb(http.StatusOK, "私钥文件格式错误")
		// }
		decryptedPrivateKey, err := x509.DecryptPEMBlock(keytBlock, PASSWORD)
		if err != nil {
			ginx.Bomb(http.StatusOK, "私钥解密失败")
		}
		privateKey, err := x509.ParsePKCS1PrivateKey(decryptedPrivateKey)
		if err != nil {
			ginx.Bomb(http.StatusOK, "私钥格式化失败")
		}
		decryptedText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, value)
		if err != nil {
			ginx.Bomb(http.StatusOK, "解密失败")
		}
		md := make(map[string]interface{})
		err = json.Unmarshal(decryptedText, &md)
		if err != nil {
			fmt.Println("额外信息解析失败：", err)
		}
		logger.Debug(md)
		license.SerialNumber = md["serialNumber"].(string)
		license.PermissionNode = int64(md["permissionNode"].(float64))
		licenses = append(licenses, license)
	}
	ginx.NewRender(c).Data(licenses, nil)
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

	crtFilePath := "etc/license/crt/xihang-" + id + ".crt"
	keyFilePath := "etc/license/key/xihang-" + id + ".key"

	err = os.Rename(crtFilePath, "etc/license/crt/xihang-"+id+"-"+strconv.FormatInt(time.Now().Unix(), 10)+".crt"+".bak")
	ginx.Dangerous(err)
	err = os.Rename(keyFilePath, "etc/license/key/xihang-"+id+"-"+strconv.FormatInt(time.Now().Unix(), 10)+".key"+".bak")
	ginx.Dangerous(err)

	//保存新文件
	err = os.Remove(crtFilePath)
	ginx.Dangerous(err)
	err = os.Remove(keyFilePath)
	ginx.Dangerous(err)

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
