package center

import (
	"context"
	"fmt"

	"github.com/ccfos/nightingale/v6/alert"
	"github.com/ccfos/nightingale/v6/alert/astats"
	"github.com/ccfos/nightingale/v6/alert/process"
	"github.com/ccfos/nightingale/v6/center/cconf"
	"github.com/ccfos/nightingale/v6/center/cstats"
	"github.com/ccfos/nightingale/v6/center/metas"
	"github.com/ccfos/nightingale/v6/center/sso"
	"github.com/ccfos/nightingale/v6/conf"
	"github.com/ccfos/nightingale/v6/dumper"
	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/models/migrate"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/httpx"
	"github.com/ccfos/nightingale/v6/pkg/i18nx"
	"github.com/ccfos/nightingale/v6/pkg/logx"
	"github.com/ccfos/nightingale/v6/pkg/version"
	"github.com/ccfos/nightingale/v6/prom"
	"github.com/ccfos/nightingale/v6/pushgw/idents"
	"github.com/ccfos/nightingale/v6/pushgw/writer"
	"github.com/ccfos/nightingale/v6/storage"
	"github.com/toolkits/pkg/logger"

	alertrt "github.com/ccfos/nightingale/v6/alert/router"
	"github.com/ccfos/nightingale/v6/center/attrs"
	centerrt "github.com/ccfos/nightingale/v6/center/router"
	providerrt "github.com/ccfos/nightingale/v6/provider/router"
	pushgwrt "github.com/ccfos/nightingale/v6/pushgw/router"
)

func Initialize(configDir string, cryptoKey string) (func(), error) {
	config, err := conf.InitConfig(configDir, cryptoKey)
	if err != nil {
		return nil, fmt.Errorf("failed to init config: %v", err)
	}

	cconf.LoadMetricsYaml(configDir, config.Center.MetricsYamlFile)
	cconf.LoadOpsYaml(configDir, config.Center.OpsYamlFile)

	logxClean, err := logx.Init(config.Log)
	if err != nil {
		return nil, err
	}

	i18nx.Init(configDir)
	cstats.Init()

	db, err := storage.New(config.DB)
	if err != nil {
		return nil, err
	}
	ctx := ctx.NewContext(context.Background(), db, true)
	models.InitRoot(ctx)
	migrate.Migrate(db)

	logLever := map[int]string{
		1: "DEBUG",
		2: "INFO",
		3: "WARNING",
		4: "ERROR",
	}

	logLeverT := map[string]int{
		"DEBUG":   1,
		"INFO":    2,
		"WARNING": 3,
		"ERROR":   4,
	}

	var lst *models.UserConfig
	err = ctx.DB.Where("id", 1).Find(&lst).Error
	if lst != nil {
		if lst.LogLever != logLeverT[config.Log.Level] {
			// 热修改logger日志等级
			logger.SetSeverity(logLever[lst.LogLever])
		}
		if lst.HttpHost != config.HTTP.Host {
			config.HTTP.Host = lst.HttpHost
		}
		if 8000 <= lst.HttpPort && lst.HttpPort <= 65535 {
			if int(lst.HttpPort) != config.HTTP.Port {
				config.HTTP.Port = int(lst.HttpPort)
			}
		}
		if (lst.Captcha != 2) != config.HTTP.ShowCaptcha.Enable {
			config.HTTP.ShowCaptcha = httpx.ShowCaptcha{Enable: lst.Captcha != 2}
		}
		if (lst.ApiService != 2) != config.HTTP.APIForService.Enable {
			config.HTTP.APIForService = httpx.BasicAuths{Enable: lst.ApiService != 2}
		}
		if lst.AccessExpired != config.HTTP.JWTAuth.AccessExpired {
			config.HTTP.JWTAuth.AccessExpired = lst.AccessExpired
		}
		if lst.RefreshExpired != config.HTTP.JWTAuth.RefreshExpired {
			config.HTTP.JWTAuth.RefreshExpired = lst.RefreshExpired
		}
		if (lst.OpenRsa != 2) != config.HTTP.RSA.OpenRSA {
			config.HTTP.RSA.OpenRSA = lst.OpenRsa != 2
		}
	}

	if err != nil {
		return nil, err
	}

	redis, err := storage.NewRedis(config.Redis)
	if err != nil {
		return nil, err
	}

	metas := metas.New(ctx, redis)
	idents := idents.New(ctx)

	syncStats := memsto.NewSyncStats()
	alertStats := astats.NewSyncStats()

	sso := sso.Init(config.Center, ctx)

	busiGroupCache := memsto.NewBusiGroupCache(ctx, syncStats)
	targetCache := memsto.NewTargetCache(ctx, syncStats, redis)
	dsCache := memsto.NewDatasourceCache(ctx, syncStats)
	alertMuteCache := memsto.NewAlertMuteCache(ctx, syncStats)
	alertRuleCache := memsto.NewAlertRuleCache(ctx, syncStats)
	notifyConfigCache := memsto.NewNotifyConfigCache(ctx)
	assetCache := memsto.NewAssetCache(ctx, syncStats)
	userCache := memsto.NewUserCache(ctx, syncStats)
	userGroupCache := memsto.NewUserGroupCache(ctx, syncStats)
	// licenseCache := memsto.NewLicenseCache(ctx, syncStats)

	promClients := prom.NewPromClient(ctx, config.Alert.Heartbeat)

	externalProcessors := process.NewExternalProcessors()
	alert.Start(config.Alert, config.Pushgw, syncStats, alertStats, externalProcessors, targetCache, busiGroupCache, alertMuteCache, alertRuleCache, notifyConfigCache, dsCache, ctx, promClients, assetCache, userCache, userGroupCache)

	writers := writer.NewWriters(config.Pushgw)

	httpx.InitRSAConfig(&config.HTTP.RSA)
	go version.GetGithubVersion()

	alertrtRouter := alertrt.New(config.HTTP, config.Alert, alertMuteCache, targetCache, busiGroupCache, alertStats, ctx, externalProcessors)
	centerRouter := centerrt.New(config.HTTP, config.Center, cconf.Operations, dsCache, notifyConfigCache, promClients, redis, sso, ctx, metas, idents, targetCache, userCache, userGroupCache, assetCache)
	pushgwRouter := pushgwrt.New(config.HTTP, config.Pushgw, targetCache, busiGroupCache, idents, writers, ctx)
	providerRouter := providerrt.New(config.HTTP, targetCache, busiGroupCache, assetCache, ctx)
	attrs.StartAttrSync(ctx, promClients, assetCache)

	r := httpx.GinEngine(config.Global.RunMode, config.HTTP)

	centerRouter.Config(r)
	alertrtRouter.Config(r)
	pushgwRouter.Config(r)
	providerRouter.Config(r)
	dumper.ConfigRouter(r)

	httpClean := httpx.Init(config.HTTP, r)

	return func() {
		logxClean()
		httpClean()
	}, nil
}
