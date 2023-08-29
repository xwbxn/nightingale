package router

import (
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/logger"

	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/httpx"
)

type Router struct {
	HTTP           httpx.Config
	TargetCache    *memsto.TargetCacheType
	BusiGroupCache *memsto.BusiGroupCacheType
	AssetCache     *memsto.AssetCacheType
	Ctx            *ctx.Context
}

func New(httpConfig httpx.Config, tc *memsto.TargetCacheType, bg *memsto.BusiGroupCacheType, ac *memsto.AssetCacheType, ctx *ctx.Context) *Router {
	return &Router{
		HTTP:           httpConfig,
		Ctx:            ctx,
		TargetCache:    tc,
		BusiGroupCache: bg,
		AssetCache:     ac,
	}
}

func (rt *Router) Config(r *gin.Engine) {
	logger.Infof("Http Provider enable: %t", rt.HTTP.Provider.Enable)

	if !rt.HTTP.Provider.Enable {
		return
	}

	registerMetrics()

	provider := r.Group("/categraf")
	if len(rt.HTTP.Provider.BasicAuth) > 0 {
		provider.Use(gin.BasicAuth(rt.HTTP.Provider.BasicAuth))
	}
	// no need basic auth
	provider.GET("/configs", rt.categrafConfigGet)
	provider.GET("/upgrade", rt.targetVersionGet)
	provider.HEAD("/upgrade", rt.targetVersionHead)

}
