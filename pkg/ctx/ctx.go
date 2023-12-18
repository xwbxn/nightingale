package ctx

import (
	"context"

	"github.com/ccfos/nightingale/v6/conf"

	"gorm.io/gorm"
)

type Context struct {
	DB        *gorm.DB
	CenterApi conf.CenterApi
	Ctx       context.Context
	IsCenter  bool
}

func NewContext(ctx context.Context, db *gorm.DB, isCenter bool, centerApis ...conf.CenterApi) *Context {
	var api conf.CenterApi
	if len(centerApis) > 0 {
		api = centerApis[0]
	}

	return &Context{
		Ctx:       ctx,
		DB:        db,
		CenterApi: api,
		IsCenter:  isCenter,
	}
}

// set db to Context
func (c *Context) SetDB(db *gorm.DB) {
	c.DB = db
}

// get context from Context
func (c *Context) GetContext() context.Context {
	return c.Ctx
}

// get db from Context
func (c *Context) GetDB() *gorm.DB {
	return c.DB
}

// 自动事务, fc返回error则rollback，返回nil则commit
func (c *Context) Transaction(fc func(ctx *Context) error) error {
	return c.DB.Transaction(func(tx *gorm.DB) error {
		ctx := NewContext(c.Ctx, tx, c.IsCenter, c.CenterApi)
		err := fc(ctx)
		return err
	})
}

// 开启手动事务, 显式控制commit和rollback
func (c *Context) BeginTransacion() *Context {
	return &Context{
		Ctx:       c.Ctx,
		DB:        c.DB.Begin(),
		CenterApi: c.CenterApi,
		IsCenter:  c.IsCenter,
	}
}

func (c *Context) Rollback() {
	c.DB.Rollback()
}

func (c *Context) Commit() {
	c.DB.Commit()
}
