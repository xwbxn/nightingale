package health

import (
	"context"
	"time"

	"github.com/ccfos/nightingale/v6/alert/aconf"
	"github.com/ccfos/nightingale/v6/alert/astats"
	"github.com/ccfos/nightingale/v6/memsto"
	"github.com/ccfos/nightingale/v6/prom"
	"github.com/ccfos/nightingale/v6/pushgw/writer"
	"github.com/toolkits/pkg/logger"
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
	assets := s.assetCache.GetAll()
	healthRules := make(map[string]*HealthRuleContext)
	for _, asset := range assets {
		// TODO: datsourceID 需要可选传入
		healthRule := NewHealthRuleContext(asset, 1, s.promClients, s.writers, s.assetCache)
		healthRules[healthRule.Hash()] = healthRule
	}

	for hash, rule := range healthRules {
		if _, has := s.healthRules[hash]; !has {
			rule.Prepare()
			rule.Start()
			logger.Debugf("add health rules: %s, %s", rule.Key(), hash)
			s.healthRules[hash] = rule
		}
	}

	for hash, rule := range s.healthRules {
		if _, has := healthRules[hash]; !has {
			rule.Stop()
			delete(s.healthRules, hash)
			logger.Debugf("delete health rules: %s, %s", rule.Key(), hash)
		}
	}
}
