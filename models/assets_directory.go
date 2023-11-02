// Package models  资产目录
// date : 2023-9-19 11:32
// desc : 资产目录
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// AssetsDirectory  资产目录。
// 说明:
// 表名:assets_directory
// group: AssetsDirectory
// version:2023-9-19 11:32
type AssetsDirectory struct {
	Id        int64          `gorm:"column:id;primaryKey" json:"id" `                          //type:BIGINT   comment:主键        version:2023-9-19 11:32
	Name      string         `gorm:"column:name" json:"name" `                                 //type:string   comment:名称        version:2023-9-19 11:32
	ParentId  int64          `gorm:"column:parent_id" json:"parent_id" `                       //type:BIGINT   comment:父节点      version:2023-9-19 11:32
	Sort      int64          `gorm:"column:sort" json:"sort" `                                 //type:*int     comment:序号        version:2023-9-19 11:32
	CreatedBy string         `gorm:"column:created_by" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-9-19 11:32
	CreatedAt int64          `gorm:"column:created_at" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-9-19 11:32
	UpdatedBy string         `gorm:"column:updated_by" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-9-19 11:32
	UpdatedAt int64          `gorm:"column:updated_at" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-9-19 11:32
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at" swaggerignore:"true"` //type:string   comment:删除时间    version:2023-9-19 11:32
}

type AssetsDirTree struct {
	Id       int64            `gorm:"column:id;primaryKey" json:"id" `    //type:BIGINT   comment:主键        version:2023-9-19 11:32
	Name     string           `gorm:"column:name" json:"name" `           //type:string   comment:名称        version:2023-9-19 11:32
	ParentId int64            `gorm:"column:parent_id" json:"parent_id" ` //type:BIGINT   comment:父节点      version:2023-9-19 11:32
	Sort     int64            `gorm:"column:sort" json:"sort" `           //type:*int     comment:序号        version:2023-9-19 11:32
	Count    int64            `gorm:"-" json:"count" `
	Children []*AssetsDirTree `gorm:"-" json:"children" `
}

// TableName 表名:assets_directory，资产目录。
// 说明:
func (a *AssetsDirectory) TableName() string {
	return "assets_directory"
}

func BuildDirTree(ctx *ctx.Context, rootId int64) (*AssetsDirTree, error) {
	queue := make([]*AssetsDirTree, 0)

	rootAssetsDirTree, err := AssetsDirectoryGetsMap(ctx, map[string]interface{}{"id": rootId})
	if err != nil || len(rootAssetsDirTree) == 0 {
		return rootAssetsDirTree[0], err
	}
	queue = append(queue, rootAssetsDirTree[0])
	forIndex := 0
	flogTure := true
	for flogTure {
		if forIndex == len(queue) {
			break
		}
		lst, err := AssetsDirectoryGetsMap(ctx, map[string]interface{}{"parent_id": queue[forIndex].Id})
		if err != nil {
			return nil, err
		}
		mapLst := make(map[int64]*AssetsDirTree)
		for _, val := range lst {
			mapLst[val.Sort] = val
		}
		var top, count int64
		top = 0
		count = 0
		flog := true
		for flog {
			if count == int64(len(lst)) {
				break
			}

			queue[forIndex].Children = append(queue[forIndex].Children, mapLst[top])
			queue = append(queue, mapLst[top])
			count++
			top = (*(mapLst[top])).Id
			// for _, val := range lst {
			// 	if val.Sort == top {
			// 		queue[forIndex].Children = append(queue[forIndex].Children, val)
			// 		queue = append(queue, val)
			// 		count++
			// 		top = val.Id
			// 		break
			// 	}
			// }

		}
		forIndex++
	}
	return rootAssetsDirTree[0], err
}

// 条件查询
func AssetsDirectoryGets(ctx *ctx.Context, query string, limit, offset int) ([]AssetsDirectory, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	// 这里使用列名的硬编码构造查询参数, 避免从前台传入造成注入风险
	if query != "" {
		q := "%" + query + "%"
		session = session.Where("id like ?", q)
	}

	var lst []AssetsDirectory
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func AssetsDirectoryGetById(ctx *ctx.Context, id int64) (*AssetsDirectory, error) {
	var obj *AssetsDirectory
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func AssetsDirectoryGetsAll(ctx *ctx.Context) ([]AssetsDirectory, error) {
	var lst []AssetsDirectory
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// map查询
func AssetsDirectoryGetsMap(ctx *ctx.Context, where map[string]interface{}) ([]*AssetsDirTree, error) {
	var lst []*AssetsDirTree
	err := DB(ctx).Model(&AssetsDirectory{}).Where(where).Find(&lst).Error

	return lst, err
}

// 增加资产目录
func (a *AssetsDirectory) Add(ctx *ctx.Context) error {
	// 这里写AssetsDirectory的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 删除资产目录
func (a *AssetsDirectory) Del(tx *gorm.DB) error {
	// 这里写AssetsDirectory的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Delete(a).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 更新资产目录
func (a *AssetsDirectory) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写AssetsDirectory的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(a).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 更新资产目录map
func AssetsDirectoryUpdate(tx *gorm.DB, where map[string]interface{}, updateFrom map[string]interface{}) error {
	// 这里写AssetsDirectory的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Model(&AssetsDirectory{}).Where(where).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 根据条件统计个数
func AssetsDirectoryCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetsDirectory{}).Where(where, args...))
}

// 统计目录下资产数量
func AssetsDirCount(ctx *ctx.Context, assetsDirTree *AssetsDirTree) (*AssetsDirTree, error) {
	if assetsDirTree.Children == nil {
		num, err := AssetsCountMap(ctx, map[string]interface{}{"directory_id": assetsDirTree.Id})
		if err != nil {
			return nil, err
		}
		assetsDirTree.Count = num
		return assetsDirTree, err
	}
	var sum int64
	sum = 0
	var err error
	for index, val := range assetsDirTree.Children {
		assetsDirTree.Children[index], err = AssetsDirCount(ctx, assetsDirTree.Children[index])
		if err != nil {
			return nil, err
		}
		sum += val.Count
	}
	assetsDirTree.Count = sum
	return assetsDirTree, err
}

// 获取目录下资产id列表
func AssetsDirIds(ctx *ctx.Context, assetsDirTree *AssetsDirTree, ids []int64) (idsR []int64, err error) {
	if assetsDirTree.Children == nil {
		ids = append(ids, assetsDirTree.Id)
		return ids, nil
	}
	for _, val := range assetsDirTree.Children {
		idsT, err := AssetsDirIds(ctx, val, ids)
		if err != nil {
			return nil, err
		}
		idsR = append(idsR, idsT...)
	}
	idsR = append(idsR, assetsDirTree.Id)
	return idsR, err
}
