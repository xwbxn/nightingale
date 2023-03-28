package router

import (
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/logger"

	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/httpx"
	"github.com/ccfos/nightingale/v6/provider/cpconf"
)

type Router struct {
	HTTP           httpx.Config
	TargetCache    *memsto.TargetCacheType
	BusiGroupCache *memsto.BusiGroupCacheType
	ProviderCache  *memsto.ProviderCacheType
	Ctx            *ctx.Context
}

func New(httpConfig httpx.Config, provider cpconf.Provider, tc *memsto.TargetCacheType, bg *memsto.BusiGroupCacheType, hpc *memsto.ProviderCacheType, ctx *ctx.Context) *Router {
	return &Router{
		HTTP:           httpConfig,
		Ctx:            ctx,
		TargetCache:    tc,
		BusiGroupCache: bg,
		ProviderCache:  hpc,
	}
}

func (rt *Router) Config(r *gin.Engine) {
	logger.Infof("Http Provider enable: %t", rt.HTTP.Provider.Enable)

	if !rt.HTTP.Provider.Enable {
		return
	}

	registerMetrics()

	if len(rt.HTTP.Provider.BasicAuth) > 0 {
		// enable basic auth
		auth := gin.BasicAuth(rt.HTTP.Provider.BasicAuth)
		r.GET("/categraf/configs", auth, rt.CategrafConfigGet)
	} else {
		// no need basic auth
		r.GET("/categraf/configs", rt.CategrafConfigGet)
	}
}
