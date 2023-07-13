package models

import (
	"reflect"
	"strings"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/toolkits/pkg/logger"
)

type Groups struct {
	Id       int64     `json:"id" gorm:"primaryKey"` // id
	Name     string    `json:"name"`                 // 组织名
	ParentId int64     `json:"parent_id"`            // 组织父id
	Path     string    `json:"path"`                 // 路径
	Children []*Groups `json:"children" gorm:"-"`
}

func tree(menus []*Groups, pid int64) []*Groups {
	//定义子节点目录
	var nodes []*Groups
	if reflect.ValueOf(menus).IsValid() {
		//循环所有一级菜单
		for _, v := range menus {
			//查询所有该菜单下的所有子菜单
			if v.ParentId == pid {
				//特别注意压入元素不是单个所有加三个 **...** 告诉切片无论多少元素一并压入
				v.Children = append(v.Children, tree(menus, v.Id)...)
				nodes = append(nodes, v)
			}

		}
	}
	return nodes
}

/**
获取表名
**/
func (e *Groups) TableName() string {
	return "groups"
}

/**
添加单条数据
**/
func (e *Groups) Add(ctx *ctx.Context) error {
	return Insert(ctx, e)
}

/**
删除单条数据
**/
func GroupsDel(ctx *ctx.Context, ids []string) error {
	if len(ids) == 0 {
		panic("ids empty")
	}
	return DB(ctx).Where("id in ?", ids).Delete(new(Groups)).Error
}

/**
更新单条数据
**/
func (u *Groups) Update(ctx *ctx.Context, selectField interface{}, selectFields ...interface{}) error {
	return DB(ctx).Model(u).Select(selectField, selectFields...).Updates(u).Error
}

func GroupsList(ctx *ctx.Context) ([]*Groups, error) {
	var lst []*Groups
	err := DB(ctx).Find(&lst).Error

	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}
	return tree(lst, 0), nil
}

/**
获取一组数据
**/
func GroupsGet(ctx *ctx.Context, where string, args ...interface{}) (*Groups, error) {
	var lst []*Groups
	err := DB(ctx).Where(where, args...).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}
	return lst[0], nil
}

func GroupsTotal(ctx *ctx.Context, prods []string, bgid, stime, etime int64, dsIds []int64, query, name string, parent_id int64, path string) (int64, error) {
	session := DB(ctx).Model(&Groups{}).Where("last_eval_time between ? and ?", stime, etime)

	if len(prods) > 0 {
		session = session.Where("rule_prod in ?", prods)
	}

	if bgid > 0 {
		session = session.Where("group_id = ?", bgid)
	}

	if len(dsIds) > 0 {
		session = session.Where("datasource_id in ?", dsIds)
	}

	if len(name) > 0 {
		session = session.Where("cate in ?", name)
	}

	if len(path) > 0 {
		session = session.Where("cate in ?", path)
	}

	if parent_id > 0 {
		session = session.Where("cate in ?", parent_id)
	}

	if query != "" {
		arr := strings.Fields(query)
		for i := 0; i < len(arr); i++ {
			qarg := "%" + arr[i] + "%"
			session = session.Where("rule_name like ? or tags like ?", qarg, qarg)
		}
	}

	return Count(session)
}

func GroupsGets(ctx *ctx.Context, prods []string, bgid, stime, etime int64, dsIds []int64, query string, name string, parent_id int64, path string, limit, offset int) ([]Groups, error) {
	session := DB(ctx).Where("last_eval_time between ? and ?", stime, etime)

	if len(prods) != 0 {
		session = session.Where("parent_id in ?", parent_id)
	}

	if bgid > 0 {
		session = session.Where("name = ?", name)
	}

	if len(dsIds) > 0 {
		session = session.Where("path in ?", path)
	}

	if len(name) > 0 {
		session = session.Where("cate in ?", name)
	}

	if parent_id > 0 {
		session = session.Where("cate in ?", parent_id)
	}

	if len(path) > 0 {
		session = session.Where("cate in ?", path)
	}

	if query != "" {
		arr := strings.Fields(query)
		for i := 0; i < len(arr); i++ {
			qarg := "%" + arr[i] + "%"
			session = session.Where("rule_name like ? or tags like ?", qarg, qarg)
		}
	}

	var lst []Groups
	err := session.Order("id desc").Limit(limit).Offset(offset).Find(&lst).Error

	return lst, err
}

func GroupsGetById(ctx *ctx.Context, id int64) (*Groups, error) {
	return GroupsGet(ctx, "id=?", id)
}

func (m *Groups) UpdateAll(ctx *ctx.Context, id int64, name string, parent_id int64, path string) error {

	modes, err := GroupsGetById(ctx, id)
	if modes == nil {
		logger.Errorf("Sorry,This id is empty, %s", err)
	}
	modes.Name = name
	modes.ParentId = parent_id
	modes.Path = path

	return DB(ctx).Model(m).Updates(modes).Error
}

func GroupsDels(ctx *ctx.Context, ids []int64, bgid ...int64) error {
	for i := 0; i < len(ids); i++ {
		session := DB(ctx).Where("id = ?", ids[i])
		if len(bgid) > 0 {
			session = session.Where("group_id = ?", bgid[0])
		}
		ret := session.Delete(&Groups{})
		if ret.Error != nil {
			return ret.Error
		}

	}

	return nil
}
