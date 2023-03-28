package router

import (
	"encoding/json"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
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

type httpRemoteProviderResponse struct {
	// version is signature/md5 of current Config, server side should deal with the Version calculate
	Version string `json:"version"`

	// ConfigMap (InputName -> Config), if version is identical, server side can set Config to nil
	Configs map[string]ConfigWithFormat `json:"configs"`
}

func (rt *Router) CategrafConfigGet(c *gin.Context) {
	ident := ginx.QueryStr(c, "ident")
	busigroup := ginx.QueryStr(c, "busigroup")

	group := rt.BusiGroupCache.GetByBusiGroupLabel(busigroup)
	if group == nil {
		ginx.Dangerous("config not found", 404)
	}

	provider := rt.ProviderCache.GetByIdentAndGroup(ident, group.Id)
	resp := convertToResponse(provider)
	c.JSON(200, resp)
}

func convertToResponse(model *models.Provider) httpRemoteProviderResponse {
	resp := httpRemoteProviderResponse{
		Version: model.Version,
	}
	err := json.Unmarshal([]byte(model.Configs), &resp.Configs)
	if err != nil {
		logger.Warningf("decode config %s fail: %s", model.Id, err.Error())
	}
	return resp
}
