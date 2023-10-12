// Package models  资产树
// date : 2023-07-21 09:50
// desc : 资产树
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// AssetTree  资产树。
// 说明:
// 表名:asset_tree
// group: AssetTree
// version:2023-07-21 09:50
type AssetTree struct {
	Id           int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-07-21 09:50
	Status       int64          `gorm:"column:STATUS" json:"status" `                             //type:*int     comment:资产状态    version:2023-07-21 09:58
	Name         string         `gorm:"column:NAME" json:"name" `                                 //type:string   comment:名称        version:2023-07-21 09:50
	ManagementIp string         `gorm:"column:MANAGEMENT_IP" json:"management_ip" `               //type:string   comment:管理IP      version:2023-08-07 09:40
	SerialNumber string         `gorm:"column:SERIAL_NUMBER" json:"serial_number" `               //type:string   comment:序列号      version:2023-08-07 09:40
	PropertyId   int64          `gorm:"column:PROPERTY_ID" json:"property_id" `                   //type:*int     comment:属性ID      version:2023-07-28 10:43
	ParentId     int64          `gorm:"column:PARENT_ID" json:"parent_id" `                       //type:*int     comment:父ID        version:2023-07-21 09:50
	Type         string         `gorm:"column:TYPE" json:"type" `                                 //type:string   comment:类型        version:2023-07-24 09:35
	Remark       string         `gorm:"column:REMARK" json:"remark" `                             //type:string   comment:备注        version:2023-07-21 09:50
	CreatedBy    string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-07-21 09:50
	CreatedAt    int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-07-21 09:50
	UpdatedBy    string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-07-21 09:50
	UpdatedAt    int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-07-21 09:50
	DeletedAt    gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*int       comment:删除时间        version:2023-9-08 16:39
}

// FrontTree  资产树。
type FrontTree struct {
	Id           int64                  `gorm:"-" json:"id" `            //type:*int     comment:主键        version:2023-07-21 09:50
	Status       int64                  `gorm:"-" json:"status" `        //type:*int     comment:资产状态    version:2023-07-21 09:58
	Name         string                 `gorm:"-" json:"name" `          //type:string   comment:名称        version:2023-07-21 09:50
	ManagementIp string                 `gorm:"-" json:"management_ip" ` //type:string   comment:管理IP      version:2023-08-07 09:40
	SerialNumber string                 `gorm:"-" json:"serial_number" ` //type:string   comment:序列号      version:2023-08-07 09:40
	PropertyId   int64                  `gorm:"-" json:"property_id" `   //type:*int     comment:属性ID      version:2023-07-28 10:43
	ParentId     int64                  `gorm:"-" json:"parent_id" `     //type:*int     comment:父ID        version:2023-07-21 09:50
	Type         string                 `gorm:"-" json:"type" `          //type:string   comment:类型        version:2023-07-24 09:35
	Remark       string                 `gorm:"-" json:"remark" `        //type:string   comment:备注        version:2023-07-21 09:50
	Query        map[string]interface{} `gorm:"-" json:"query" `
	Children     []*FrontTree           `gorm:"-" json:"children" `
}

// TableName 表名:asset_tree，资产树。
// 说明:
func (a *AssetTree) TableName() string {
	return "asset_tree"
}

//默认树：设备状态-设备类型-设备厂商
func BuildAssetTree(ctx *ctx.Context, query map[string]interface{}) ([]*FrontTree, error) {

	var tree []*AssetTree
	err := DB(ctx).Select("ID", "NAME", "PARENT_ID, TYPE, PROPERTY_ID, MANAGEMENT_IP, SERIAL_NUMBER").Where(query).Find(&tree).Error
	if err != nil {
		return nil, err
	}

	frontTree := PackagingFrontTree(tree, query["status"].(int64))

	return getTreeRecursive(frontTree, 0), err
}

//部分资产树
func BuildPartAssetTree(ctx *ctx.Context, query map[string]interface{}) ([]*FrontTree, error) {
	var tree []*AssetTree
	err := DB(ctx).Select("ID", "NAME", "PARENT_ID, TYPE, PROPERTY_ID, MANAGEMENT_IP, SERIAL_NUMBER").Where(query).Find(&tree).Error
	if err != nil || len(tree) == 0 {
		return nil, err
	}
	var proTree []*AssetTree
	err = DB(ctx).Select("ID", "NAME", "PARENT_ID, TYPE, PROPERTY_ID, MANAGEMENT_IP, SERIAL_NUMBER").Where("PARENT_ID = ?", tree[0].Id).Find(&proTree).Error
	if err != nil {
		return nil, err
	}
	tree = append(tree, proTree...)
	var devTree []*AssetTree
	for _, val := range proTree {
		err = DB(ctx).Select("ID", "NAME", "PARENT_ID, TYPE, PROPERTY_ID, MANAGEMENT_IP, SERIAL_NUMBER").Where("PARENT_ID = ?", val.Id).Find(&devTree).Error
		if err != nil {
			return nil, err
		}
		tree = append(tree, devTree...)
	}

	frontTree := PackagingFrontTree(tree, query["status"].(int64))

	return getTreeRecursive(frontTree, 0), err
}

//构建树
func getTreeRecursive(list []*FrontTree, parentId int64) []*FrontTree {
	res := make([]*FrontTree, 0)
	for _, v := range list {
		if v.ParentId == parentId {
			v.Children = getTreeRecursive(list, v.Id)
			res = append(res, v)
		}
	}
	return res
}

// 条件查询
func AssetTreeGets(ctx *ctx.Context, query string, limit, offset int) ([]AssetTree, error) {
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

	var lst []AssetTree
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func AssetTreeGetById(ctx *ctx.Context, id int64) (*AssetTree, error) {
	var obj *AssetTree
	err := DB(ctx).Debug().Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按id批量查询
func AssetTreeBatchGetById(ctx *ctx.Context, ids []int64) ([]AssetTree, error) {
	var obj []AssetTree
	err := DB(ctx).Debug().Where("PROPERTY_ID IN ? AND TYPE = ?", ids, "asset").Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按type和assetId查询
func AssetTreeGetByMap(ctx *ctx.Context, m map[string]interface{}) ([]AssetTree, error) {
	var obj []AssetTree
	err := DB(ctx).Debug().Where(m).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func AssetTreeGetsAll(ctx *ctx.Context) ([]AssetTree, error) {
	var lst []AssetTree
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加资产树
func (a *AssetTree) Add(ctx *ctx.Context) error {
	// 这里写AssetTree的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 删除资产树
func (a *AssetTree) Del(tx *gorm.DB) error {
	// 这里写AssetTree的业务逻辑，通过error返回错误

	var tree []*AssetTree
	err := tx.Debug().Where("PARENT_ID = ?", a.Id).Find(&tree).Error
	if err != nil {
		tx.Rollback()
	}

	if len(tree) != 0 {
		for _, val := range tree {
			err = val.Del(tx)
			if err != nil {
				tx.Rollback()
			}
		}

	}
	err = tx.Debug().Delete(&AssetTree{}, a.Id).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 删除资产详情
func AssetTreeBatchDel(tx *gorm.DB, ids []int64) error {
	// 这里写AssetTree的业务逻辑，通过error返回错误

	//删除资产详情
	err := tx.Where("PROPERTY_ID IN ? AND TYPE = ?", ids, "asset").Delete(&AssetTree{}).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 更新资产树
func (a *AssetTree) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写AssetTree的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(a).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func AssetTreeCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetTree{}).Where(where, args...))
}

func (a *AssetTree) AssetCountGet(ctx *ctx.Context) (num int64, err error) {
	var tree []*AssetTree
	var numVal int64
	err = DB(ctx).Where("PARENT_ID = ?", a.Id).Find(&tree).Error
	if tree == nil {
		return 0, err
	} else {
		for _, val := range tree {
			numVal, err = val.AssetCountGet(ctx)
			num += numVal
			if val.Type == "asset" {
				num++
			}
		}
	}
	return num, err
}

//封装返回值
func PackagingFrontTree(tree []*AssetTree, status int64) []*FrontTree {

	var frontTree []*FrontTree
	//封装资产树查询资产详情条件
	for i := 0; i < len(tree); i++ {
		m := make(map[string]interface{})
		m["DEVICE_STATUS"] = status

		if tree[i].Type == "type" {
			m["DEVICE_TYPE"] = tree[i].PropertyId
		}
		if tree[i].Type == "producer" {
			for _, val := range tree {
				if val.Id == tree[i].ParentId {
					m["DEVICE_TYPE"] = val.PropertyId
					break
				}
			}
			m["DEVICE_PRODUCER"] = tree[i].PropertyId
		}
		if tree[i].Type == "asset" {
			var froId int64
			for _, val := range tree {
				if val.Id == tree[i].ParentId {
					m["DEVICE_PRODUCER"] = val.PropertyId
					froId = val.ParentId
					break
				}
			}
			for _, val := range tree {
				if val.Id == froId {
					m["DEVICE_TYPE"] = val.PropertyId
					break
				}
			}

			m["DEVICE_NAME"] = tree[i].Name

			//增加资产ID
			m["ID"] = tree[i].PropertyId
		}

		frontTree = append(frontTree, &FrontTree{
			Id:           tree[i].Id,
			Status:       status,
			Name:         tree[i].Name,
			ManagementIp: tree[i].ManagementIp,
			SerialNumber: tree[i].SerialNumber,
			ParentId:     tree[i].ParentId,
			Type:         tree[i].Type,
			Query:        m,
		})
	}
	return frontTree
}

// 更新资产树
func UpdateTree(ctx *ctx.Context, ids []int64, m map[string]interface{}) error {
	// 这里写AssetTree的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(&AssetTree{}).Where("ID IN ?", ids).Updates(m).Error
}

// 更新资产树
func UpdateTxTree(tx *gorm.DB, where map[string]interface{}, m map[string]interface{}) error {
	// 这里写AssetTree的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Debug().Model(&AssetTree{}).Where(where).Updates(m).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

//判断存在几个子节点
func ChildrenNum(ctx *ctx.Context, parent int64) (num int64, err error) {
	err = DB(ctx).Debug().Model(&AssetTree{}).Where("PARENT_ID = ?", parent).Count(&num).Error
	return num, err
}

//根据目录清除树
func AssetTreeDelById(ctx *ctx.Context, id int64) error {

	flog := true
	for flog {
		childrenNum, err := ChildrenNum(ctx, id)
		if err != nil {
			return err
		}
		if childrenNum > 0 || id == 0 {
			return err
		}
		assetTree, err := AssetTreeGetById(ctx, id)
		if err != nil {
			return err
		}
		id = assetTree.ParentId
		err = DB(ctx).Debug().Where("ID = ?", assetTree.Id).Delete(&AssetTree{}).Error
		if id == 0 {
			return err
		}
	}
	return nil
}
