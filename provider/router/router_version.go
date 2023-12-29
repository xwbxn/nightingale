package router

import (
	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func (rt *Router) categrafVersionHead(c *gin.Context) {
	new_version := getNewVersion(c, rt)
	c.Header("Client-Version", new_version.Version)
}

func (rt *Router) categrafVersionGet(c *gin.Context) {
	new_version := getNewVersion(c, rt)
	c.File(new_version.Path)
}

func getNewVersion(c *gin.Context, rt *Router) *models.TargetVersion {
	host := ginx.QueryStr(c, "id", "")
	version := ginx.QueryStr(c, "version", "")
	os := ginx.QueryStr(c, "os", "linux")
	arch := ginx.QueryStr(c, "arch", "amd64")

	if host == "" || version == "" {
		ginx.Bomb(404, "not found")
	}

	target, has := rt.TargetCache.Get(host)
	if !has {
		ginx.Bomb(404, "not found")
	}

	if target.AgentVersion == "" {
		ginx.Bomb(200, "not upgradable")
	}

	if target.AgentVersion == version {
		ginx.Bomb(200, "already latest")
	}

	new_version, err := models.TargetVersionGet(rt.Ctx, target.AgentVersion, os, arch)
	if err != nil {
		ginx.Bomb(500, "latest version error")
	}
	return new_version
}
