package router

import (
	"encoding/json"
	"strconv"
	"time"

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
// @Router       /api/n9e/xh/api-auth/add [post]
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

	data[f.Ip] = f.passWord
	// 将数组转换为JSON字符串
	newData, err := json.Marshal(data)
	ginx.Dangerous(err)
	err = rt.Redis.Set(rt.Ctx.Ctx, "apiAuth", newData, 0).Err()
	ginx.NewRender(c).Message(err)
}

// @Summary      获取密码
// @Description  获取密码
// @Tags         第三方接口
// @Accept       json
// @Produce      json
// @Param        pwd  query   string     true  "add"
// @Success      200
// @Router       /api/n9e/xh/api-auth/getpwd [get]
// @Security     ApiKeyAuth
func (rt *Router) getApiAuth(c *gin.Context) {
	pwd := ginx.QueryStr(c, "pwd", "")
	dataByte, err := json.Marshal(pwd + strconv.FormatInt(time.Now().Unix(), 10))
	ginx.Dangerous(err)
	key, _ := json.Marshal("it`s api auth!")

	dataPwd, err := secu.AesEncrypt(dataByte, key)
	ginx.Dangerous(err)

	res := secu.BASE64StdEncode(dataPwd)

	ginx.NewRender(c).Data(res, err)
}
