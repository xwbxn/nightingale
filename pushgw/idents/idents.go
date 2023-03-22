package idents

import (
	"sync"
	"time"

	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/slice"
	"gorm.io/gorm"
)

type IdentProps struct {
	BusiGroup string
}

type Set struct {
	sync.Mutex
	items map[string]IdentProps
	db    *gorm.DB
}

func New(db *gorm.DB) *Set {
	set := &Set{
		items: make(map[string]IdentProps),
		db:    db,
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
			if err := s.updateTargets(lst, now); err != nil {
				logger.Errorf("failed to update targets: %v", err)
			}
			lst = make(map[string]IdentProps, 100)
			num = 0
		}
	}

	if err := s.updateTargets(lst, now); err != nil {
		logger.Errorf("failed to update targets: %v", err)
	}
}

func (s *Set) updateTargets(lst map[string]IdentProps, now int64) error {
	count := int64(len(lst))
	if count == 0 {
		return nil
	}
	namelist := make([]string, count)
	for ident := range lst {
		namelist = append(namelist, ident)
	}

	ret := s.db.Table("target").Where("ident in ?", namelist).Update("update_at", now)
	if ret.Error != nil {
		return ret.Error
	}

	if ret.RowsAffected == count {
		return nil
	}

	// there are some idents not found in db, so insert them
	var exists []string
	err := s.db.Table("target").Where("ident in ?", namelist).Pluck("ident", &exists).Error
	if err != nil {
		return err
	}

	news := slice.SubString(namelist, exists)
	for i := 0; i < len(news); i++ {
		busigroup := lst[news[i]].BusiGroup
		err = s.db.Exec("INSERT INTO target(ident, busigroup, update_at) VALUES(?, ?)", news[i], busigroup, now).Error
		if err != nil {
			logger.Error("failed to insert target:", news[i], "error:", err)
		}
	}

	return nil
}
