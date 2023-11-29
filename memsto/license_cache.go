package memsto

import (
	"fmt"
	"sync"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/pkg/errors"
	"github.com/toolkits/pkg/logger"
)

type LicenseCache struct {
	statTotal       int64
	statLastUpdated int64
	ctx             *ctx.Context
	stats           *Stats

	sync.RWMutex
	config map[int64]*models.LicenseConfig
}

func NewLicenseCache(ctx *ctx.Context, stats *Stats) *LicenseCache {
	cache := &LicenseCache{
		statTotal:       -1,
		statLastUpdated: -1,
		ctx:             ctx,
		stats:           stats,
		config:          make(map[int64]*models.LicenseConfig),
	}
	cache.SyncLicense()
	return cache
}

func (cache *LicenseCache) GetByLicenseId(id int64) *models.LicenseConfig {
	cache.RLock()
	defer cache.RUnlock()
	return cache.config[id]
}

func (cache *LicenseCache) SyncLicense() {
	err := cache.syncLicense()
	if err != nil {
		fmt.Println("failed to sync assets:", err)
		exit(1)
	}

	go cache.loopSyncLicense()
}

func (cache *LicenseCache) loopSyncLicense() {
	duration := time.Duration(9000) * time.Millisecond
	for {
		time.Sleep(duration)
		if err := cache.syncLicense(); err != nil {
			logger.Warning("failed to sync assets:", err)
		}
	}
}

func (cache *LicenseCache) syncLicense() error {
	start := time.Now()

	stat, err := models.LicenseStatistics(cache.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec LicenseStatistics")
	}

	if !cache.StatChanged(stat.Total, stat.LastUpdated) {
		cache.stats.GaugeCronDuration.WithLabelValues("sync_assets").Set(0)
		cache.stats.GaugeSyncNumber.WithLabelValues("sync_assets").Set(0)
		logger.Debug("License not changed")
		return nil
	}

	lst, err := models.LicenseConfigGetsAll(cache.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec LicenseGetsAll")
	}

	m := make(map[int64]*models.LicenseConfig)
	for i := 0; i < len(lst); i++ {
		m[lst[i].Id] = lst[i]
	}

	cache.Set(m, stat.Total, stat.LastUpdated)

	ms := time.Since(start).Milliseconds()
	cache.stats.GaugeCronDuration.WithLabelValues("sync_assets").Set(float64(ms))
	cache.stats.GaugeSyncNumber.WithLabelValues("sync_assets").Set(float64(len(m)))

	logger.Infof("timer: sync assets done, cost: %dms, number: %d", ms, len(m))

	return nil
}

func (cache *LicenseCache) StatChanged(total, lastUpdated int64) bool {
	if cache.statTotal == total && cache.statLastUpdated == lastUpdated {
		return false
	}

	return true
}

func (cache *LicenseCache) Set(licenses map[int64]*models.LicenseConfig, total, lastUpdated int64) {
	cache.Lock()
	cache.config = licenses
	cache.Unlock()

	// only one goroutine used, so no need lock
	cache.statTotal = total
	cache.statLastUpdated = lastUpdated
}
