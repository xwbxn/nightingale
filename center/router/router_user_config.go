// Package models  用户配置
// date : 2023-10-13 10:02
// desc : 用户配置
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	picture "github.com/ccfos/nightingale/v6/pkg/picture"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取用户配置
// @Description  根据主键获取用户配置
// @Tags         用户配置
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.UserConfig
// @Router       /api/n9e/user-config/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) userConfigGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	userConfig, err := models.UserConfigGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if userConfig == nil {
		ginx.Bomb(404, "No such user_config")
	}

	ginx.NewRender(c).Data(userConfig, nil)
}

// @Summary      查询用户配置
// @Description  根据条件查询用户配置
// @Tags         用户配置
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.UserConfig
// @Router       /api/n9e/user-config/ [get]
// @Security     ApiKeyAuth
func (rt *Router) userConfigGets(c *gin.Context) {

	lst, err := models.UserConfigGet(rt.Ctx)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(lst, nil)
}

// @Summary      查询用户登录标题
// @Description  查询用户登录标题
// @Tags         用户配置
// @Accept       json
// @Produce      json
// @Success      200  {array}  string
// @Router       /api/n9e/user-config/login/title [get]
// @Security     ApiKeyAuth
func (rt *Router) loginTitleGet(c *gin.Context) {

	ginx.NewRender(c).Data(models.LogoPathGet(rt.Ctx))
}

// @Summary      创建用户配置
// @Description  创建用户配置
// @Tags         用户配置
// @Accept       json
// @Produce      json
// @Param        body  body   models.UserConfig true "add userConfig"
// @Success      200
// @Router       /api/n9e/user-config/ [post]
// @Security     ApiKeyAuth
func (rt *Router) userConfigAdd(c *gin.Context) {
	var f models.UserConfig
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新用户配置
// @Description  更新用户配置
// @Tags         用户配置
// @Accept       json
// @Produce      json
// @Param        body  body   models.UserConfig true "update userConfig"
// @Success      200
// @Router       /api/n9e/user-config/ [put]
// @Security     ApiKeyAuth
func (rt *Router) userConfigPut(c *gin.Context) {
	var f models.UserConfig
	ginx.BindJSON(c, &f)

	old, err := models.UserConfigGet(rt.Ctx)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "user_config not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	if f.LogLever != old.LogLever {
		old.LogLever = f.LogLever
	}
	if f.AccessExpired != old.AccessExpired {
		old.AccessExpired = f.AccessExpired
	}
	if f.RefreshExpired != old.RefreshExpired {
		old.RefreshExpired = f.RefreshExpired
	}
	if f.HttpHost != old.HttpHost {
		old.HttpHost = f.HttpHost
	}
	if 8000 <= f.HttpPort && f.HttpPort <= 65535 {
		if old.HttpPort == f.HttpPort {
			old.HttpPort = f.HttpPort
		}
	} else {
		ginx.Bomb(http.StatusOK, "Http端口不合法!")
	}
	if f.OpenRsa != old.OpenRsa {
		old.OpenRsa = f.OpenRsa
	}
	if f.ApiService != old.ApiService {
		old.ApiService = f.ApiService
	}
	if f.Captcha != old.Captcha {
		old.Captcha = f.Captcha
	}

	err = old.Update(rt.Ctx, "LOG_LEVER", "HTTP_HOST", "HTTP_PORT", "CAPTCHA", "API_SERVICE", "ACCESS_EXPIRED", "REFRESH_EXPIRED", "OPEN_RSA")
	ginx.NewRender(c).Message(err)
}

// @Summary      删除用户配置
// @Description  根据主键删除用户配置
// @Tags         用户配置
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/user-config/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) userConfigDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	userConfig, err := models.UserConfigGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if userConfig == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(userConfig.Del(rt.Ctx))
}

// @Summary      选择图片上传
// @Description  选择图片上传
// @Tags         用户配置
// @Accept       json
// @Produce      json
// @Param        picture  formData file false "图片"
// @Param        logoName query    string  false  "图片名称"
// @Success      200
// @Router       /api/n9e/user-config/picture [post]
// @Security     ApiKeyAuth
func (rt *Router) userPictureAdd(c *gin.Context) {
	_, fileHeader, err := c.Request.FormFile("picture")
	filename := ginx.QueryStr(c, "logoName")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "文件上传失败")
	}

	suffix, err := picture.VerifyPicture(fileHeader)
	ginx.Dangerous(err)

	filePath, err, fakePath := picture.GenerateLogoName(filename, suffix)
	ginx.Dangerous(err)

	c.SaveUploadedFile(fileHeader, filePath)

	ginx.NewRender(c).Data(fakePath, err)
}

// @Summary      更新用户Logo
// @Description  更新用户Logo
// @Tags         用户配置
// @Accept       json
// @Produce      json
// @Param        body  body   models.LogoConfig true "更新Logo信息"
// @Success      200
// @Router       /api/n9e/user-config/logo [put]
// @Security     ApiKeyAuth
func (rt *Router) logoPut(c *gin.Context) {
	var f models.LogoConfig
	ginx.BindJSON(c, &f)

	old, err := models.UserLogoConfigGet(rt.Ctx)

	if f.LoginTitle != old.LoginTitle {
		old.LoginTitle = f.LoginTitle
	}
	if f.LogoTop != old.LogoTop {
		old.LogoTop = f.LogoTop
	}
	if f.LogoTitle != old.LogoTitle {
		old.LogoTitle = f.LogoTitle
	}

	err = old.UserLogoUpdate(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      查询logo信息
// @Description  查询logo信息
// @Tags         用户配置
// @Accept       json
// @Produce      json
// @Success      200 {array}  models.LogoConfig
// @Router       /api/n9e/user-config/getInfo [get]
// @Security     ApiKeyAuth
func (rt *Router) userLogoGet(c *gin.Context) {
	//pictureTag := ginx.QueryInt(c, "tag", 0)
	//设置文件类型

	res, err := models.LogoPathGet(rt.Ctx)

	ginx.NewRender(c).Data(res, err)

}
