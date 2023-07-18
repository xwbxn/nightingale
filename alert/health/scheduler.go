package health

import (
	"context"
	"fmt"
	"time"

	"github.com/ccfos/nightingale/v6/alert/aconf"
	"github.com/ccfos/nightingale/v6/alert/astats"
	"github.com/ccfos/nightingale/v6/alert/naming"
	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/prom"
	"github.com/ccfos/nightingale/v6/pushgw/writer"
)

type Scheduler struct {
	// key: hash
	healthRules map[string]*HealthRuleContext

	aconf aconf.Alert

	assetCache *memsto.AssetCacheType

	promClients *prom.PromClientMap
	writers     *writer.WritersType

	stats *astats.Stats
}

func NewScheduler(aconf aconf.Alert, ac *memsto.AssetCacheType, promClients *prom.PromClientMap, writers *writer.WritersType, stats *astats.Stats) *Scheduler {
	scheduler := &Scheduler{
		aconf:       aconf,
		healthRules: make(map[string]*HealthRuleContext),

		assetCache: ac,

		promClients: promClients,
		writers:     writers,

		stats: stats,
	}

	go scheduler.LoopSyncHealthRules(context.Background())
	return scheduler
}

func (s *Scheduler) LoopSyncHealthRules(ctx context.Context) {
	time.Sleep(time.Duration(s.aconf.EngineDelay) * time.Second)
	duration := 9000 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			s.syncHealthRules()
		}
	}
}

func (s *Scheduler) syncHealthRules() {
	ids := s.assetCache.GetTypeIds()
	healthRules := make(map[string]*HealthRuleContext)
	for _, id := range ids {
		atype, has := s.assetCache.GetType(id)
		if !has {
			continue
		}

		// datasourceIds := s.promClients.Hit(asset.DatasourceIdsJson)
		datasourceIds := s.promClients.Hit([]int64{0})
		for _, dsId := range datasourceIds {
			if !naming.DatasourceHashRing.IsHit(dsId, fmt.Sprintf("%s", atype.Name), s.aconf.Heartbeat.Endpoint) {
				continue
			}

			healthRule := NewHealthRuleContext(atype, dsId, s.promClients, s.writers, s.assetCache)
			healthRules[healthRule.Hash()] = healthRule
		}
	}

	for hash, rule := range healthRules {
		if _, has := s.healthRules[hash]; !has {
			rule.Prepare()
			rule.Start()
			s.healthRules[hash] = rule
		}
	}

	for hash, rule := range s.healthRules {
		if _, has := healthRules[hash]; !has {
			rule.Stop()
			delete(s.healthRules, hash)
		}
	}
}
