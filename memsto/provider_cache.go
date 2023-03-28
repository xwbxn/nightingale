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
	hps map[int64]*models.Provider // key: id
}

func NewProviderCache(ctx *ctx.Context, stats *Stats) *ProviderCacheType {
	hpc := &ProviderCacheType{
		statTotal:       -1,
		statLastUpdated: -1,
		ctx:             ctx,
		stats:           stats,
		hps:             make(map[int64]*models.Provider),
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

func (hpc *ProviderCacheType) Set(hps map[int64]*models.Provider, total, lastUpdated int64) {
	hpc.Lock()
	hpc.hps = hps
	hpc.Unlock()

	// only one goroutine used, so no need lock
	hpc.statTotal = total
	hpc.statLastUpdated = lastUpdated
}

func (hpc *ProviderCacheType) GetByIdentAndGroup(ident string, groupId int64) *models.Provider {
	hpc.RLock()
	defer hpc.RUnlock()

	for _, item := range hpc.hps {
		if item.Ident == ident && item.GroupId == groupId {
			return item
		}
	}

	return nil
}

func (hpc *ProviderCacheType) SyncProviders() {
	err := hpc.syncProviders()
	if err != nil {
		fmt.Println("failed to sync http providers:", err)
		exit(1)
	}

	go hpc.loopSyncProviders()
}

func (hpc *ProviderCacheType) loopSyncProviders() {
	duration := time.Duration(9000) * time.Millisecond
	for {
		time.Sleep(duration)
		if err := hpc.syncProviders(); err != nil {
			logger.Warning("failed to sync http providers:", err)
		}
	}
}

func (hpc *ProviderCacheType) syncProviders() error {
	start := time.Now()

	stat, err := models.ProviderStatistics(hpc.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec HttpProviderStatistics")
	}

	// 因为有内容变化，因此全部同步，不进行状态判定
	lst, err := models.ProviderGetAll(hpc.ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to exec HttpProviderGetAll")
	}

	m := make(map[int64]*models.Provider)
	for i := 0; i < len(lst); i++ {
		m[lst[i].Id] = lst[i]
	}

	hpc.Set(m, stat.Total, stat.LastUpdated)

	ms := time.Since(start).Milliseconds()
	hpc.stats.GaugeCronDuration.WithLabelValues("sync_providers").Set(float64(ms))
	hpc.stats.GaugeSyncNumber.WithLabelValues("sync_providers").Set(float64(len(m)))

	logger.Infof("timer: sync providers done, cost: %dms, number: %d", ms, len(m))

	return nil
}
