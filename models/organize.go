package models

import (
	"reflect"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/toolkits/pkg/logger"
)

type Organize struct {
	Id       int64       `json:"id" gorm:"primaryKey"` // id
	Name     string      `json:"name"`                 // 组织名
	ParentId int64       `json:"parent_id"`            // 组织父id
	Path     string      `json:"-"`                    // 路径
	Children []*Organize `json:"children" gorm:"-"`
	Son      int64       `json:"-"`
}

func tree(menus []*Organize, pid int64) []*Organize {
	//定义子节点目录
	var nodes []*Organize
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

/*
*
获取表名
*
*/
func (e *Organize) TableName() string {
	return "organize"
}

/*
*
添加单条数据
*
*/
func (e *Organize) Add(ctx *ctx.Context) error {
	return Insert(ctx, e)
}

/*
*
删除单条数据
*
*/
func OrganizeDel(ctx *ctx.Context, ids []string) error {
	if len(ids) == 0 {
		panic("ids empty")
	}
	return DB(ctx).Where("id in ?", ids).Delete(new(Organize)).Error
}

/*
*
更新单条数据
*
*/
func (u *Organize) Update(ctx *ctx.Context, selectField interface{}, selectFields ...interface{}) error {
	return DB(ctx).Model(u).Select(selectField, selectFields...).Updates(u).Error
}

func OrganizeList(ctx *ctx.Context) ([]*Organize, error) {
	var lst []*Organize
	err := DB(ctx).Find(&lst).Error

	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}
	return tree(lst, 0), nil
}

/*
*
获取一组数据
*
*/
func OrganizeGet(ctx *ctx.Context, where string, args ...interface{}) (*Organize, error) {
	var lst []*Organize
	err := DB(ctx).Where(where, args...).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}
	return lst[0], nil
}

func OrganizeGets(ctx *ctx.Context) ([]Organize, error) {
	session := DB(ctx)

	var lst []Organize
	err := session.Find(&lst).Error

	return lst, err
}

func OrganizeGetById(ctx *ctx.Context, id int64) (*Organize, error) {
	return OrganizeGet(ctx, "id=?", id)
}

func (m *Organize) UpdateAll(ctx *ctx.Context, id int64, name string, parent_id int64, path string) error {

	modes, err := OrganizeGetById(ctx, id)
	if modes == nil {
		logger.Errorf("Sorry,This id is empty, %s", err)
	}
	modes.Name = name
	modes.ParentId = parent_id
	modes.Path = path

	return DB(ctx).Model(m).Updates(modes).Error
}

func OrganizeDels(ctx *ctx.Context, ids []int64) error {
	for i := 0; i < len(ids); i++ {
		session := DB(ctx).Where("id = ?", ids[i])
		ret := session.Delete(&Organize{})
		if ret.Error != nil {
			return ret.Error
		}

	}

	return nil
}

func (m *Organize) UpdatesAll(ctx *ctx.Context, ids []int64, name string, parent_id int64, path string) error {
	for i := 0; i < len(ids); i++ {
		modes, err := OrganizeGetById(ctx, ids[i])
		if modes == nil {
			logger.Errorf("Sorry,This id is empty, %s", err)
		}
		modes.Name = name
		modes.ParentId = parent_id
		modes.Path = path

		return DB(ctx).Model(m).Updates(modes).Error

	}
	return nil
}

// 前端输出组织接口
// 前端结构定义
type FeOrg struct {
	Id       int64    `json:"id"`   // id
	Name     string   `json:"name"` // 组织名
	ParentId int64    `json:"-"`    // 组织父id
	Children []*FeOrg `json:"children"`
}

// 前端数组织生成
func fetree(menus []*FeOrg, pid int64) []*FeOrg {
	//定义子节点目录
	var nodes []*FeOrg
	if reflect.ValueOf(menus).IsValid() {
		//循环所有一级菜单
		for _, v := range menus {
			//查询所有该菜单下的所有子菜单
			if v.ParentId == pid {
				//特别注意压入元素不是单个所有加三个 **...** 告诉切片无论多少元素一并压入
				v.Children = append(v.Children, fetree(menus, v.Id)...)
				nodes = append(nodes, v)
			}

		}
	}
	return nodes
}

// 前端组织树
func OrgList(ctx *ctx.Context) ([]*FeOrg, error) {
	var lst []*Organize
	var felst []*FeOrg
	err := DB(ctx).Find(&lst).Error
	for i := 0; i < len(lst); i++ {
		felst = append(felst, &FeOrg{
			Id:       lst[i].Id,
			Name:     lst[i].Name,
			ParentId: lst[i].ParentId,
		})

	}
	x := fetree(felst, 0)
	return x, err
}
