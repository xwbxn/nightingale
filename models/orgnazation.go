package models

import (
	"reflect"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

type Organization struct {
	Id          int64           `json:"id" gorm:"primaryKey"` // id
	Name        string          `json:"name"`                 // 组织名
	ParentId    int64           `json:"parent_id"`            // 组织父id
	Path        string          `json:"-"`                    // 路径
	Children    []*Organization `json:"children" gorm:"-"`
	City        string          `json:"city"`
	Manger      string          `json:"manger"`
	Phone       string          `json:"phone"`
	Address     string          `json:"address"`
	Description string          `json:"description"`
}

func tree(menus []*Organization, pid int64) []*Organization {
	//定义子节点目录
	var nodes []*Organization
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
func (e *Organization) TableName() string {
	return "organization"
}

/*
*
添加单条数据
*
*/
func (e *Organization) Add(ctx *ctx.Context) error {
	return Insert(ctx, e)
}

/*
*
删除单条数据
*
*/
func OrganizationDel(ctx *ctx.Context, ids []string) error {
	if len(ids) == 0 {
		panic("ids empty")
	}
	return DB(ctx).Where("id in ?", ids).Delete(new(Organization)).Error
}

func OrganizationList(ctx *ctx.Context) ([]*Organization, error) {
	var lst []*Organization
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
func OrganizationGet(ctx *ctx.Context, where string, args ...interface{}) (*Organization, error) {
	var lst []*Organization
	err := DB(ctx).Where(where, args...).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	if len(lst) == 0 {
		return nil, nil
	}
	return lst[0], nil
}

func OrganizationGets(ctx *ctx.Context) ([]Organization, error) {
	session := DB(ctx)

	var lst []Organization
	err := session.Find(&lst).Error

	return lst, err
}

func OrganizationGetsByIds(ctx *ctx.Context, ids []int64) ([]Organization, error) {
	session := DB(ctx)

	var lst []Organization
	err := session.Where("ID IN ?", ids).Find(&lst).Error

	return lst, err
}

func OrganizationGetById(ctx *ctx.Context, id int64) (*Organization, error) {
	return OrganizationGet(ctx, "id=?", id)
}

func (m *Organization) Update(ctx *ctx.Context, data interface{}) error {
	return DB(ctx).Model(m).Select("name", "parent_id", "path").Updates(data).Error
}

func OrganizationDels(ctx *ctx.Context, ids []int64) error {
	for i := 0; i < len(ids); i++ {
		session := DB(ctx).Where("id = ?", ids[i])
		ret := session.Delete(&Organization{})
		if ret.Error != nil {
			return ret.Error
		}
	}
	return nil
}

// 前端输出组织接口
// 前端结构定义
type feOrg struct {
	Id          int64    `json:"id"`   // id
	Name        string   `json:"name"` // 组织名
	ParentId    int64    `json:"-"`    // 组织父id
	Children    []*feOrg `json:"children"`
	City        string   `json:"city"`
	Manger      string   `json:"manger"`
	Phone       string   `json:"phone"`
	Address     string   `json:"address"`
	Description string   `json:"description"`
}

// 前端数组织生成
func feTree(menus []*feOrg, pid int64) []*feOrg {
	//定义子节点目录
	var nodes []*feOrg
	if reflect.ValueOf(menus).IsValid() {
		//循环所有一级菜单
		for _, v := range menus {
			//查询所有该菜单下的所有子菜单
			if v.ParentId == pid {
				//特别注意压入元素不是单个所有加三个 **...** 告诉切片无论多少元素一并压入
				v.Children = append(v.Children, feTree(menus, v.Id)...)
				nodes = append(nodes, v)
			}

		}
	}
	return nodes
}

// 前端组织树
func OrganizationTreeGetsFE(ctx *ctx.Context) ([]*feOrg, error) {
	var lst []*Organization
	var felst []*feOrg
	err := DB(ctx).Find(&lst).Error
	for i := 0; i < len(lst); i++ {
		felst = append(felst, &feOrg{
			Id:          lst[i].Id,
			Name:        lst[i].Name,
			ParentId:    lst[i].ParentId,
			City:        lst[i].City,
			Manger:      lst[i].Manger,
			Phone:       lst[i].Phone,
			Address:     lst[i].Address,
			Description: lst[i].Description,
		})

	}
	x := feTree(felst, 0)
	return x, err
}

// 根据条件统计个数
func OrganizationCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&Organization{}).Where(where, args...))
}

//根据name模糊查询id
func OrgIdByName(ctx *ctx.Context, name string) ([]int64, error) {
	ids := make([]int64, 0)
	name = "%" + name + "%"
	err := DB(ctx).Debug().Model(&Organization{}).Where("name like ?", name).Pluck("id", &ids).Error
	return ids, err
}

func IsContain(items []int64, item int64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
