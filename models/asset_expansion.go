// Package models  资产扩展
// date : 2023-07-23 09:04
// desc : 资产扩展
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// AssetExpansion  资产扩展。
// 说明:
// 表名:asset_expansion
// group: AssetExpansion
// version:2023-07-23 09:04
type AssetExpansion struct {
	Id               int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键                                          version:2023-07-23 09:04
	AssetId          int64          `gorm:"column:ASSET_ID" json:"asset_id" `                         //type:*int     comment:资产ID                                        version:2023-07-28 16:17
	ConfigCategory   string         `gorm:"column:CONFIG_CATEGORY" json:"config_category" `           //type:string   comment:配置类别(1:基本信息,2:硬件配置,3:网络配置)    version:2023-08-05 14:41
	PropertyCategory string         `gorm:"column:PROPERTY_CATEGORY" json:"property_category" `       //type:string   comment:属性类别                                      version:2023-07-23 09:04
	GroupId          string         `gorm:"column:GROUP_ID" json:"group_id" `                         //type:string   comment:分组ID                                        version:2023-07-23 09:04
	PropertyNameCn   string         `gorm:"column:PROPERTY_NAME_CN" json:"property_name_cn" `         //type:string   comment:属性名称                                      version:2023-07-28 17:32
	PropertyName     string         `gorm:"column:PROPERTY_NAME" json:"property_name" `               //type:string   comment:英文名称                                      version:2023-07-28 17:32
	PropertyValue    string         `gorm:"column:PROPERTY_VALUE" json:"property_value" `             //type:string   comment:属性值                                        version:2023-07-23 09:04
	AssociatedTable  string         `gorm:"column:ASSOCIATED_TABLE" json:"associated_table" `         //type:string   comment:关联表名                                      version:2023-07-23 09:04
	CreatedBy        string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人                                        version:2023-07-23 09:04
	CreatedAt        int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间                                      version:2023-07-23 09:04
	UpdatedBy        string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人                                        version:2023-07-23 09:04
	UpdatedAt        int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间                                      version:2023-07-23 09:04
	DeletedAt        gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*int       comment:删除时间        version:2023-9-08 16:39
}

type AssetNetWork struct {
	ManagementIp     string `gorm:"-" json:"management_ip" cn:"设备IP"`  //type:string   comment:管理IP                                   version:2023-07-21 08:45
	SerialNumber     string `gorm:"-" json:"serial_number" cn:"设备序列号"` //type:string   comment:序列号                                   version:2023-07-21 08:45
	DeviceName       string `gorm:"-" json:"device_name" cn:"设备名称"`    //type:string   comment:设备名称                                 version:2023-07-21 08:45
	Type             int64  `gorm:"-" json:"type" cn:"类型" validate:"omitempty,oneof=0 1" source:"type=option,value=[带外IP;生产IP]"`
	IP               string `gorm:"-" json:"ip" cn:"IP" `
	SubnetMask       string `gorm:"-" json:"subnet_mask" cn:"子网掩码" prop:"ext_mask"`
	Gateway          string `gorm:"-" json:"gateway" cn:"网关" prop:"ext_gateway"`
	MacAddress       string `gorm:"-" json:"mac_address" cn:"MAC地址" prop:"ext_mac"`
	SwitchIp         string `gorm:"-" json:"switch_ip" cn:"交换机IP" prop:"ext_switch_ip"`
	SwitchName       string `gorm:"-" json:"switch_name" cn:"交换机名称" prop:"ext_switch_name"`
	SwitchPort       string `gorm:"-" json:"switch_port" cn:"交换机端口" prop:"ext_corres_port"`
	SwitchMacAddress string `gorm:"-" json:"switch_mac_address" cn:"交换机MAC地址" prop:"ext_switch_mac"`
	DisName          string `gorm:"-" json:"dis_name" cn:"配线架名称" prop:"ext_distr_frame"`  //type:string   comment:配线架名称    version:2023-07-16 10:14
	DisPort          string `gorm:"-" json:"dis_port" cn:"配线架端口" prop:"ext_docking_port"` //type:string   comment:配线架名称    version:2023-07-16 10:14
	OS               string `gorm:"-" json:"os" cn:"操作系统" prop:"ext_operate_system"`
	ConnectionName   string `gorm:"-" json:"connection_name" cn:"连接用户名" prop:"ext_username"`
	ConnectionPwd    string `gorm:"-" json:"connection_pwd" cn:"连接密码" prop:"ext_password"`
	ConnectionMode   string `gorm:"-" json:"connection_mode" cn:"连接方式" prop:"ext_link_method"`
	ConnectionPort   string `gorm:"-" json:"connection_port" cn:"连接端口" prop:"ext_link_port"`
	RemoteDesktop    string `gorm:"-" json:"Remote_desktop" cn:"远程桌面" prop:"ext_remote_desktop"`
	RemotePort       string `gorm:"-" json:"Remote_port" cn:"远程端口" prop:"ext_port"`
}

// TableName 表名:asset_expansion，资产扩展。
// 说明:
func (a *AssetExpansion) TableName() string {
	return "asset_expansion"
}

// 条件查询
func AssetExpansionGets(ctx *ctx.Context, query string, limit, offset int) ([]AssetExpansion, error) {
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

	var lst []AssetExpansion
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func AssetExpansionGetById(ctx *ctx.Context, id int64) (*AssetExpansion, error) {
	var obj *AssetExpansion
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按map查询
func AssetExpansionGetByMap(ctx *ctx.Context, query map[string]interface{}) ([]AssetExpansion, error) {
	var lst []AssetExpansion
	err := DB(ctx).Where(query).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	return lst, nil
}

// 查询网络配置
func AssetNetConfigGetByAssetId(ctx *ctx.Context, assetIds []int64) ([]AssetExpansion, error) {
	var lst []AssetExpansion
	err := DB(ctx).Debug().Where("ASSET_ID IN ? AND config_category = ? AND (property_category = ? OR property_category = ?)", assetIds, "form-netconfig", "group_management", "group_ip").Find(&lst).Error
	if err != nil {
		return nil, err
	}

	return lst, nil
}

// 按map查询GroupId
func GroupIdGetByMap(ctx *ctx.Context, query map[string]interface{}) ([]string, error) {
	var lst []string
	err := DB(ctx).Model(&AssetExpansion{}).Select("GROUP_ID").Where(query).Find(&lst).Error
	if err != nil {
		return nil, err
	}

	return lst, nil
}

// 查询所有
func AssetExpansionGetsAll(ctx *ctx.Context) ([]AssetExpansion, error) {
	var lst []AssetExpansion
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加资产扩展
func (a *AssetExpansion) Add(ctx *ctx.Context) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 增加资产扩展
func (a *AssetExpansion) AddTx(tx *gorm.DB) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	// logger.Debug(a)
	err := tx.Create(&a).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 批量增加资产扩展
func AssetExpansionTxBatchAdd(tx *gorm.DB, a []AssetExpansion) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Create(&a).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 批量增加资产扩展
func AssetExpansionBatchAdd(ctx *ctx.Context, a []AssetExpansion) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(&a).Error
}

// 删除资产扩展
func (a *AssetExpansion) Del(ctx *ctx.Context) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(a).Error
}

// 删除资产扩展Tx
func (a *AssetExpansion) TxDel(tx *gorm.DB) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Delete(a).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 批量删除资产扩展
func AssetExpansionBatchDel(tx *gorm.DB, ids []int64) error {
	//删除资产扩展
	err := tx.Where("ASSET_ID IN ?", ids).Delete(&AssetExpansion{}).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 根据map删除资产扩展Tx
func MapTxDel(tx *gorm.DB, m map[string]interface{}) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Where(m).Delete(&AssetExpansion{}).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 根据groupId删除资产扩展
func Del(ctx *ctx.Context, groupId string) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Where("GROUP_ID = ?", groupId).Delete(AssetExpansion{}).Error
}

// 更新资产扩展
func (a *AssetExpansion) Update(tx *gorm.DB, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Model(a).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 更新资产扩展
func UpdateAssetExpansionMap(tx *gorm.DB, where map[string]interface{}, m map[string]interface{}) error {
	// 这里写AssetExpansion的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Model(&AssetExpansion{}).Where(where).Updates(m).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 根据条件统计个数
func AssetExpansionCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetExpansion{}).Where(where, args...))
}

//以group处理新数据
func UpdateAssetExpansionGroup(ctx *ctx.Context, m map[string]interface{}, f []AssetExpansion, name string) error {
	//查询具体的group
	groupIds, err := GroupIdGetByMap(ctx, m)
	if err != nil {
		return err
	}

	mdata := make(map[string][]AssetExpansion)
	mNewdata := make(map[string][]AssetExpansion)

	for index, val := range f {
		//将新数据单独存储（id=0）
		if val.Id == 0 {
			f[index].CreatedBy = name
			assetExpansions := mNewdata[val.GroupId]
			assetExpansions = append(assetExpansions, f[index])
			mNewdata[val.GroupId] = assetExpansions
		} else {
			//将需更新或删除的数据存储
			assetExpansions := mdata[val.GroupId]
			assetExpansions = append(assetExpansions, val)
			mdata[val.GroupId] = assetExpansions
		}

	}

	//启动事务
	tx := DB(ctx).Begin()

	for _, groupId := range groupIds {
		//判断旧数据是否依然存在
		assetExpansions, ok := mdata[groupId]
		if ok {
			//旧数据存在，更新数据
			err = UpdateAssetExpansion(tx, ctx, assetExpansions, name)
			if err != nil {
				tx.Rollback()
			}
		} else {
			//旧数据不存在，删除旧数据
			err = tx.Debug().Model(AssetExpansion{}).Where("GROUP_ID = ?", groupId).Delete(assetExpansions).Error
			if err != nil {
				tx.Rollback()
			}
		}
	}

	//将新数据插入表中
	for groupId := range mNewdata {

		assetExp := mNewdata[groupId]
		err = tx.Create(&assetExp).Error
		if err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
	return err
}

//以属性处理数据
func UpdateAssetExpansion(tx *gorm.DB, ctx *ctx.Context, f []AssetExpansion, name string) error {

	//通过groupId查询一组扩展数据
	old, err := AssetExpansionGetByMap(ctx, map[string]interface{}{"GROUP_ID": f[0].GroupId})
	if err != nil {
		return err
	}

	for _, oldVal := range old {
		var tag = false
		for index := range f {
			//新旧数据id匹配，更新数据
			if f[index].Id == oldVal.Id {
				f[index].UpdatedBy = name
				err = oldVal.Update(tx, f[index], "*")
				if err != nil {
					tx.Rollback()
				}
				tag = true
				break
			}
		}
		//新数据不存在，旧数据存在，删除该记录
		if !tag {
			err = tx.Delete(&oldVal).Error
			if err != nil {
				tx.Rollback()
			}
		}
	}
	//新数据存在，原数据不存在，新增该记录
	for index := range f {
		if f[index].Id == 0 {
			f[index].CreatedBy = name
			assetExp := f[index]
			err = tx.Create(&assetExp).Error
			if err != nil {
				tx.Rollback()
			}
		}
	}
	return err
}
