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

type ProviderCacheType struct {
	statTotal       int64
	statLastUpdated int64
	ctx             *ctx.Context
	stats           *Stats

	sync.RWMutex
	hps map[int64]*models.Asset // key: id
}

func NewProviderCache(ctx *ctx.Context, stats *Stats) *ProviderCacheType {
	hpc := &ProviderCacheType{
		statTotal:       -1,
		statLastUpdated: -1,
		ctx:             ctx,
		stats:           stats,
		hps:             make(map[int64]*models.Asset),
	}
	hpc.SyncProviders()
	return hpc
}

func (hpc *ProviderCacheType) StatChanged(total, lastUpdated int64) bool {
	if hpc.statTotal == total && hpc.statLastUpdated == lastUpdated {
		return false
	}

	return true
}

func (hpc *ProviderCacheType) Set(hps map[int64]*models.Asset, total, lastUpdated int64) {
	hpc.Lock()
	hpc.hps = hps
	hpc.Unlock()

	// only one goroutine used, so no need lock
	hpc.statTotal = total
	hpc.statLastUpdated = lastUpdated
}

func (hpc *ProviderCacheType) GetByIdent(ident string) []*models.Asset {
	hpc.RLock()
	defer hpc.RUnlock()

	items := []*models.Asset{}
	for _, item := range hpc.hps {
		if item.Ident == ident {
			items = append(items, item)
		}
	}

	return items
}

func (hpc *ProviderCacheType) SyncProviders() {
	err := hpc.syncProviders()
	if err != nil {
		fmt.Println("failed to sync assets:", err)
		exit(1)
	}

	go hpc.loopSyncProviders()
}

func (hpc *ProviderCacheType) loopSyncProviders() {
	duration := time.Duration(9000) * time.Millisecond
	for {
		time.Sleep(duration)
		if err := hpc.syncProviders(); err != nil {
			logger.Warning("failed to sync assets:", err)
		}
	}
}

func (hpc *ProviderCacheType) syncProviders() error {
	start := time.Now()

	stat, err := models.AssetStatistics(hpc.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec AssetsStatistics")
	}

	if !hpc.StatChanged(stat.Total, stat.LastUpdated) {
		hpc.stats.GaugeCronDuration.WithLabelValues("sync_assets").Set(0)
		hpc.stats.GaugeSyncNumber.WithLabelValues("sync_assets").Set(0)
		logger.Debug("assets not changed")
		return nil
	}

	lst, err := models.AssetGetAll(hpc.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec HttpProviderGetAll")
	}

	m := make(map[int64]*models.Asset)
	for i := 0; i < len(lst); i++ {
		m[lst[i].Id] = lst[i]
	}

	hpc.Set(m, stat.Total, stat.LastUpdated)

	ms := time.Since(start).Milliseconds()
	hpc.stats.GaugeCronDuration.WithLabelValues("sync_assets").Set(float64(ms))
	hpc.stats.GaugeSyncNumber.WithLabelValues("sync_assets").Set(float64(len(m)))

	logger.Infof("timer: sync assets done, cost: %dms, number: %d", ms, len(m))

	return nil
}
