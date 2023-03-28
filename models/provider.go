package models

import (
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Provider struct {
	Id       int64  `json:"id" gorm:"primaryKey"`
	Version  string `json:"version"`
	Ident    string `json:"ident"`
	GroupId  int64  `json:"group_id"`
	Configs  string `json:"configs"`
	CreateAt int64  `json:"create_at"`
	CreateBy string `json:"create_by"`
	UpdateAt int64  `json:"update_at"`
	UpdateBy string `json:"update_by"`
}

func (hp *Provider) TableName() string {
	return "http_provider"
}

func (hp *Provider) Verify() error {
	return nil
}

func (hp *Provider) Update(ctx *ctx.Context, selectField interface{}, selectFields ...interface{}) error {
	if err := hp.Verify(); err != nil {
		return err
	}

	return DB(ctx).Model(hp).Select(selectField, selectFields...).Updates(hp).Error
}

func ProviderCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&Provider{}).Where(where, args...))
}

func (hp *Provider) Add(ctx *ctx.Context) error {
	if err := hp.Verify(); err != nil {
		return err
	}

	num, err := ProviderCount(ctx, "ident=? and group_id=?", hp.Ident, hp.GroupId)
	if err != nil {
		return errors.WithMessage(err, "failed to count http-provider")
	}

	if num > 0 {
		return errors.New("Provider already exists")
	}

	now := time.Now().Unix()
	hp.CreateAt = now
	hp.UpdateAt = now
	return Insert(ctx, hp)
}

func (hp *Provider) Del(ctx *ctx.Context) error {
	return DB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id=?", hp.Id).Delete(&Provider{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func ProviderGet(ctx *ctx.Context, where string, args ...interface{}) (*Provider, error) {
	var lst []*Provider
	err := DB(ctx).Where(where, args...).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}

	return lst[0], nil
}

func ProviderGetById(ctx *ctx.Context, id int64) (*Provider, error) {
	return ProviderGet(ctx, "id = ?", id)
}

func ProviderGetAll(ctx *ctx.Context) ([]*Provider, error) {
	var lst []*Provider
	err := DB(ctx).Find(&lst).Error
	return lst, err
}

func ProviderStatistics(ctx *ctx.Context) (*Statistics, error) {
	session := DB(ctx).Model(&Provider{}).Select("count(*) as total", "max(update_at) as last_updated")

	var stats []*Statistics
	err := session.Find(&stats).Error
	if err != nil {
		return nil, err
	}

	return stats[0], nil
}
