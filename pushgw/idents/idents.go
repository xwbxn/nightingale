package idents

import (
	"sync"
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/poster"

	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/slice"
)

type IdentProps struct {
	BusiGroup string
}

type Set struct {
	sync.Mutex
	items map[string]struct{}
	ctx   *ctx.Context
}

func New(ctx *ctx.Context) *Set {
	set := &Set{
		items: make(map[string]struct{}),
		ctx:   ctx,
	}

	set.Init()
	return set
}

func (s *Set) Init() {
	go s.LoopPersist()
}

func (s *Set) MSet(items map[string]IdentProps) {
	s.Lock()
	defer s.Unlock()
	for ident, props := range items {
		s.items[ident] = props
	}
}

func (s *Set) LoopPersist() {
	for {
		time.Sleep(time.Second)
		s.persist()
	}
}

func (s *Set) persist() {
	var items map[string]IdentProps

	s.Lock()
	if len(s.items) == 0 {
		s.Unlock()
		return
	}

	items = s.items
	s.items = make(map[string]IdentProps)
	s.Unlock()

	s.updateTimestamp(items)
}

func (s *Set) updateTimestamp(items map[string]IdentProps) {
	lst := make(map[string]IdentProps, 100)
	now := time.Now().Unix()
	num := 0
	for ident, props := range items {
		lst[ident] = props
		num++
		if num == 100 {
			if err := s.UpdateTargets(lst, now); err != nil {
				logger.Errorf("failed to update targets: %v", err)
			}
			lst = make(map[string]IdentProps, 100)
			num = 0
		}
	}

	if err := s.UpdateTargets(lst, now); err != nil {
		logger.Errorf("failed to update targets: %v", err)
	}
}

type TargetUpdate struct {
	Lst []string `json:"lst"`
	Now int64    `json:"now"`
}

func (s *Set) UpdateTargets(lst []string, now int64) error {
	if !s.ctx.IsCenter {
		t := TargetUpdate{
			Lst: lst,
			Now: now,
		}
		err := poster.PostByUrls(s.ctx, "/v1/n9e/target-update", t)
		return err
	}

	count := int64(len(lst))
	if count == 0 {
		return nil
	}
	namelist := []string{}
	for ident := range lst {
		namelist = append(namelist, ident)
	}

	ret := s.ctx.DB.Table("target").Where("ident in ?", lst).Update("update_at", now)
	if ret.Error != nil {
		return ret.Error
	}

	if ret.RowsAffected == count {
		return nil
	}

	// there are some idents not found in db, so insert them
	var exists []string
	err := s.ctx.DB.Table("target").Where("ident in ?", lst).Pluck("ident", &exists).Error
	if err != nil {
		return err
	}

	news := slice.SubString(namelist, exists)
	for i := 0; i < len(news); i++ {
		err = s.ctx.DB.Exec("INSERT INTO target(ident, update_at) VALUES(?, ?)", news[i], now).Error
		if err != nil {
			logger.Error("failed to insert target:", news[i], "error:", err)
		}
	}

	return nil
}
