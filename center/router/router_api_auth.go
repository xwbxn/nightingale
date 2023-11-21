package router

import (
	"encoding/json"

	"github.com/ccfos/nightingale/v6/pkg/secu"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

type ApiAuth struct {
	Ip       string
	passWord string
}

// @Summary      创建接口管理
// @Description  创建接口管理
// @Tags         第三方接口
// @Accept       json
// @Produce      json
// @Param        body  body   ApiAuth true "add query"
// @Success      200
// @Router       /api/n9e/xh/assets/filter [post]
// @Security     ApiKeyAuth
func (rt *Router) apiAuthAdd(c *gin.Context) {
	var f ApiAuth
	ginx.BindJSON(c, &f)

	jsonData, err := rt.Redis.Get(rt.Ctx.Ctx, "apiAuth").Bytes()
	if err != nil {
		if err.Error() != "redis: nil" {
			ginx.Dangerous(err)
		}
	}
	data := make(map[string]string)
	err = json.Unmarshal(jsonData, &data)
	ginx.Dangerous(err)

	dataByte, err := json.Marshal(data)
	ginx.Dangerous(err)
	key, _ := json.Marshal("it`s api auth!")

	dataPwd, err := secu.AesEncrypt(dataByte, key)
	ginx.Dangerous(err)

	data[f.Ip] = secu.BASE64StdEncode(dataPwd)

}

// @Summary      获取密码
// @Description  获取密码
// @Tags         第三方接口
// @Accept       json
// @Produce      json
// @Param        body  body   ApiAuth true "add query"
// @Success      200
// @Router       /api/n9e/xh/api-auth/getpwd [get]
// @Security     ApiKeyAuth
func (rt *Router) getApiAuth(c *gin.Context) {

}
