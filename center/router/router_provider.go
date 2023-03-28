package router

import (
	"encoding/json"
	"net/http"

	"github.com/ccfos/nightingale/v6/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toolkits/pkg/ginx"
)

type ConfigFormat string

const (
	YamlFormat ConfigFormat = "yaml"
	TomlFormat ConfigFormat = "toml"
	JsonFormat ConfigFormat = "json"
)

type ConfigWithFormat struct {
	Config string       `json:"config"`
	Format ConfigFormat `json:"format"`
}

func (rt *Router) providerGet(c *gin.Context) {
	groupId := ginx.UrlParamInt64(c, "id")
	provider, err := models.ProviderGetById(rt.Ctx, groupId)
	vm := providerModel{
		Id:      provider.Id,
		Ident:   provider.Ident,
		GroupId: provider.GroupId,
	}
	err = json.Unmarshal([]byte(provider.Configs), &vm.Configs)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(vm, err)
}

type providerModel struct {
	Id      int64                       `json:"id"`
	Ident   string                      `json:"ident"`
	GroupId int64                       `json:"group_id"`
	Configs map[string]ConfigWithFormat `json:"configs"`
}

func (rt *Router) providerAdd(c *gin.Context) {
	var f providerModel
	ginx.BindJSON(c, &f)

	var provider = models.Provider{
		Version: newVersion(),
		Ident:   f.Ident,
		GroupId: f.GroupId,
	}
	configs, err := json.Marshal(f.Configs)
	ginx.Dangerous(err)

	provider.Configs = string(configs)
	err = provider.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

func (rt *Router) providerPut(c *gin.Context) {
	var f providerModel
	ginx.BindJSON(c, &f)
	oldConfig, err := models.ProviderGet(rt.Ctx, "id=?", f.Id)
	ginx.Dangerous(err)

	if oldConfig == nil {
		ginx.Bomb(http.StatusOK, "provider not found")
	}

	if oldConfig.Ident != f.Ident && oldConfig.GroupId != f.GroupId {
		// check duplication
		num, err := models.RoleCount(rt.Ctx, "ident=? and group_id=? and id<>?", f.Ident, f.GroupId, oldConfig.Id)
		ginx.Dangerous(err)

		if num > 0 {
			ginx.Bomb(http.StatusOK, "provider name already exists")
		}
	}

	oldConfig.Ident = f.Ident
	oldConfig.GroupId = f.GroupId
	oldConfig.Version = newVersion()
	configs, err := json.Marshal(f.Configs)
	ginx.Dangerous(err)
	oldConfig.Configs = string(configs)

	ginx.NewRender(c).Message(oldConfig.Update(rt.Ctx, "ident", "group_id", "version", "configs"))
}

func (rt *Router) providerDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	target, err := models.ProviderGet(rt.Ctx, "id=?", id)
	ginx.Dangerous(err)

	if target == nil {
		ginx.NewRender(c).Message(nil)
		return
	}

	ginx.NewRender(c).Message(target.Del(rt.Ctx))
}

func newVersion() string {
	u1, _ := uuid.NewUUID()
	return u1.String()
}
