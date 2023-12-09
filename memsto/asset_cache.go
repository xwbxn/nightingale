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

type AssetCacheType struct {
	// 资产
	statTotal       int64
	statLastUpdated int64
	// 监控
	mStatTotal       int64
	mStatLastUpdated int64

	ctx   *ctx.Context
	stats *Stats

	sync.RWMutex
	assets map[int64]*models.Asset // key: id
	types  map[string]*models.AssetType
}

func NewAssetCache(ctx *ctx.Context, stats *Stats) *AssetCacheType {
	cache := &AssetCacheType{
		statTotal:        -1,
		statLastUpdated:  -1,
		ctx:              ctx,
		stats:            stats,
		assets:           make(map[int64]*models.Asset),
		mStatTotal:       -1,
		mStatLastUpdated: -1,
	}
	cache.SyncAssets()
	cache.SyncAssetType()
	return cache
}

func (cache *AssetCacheType) StatChanged(total, lastUpdated int64) bool {
	if cache.statTotal == total && cache.statLastUpdated == lastUpdated {
		return false
	}
	return true
}

func (cache *AssetCacheType) MStatChanged(total, lastUpdated int64) bool {
	if cache.mStatTotal == total && cache.mStatLastUpdated == lastUpdated {
		return false
	}
	return true
}

func (cache *AssetCacheType) Set(assets map[int64]*models.Asset, total, lastUpdated, mTotal, mLastUpdated int64) {
	cache.Lock()
	cache.assets = assets
	cache.Unlock()

	// only one goroutine used, so no need lock
	cache.statTotal = total
	cache.statLastUpdated = lastUpdated
	cache.mStatTotal = mTotal
	cache.mStatLastUpdated = mLastUpdated
}

func (cache *AssetCacheType) SetTypes(types map[string]*models.AssetType) {
	cache.Lock()
	cache.types = types
	cache.Unlock()
}

func (cache *AssetCacheType) Get(id int64) (*models.Asset, bool) {
	cache.RLock()
	defer cache.RUnlock()
	val, has := cache.assets[id]
	return val, has
}

func (cache *AssetCacheType) GetAll() []*models.Asset {
	cache.RLock()
	defer cache.RUnlock()

	items := []*models.Asset{}
	for _, item := range cache.assets {
		items = append(items, item)
	}

	return items
}

func (cache *AssetCacheType) GetByType(atype string) []*models.Asset {
	cache.RLock()
	defer cache.RUnlock()

	items := []*models.Asset{}
	for _, item := range cache.assets {
		if item.Type == atype {
			items = append(items, item)
		}
	}

	return items
}

func (cache *AssetCacheType) GetType(id string) (*models.AssetType, bool) {
	cache.RLock()
	defer cache.RUnlock()
	val, has := cache.types[id]
	return val, has
}

func (cache *AssetCacheType) GetByIdent(ident string) []*models.Asset {
	cache.RLock()
	defer cache.RUnlock()

	items := []*models.Asset{}
	for _, item := range cache.assets {
		if item.Ident == ident {
			items = append(items, item)
		}
	}

	return items
}

func (cache *AssetCacheType) GetTypeIds() []string {
	cache.RLock()
	defer cache.RUnlock()

	count := len(cache.types)
	list := make([]string, 0, count)
	for id := range cache.types {
		list = append(list, id)
	}

	return list
}

func (cache *AssetCacheType) SyncAssets() {
	err := cache.syncAssets()
	if err != nil {
		fmt.Println("failed to sync assets:", err)
		exit(1)
	}

	go cache.loopSyncAssets()
}

func (cache *AssetCacheType) loopSyncAssets() {
	duration := time.Duration(9000) * time.Millisecond
	for {
		time.Sleep(duration)
		if err := cache.syncAssets(); err != nil {
			logger.Warning("failed to sync assets:", err)
		}
	}
}

func (cache *AssetCacheType) syncAssets() error {
	start := time.Now()

	stat, err := models.AssetStatistics(cache.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec AssetsStatistics")
	}
	mStat, err := models.MonitoringStatistics(cache.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec MonitoringStatistics")
	}

	// TODO: 这里经常不刷新，暂时改成每次都刷新
	if !cache.StatChanged(stat.Total, stat.LastUpdated) && !cache.MStatChanged(mStat.Total, mStat.LastUpdated) {
		cache.stats.GaugeCronDuration.WithLabelValues("sync_assets").Set(0)
		cache.stats.GaugeSyncNumber.WithLabelValues("sync_assets").Set(0)
		logger.Debug("assets not changed")
		return nil
	}

	lst, err := models.AssetGetsAll(cache.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec AssetGetsAll")
	}
	mlst, err := models.MonitoringGetsAll(cache.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec MonitoringGetsAll")
	}

	m := make(map[int64]*models.Asset)
	for i := 0; i < len(lst); i++ {
		if old, has := cache.assets[lst[i].Id]; has {
			// 刷新缓存时，需要保留原有的监控记录
			lst[i].Health = old.Health
			lst[i].HealthAt = old.HealthAt
			lst[i].Metrics = old.Metrics
		} else {
			//默认值
			lst[i].Health = 2
			lst[i].Metrics = []*models.AssetMetric{}
		}

		lst[i].Monitorings = make([]*models.Monitoring, 0)
		for _, m := range mlst {
			if m.AssetId == lst[i].Id {
				lst[i].Monitorings = append(lst[i].Monitorings, m)
			}
		}
		m[lst[i].Id] = lst[i]
	}

	cache.Set(m, stat.Total, stat.LastUpdated, mStat.Total, mStat.LastUpdated)

	ms := time.Since(start).Milliseconds()
	cache.stats.GaugeCronDuration.WithLabelValues("sync_assets").Set(float64(ms))
	cache.stats.GaugeSyncNumber.WithLabelValues("sync_assets").Set(float64(len(m)))

	logger.Infof("timer: sync assets done, cost: %dms, number: %d", ms, len(m))

	return nil
}

func (cache *AssetCacheType) SyncAssetType() {
	err := cache.syncAssetType()
	if err != nil {
		fmt.Println("failed to sync asset_types:", err)
		exit(1)
	}

	go cache.loopSyncAssetType()
}

func (cache *AssetCacheType) loopSyncAssetType() error {
	duration := time.Duration(9000) * time.Millisecond
	for {
		time.Sleep(duration)
		if err := cache.syncAssetType(); err != nil {
			logger.Warning("failed to sync asset_types:", err)
		}
	}
}

func (cache *AssetCacheType) syncAssetType() error {
	start := time.Now()

	lst, err := models.AssetTypeGetsAll()
	if err != nil {
		return errors.WithMessage(err, "failed to exec AssetTypeGetsAll")
	}

	m := make(map[string]*models.AssetType)
	for i := 0; i < len(lst); i++ {
		m[lst[i].Name] = lst[i]
	}

	cache.SetTypes(m)

	ms := time.Since(start).Milliseconds()
	cache.stats.GaugeCronDuration.WithLabelValues("sync_asset_types").Set(float64(ms))
	cache.stats.GaugeSyncNumber.WithLabelValues("sync_asset_types").Set(float64(len(m)))

	logger.Infof("timer: sync asset_types done, cost: %dms, number: %d", ms, len(m))

	return nil
}
