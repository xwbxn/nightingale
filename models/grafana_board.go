package models

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/toolkits/pkg/str"
	"gorm.io/gorm"
)

type GrafanaBoard struct {
	Id         int64  `json:"id" gorm:"primaryKey"`
	GroupId    int64  `json:"group_id"`
	Name       string `json:"name"`
	Ident      string `json:"ident"`
	Tags       string `json:"tags"`
	CreateAt   int64  `json:"create_at"`
	CreateBy   string `json:"create_by"`
	UpdateAt   int64  `json:"update_at"`
	UpdateBy   string `json:"update_by"`
	Configs    string `json:"configs" gorm:"-"`
	Public     int    `json:"public"` // 0: false, 1: true
	GrafanaId  int64  `json:"grafana_id"`
	GrafanaUrl string `json:"grafana_url"`
}

func (g *GrafanaBoard) TableName() string {
	return "grafana_board"
}

func (g *GrafanaBoard) Verify() error {
	if g.Name == "" {
		return errors.New("Name is blank")
	}

	if str.Dangerous(g.Name) {
		return errors.New("Name has invalid characters")
	}

	return nil
}

func (g *GrafanaBoard) GrafanaBoardCanRenameIdent(ident string) (bool, error) {
	if ident == "" {
		return true, nil
	}

	cnt, err := Count(DB().Model(g).Where("ident=? and id <> ?", ident, g.Id))
	if err != nil {
		return false, err
	}

	return cnt == 0, nil
}

func (g *GrafanaBoard) Add() error {
	if err := g.Verify(); err != nil {
		return err
	}

	now := time.Now().Unix()
	g.CreateAt = now
	g.UpdateAt = now

	return Insert(g)
}

func (g *GrafanaBoard) Update(selectField interface{}, selectFields ...interface{}) error {
	if err := g.Verify(); err != nil {
		return err
	}

	return DB().Model(g).Select(selectField, selectFields...).Updates(g).Error
}

func (g *GrafanaBoard) Del() error {
	return DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id=?", g.Id).Delete(&GrafanaBoard{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func GrafanaBoardGetByID(id int64) (*GrafanaBoard, error) {
	var lst []*GrafanaBoard
	err := DB().Where("id = ?", id).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}

	return lst[0], nil
}

// BoardGet for detail page
func GrafanaBoardGet(where string, args ...interface{}) (*GrafanaBoard, error) {
	var lst []*GrafanaBoard
	err := DB().Where(where, args...).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}

	return lst[0], nil
}

func GrafanaBoardCount(where string, args ...interface{}) (num int64, err error) {
	return Count(DB().Model(&GrafanaBoard{}).Where(where, args...))
}

func GrafanaBoardExists(where string, args ...interface{}) (bool, error) {
	num, err := GrafanaBoardCount(where, args...)
	return num > 0, err
}

// BoardGets for list page
func GrafanaBoardGets(groupId int64, query string) ([]GrafanaBoard, error) {
	session := DB().Where("group_id=?", groupId).Order("name")

	arr := strings.Fields(query)
	if len(arr) > 0 {
		for i := 0; i < len(arr); i++ {
			if strings.HasPrefix(arr[i], "-") {
				q := "%" + arr[i][1:] + "%"
				session = session.Where("name not like ? and tags not like ?", q, q)
			} else {
				q := "%" + arr[i] + "%"
				session = session.Where("(name like ? or tags like ?)", q, q)
			}
		}
	}

	var objs []GrafanaBoard
	err := session.Find(&objs).Error
	return objs, err
}
