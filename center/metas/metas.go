package metas

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/storage"

	"github.com/toolkits/pkg/logger"
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
		m[meta.IpAddress] = meta
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

	// there are some asset not updated in db, so insert them or update
	for k := range m {
		exists, err := models.AssetGet(s.ctx, "ip = ? and plugin = 'host'", k)
		if err != nil {
			return err
		}
		new := m[k]
		if exists == nil {
			//new
			asset := models.Asset{
				Ident:   new.Hostname,
				Name:    new.Hostname,
				Label:   new.IpAddress,
				GroupId: 1,
				Type:    "主机设备",
				Memo:    "自动注册",
				Plugin:  "host",
				Ip:      new.IpAddress,
			}
			params, _ := json.Marshal(asset)
			asset.Params = string(params)
			// insert
			_, err = asset.AddXH(s.ctx)
			if err != nil {
				logger.Error("failed to insert assets:", k, "error:", err)
			}
		} else if exists.Ident == "" || exists.Configs == "" {
			exists.Ident = new.Hostname
			exists.UpdateAt = time.Now().Unix()
			err := exists.Update(s.ctx, "ident", "configs", "update_at")
			if err != nil {
				logger.Error("failed to update assets:", k, "error:", err)
			}
		}
	}

	newMap := make(map[string]interface{}, count)
	extendMap := make(map[string]interface{})
	for ident, meta := range m {
		if meta.ExtendInfo != nil {
			extendMeta := meta.ExtendInfo
			meta.ExtendInfo = make(map[string]interface{})
			extendMetaStr, err := json.Marshal(extendMeta)
			if err != nil {
				return err
			}
			extendMap[models.WrapExtendIdent(ident)] = extendMetaStr
		}
		newMap[models.WrapIdent(ident)] = meta
	}
	err := storage.MSet(context.Background(), s.redis, newMap)
	if err != nil {
		return err
	}

	if len(extendMap) > 0 {
		err = storage.MSet(context.Background(), s.redis, extendMap)
		if err != nil {
			return err
		}
	}

	return err
}
