package memsto

import (
	"fmt"
	"sync"
	"time"

	"github.com/ccfos/nightingale/v6/center/ws"
	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/version"
	"github.com/pkg/errors"
	"github.com/toolkits/pkg/logger"
)

type LicenseCache struct {
	statTotal       int64
	statLastUpdated int64
	ctx             *ctx.Context
	stats           *Stats

	sync.RWMutex
	config   map[int64]*models.LicenseConfig
	license  map[int64]*models.License
	lastTime time.Time
}

func NewLicenseCache(ctx *ctx.Context, stats *Stats) *LicenseCache {
	cache := &LicenseCache{
		statTotal:       -1,
		statLastUpdated: -1,
		ctx:             ctx,
		stats:           stats,
		config:          make(map[int64]*models.LicenseConfig),
		license:         make(map[int64]*models.License),
		lastTime:        time.Now(),
	}
	cache.SyncLicenseConfig()
	cache.SyncLicense()
	return cache
}

func (cache *LicenseCache) GetByLicenseConfig() *models.LicenseConfig {
	cache.RLock()
	defer cache.RUnlock()
	for key := range cache.config {
		return cache.config[key]
	}
	return nil
}

func (cache *LicenseCache) GetByLicense() *models.License {
	cache.RLock()
	defer cache.RUnlock()
	max := int64(-1)
	now := time.Now().Unix()
	for key, val := range cache.license {
		if val.StartTime < now && val.EndTime > now && key > max {
			max = key
		}
	}
	if len(cache.license) == 0 {
		return nil
	}
	return cache.license[max]
}

func (cache *LicenseCache) GetLicenseAll() []*models.License {
	cache.RLock()
	defer cache.RUnlock()
	licenses := make([]*models.License, 0)
	for _, val := range cache.license {
		licenses = append(licenses, val)
	}
	return licenses
}

func (cache *LicenseCache) SyncLicenseConfig() {
	err := cache.syncLicenseConfig()
	if err != nil {
		fmt.Println("failed to sync licenseConfig:", err)
		exit(1)
	}

	go cache.loopSyncLicenseConfig()
}

func (cache *LicenseCache) SyncLicense() {
	err := cache.syncLicense()
	if err != nil {
		fmt.Println("failed to sync license:", err)
		exit(1)
	}

	go cache.loopSyncLicense()
}

func (cache *LicenseCache) loopSyncLicenseConfig() {
	duration := time.Duration(9000) * time.Millisecond
	for {
		time.Sleep(duration)
		if err := cache.syncLicenseConfig(); err != nil {
			logger.Warning("failed to sync licenseConfig:", err)
		}
	}
}

func (cache *LicenseCache) loopSyncLicense() {
	duration := time.Duration(9000) * time.Millisecond
	for {
		time.Sleep(duration)
		if err := cache.syncLicense(); err != nil {
			logger.Warning("failed to sync license:", err)
		}
	}
}

func (cache *LicenseCache) syncLicenseConfig() error {
	start := time.Now()

	stat, err := models.LicenseConfigStatistics(cache.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec LicenseConfigStatistics")
	}

	if !cache.StatChanged(stat.Total, stat.LastUpdated) {
		cache.stats.GaugeCronDuration.WithLabelValues("sync_licenseConfig").Set(0)
		cache.stats.GaugeSyncNumber.WithLabelValues("sync_licenseConfig").Set(0)
		logger.Debug("licenseConfig not changed")
		return nil
	}

	lst, err := models.LicenseConfigGetsAll(cache.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec licenseConfigGetsAll")
	}

	m := make(map[int64]*models.LicenseConfig)
	for i := 0; i < len(lst); i++ {
		m[lst[i].Id] = lst[i]
	}

	cache.SetLicenseConfig(m, stat.Total, stat.LastUpdated)

	ms := time.Since(start).Milliseconds()
	cache.stats.GaugeCronDuration.WithLabelValues("sync_licenseConfig").Set(float64(ms))
	cache.stats.GaugeSyncNumber.WithLabelValues("sync_licenseConfig").Set(float64(len(m)))

	logger.Infof("timer: sync licenseConfig done, cost: %dms, number: %d", ms, len(m))

	return nil
}

func (cache *LicenseCache) syncLicense() error {
	start := time.Now()

	stat, err := models.LicenseStatistics()
	if err != nil {
		logger.Errorf("failed to exec LicenseStatistics: %s", err.Error())
		return nil
	}

	if !cache.StatChanged(stat.Total, stat.LastUpdated) {
		cache.stats.GaugeCronDuration.WithLabelValues("sync_license").Set(0)
		cache.stats.GaugeSyncNumber.WithLabelValues("sync_license").Set(0)
		logger.Debug("License not changed")
		return nil
	}

	//获取探针数量
	usedNode, err := models.TargetTotalCount(cache.ctx)
	if err != nil {
		logger.Errorf("failed to exec TargetTotalCount: %s", err.Error())
		return nil
	}

	lst, err := models.LicenseCacheGets()
	if err != nil {
		logger.Errorf("failed to exec LicenseGetsAll: %s", err.Error())
		return nil
	}

	max := int64(-1)
	m := make(map[int64]*models.License)
	for i := 0; i < len(lst); i++ {
		lst[i].UsedNode = usedNode
		lst[i].TargetVersion = version.Version
		m[lst[i].Id] = lst[i]
		if max < lst[i].Id {
			max = lst[i].Id
		}
	}
	if m[max].PermissionNode < usedNode && time.Since(cache.lastTime).Hours() > 6 {
		ws.SetMessage(758493, "已用节点数超出许可节点数")
		cache.lastTime = time.Now()
	}

	cache.SetLicense(m, stat.Total, stat.LastUpdated)

	ms := time.Since(start).Milliseconds()
	cache.stats.GaugeCronDuration.WithLabelValues("sync_license").Set(float64(ms))
	cache.stats.GaugeSyncNumber.WithLabelValues("sync_license").Set(float64(len(m)))

	logger.Infof("timer: sync license done, cost: %dms, number: %d", ms, len(m))

	return nil
}

func (cache *LicenseCache) StatChanged(total, lastUpdated int64) bool {
	if cache.statTotal == total && cache.statLastUpdated == lastUpdated {
		return false
	}

	return true
}

func (cache *LicenseCache) SetLicenseConfig(licenses map[int64]*models.LicenseConfig, total, lastUpdated int64) {
	cache.Lock()
	cache.config = licenses
	cache.Unlock()

	// only one goroutine used, so no need lock
	cache.statTotal = total
	cache.statLastUpdated = lastUpdated
}

func (cache *LicenseCache) SetLicense(licenses map[int64]*models.License, total, lastUpdated int64) {
	cache.Lock()
	cache.license = licenses
	cache.Unlock()

	// only one goroutine used, so no need lock
	cache.statTotal = total
	cache.statLastUpdated = lastUpdated
}
