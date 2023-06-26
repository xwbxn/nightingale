package metas

import (
	"context"
	"sync"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/storage"

	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/slice"
)

type Set struct {
	sync.RWMutex
	items map[string]models.HostMeta
	redis storage.Redis
	ctx   *ctx.Context
}

func New(ctx *ctx.Context, redis storage.Redis) *Set {
	set := &Set{
		items: make(map[string]models.HostMeta),
		redis: redis,
		ctx:   ctx,
	}

	set.Init()
	return set
}

func (s *Set) Init() {
	go s.LoopPersist()
}

func (s *Set) MSet(items map[string]models.HostMeta) {
	s.Lock()
	defer s.Unlock()
	for ident, meta := range items {
		s.items[ident] = meta
	}
}

func (s *Set) Set(ident string, meta models.HostMeta) {
	s.Lock()
	defer s.Unlock()
	s.items[ident] = meta
}

func (s *Set) LoopPersist() {
	for {
		time.Sleep(time.Second)
		s.persist()
	}
}

func (s *Set) persist() {
	var items map[string]models.HostMeta

	s.Lock()
	if len(s.items) == 0 {
		s.Unlock()
		return
	}

	items = s.items
	s.items = make(map[string]models.HostMeta)
	s.Unlock()

	s.updateMeta(items)
}

func (s *Set) updateMeta(items map[string]models.HostMeta) {
	m := make(map[string]models.HostMeta, 100)
	num := 0

	for _, meta := range items {
		m[meta.Hostname] = meta
		num++
		if num == 100 {
			if err := s.updateTargets(m); err != nil {
				logger.Errorf("failed to update targets: %v", err)
			}
			m = make(map[string]models.HostMeta, 100)
			num = 0
		}
	}

	if err := s.updateTargets(m); err != nil {
		logger.Errorf("failed to update targets: %v", err)
	}
}

func (s *Set) updateTargets(m map[string]models.HostMeta) error {
	count := int64(len(m))
	if count == 0 {
		return nil
	}

	// there are some idents not found in db, so insert them
	var exists []string
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	err := s.ctx.DB.Table("assets").Where("ident in ? and plugin = 'host'", keys).Pluck("ident", &exists).Error
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	news := slice.SubString(keys, exists)
	for i := 0; i < len(news); i++ {
		conf, _ := models.AssetGenConfig("主机", make(map[string]interface{}))

		new := m[news[i]]
		err = s.ctx.DB.Exec("INSERT INTO assets(configs, group_id, ident, name, label, type, memo, create_at, create_by, plugin) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?,?)",
			conf.String(), 0, new.Hostname, new.Hostname, new.RemoteAddr, "主机", "auto", now, "system", "host").Error
		if err != nil {
			logger.Error("failed to insert assets:", news[i], "error:", err)
		}
	}

	newMap := make(map[string]interface{}, count)
	for ident, meta := range m {
		newMap[models.WrapIdent(ident)] = meta
	}
	err = storage.MSet(context.Background(), s.redis, newMap)
	return err
}
