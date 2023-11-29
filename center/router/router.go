package router

import (
	"fmt"
	"net/http"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/center/cconf"
	"github.com/ccfos/nightingale/v6/center/cstats"
	"github.com/ccfos/nightingale/v6/center/metas"
	"github.com/ccfos/nightingale/v6/center/sso"
	_ "github.com/ccfos/nightingale/v6/front/statik"
	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/pkg/aop"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/httpx"
	"github.com/ccfos/nightingale/v6/pkg/version"
	"github.com/ccfos/nightingale/v6/prom"
	"github.com/ccfos/nightingale/v6/pushgw/idents"
	"github.com/ccfos/nightingale/v6/storage"

	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/runner"

	_ "github.com/ccfos/nightingale/v6/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	HTTP              httpx.Config
	Center            cconf.Center
	Operations        cconf.Operation
	DatasourceCache   *memsto.DatasourceCacheType
	NotifyConfigCache *memsto.NotifyConfigCacheType
	PromClients       *prom.PromClientMap
	Redis             storage.Redis
	MetaSet           *metas.Set
	IdentSet          *idents.Set
	TargetCache       *memsto.TargetCacheType
	Sso               *sso.SsoClient
	UserCache         *memsto.UserCacheType
	UserGroupCache    *memsto.UserGroupCacheType
	Ctx               *ctx.Context
	assetCache        *memsto.AssetCacheType

	DatasourceCheckHook func(*gin.Context) bool
}

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(httpConfig httpx.Config, center cconf.Center, operations cconf.Operation, ds *memsto.DatasourceCacheType, ncc *memsto.NotifyConfigCacheType,
	pc *prom.PromClientMap, redis storage.Redis, sso *sso.SsoClient, ctx *ctx.Context, metaSet *metas.Set, idents *idents.Set, tc *memsto.TargetCacheType,
	uc *memsto.UserCacheType, ugc *memsto.UserGroupCacheType, ac *memsto.AssetCacheType) *Router {
	return &Router{
		HTTP:              httpConfig,
		Center:            center,
		Operations:        operations,
		DatasourceCache:   ds,
		NotifyConfigCache: ncc,
		PromClients:       pc,
		Redis:             redis,
		MetaSet:           metaSet,
		IdentSet:          idents,
		TargetCache:       tc,
		Sso:               sso,
		UserCache:         uc,
		UserGroupCache:    ugc,
		Ctx:               ctx,
		assetCache:        ac,

		DatasourceCheckHook: func(ctx *gin.Context) bool { return false },
	}
}

func stat() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		code := fmt.Sprintf("%d", c.Writer.Status())
		method := c.Request.Method
		labels := []string{cstats.Service, code, c.FullPath(), method}

		cstats.RequestCounter.WithLabelValues(labels...).Inc()
		cstats.RequestDuration.WithLabelValues(labels...).Observe(float64(time.Since(start).Seconds()))
	}
}

func languageDetector(i18NHeaderKey string) gin.HandlerFunc {
	headerKey := i18NHeaderKey
	return func(c *gin.Context) {
		if headerKey != "" {
			lang := c.GetHeader(headerKey)
			if lang != "" {
				if strings.HasPrefix(lang, "zh") {
					c.Request.Header.Set("X-Language", "zh")
				} else if strings.HasPrefix(lang, "en") {
					c.Request.Header.Set("X-Language", "en")
				} else {
					c.Request.Header.Set("X-Language", lang)
				}
			} else {
				c.Request.Header.Set("X-Language", "en")
			}
		}
		c.Next()
	}
}

func (rt *Router) configNoRoute(r *gin.Engine, fs *http.FileSystem) {
	r.NoRoute(func(c *gin.Context) {
		arr := strings.Split(c.Request.URL.Path, ".")
		suffix := arr[len(arr)-1]

		if strings.HasPrefix(arr[0], "/prod-api") {
			cwdarr := []string{"/"}
			if runtime.GOOS == "windows" {
				cwdarr[0] = ""
			}
			cwdarr = append(cwdarr, strings.Split(runner.Cwd, "/")...)
			cwdarr = append(cwdarr, "public")
			filename := strings.ReplaceAll(c.Request.URL.Path, "/prod-api", "")
			if filename == "" || filename == "/" {
				filename = "/index.html"
			}
			cwdarr = append(cwdarr, strings.Split(filename, "/")...)
			c.File(path.Join(cwdarr...))
		} else {
			switch suffix {
			case "png", "jpeg", "jpg", "svg", "ico", "gif", "css", "js", "html", "htm", "gz", "zip", "map", "ttf":
				if !rt.Center.UseFileAssets {
					c.FileFromFS(c.Request.URL.Path, *fs)
				} else {
					cwdarr := []string{"/"}
					if runtime.GOOS == "windows" {
						cwdarr[0] = ""
					}
					cwdarr = append(cwdarr, strings.Split(runner.Cwd, "/")...)
					cwdarr = append(cwdarr, "pub")
					cwdarr = append(cwdarr, strings.Split(c.Request.URL.Path, "/")...)
					c.File(path.Join(cwdarr...))
				}
			default:
				if !rt.Center.UseFileAssets {
					c.FileFromFS("/", *fs)
				} else {
					cwdarr := []string{"/"}
					if runtime.GOOS == "windows" {
						cwdarr[0] = ""
					}
					cwdarr = append(cwdarr, strings.Split(runner.Cwd, "/")...)
					cwdarr = append(cwdarr, "pub")
					cwdarr = append(cwdarr, "index.html")
					c.File(path.Join(cwdarr...))
				}
			}
		}
	})
}

func (rt *Router) Config(r *gin.Engine) {

	r.Use(stat())
	r.Use(languageDetector(rt.Center.I18NHeaderKey))
	r.Use(aop.Recovery())
	r.Use(aop.LoggerToFile())

	statikFS, err := fs.New()
	if err != nil {
		logger.Errorf("cannot create statik fs: %v", err)
	}

	if !rt.Center.UseFileAssets {
		r.StaticFS("/pub", statikFS)
	}

	pagesPrefix := "/api/n9e"
	pages := r.Group(pagesPrefix)
	{

		if rt.Center.AnonymousAccess.PromQuerier {
			pages.Any("/proxy/:id/*url", rt.dsProxy)
			pages.POST("/query-range-batch", rt.promBatchQueryRange)
			pages.POST("/query-instant-batch", rt.promBatchQueryInstant)
			pages.GET("/datasource/brief", rt.datasourceBriefs)
		} else {
			pages.Any("/proxy/:id/*url", rt.auth(), rt.dsProxy)
			pages.POST("/query-range-batch", rt.auth(), rt.promBatchQueryRange)
			pages.POST("/query-instant-batch", rt.auth(), rt.promBatchQueryInstant)
			pages.GET("/datasource/brief", rt.auth(), rt.datasourceBriefs)
		}

		pages.POST("/auth/login", rt.jwtMock(), rt.loginPost)
		pages.POST("/auth/logout", rt.jwtMock(), rt.auth(), rt.logoutPost)
		pages.POST("/auth/refresh", rt.jwtMock(), rt.refreshPost)
		pages.POST("/auth/captcha", rt.jwtMock(), rt.generateCaptcha)
		pages.POST("/auth/captcha-verify", rt.jwtMock(), rt.captchaVerify)
		pages.GET("/auth/ifshowcaptcha", rt.ifShowCaptcha)

		pages.GET("/auth/sso-config", rt.ssoConfigNameGet)
		pages.GET("/auth/rsa-config", rt.rsaConfigGet)
		pages.GET("/auth/redirect", rt.loginRedirect)
		pages.GET("/auth/redirect/cas", rt.loginRedirectCas)
		pages.GET("/auth/redirect/oauth", rt.loginRedirectOAuth)
		pages.GET("/auth/callback", rt.loginCallback)
		pages.GET("/auth/callback/cas", rt.loginCallbackCas)
		pages.GET("/auth/callback/oauth", rt.loginCallbackOAuth)
		pages.GET("/auth/perms", rt.allPerms)

		pages.GET("/metrics/desc", rt.metricsDescGetFile)
		pages.POST("/metrics/desc", rt.metricsDescGetMap)

		pages.GET("/notify-channels", rt.notifyChannelsGets)
		pages.GET("/contact-keys", rt.contactKeysGets)

		pages.GET("/self/perms", rt.auth(), rt.user(), rt.permsGets)
		pages.GET("/self/profile", rt.auth(), rt.user(), rt.selfProfileGet)
		pages.PUT("/self/profile", rt.auth(), rt.user(), rt.selfProfilePut)
		pages.PUT("/self/password", rt.auth(), rt.user(), rt.selfPasswordPut)

		pages.GET("/users", rt.auth(), rt.admin(), rt.perm("/users"), rt.userGets)
		pages.GET("/users/xh", rt.auth(), rt.admin(), rt.userGetsXH)
		pages.POST("/users", rt.auth(), rt.admin(), rt.userAddPost)
		pages.GET("/user/:id/profile", rt.auth(), rt.userProfileGet)
		pages.GET("/user/getNames", rt.auth(), rt.userNameGets)
		pages.PUT("/user/:id/profile", rt.auth(), rt.admin(), rt.userProfilePut)
		pages.PUT("/user/:id/password", rt.auth(), rt.admin(), rt.userPasswordPut)
		pages.DELETE("/user/:id", rt.auth(), rt.admin(), rt.userDel)
		pages.POST("/users/update-property", rt.auth(), rt.admin(), rt.userPropertyUpdate) //批量修改用户状态/组织
		pages.POST("/users/batch-del", rt.auth(), rt.admin(), rt.userDels)                 //批量删除用户
		pages.POST("/xh/users/templet", rt.auth(), rt.user(), rt.templeUserXH)
		pages.POST("/xh/users/import-xls", rt.auth(), rt.user(), rt.importUserXH)

		pages.GET("/metric-views", rt.auth(), rt.metricViewGets)
		pages.DELETE("/metric-views", rt.auth(), rt.user(), rt.metricViewDel)
		pages.POST("/metric-views", rt.auth(), rt.user(), rt.metricViewAdd)
		pages.PUT("/metric-views", rt.auth(), rt.user(), rt.metricViewPut)

		pages.GET("/user-groups", rt.auth(), rt.user(), rt.userGroupGets)
		pages.POST("/user-groups", rt.auth(), rt.user(), rt.perm("/user-groups/add"), rt.userGroupAdd)
		pages.GET("/user-group/:id", rt.auth(), rt.user(), rt.userGroupGet)
		pages.PUT("/user-group/:id", rt.auth(), rt.user(), rt.perm("/user-groups/put"), rt.userGroupWrite(), rt.userGroupPut)
		pages.DELETE("/user-group/:id", rt.auth(), rt.user(), rt.perm("/user-groups/del"), rt.userGroupWrite(), rt.userGroupDel)
		pages.POST("/user-group/:id/members", rt.auth(), rt.user(), rt.perm("/user-groups/put"), rt.userGroupWrite(), rt.userGroupMemberAdd)
		pages.DELETE("/user-group/:id/members", rt.auth(), rt.user(), rt.perm("/user-groups/put"), rt.userGroupWrite(), rt.userGroupMemberDel)

		pages.GET("/busi-groups", rt.auth(), rt.user(), rt.busiGroupGets)
		pages.POST("/busi-groups", rt.auth(), rt.user(), rt.perm("/busi-groups/add"), rt.busiGroupAdd)
		pages.GET("/busi-groups/alertings", rt.auth(), rt.busiGroupAlertingsGets)
		pages.GET("/busi-group/:id", rt.auth(), rt.user(), rt.bgro(), rt.busiGroupGet)
		pages.PUT("/busi-group/:id", rt.auth(), rt.user(), rt.perm("/busi-groups/put"), rt.bgrw(), rt.busiGroupPut)
		pages.POST("/busi-group/:id/members", rt.auth(), rt.user(), rt.perm("/busi-groups/put"), rt.bgrw(), rt.busiGroupMemberAdd)
		pages.DELETE("/busi-group/:id/members", rt.auth(), rt.user(), rt.perm("/busi-groups/put"), rt.bgrw(), rt.busiGroupMemberDel)
		pages.DELETE("/busi-group/:id", rt.auth(), rt.user(), rt.perm("/busi-groups/del"), rt.bgrw(), rt.busiGroupDel)
		pages.GET("/busi-group/:id/perm/:perm", rt.auth(), rt.user(), rt.checkBusiGroupPerm)

		pages.GET("/targets", rt.auth(), rt.user(), rt.targetGets)
		pages.POST("/target/list", rt.auth(), rt.user(), rt.targetGetsByHostFilter)
		pages.DELETE("/targets", rt.auth(), rt.user(), rt.perm("/targets/del"), rt.targetDel)
		pages.GET("/targets/tags", rt.auth(), rt.user(), rt.targetGetTags)
		pages.POST("/targets/tags", rt.auth(), rt.user(), rt.perm("/targets/put"), rt.targetBindTagsByFE)
		pages.DELETE("/targets/tags", rt.auth(), rt.user(), rt.perm("/targets/put"), rt.targetUnbindTagsByFE)
		pages.PUT("/targets/note", rt.auth(), rt.user(), rt.perm("/targets/put"), rt.targetUpdateNote)
		pages.PUT("/targets/bgid", rt.auth(), rt.user(), rt.perm("/targets/put"), rt.targetUpdateBgid)

		pages.POST("/builtin-cate-favorite", rt.auth(), rt.user(), rt.builtinCateFavoriteAdd)
		pages.DELETE("/builtin-cate-favorite/:name", rt.auth(), rt.user(), rt.builtinCateFavoriteDel)

		pages.GET("/builtin-boards", rt.builtinBoardGets)
		pages.GET("/builtin-board/:name", rt.builtinBoardGet)
		pages.GET("/dashboards/builtin/list", rt.builtinBoardGets)
		pages.GET("/builtin-boards-cates", rt.auth(), rt.user(), rt.builtinBoardCateGets)
		pages.POST("/builtin-boards-detail", rt.auth(), rt.user(), rt.builtinBoardDetailGets)
		pages.GET("/integrations/icon/:cate/:name", rt.builtinIcon)
		pages.GET("/integrations/makedown/:cate", rt.builtinMarkdown)

		pages.GET("/busi-group/:id/boards", rt.auth(), rt.user(), rt.perm("/dashboards"), rt.bgro(), rt.boardGets)
		pages.POST("/busi-group/:id/boards", rt.auth(), rt.user(), rt.perm("/dashboards/add"), rt.bgrw(), rt.boardAdd)
		pages.POST("/busi-group/:id/board/:bid/clone", rt.auth(), rt.user(), rt.perm("/dashboards/add"), rt.bgrw(), rt.boardClone)

		pages.GET("/board/:bid", rt.boardGet)
		pages.GET("/board/:bid/pure", rt.boardPureGet)
		pages.PUT("/board/:bid", rt.auth(), rt.user(), rt.perm("/dashboards/put"), rt.boardPut)
		pages.PUT("/board/:bid/configs", rt.auth(), rt.user(), rt.perm("/dashboards/put"), rt.boardPutConfigs)
		pages.PUT("/board/:bid/public", rt.auth(), rt.user(), rt.perm("/dashboards/put"), rt.boardPutPublic)
		pages.DELETE("/boards", rt.auth(), rt.user(), rt.perm("/dashboards/del"), rt.boardDel)

		pages.GET("/share-charts", rt.chartShareGets)
		pages.POST("/share-charts", rt.auth(), rt.chartShareAdd)

		pages.GET("/alert-rules/builtin/alerts-cates", rt.auth(), rt.user(), rt.builtinAlertCateGets)
		pages.GET("/alert-rules/builtin/list", rt.auth(), rt.user(), rt.builtinAlertRules)

		pages.GET("/busi-group/:id/alert-rules", rt.auth(), rt.user(), rt.perm("/alert-rules"), rt.alertRuleGets)
		pages.GET("/busi-group/:id/alert-rules/xh", rt.auth(), rt.user(), rt.alertRuleGetsXH)
		pages.POST("/busi-group/:id/alert-rules", rt.auth(), rt.user(), rt.perm("/alert-rules/add"), rt.bgrw(), rt.alertRuleAddByFE)
		pages.POST("/busi-group/:id/alert-rules/import", rt.auth(), rt.user(), rt.perm("/alert-rules/add"), rt.bgrw(), rt.alertRuleAddByImport)
		pages.DELETE("/busi-group/:id/alert-rules", rt.auth(), rt.user(), rt.perm("/alert-rules/del"), rt.bgrw(), rt.alertRuleDel)
		pages.PUT("/busi-group/:id/alert-rules/fields", rt.auth(), rt.user(), rt.perm("/alert-rules/put"), rt.bgrw(), rt.alertRulePutFields)
		pages.PUT("/busi-group/:id/alert-rule/:arid", rt.auth(), rt.user(), rt.perm("/alert-rules/put"), rt.alertRulePutByFE)
		pages.GET("/alert-rule/:arid", rt.auth(), rt.user(), rt.perm("/alert-rules"), rt.alertRuleGet)
		pages.POST("/alert-rule/arid", rt.auth(), rt.user(), rt.alertRuleGetsByIds)
		pages.PUT("/busi-group/:id/alert-rule/:arid/validate", rt.auth(), rt.user(), rt.perm("/alert-rules/put"), rt.alertRuleValidation)

		pages.GET("/busi-group/:id/recording-rules", rt.auth(), rt.user(), rt.perm("/recording-rules"), rt.recordingRuleGets)
		pages.POST("/busi-group/:id/recording-rules", rt.auth(), rt.user(), rt.perm("/recording-rules/add"), rt.bgrw(), rt.recordingRuleAddByFE)
		pages.DELETE("/busi-group/:id/recording-rules", rt.auth(), rt.user(), rt.perm("/recording-rules/del"), rt.bgrw(), rt.recordingRuleDel)
		pages.PUT("/busi-group/:id/recording-rule/:rrid", rt.auth(), rt.user(), rt.perm("/recording-rules/put"), rt.bgrw(), rt.recordingRulePutByFE)
		pages.GET("/recording-rule/:rrid", rt.auth(), rt.user(), rt.perm("/recording-rules"), rt.recordingRuleGet)
		pages.PUT("/busi-group/:id/recording-rules/fields", rt.auth(), rt.user(), rt.perm("/recording-rules/put"), rt.recordingRulePutFields)

		pages.GET("/busi-group/:id/alert-mutes", rt.auth(), rt.user(), rt.perm("/alert-mutes"), rt.bgro(), rt.alertMuteGetsByBG)
		pages.POST("/busi-group/:id/alert-mutes", rt.auth(), rt.user(), rt.perm("/alert-mutes/add"), rt.bgrw(), rt.alertMuteAdd)
		pages.DELETE("/busi-group/:id/alert-mutes", rt.auth(), rt.user(), rt.perm("/alert-mutes/del"), rt.bgrw(), rt.alertMuteDel)
		pages.PUT("/busi-group/:id/alert-mute/:amid", rt.auth(), rt.user(), rt.perm("/alert-mutes/put"), rt.alertMutePutByFE)
		pages.PUT("/busi-group/:id/alert-mutes/fields", rt.auth(), rt.user(), rt.perm("/alert-mutes/put"), rt.bgrw(), rt.alertMutePutFields)

		pages.GET("/busi-group/:id/alert-subscribes", rt.auth(), rt.user(), rt.perm("/alert-subscribes"), rt.bgro(), rt.alertSubscribeGets)
		pages.GET("/alert-subscribe/:sid", rt.auth(), rt.user(), rt.perm("/alert-subscribes"), rt.alertSubscribeGet)
		pages.POST("/busi-group/:id/alert-subscribes", rt.auth(), rt.user(), rt.perm("/alert-subscribes/add"), rt.bgrw(), rt.alertSubscribeAdd)
		pages.PUT("/busi-group/:id/alert-subscribes", rt.auth(), rt.user(), rt.perm("/alert-subscribes/put"), rt.bgrw(), rt.alertSubscribePut)
		pages.DELETE("/busi-group/:id/alert-subscribes", rt.auth(), rt.user(), rt.perm("/alert-subscribes/del"), rt.bgrw(), rt.alertSubscribeDel)

		if rt.Center.AnonymousAccess.AlertDetail {
			pages.GET("/alert-cur-event/:eid", rt.alertCurEventGet)
			pages.GET("/alert-his-event/:eid", rt.alertHisEventGet)
		} else {
			pages.GET("/alert-cur-event/:eid", rt.auth(), rt.alertCurEventGet)
			pages.GET("/alert-his-event/:eid", rt.auth(), rt.alertHisEventGet)
		}

		// card logic
		pages.GET("/alert-cur-events/list", rt.auth(), rt.alertCurEventsList)
		pages.GET("/alert-events/list/xh", rt.auth(), rt.alertEventsListXH)
		pages.GET("/alert-cur-events/card", rt.auth(), rt.alertCurEventsCard)
		pages.GET("/alert-cur-events/card/xh", rt.auth(), rt.alertCurEventsCardXH)
		pages.POST("/alert-cur-events/card/details", rt.auth(), rt.alertCurEventsCardDetails)
		pages.GET("/alert-his-events/list", rt.auth(), rt.alertHisEventsList)
		// pages.GET("/alert-his-events/list/xh", rt.auth(), rt.alertHisEventsListXH)
		pages.POST("/alert-his-events/batch-del", rt.auth(), rt.admin(), rt.alertHisEventBatchDel)
		pages.POST("/alert-cur-events/batch-del", rt.auth(), rt.admin(), rt.alertCurEventBatchDel)
		pages.POST("/alert-events/export-xls", rt.auth(), rt.admin(), rt.exportEventXH)
		// pages.POST("/alert-his-events/export-xls", rt.auth(), rt.admin(), rt.exportHisEventXH)
		pages.DELETE("/alert-cur-events", rt.auth(), rt.user(), rt.perm("/alert-cur-events/del"), rt.alertCurEventDel)

		pages.POST("/alert-his-event/solve/:eid", rt.auth(), rt.user(), rt.alertHisEventSolve) //人工解决异常接口
		pages.POST("/alert-his-event/close/:eid", rt.auth(), rt.user(), rt.alertHisEventClose) //人工关闭异常接口

		pages.GET("/alert-aggr-views", rt.auth(), rt.alertAggrViewGets)
		pages.DELETE("/alert-aggr-views", rt.auth(), rt.user(), rt.alertAggrViewDel)
		pages.POST("/alert-aggr-views", rt.auth(), rt.user(), rt.alertAggrViewAdd)
		pages.PUT("/alert-aggr-views", rt.auth(), rt.user(), rt.alertAggrViewPut)

		pages.GET("/busi-group/:id/task-tpls", rt.auth(), rt.user(), rt.perm("/job-tpls"), rt.bgro(), rt.taskTplGets)
		pages.POST("/busi-group/:id/task-tpls", rt.auth(), rt.user(), rt.perm("/job-tpls/add"), rt.bgrw(), rt.taskTplAdd)
		pages.DELETE("/busi-group/:id/task-tpl/:tid", rt.auth(), rt.user(), rt.perm("/job-tpls/del"), rt.bgrw(), rt.taskTplDel)
		pages.POST("/busi-group/:id/task-tpls/tags", rt.auth(), rt.user(), rt.perm("/job-tpls/put"), rt.bgrw(), rt.taskTplBindTags)
		pages.DELETE("/busi-group/:id/task-tpls/tags", rt.auth(), rt.user(), rt.perm("/job-tpls/put"), rt.bgrw(), rt.taskTplUnbindTags)
		pages.GET("/busi-group/:id/task-tpl/:tid", rt.auth(), rt.user(), rt.perm("/job-tpls"), rt.bgro(), rt.taskTplGet)
		pages.PUT("/busi-group/:id/task-tpl/:tid", rt.auth(), rt.user(), rt.perm("/job-tpls/put"), rt.bgrw(), rt.taskTplPut)

		pages.GET("/busi-group/:id/tasks", rt.auth(), rt.user(), rt.perm("/job-tasks"), rt.bgro(), rt.taskGets)
		pages.POST("/busi-group/:id/tasks", rt.auth(), rt.user(), rt.perm("/job-tasks/add"), rt.bgrw(), rt.taskAdd)
		pages.GET("/busi-group/:id/task/*url", rt.auth(), rt.user(), rt.perm("/job-tasks"), rt.taskProxy)
		pages.PUT("/busi-group/:id/task/*url", rt.auth(), rt.user(), rt.perm("/job-tasks/put"), rt.bgrw(), rt.taskProxy)

		pages.GET("/servers", rt.auth(), rt.admin(), rt.serversGet)
		pages.GET("/server-clusters", rt.auth(), rt.admin(), rt.serverClustersGet)

		pages.POST("/datasource/list", rt.auth(), rt.datasourceList)
		pages.POST("/datasource/plugin/list", rt.auth(), rt.pluginList)
		pages.POST("/datasource/upsert", rt.auth(), rt.admin(), rt.datasourceUpsert)
		pages.POST("/datasource/desc", rt.auth(), rt.admin(), rt.datasourceGet)
		pages.POST("/datasource/status/update", rt.auth(), rt.admin(), rt.datasourceUpdataStatus)
		pages.DELETE("/datasource/", rt.auth(), rt.admin(), rt.datasourceDel)

		pages.GET("/roles", rt.auth(), rt.admin(), rt.roleGets)
		pages.POST("/roles", rt.auth(), rt.admin(), rt.roleAdd)
		pages.PUT("/roles", rt.auth(), rt.admin(), rt.rolePut)
		pages.DELETE("/role/:id", rt.auth(), rt.admin(), rt.roleDel)

		pages.GET("/role/:id/ops", rt.auth(), rt.admin(), rt.operationOfRole)
		pages.PUT("/role/:id/ops", rt.auth(), rt.admin(), rt.roleBindOperation)
		pages.GET("/operation", rt.operations)

		pages.GET("/notify-tpls", rt.auth(), rt.admin(), rt.notifyTplGets)
		pages.PUT("/notify-tpl/content", rt.auth(), rt.admin(), rt.notifyTplUpdateContent)
		pages.PUT("/notify-tpl", rt.auth(), rt.admin(), rt.notifyTplUpdate)
		pages.POST("/notify-tpl", rt.auth(), rt.admin(), rt.notifyTplAdd)
		pages.DELETE("/notify-tpl/:id", rt.auth(), rt.admin(), rt.notifyTplDel)
		pages.POST("/notify-tpl/preview", rt.auth(), rt.admin(), rt.notifyTplPreview)

		pages.GET("/sso-configs", rt.auth(), rt.admin(), rt.ssoConfigGets)
		pages.PUT("/sso-config", rt.auth(), rt.admin(), rt.ssoConfigUpdate)

		pages.GET("/webhooks", rt.auth(), rt.admin(), rt.webhookGets)
		pages.PUT("/webhooks", rt.auth(), rt.admin(), rt.webhookPuts)

		pages.GET("/notify-script", rt.auth(), rt.admin(), rt.notifyScriptGet)
		pages.PUT("/notify-script", rt.auth(), rt.admin(), rt.notifyScriptPut)

		pages.GET("/notify-channel", rt.auth(), rt.admin(), rt.notifyChannelGets)
		pages.PUT("/notify-channel", rt.auth(), rt.admin(), rt.notifyChannelPuts)

		pages.GET("/notify-contact", rt.auth(), rt.admin(), rt.notifyContactGets)
		pages.PUT("/notify-contact", rt.auth(), rt.admin(), rt.notifyContactPuts)

		pages.GET("/notify-config", rt.auth(), rt.admin(), rt.notifyConfigGet)
		pages.PUT("/notify-config", rt.auth(), rt.admin(), rt.notifyConfigPut)

		// 资产管理
		pages.GET("/assets/:id", rt.auth(), rt.admin(), rt.assetsGet)
		pages.GET("/assets", rt.auth(), rt.admin(), rt.assetsGets)
		pages.POST("/assets", rt.auth(), rt.admin(), rt.assetsAdd)
		pages.PUT("/assets", rt.auth(), rt.admin(), rt.assetPut)
		pages.PUT("/assets/optmetrics", rt.auth(), rt.admin(), rt.putOptionalMetrics)
		pages.DELETE("/assets", rt.auth(), rt.user(), rt.perm("/assets/del"), rt.assetDel)
		pages.POST("/assets/config/default/:type", rt.auth(), rt.admin(), rt.assetDefaultConfigGet)
		pages.GET("/assets/idents", rt.auth(), rt.admin(), rt.assetIdentGetAll)
		pages.GET("/assets/types", rt.auth(), rt.admin(), rt.assetGetTypeList)
		pages.GET("/assets/tags", rt.auth(), rt.user(), rt.assetGetTags)
		pages.POST("/assets/tags", rt.auth(), rt.user(), rt.perm("/assets/put"), rt.assetBindTagsByFE)
		pages.DELETE("/assets/tags", rt.auth(), rt.user(), rt.perm("/assets/put"), rt.assetUnbindTagsByFE)
		pages.PUT("/assets/bgid", rt.auth(), rt.user(), rt.perm("/assets/put"), rt.assetUpdateBgid)
		pages.PUT("/assets/note", rt.auth(), rt.user(), rt.perm("/assets/note"), rt.assetUpdateNote)
		pages.PUT("/assets/orgnazation", rt.auth(), rt.user(), rt.assetUpdateOrganization) // 批量修改资产组织ID

		// 大屏接口管理
		pages.GET("/api-service", rt.auth(), rt.admin(), rt.apiServiceGets)
		pages.GET("/api-service/:id", rt.auth(), rt.admin(), rt.apiServiceGet)
		pages.GET("/api-service/:id/execute", rt.auth(), rt.admin(), rt.apiServiceExecute)
		pages.POST("/api-service", rt.auth(), rt.admin(), rt.apiServiceAdd)
		pages.PUT("/api-service", rt.auth(), rt.admin(), rt.apiServicePut)
		pages.DELETE("/api-service/:id", rt.auth(), rt.admin(), rt.apiServiceDel)
		pages.GET("/api-service/options", rt.auth(), rt.admin(), rt.apiServiceGetsOptions)

		// 大屏管理
		pages.GET("/bigscreen", rt.auth(), rt.admin(), rt.bigscreenGets)
		pages.GET("/bigscreen/:id", rt.auth(), rt.admin(), rt.bigscreenGet)
		pages.POST("/bigscreen", rt.auth(), rt.admin(), rt.bigscreenAdd)
		pages.PUT("/bigscreen", rt.auth(), rt.admin(), rt.bigscreenPut)
		pages.DELETE("/bigscreen/:id", rt.auth(), rt.admin(), rt.bigscreenDel)

		//西航
		pages.POST("/xh/assets", rt.auth(), rt.user(), rt.assetsAddXH)
		pages.PUT("/xh/assets", rt.auth(), rt.user(), rt.assetPutXH)
		pages.PUT("/xh/assets/batch-update", rt.auth(), rt.user(), rt.assetUpdateXH)
		pages.POST("/xh/assets/batch-del", rt.auth(), rt.user(), rt.assetDelXH)
		pages.GET("/xh/assets/filter", rt.auth(), rt.user(), rt.assetGetFilter)
		pages.GET("/xh/assets/xml", rt.auth(), rt.user(), rt.xmlceshi)
		pages.GET("/xh/assets/id", rt.auth(), rt.user(), rt.assetGetById)
		pages.POST("/xh/asset/import-xls", rt.auth(), rt.user(), rt.importAssetXH)
		pages.POST("/xh/asset/export-xls", rt.auth(), rt.user(), rt.exportAssetXH)
		pages.POST("/xh/asset/templet", rt.auth(), rt.user(), rt.templeAssetXH)

		pages.GET("/organization/:id", rt.auth(), rt.user(), rt.organizationGet)    // 依据id获取组织信息
		pages.GET("/organization", rt.auth(), rt.user(), rt.organizationGets)       // 获取组织树
		pages.PUT("/organization", rt.auth(), rt.user(), rt.organizationPut)        // 修改组织信息
		pages.DELETE("/organization/:id", rt.auth(), rt.user(), rt.organizationDel) // 删除组织信息
		pages.POST("/organization", rt.auth(), rt.user(), rt.organizationAdd)       // 上传组织信息
		pages.POST("/organization/name", rt.auth(), rt.user(), rt.organizationGetsByIds)

		pages.GET("/es-index-pattern", rt.auth(), rt.esIndexPatternGet)
		pages.GET("/es-index-pattern-list", rt.auth(), rt.esIndexPatternGetList)
		pages.POST("/es-index-pattern", rt.auth(), rt.admin(), rt.esIndexPatternAdd)
		pages.PUT("/es-index-pattern", rt.auth(), rt.admin(), rt.esIndexPatternPut)
		pages.DELETE("/es-index-pattern", rt.auth(), rt.admin(), rt.esIndexPatternDel)

		//设备类型管理
		pages.GET("/device-type", rt.auth(), rt.user(), rt.deviceTypeGets)
		pages.GET("/device-type/:id", rt.auth(), rt.user(), rt.deviceTypeGet)
		pages.POST("/device-type", rt.auth(), rt.user(), rt.deviceTypeAdd)
		pages.POST("/device-type/import-xls", rt.auth(), rt.user(), rt.importsDeviceType)
		pages.POST("/device-type/templet", rt.auth(), rt.user(), rt.templetDeviceType)
		pages.PUT("/device-type", rt.auth(), rt.user(), rt.deviceTypePut)
		pages.POST("/device-type/batch-del", rt.auth(), rt.user(), rt.deviceTypeBatchDel)

		//设备模型管理-
		pages.GET("/device-model/getmodel", rt.auth(), rt.user(), rt.deviceModelGets)
		pages.GET("/device-model/:id", rt.auth(), rt.user(), rt.deviceModelGet)
		pages.POST("/device-model", rt.auth(), rt.user(), rt.deviceModelAdd)
		pages.PUT("/device-model", rt.auth(), rt.user(), rt.deviceModelPut)
		pages.POST("/device-model/batch-del", rt.auth(), rt.user(), rt.deviceModelBatchDel)
		pages.POST("/device-model/batch", rt.auth(), rt.user(), rt.deviceModelBatchGets)
		pages.POST("/device-model/import-xls", rt.auth(), rt.user(), rt.importDeviceModels)
		pages.POST("/device-model/export-xls", rt.auth(), rt.user(), rt.exportDeviceModels)
		pages.POST("/device-model/picture", rt.auth(), rt.user(), rt.DeviceModelpictureAdd)
		pages.POST("/device-model/templet", rt.auth(), rt.user(), rt.templetDeviceModels)

		//设备厂商管理
		pages.GET("/device-producer/list", rt.auth(), rt.user(), rt.deviceProducerGets)
		pages.GET("/device-producer/:id", rt.auth(), rt.user(), rt.deviceProducerGet)
		pages.GET("/device-producer/getName", rt.auth(), rt.user(), rt.deviceProducerGetNames)
		pages.POST("/device-producer", rt.auth(), rt.user(), rt.deviceProducerAdd)
		pages.PUT("/device-producer", rt.auth(), rt.user(), rt.deviceProducerPut)
		pages.DELETE("/device-producer/:id", rt.auth(), rt.user(), rt.deviceProducerDel)
		pages.POST("/device-producer/batch-del", rt.auth(), rt.user(), rt.deviceProducerBatchDel)
		pages.POST("/device-producer/import-xls", rt.auth(), rt.user(), rt.importDeviceProducer)
		pages.POST("/device-producer/export-xls", rt.auth(), rt.user(), rt.downloadDeviceProducer)
		pages.POST("/device-producer/templet", rt.auth(), rt.user(), rt.templetDeviceProducer)

		//数据中心管理
		pages.GET("/datacenter/list", rt.auth(), rt.user(), rt.datacenterGets)
		pages.GET("/datacenter/:id", rt.auth(), rt.user(), rt.datacenterGet)
		pages.POST("/datacenter", rt.auth(), rt.user(), rt.datacenterAdd)
		pages.PUT("/datacenter", rt.auth(), rt.user(), rt.datacenterPut)
		pages.DELETE("/datacenter/:id", rt.auth(), rt.user(), rt.datacenterDel)

		//机房信息管理
		pages.GET("/computer-room/list", rt.auth(), rt.user(), rt.computerRoomGets)
		pages.GET("/computer-room/:id", rt.auth(), rt.user(), rt.computerRoomGet)
		pages.GET("/computer-room/datacenterId", rt.auth(), rt.user(), rt.computerRoomNameGet)
		pages.POST("/computer-room", rt.auth(), rt.user(), rt.computerRoomAdd)
		pages.POST("/computer-room/picture", rt.auth(), rt.user(), rt.computerRoompictureAdd)
		pages.PUT("/computer-room", rt.auth(), rt.user(), rt.computerRoomPut)
		pages.DELETE("/computer-room/:id", rt.auth(), rt.user(), rt.computerRoomDel)

		//机房分区管理
		pages.GET("/room-partition/room-id", rt.auth(), rt.user(), rt.roomPartitionGets)
		// pages.GET("/room-partition", rt.auth(), rt.admin(), rt.roomPartitionGets)
		pages.GET("/room-partition/:id", rt.auth(), rt.user(), rt.roomPartitionGet)
		pages.POST("/room-partition", rt.auth(), rt.user(), rt.roomPartitionAdd)
		pages.PUT("/room-partition", rt.auth(), rt.user(), rt.roomPartitionPut)
		pages.DELETE("/room-partition/:id", rt.auth(), rt.user(), rt.roomPartitionDel)

		//机柜信息管理
		pages.GET("/device-cabinet/list", rt.auth(), rt.user(), rt.deviceCabinetGets)
		pages.GET("/device-cabinet/:id", rt.auth(), rt.user(), rt.deviceCabinetGet)
		pages.GET("/device-cabinet/getNames", rt.auth(), rt.user(), rt.deviceCabinetNameGet)
		pages.POST("/device-cabinet", rt.auth(), rt.user(), rt.deviceCabinetAdd)
		pages.POST("/device-cabinet/import-xls", rt.auth(), rt.user(), rt.importDeviceCabinet)
		pages.POST("/device-cabinet/export-xls", rt.auth(), rt.user(), rt.downloadDeviceCabinet)
		pages.POST("/device-cabinet/templet", rt.auth(), rt.user(), rt.templetDeviceCabinet)
		pages.PUT("/device-cabinet", rt.auth(), rt.user(), rt.deviceCabinetPut)
		pages.DELETE("/device-cabinet/:id", rt.auth(), rt.user(), rt.deviceCabinetDel)

		//机柜组信息管理
		pages.GET("/cabinet-group/list", rt.auth(), rt.user(), rt.cabinetGroupGets)
		pages.GET("/cabinet-group/:id", rt.auth(), rt.user(), rt.cabinetGroupGet)
		pages.POST("/cabinet-group", rt.auth(), rt.user(), rt.cabinetGroupAdd)
		pages.PUT("/cabinet-group", rt.auth(), rt.user(), rt.cabinetGroupPut)
		pages.DELETE("/cabinet-group/:id", rt.auth(), rt.user(), rt.cabinetGroupDel)

		//配线架信息管理
		pages.GET("/distribution-frame/list", rt.auth(), rt.user(), rt.distributionFrameGets)
		pages.GET("/distribution-frame/:id", rt.auth(), rt.user(), rt.distributionFrameGet)
		pages.GET("/distribution-frame/query", rt.auth(), rt.user(), rt.distributionFrameQuery)
		pages.POST("/distribution-frame", rt.auth(), rt.user(), rt.distributionFrameAdd)
		pages.POST("/distribution-frame/import-xls", rt.auth(), rt.user(), rt.importDistributionFrame)
		pages.POST("/distribution-frame/export-xls", rt.auth(), rt.user(), rt.downloadDistributionFrame)
		pages.POST("/distribution-frame/templet", rt.auth(), rt.user(), rt.templetDistributionFrame)
		pages.PUT("/distribution-frame", rt.auth(), rt.user(), rt.distributionFramePut)
		pages.DELETE("/distribution-frame/:id", rt.auth(), rt.user(), rt.distributionFrameDel)

		//PDU信息管理
		pages.GET("/pdu/list", rt.auth(), rt.user(), rt.pduGets)
		pages.GET("/pdu/:id", rt.auth(), rt.user(), rt.pduGet)
		pages.POST("/pdu", rt.auth(), rt.user(), rt.pduAdd)
		pages.POST("/pdu/import-xls", rt.auth(), rt.user(), rt.importpdu)
		pages.POST("/pdu/export-xls", rt.auth(), rt.user(), rt.downloadpdu)
		pages.POST("/pdu/templet", rt.auth(), rt.user(), rt.templetpdu)
		pages.PUT("/pdu", rt.auth(), rt.user(), rt.pduPut)
		pages.DELETE("/pdu/:id", rt.auth(), rt.user(), rt.pduDel)

		//字典表类型管理
		pages.GET("/dict-type/list", rt.auth(), rt.user(), rt.dictTypeGets)
		pages.GET("/dict-type/:id", rt.auth(), rt.user(), rt.dictTypeGet)
		pages.POST("/dict-type", rt.auth(), rt.user(), rt.dictTypeAdd)
		pages.PUT("/dict-type", rt.auth(), rt.user(), rt.dictTypePut)
		pages.DELETE("/dict-type/:id", rt.auth(), rt.user(), rt.dictTypeDel)

		//字典表数据管理
		pages.GET("/dict-data/list", rt.auth(), rt.user(), rt.dictDataGets)
		pages.GET("/dict-data/:code", rt.auth(), rt.user(), rt.dictDataGet)
		pages.GET("/dict-data/exp", rt.auth(), rt.user(), rt.dictDataGetExp)
		pages.POST("/dict-data", rt.auth(), rt.user(), rt.dictDataAdd)
		pages.POST("/dict-data/one", rt.auth(), rt.user(), rt.dictDataOneAdd)
		pages.POST("/dict-data/asset-batch", rt.auth(), rt.user(), rt.dictDataBtachDel)
		pages.PUT("/dict-data", rt.auth(), rt.user(), rt.dictDataPut)
		pages.PUT("/dict-data/one", rt.auth(), rt.user(), rt.dictDataOnePut)
		pages.DELETE("/dict-data/:id", rt.auth(), rt.user(), rt.dictDataDel)

		//资产树管理
		pages.GET("/asset-tree/list", rt.auth(), rt.user(), rt.assetTreeGets)
		pages.GET("/asset-tree/data", rt.auth(), rt.user(), rt.assetTreeGetALL)
		pages.GET("/asset-tree/part", rt.auth(), rt.user(), rt.assetTreeGetPart)
		pages.GET("/asset-tree/count", rt.auth(), rt.user(), rt.assetTreeGetCount)
		pages.GET("/asset-tree/:id", rt.auth(), rt.user(), rt.assetTreeGet)
		pages.POST("/asset-tree", rt.auth(), rt.user(), rt.assetTreeAdd)
		pages.POST("/asset-tree/parent/:id", rt.auth(), rt.user(), rt.assetTreeUpdate)
		pages.POST("/asset-tree/transfer", rt.auth(), rt.user(), rt.assetTreeTransfer)
		pages.POST("/asset-tree/asset", rt.auth(), rt.user(), rt.assetTreeGetsByAssetId)
		pages.PUT("/asset-tree", rt.auth(), rt.user(), rt.assetTreePut)
		pages.DELETE("/asset-tree/:id", rt.auth(), rt.user(), rt.assetTreeDel)

		//资产详情管理
		// pages.POST("/asset-basic/list", rt.auth(), rt.user(), rt.assetBasicGetsByTree)
		pages.GET("/asset-basic/list", rt.auth(), rt.user(), rt.assetBasicGets)
		pages.GET("/asset-basic/copy", rt.auth(), rt.user(), rt.assetCopyGetsByAssetId)
		pages.GET("/asset-basic/:id", rt.auth(), rt.user(), rt.assetBasicGet)
		// pages.GET("/asset-basic/statistics", rt.auth(), rt.admin(), rt.assetStatisticsAll)
		pages.POST("/asset-basic", rt.auth(), rt.user(), rt.assetBasicAdd)
		pages.POST("/asset-basic/copy", rt.auth(), rt.user(), rt.assetCopyAdd)
		pages.POST("/asset-basic/table", rt.auth(), rt.user(), rt.assetFieldGetsByTableName)
		pages.POST("/asset-basic/del", rt.auth(), rt.user(), rt.assetBatchDel)
		pages.POST("/asset-basic/status", rt.auth(), rt.user(), rt.assetBasicStatusUpdateByAssetId)
		pages.POST("/asset-basic/templet", rt.auth(), rt.user(), rt.templeAsset)
		pages.POST("/asset-basic/import-xls", rt.auth(), rt.user(), rt.importAsset)
		pages.POST("/asset-basic/export-xls", rt.auth(), rt.user(), rt.exportAsset)
		pages.POST("/asset-basic/batch-update", rt.auth(), rt.user(), rt.assetBasicsUpdate)
		pages.POST("/asset-basic/list/query", rt.auth(), rt.user(), rt.assetBasicBatchGets)
		pages.PUT("/asset-basic", rt.auth(), rt.user(), rt.assetBasicPut)
		pages.DELETE("/asset-basic/:id", rt.auth(), rt.user(), rt.assetBasicDel)

		//资产扩展管理
		pages.GET("/asset-expansion/list", rt.auth(), rt.user(), rt.assetExpansionGets)
		pages.GET("/asset-expansion/:id", rt.auth(), rt.user(), rt.assetExpansionGet)
		pages.GET("/asset-expansion/asset", rt.auth(), rt.user(), rt.assetExpansionGetByAssetId)
		pages.POST("/asset-expansion", rt.auth(), rt.user(), rt.assetExpansionAdd)
		pages.POST("/asset-expansion/batch", rt.auth(), rt.user(), rt.assetExpansionBatchAdd)
		pages.POST("/asset-expansion/map", rt.auth(), rt.user(), rt.assetExpansionGetByMap)
		pages.POST("/asset-expansion/netconfig/templet", rt.auth(), rt.user(), rt.templeAssetNetWork)
		pages.POST("/asset-expansion/netconfig/import-xls", rt.auth(), rt.user(), rt.importAssetNetWork)
		pages.POST("/asset-expansion/netconfig/export-xls", rt.auth(), rt.user(), rt.exportAssetNetWork)
		pages.PUT("/asset-expansion", rt.auth(), rt.user(), rt.assetExpansionPut)
		pages.DELETE("/asset-expansion/:id", rt.auth(), rt.user(), rt.assetExpansionDel)

		//资产维保管理
		pages.GET("/asset-maintenance", rt.auth(), rt.user(), rt.assetMaintenanceGets)
		pages.GET("/asset-maintenance/:id", rt.auth(), rt.user(), rt.assetMaintenanceGet)
		pages.GET("/asset-maintenance/asset", rt.auth(), rt.user(), rt.assetMaintenanceGetAssetId)
		pages.POST("/asset-maintenance", rt.auth(), rt.user(), rt.assetMaintenanceAdd)
		pages.PUT("/asset-maintenance", rt.auth(), rt.user(), rt.assetMaintenancePut)
		pages.DELETE("/asset-maintenance/:id", rt.auth(), rt.user(), rt.assetMaintenanceDel)

		//资产管理
		pages.GET("/asset-management", rt.auth(), rt.user(), rt.assetManagementGets)
		pages.GET("/asset-management/:id", rt.auth(), rt.user(), rt.assetManagementGet)
		pages.GET("/asset-management/asset", rt.auth(), rt.user(), rt.assetManagementGetAI)
		pages.POST("/asset-management", rt.auth(), rt.user(), rt.assetManagementAdd)
		pages.PUT("/asset-management", rt.auth(), rt.user(), rt.assetManagementPut)
		pages.DELETE("/asset-management/:id", rt.auth(), rt.user(), rt.assetManagementDel)

		//资产变更管理
		pages.GET("/asset-alter", rt.auth(), rt.user(), rt.assetAlterGets)
		pages.GET("/asset-alter/:id", rt.auth(), rt.user(), rt.assetAlterGet)
		pages.GET("/asset-alter/asset", rt.auth(), rt.user(), rt.assetAlterGetByAssetId)
		pages.POST("/asset-alter", rt.auth(), rt.user(), rt.assetAlterAdd)
		pages.POST("/asset-alter/export-xls", rt.auth(), rt.user(), rt.downloadAssetAlter)
		pages.PUT("/asset-alter", rt.auth(), rt.user(), rt.assetAlterPut)
		pages.DELETE("/asset-alter/:id", rt.auth(), rt.user(), rt.assetAlterDel)

		//探针版本管理
		pages.POST("/target/version", rt.auth(), rt.user(), rt.importNewVersion)

		pages.GET("/target/:ident/version", rt.auth(), rt.user(), rt.UsableVersionGet)
		pages.PUT("/target/:ident/version", rt.auth(), rt.user(), rt.targetVersionPut)
		//更新包管理
		pages.POST("/server/update", rt.auth(), rt.user(), rt.importUpgradePack)

		//设备类型表单配置表管理
		pages.GET("/device-type_config", rt.auth(), rt.user(), rt.deviceTypeConfigGets)
		pages.GET("/device-type_config/:id", rt.auth(), rt.user(), rt.deviceTypeConfigGet)
		pages.POST("/device-type_config", rt.auth(), rt.user(), rt.deviceTypeConfigAdd)
		pages.PUT("/device-type_config", rt.auth(), rt.user(), rt.deviceTypeConfigPut)
		pages.DELETE("/device-type_config/:id", rt.auth(), rt.user(), rt.deviceTypeConfigDel)

		//备件基础数据管理
		pages.GET("/spare-part_basic", rt.auth(), rt.user(), rt.sparePartBasicGets)
		pages.GET("/spare-part_basic/:id", rt.auth(), rt.user(), rt.sparePartBasicGet)
		pages.POST("/spare-part_basic", rt.auth(), rt.user(), rt.sparePartBasicAdd)
		pages.POST("/spare-part_basic/picture", rt.auth(), rt.user(), rt.sparePartBasicPictureAdd)
		pages.POST("/spare-part_basic/import-xls", rt.auth(), rt.user(), rt.importsSparePartBasic)
		pages.POST("/spare-part_basic/export-xls", rt.auth(), rt.user(), rt.exportSparePartBasic)
		pages.POST("/spare-part_basic/templet", rt.auth(), rt.user(), rt.templetSparePartBasic)
		pages.PUT("/spare-part_basic", rt.auth(), rt.user(), rt.sparePartBasicPut)
		pages.POST("/spare-part_basic/batch-del", rt.auth(), rt.user(), rt.sparePartBasicBatchDel)

		//部件类型管理
		pages.GET("/component-type", rt.auth(), rt.user(), rt.componentTypeGets)
		pages.GET("/component-type/:id", rt.auth(), rt.user(), rt.componentTypeGet)
		pages.POST("/component-type", rt.auth(), rt.user(), rt.componentTypeAdd)
		pages.POST("/component-type/picture", rt.auth(), rt.user(), rt.componentTypePictureAdd)
		pages.PUT("/component-type", rt.auth(), rt.user(), rt.componentTypePut)
		pages.POST("/component-type/batch-del", rt.auth(), rt.user(), rt.componentTypeBatchDel)

		//库房信息管理
		pages.GET("/storeroom-management", rt.auth(), rt.user(), rt.storeroomManagementGets)
		pages.GET("/storeroom-management/:id", rt.auth(), rt.user(), rt.storeroomManagementGet)
		pages.POST("/storeroom-management", rt.auth(), rt.user(), rt.storeroomManagementAdd)
		pages.POST("/storeroom-management/import-xls", rt.auth(), rt.user(), rt.importsStoreroomManagement)
		pages.POST("/storeroom-management/templet", rt.auth(), rt.user(), rt.templetStoreroomManagement)
		pages.PUT("/storeroom-management", rt.auth(), rt.user(), rt.storeroomManagementPut)
		pages.POST("/storeroom-management/batch-del", rt.auth(), rt.user(), rt.storeroomManagementBatchDel)

		//设备上下线
		pages.GET("/device-online", rt.auth(), rt.user(), rt.deviceOnlineGets)
		pages.GET("/device-online/:id", rt.auth(), rt.user(), rt.deviceOnlineGet)
		pages.POST("/device-online", rt.auth(), rt.user(), rt.deviceOnlineAdd)
		pages.PUT("/device-online", rt.auth(), rt.user(), rt.deviceOnlinePut)
		pages.DELETE("/device-online/:id", rt.auth(), rt.user(), rt.deviceOnlineDel)

		//设备报废管理
		pages.GET("/device-scrap", rt.auth(), rt.user(), rt.deviceScrapGets)
		pages.GET("/device-scrap/:id", rt.auth(), rt.user(), rt.deviceScrapGet)
		pages.POST("/device-scrap", rt.auth(), rt.user(), rt.deviceScrapAdd)
		pages.POST("/device-scrap/export-xls", rt.auth(), rt.user(), rt.exportDeviceScraps)
		// pages.PUT("/device-scrap", rt.auth(), rt.admin(), rt.deviceScrapPut)
		pages.POST("/device-scrap/batch", rt.auth(), rt.user(), rt.deviceScrapBatchDel)

		// pages.GET("/dashboard/assets/statistics", rt.getDashboardAssetStatistics) //资产统计接口
		pages.GET("/dashboard/count/ceshi", rt.getDashboardDataCount) //首页数据统计接口
		pages.POST("/dashboard/asset/details/ceshi", rt.AssetDetails) //资产详情接口
		// pages.GET("/dashboard/alarm/details", rt.AlarmDetails)                    //告警详情接口
		pages.GET("/dashboard/his-alarms/his-query/ceshi", rt.AlarmHisQueryGet)                      //历史告警搜索记录查询
		pages.DELETE("/dashboard/his-alarms/his-query/ceshi", rt.AlarmHisDel)                        //历史告警搜索记录删除
		pages.GET("/dashboard/his-alarm/filter/ceshi", rt.AlarmHisFilter)                            //历史告警过滤条件查询
		pages.GET("/dashboard/his-alarms/ceshi", rt.AlarmHisGet)                                     //历史告警查询
		pages.GET("/dashboard/data/ceshi", rt.DataBoardPerson)                                       //列表数据测试
		pages.GET("/dashboard/user/list/ceshi", rt.DataBoardAssetList)                               //数据看板筛选资产
		pages.POST("/dashboard/user/add", rt.auth(), rt.user(), rt.DataBoardAsset)                   //创建/删除用户看板
		pages.GET("/dashboard/user/data", rt.auth(), rt.user(), rt.DataBoardAssetsGet)               //数据面板
		pages.GET("/dashboard/user/page-name/get", rt.auth(), rt.user(), rt.DataBoardPageNameGets)   //获取pageName列表
		pages.DELETE("/dashboard/user/page-name/del", rt.auth(), rt.user(), rt.DataBoardPageNameDel) //获取pageName列表
		pages.GET("/dashboard/alert-cur-events/ceshi", rt.getAlertListByFE)                          // 告警列表接口前端接口返回

		//资产目录
		pages.GET("/asset-directory/tree", rt.auth(), rt.user(), rt.assetsDirTreeGet)
		pages.GET("/asset-directory/move", rt.auth(), rt.user(), rt.assetsDirectoryMove)
		pages.POST("/asset-directory", rt.auth(), rt.user(), rt.assetsDirectoryAdd)
		pages.PUT("/asset-directory", rt.auth(), rt.user(), rt.assetsDirectoryPut)
		pages.DELETE("/asset-directory/:id", rt.auth(), rt.user(), rt.assetsDirectoryDel)

		//资产扩展（西航）
		pages.PUT("/xh/assets-expansion", rt.auth(), rt.user(), rt.assetsExpansionPut)

		//监控（西航）
		pages.GET("/xh/monitoring", rt.auth(), rt.user(), rt.monitoringAllGets)
		pages.GET("/xh/monitoring/filter", rt.auth(), rt.user(), rt.monitoringGets)
		pages.GET("/xh/monitoring/:id", rt.auth(), rt.user(), rt.monitoringGet)
		pages.GET("/xh/monitoring/asset", rt.auth(), rt.user(), rt.monitoringGetByAssetId)
		pages.POST("/xh/monitoring", rt.auth(), rt.user(), rt.monitoringAdd)
		pages.PUT("/xh/monitoring", rt.auth(), rt.user(), rt.monitoringPut)
		pages.DELETE("/xh/monitoring/:id", rt.auth(), rt.user(), rt.monitoringDel)
		pages.POST("/xh/monitoring/status", rt.auth(), rt.user(), rt.monitoringStatusUpdate)
		pages.POST("/xh/monitoring/data", rt.auth(), rt.user(), rt.monitoringData)
		pages.GET("/monitoring/options", rt.auth(), rt.user(), rt.monitoringGetOptions)

		//用户配置
		pages.GET("/user-config", rt.auth(), rt.admin(), rt.userConfigGets)
		//pages.GET("/user-config/:id", rt.auth(), rt.admin(), rt.userConfigGet)
		pages.GET("/user-config/getInfo", rt.userLogoGet)
		pages.POST("/user-config/picture", rt.auth(), rt.admin(), rt.userPictureAdd)
		pages.PUT("/user-config/logo", rt.auth(), rt.admin(), rt.logoPut)
		pages.PUT("/user-config", rt.auth(), rt.admin(), rt.userConfigPut)
		pages.DELETE("/user-config/:id", rt.auth(), rt.admin(), rt.userConfigDel)

		//操作日志
		pages.GET("/operation-log", rt.auth(), rt.admin(), rt.operationLogGets)
		pages.GET("/operation-log/:id", rt.auth(), rt.admin(), rt.operationLogGet)
		pages.POST("/operation-log", rt.auth(), rt.admin(), rt.operationLogAdd)
		pages.PUT("/operation-log", rt.auth(), rt.admin(), rt.operationLogPut)
		pages.DELETE("/operation-log/:id", rt.auth(), rt.admin(), rt.operationLogDel)

		//静态图片
		pages.Static("/images", "./etc/picture/logo")

		//接口访问设置
		pages.GET("/xh/assets/out", rt.auth(), rt.thirdUser(), rt.assetGetAll)

		//系统日志
		pages.GET("/xh/sys-log/filter", rt.auth(), rt.user(), rt.systemLogGets)
		pages.POST("/xh/sys-log/export-xls", rt.auth(), rt.admin(), rt.exportSystemLogXH)

		//许可管理
		pages.POST("/xh/license/add-file", rt.auth(), rt.admin(), rt.licenseAddXH)
		pages.GET("/xh/license/list", rt.auth(), rt.admin(), rt.licenseGetsXH)
		pages.PUT("/xh/license/update", rt.auth(), rt.admin(), rt.licenseUpdateXH)
		pages.POST("/xh/license-config", rt.auth(), rt.admin(), rt.licenseConfigAdd)
		pages.PUT("/xh/license-config", rt.auth(), rt.admin(), rt.licenseConfigPut)
	}

	r.GET("/api/n9e/versions", func(c *gin.Context) {
		v := version.Version
		lastIndex := strings.LastIndex(version.Version, "-")
		if lastIndex != -1 {
			v = version.Version[:lastIndex]
		}

		ginx.NewRender(c).Data(gin.H{"version": v, "github_verison": version.GithubVersion.Load().(string)}, nil)
	})

	if rt.HTTP.APIForService.Enable {
		service := r.Group("/v1/n9e")
		if len(rt.HTTP.APIForService.BasicAuth) > 0 {
			service.Use(gin.BasicAuth(rt.HTTP.APIForService.BasicAuth))
		}
		{
			service.Any("/prometheus/*url", rt.dsProxy)
			service.POST("/users", rt.userAddPost)
			service.GET("/users", rt.userFindAll)

			service.GET("/user-groups", rt.userGroupGetsByService)
			service.GET("/user-group-members", rt.userGroupMemberGetsByService)

			service.GET("/targets", rt.targetGetsByService)
			service.GET("/targets/tags", rt.targetGetTags)
			service.POST("/targets/tags", rt.targetBindTagsByService)
			service.DELETE("/targets/tags", rt.targetUnbindTagsByService)
			service.PUT("/targets/note", rt.targetUpdateNoteByService)

			service.POST("/alert-rules", rt.alertRuleAddByService)
			service.DELETE("/alert-rules", rt.alertRuleDelByService)
			service.PUT("/alert-rule/:arid", rt.alertRulePutByService)
			service.GET("/alert-rule/:arid", rt.alertRuleGet)
			service.GET("/alert-rules", rt.alertRulesGetByService)

			service.GET("/alert-subscribes", rt.alertSubscribeGetsByService)

			service.GET("/busi-groups", rt.busiGroupGetsByService)

			service.GET("/datasources", rt.datasourceGetsByService)
			service.GET("/datasource-ids", rt.getDatasourceIds)
			service.POST("/server-heartbeat", rt.serverHeartbeat)
			service.GET("/servers-active", rt.serversActive)

			service.GET("/recording-rules", rt.recordingRuleGetsByService)

			service.GET("/alert-mutes", rt.alertMuteGets)
			service.POST("/alert-mutes", rt.alertMuteAddByService)
			service.DELETE("/alert-mutes", rt.alertMuteDel)

			service.GET("/alert-cur-events", rt.alertCurEventsList)
			service.GET("/alert-cur-events-get-by-rid", rt.alertCurEventsGetByRid)
			service.GET("/alert-his-events", rt.alertHisEventsList)
			service.GET("/alert-his-event/:eid", rt.alertHisEventGet)

			service.GET("/task-tpl/:tid", rt.taskTplGetByService)

			service.GET("/config/:id", rt.configGet)
			service.GET("/configs", rt.configsGet)
			service.GET("/config", rt.configGetByKey)
			service.PUT("/configs", rt.configsPut)
			service.POST("/configs", rt.configsPost)
			service.DELETE("/configs", rt.configsDel)

			service.POST("/conf-prop/encrypt", rt.confPropEncrypt)
			service.POST("/conf-prop/decrypt", rt.confPropDecrypt)

			service.GET("/statistic", rt.statistic)

			service.GET("/notify-tpls", rt.notifyTplGets)

			service.POST("/task-record-add", rt.taskRecordAdd)

			//前端大屏接口
			service.GET("/dashboard/assets", rt.getDashboardAssetsByFE)                 // 资产清单
			service.GET("/dashboard/organization-tree", rt.getOrganizationTreeByFE)     // 提供前端组织树接口
			service.GET("/dashboard/alert-cur-events", rt.getAlertListByFE)             // 告警列表接口前端接口返回
			service.GET("/dashboard/assets/statistics", rt.getDashboardAssetStatistics) //资产统计接口
			service.GET("/dashboard/count", rt.getDashboardDataCount)                   //首页数据统计接口
			service.POST("/dashboard/asset/details", rt.AssetDetails)                   //资产详情接口
			service.GET("/dashboard/alarm/details", rt.AlarmDetails)                    //告警详情接口
			service.GET("/dashboard/his-alarms/his-query", rt.AlarmHisQueryGet)         //历史告警搜索记录查询
			service.DELETE("/dashboard/his-alarms/his-query", rt.AlarmHisDel)           //历史告警搜索删除
			service.GET("/dashboard/his-alarm/filter", rt.AlarmHisFilter)               //历史告警过滤条件查询
			service.GET("/dashboard/his-alarms", rt.AlarmHisGet)                        //历史告警搜索记录
			service.GET("/dashboard/user/list", rt.DataBoardAssetList)                  //数据看板筛选资产
			service.POST("/dashboard/user/add", rt.DataBoardAsset)                      //创建/编辑/删除用户看板
			service.GET("/dashboard/user/data", rt.DataBoardAssetsGet)                  //数据面板

		}
	}

	if rt.HTTP.APIForAgent.Enable {
		heartbeat := r.Group("/v1/n9e")
		{
			if len(rt.HTTP.APIForAgent.BasicAuth) > 0 {
				heartbeat.Use(gin.BasicAuth(rt.HTTP.APIForAgent.BasicAuth))
			}
			heartbeat.POST("/heartbeat", rt.heartbeat)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	rt.configNoRoute(r, &statikFS)
}

func Render(c *gin.Context, data, msg interface{}) {
	if msg == nil {
		if data == nil {
			data = struct{}{}
		}
		c.JSON(http.StatusOK, gin.H{"data": data, "error": ""})
	} else {
		c.JSON(http.StatusOK, gin.H{"error": gin.H{"message": msg}})
	}
}

func Dangerous(c *gin.Context, v interface{}, code ...int) {
	if v == nil {
		return
	}

	switch t := v.(type) {
	case string:
		if t != "" {
			c.JSON(http.StatusOK, gin.H{"error": v})
		}
	case error:
		c.JSON(http.StatusOK, gin.H{"error": t.Error()})
	}
}
