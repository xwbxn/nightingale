// Package models  资产详情
// date : 2023-07-21 08:45
// desc : 资产详情
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// AssetBasic  资产详情。
// 说明:
// 表名:asset_basic
// group: AssetBasic
// version:2023-07-21 08:45
type AssetBasic struct {
	Id                     int64  `gorm:"column:ID;primaryKey" json:"id" `                                          //type:*int     comment:主键                                     version:2023-07-21 08:45
	DeviceType             int64  `gorm:"column:DEVICE_TYPE" json:"device_type" `                                   //type:*int   comment:设备类型                                 version:2023-07-21 08:45
	ManagementIp           string `gorm:"column:MANAGEMENT_IP" json:"management_ip" example:"管理IP"`                 //type:string   comment:管理IP                                   version:2023-07-21 08:45
	DeviceName             string `gorm:"column:DEVICE_NAME" json:"device_name" example:"设备名称"`                     //type:string   comment:设备名称                                 version:2023-07-21 08:45
	SerialNumber           string `gorm:"column:SERIAL_NUMBER" json:"serial_number" example:"序列号"`                  //type:string   comment:序列号                                   version:2023-07-21 08:45
	DeviceStatus           int64  `gorm:"column:DEVICE_STATUS" json:"device_status" `                               //type:*int     comment:状态(1:待上线,2:已上线,3:下线,4:报废)    version:2023-07-21 08:45
	ManagedState           int64  `gorm:"column:MANAGED_STATE" json:"managed_state" `                               //type:*int     comment:纳管状态                                 version:2023-07-28 14:36
	DeviceProducer         int64  `gorm:"column:DEVICE_PRODUCER" json:"device_producer" `                           //type:*int     comment:厂商                                     version:2023-07-25 16:30
	DeviceModel            int64  `gorm:"column:DEVICE_MODEL" json:"device_model" `                                 //type:string   comment:型号                                     version:2023-07-21 08:45
	Subtype                string `gorm:"column:SUBTYPE" json:"subtype" example:"子类型"`                              //type:string   comment:子类型                                   version:2023-07-21 08:45
	OutlineStructure       string `gorm:"column:OUTLINE_STRUCTURE" json:"outline_structure" example:"外形结构"`         //type:string   comment:外形结构                                 version:2023-07-21 08:45
	Specifications         string `gorm:"column:SPECIFICATIONS" json:"specifications" example:"规格"`                 //type:string   comment:规格                                     version:2023-07-21 08:45
	UNumber                int64  `gorm:"column:U_NUMBER" json:"u_number" `                                         //type:*int     comment:U数                                      version:2023-07-21 08:45
	UseStorage             string `gorm:"column:USE_STORAGE" json:"use_storage" example:"使用存储"`                     //type:string   comment:使用存储                                 version:2023-07-21 08:45
	DeviceManagerOne       string `gorm:"column:DEVICE_MANAGER_ONE" json:"device_manager_one" example:"设备负责人1"`     //type:string   comment:设备负责人1                              version:2023-07-21 08:45
	DeviceManagerTwo       string `gorm:"column:DEVICE_MANAGER_TWO" json:"device_manager_two" example:"设备负责人2"`     //type:string   comment:设备负责人2                              version:2023-07-21 08:45
	BusinessManagerOne     string `gorm:"column:BUSINESS_MANAGER_ONE" json:"business_manager_one" example:"业务负责人1"` //type:string   comment:业务负责人1                              version:2023-07-21 08:45
	BusinessManagerTwo     string `gorm:"column:BUSINESS_MANAGER_TWO" json:"business_manager_two" example:"业务负责人2"` //type:string   comment:业务负责人2                              version:2023-07-21 08:45
	OperatingSystem        string `gorm:"column:OPERATING_SYSTEM" json:"operating_system" example:"操作系统"`           //type:string   comment:操作系统                                 version:2023-07-21 08:45
	Remark                 string `gorm:"column:REMARK" json:"remark" example:"备注"`                                 //type:string   comment:备注                                     version:2023-07-21 08:45
	AffiliatedOrganization int64  `gorm:"column:AFFILIATED_ORGANIZATION" json:"affiliated_organization" `           //type:string   comment:所属组织机构                             version:2023-07-21 08:45
	EquipmentRoom          int64  `gorm:"column:EQUIPMENT_ROOM" json:"equipment_room" `                             //type:string   comment:所在机房                                 version:2023-07-21 08:45
	OwningCabinet          int64  `gorm:"column:OWNING_CABINET" json:"owning_cabinet" `                             //type:string   comment:所在机柜                                 version:2023-07-21 08:45
	Region                 string `gorm:"column:REGION" json:"region" example:"所在区域"`                               //type:string   comment:所在区域                                 version:2023-07-21 08:45
	CabinetLocation        int64  `gorm:"column:CABINET_LOCATION" json:"cabinet_location" `                         //type:*int     comment:机柜位置                                 version:2023-07-21 08:45
	Abreast                int64  `gorm:"column:ABREAST" json:"abreast" `                                           //type:*int     comment:并排放置(0:否,1:是)                      version:2023-07-21 08:45
	LocationDescription    string `gorm:"column:LOCATION_DESCRIPTION" json:"location_description" example:"位置描述"`   //type:string   comment:位置描述                                 version:2023-07-21 08:45
	CreatedBy              string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                 //type:string   comment:创建人                                   version:2023-07-21 08:45
	CreatedAt              int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                 //type:*int     comment:创建时间                                 version:2023-07-21 08:45
	UpdatedBy              string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                 //type:string   comment:更新人                                   version:2023-07-21 08:45
	UpdatedAt              int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                 //type:*int     comment:更新时间                                 version:2023-07-21 08:45
}

type AssetBasicFindVo struct {
	Id             int64 `gorm:"column:ID;primaryKey" json:"id" `                //type:*int     comment:主键                                     version:2023-07-21 08:45
	DeviceType     int64 `gorm:"column:DEVICE_TYPE" json:"device_type" `         //type:*int   comment:设备类型                                 version:2023-07-21 08:45
	DeviceStatus   int64 `gorm:"column:DEVICE_STATUS" json:"device_status" `     //type:*int     comment:状态(1:待上线,2:已上线,3:下线,4:报废)    version:2023-07-21 08:45
	DeviceProducer int64 `gorm:"column:DEVICE_PRODUCER" json:"device_producer" ` //type:*int     comment:厂商                                     version:2023-07-25 16:30
}

type AssetBasicDetailsVo struct {
	Id                     int64  `gorm:"column:ID;primaryKey" json:"id" `        //type:*int     comment:主键                                     version:2023-07-21 08:45
	DeviceType             int64  `gorm:"column:DEVICE_TYPE" json:"device_type" ` //type:*int   comment:设备类型                                 version:2023-07-21 08:45
	DeviceTypeName         string `gorm:"-" json:"device_type_name" `
	ManagementIp           string `gorm:"column:MANAGEMENT_IP" json:"management_ip" example:"管理IP"` //type:string   comment:管理IP                                   version:2023-07-21 08:45
	DeviceName             string `gorm:"column:DEVICE_NAME" json:"device_name" example:"设备名称"`     //type:string   comment:设备名称                                 version:2023-07-21 08:45
	SerialNumber           string `gorm:"column:SERIAL_NUMBER" json:"serial_number" example:"序列号"`  //type:string   comment:序列号                                   version:2023-07-21 08:45
	DeviceStatus           int64  `gorm:"column:DEVICE_STATUS" json:"device_status" `               //type:*int     comment:状态(1:待上线,2:已上线,3:下线,4:报废)    version:2023-07-21 08:45
	ManagedState           int64  `gorm:"column:MANAGED_STATE" json:"managed_state" `               //type:*int     comment:纳管状态                                 version:2023-07-28 14:36
	DeviceProducer         int64  `gorm:"column:DEVICE_PRODUCER" json:"device_producer" `           //type:*int     comment:厂商                                     version:2023-07-25 16:30
	DeviceProducerName     string `gorm:"-" json:"device_producer_name" `
	DeviceModel            int64  `gorm:"column:DEVICE_MODEL" json:"device_model" ` //type:string   comment:型号                                     version:2023-07-21 08:45
	DeviceModelName        string `gorm:"-" json:"device_model_name" `
	Subtype                string `gorm:"column:SUBTYPE" json:"subtype" example:"子类型"`                              //type:string   comment:子类型                                   version:2023-07-21 08:45
	OutlineStructure       string `gorm:"column:OUTLINE_STRUCTURE" json:"outline_structure" example:"外形结构"`         //type:string   comment:外形结构                                 version:2023-07-21 08:45
	Specifications         string `gorm:"column:SPECIFICATIONS" json:"specifications" example:"规格"`                 //type:string   comment:规格                                     version:2023-07-21 08:45
	UNumber                int64  `gorm:"column:U_NUMBER" json:"u_number" `                                         //type:*int     comment:U数                                      version:2023-07-21 08:45
	UseStorage             string `gorm:"column:USE_STORAGE" json:"use_storage" example:"使用存储"`                     //type:string   comment:使用存储                                 version:2023-07-21 08:45
	DeviceManagerOne       string `gorm:"column:DEVICE_MANAGER_ONE" json:"device_manager_one" example:"设备负责人1"`     //type:string   comment:设备负责人1                              version:2023-07-21 08:45
	DeviceManagerTwo       string `gorm:"column:DEVICE_MANAGER_TWO" json:"device_manager_two" example:"设备负责人2"`     //type:string   comment:设备负责人2                              version:2023-07-21 08:45
	BusinessManagerOne     string `gorm:"column:BUSINESS_MANAGER_ONE" json:"business_manager_one" example:"业务负责人1"` //type:string   comment:业务负责人1                              version:2023-07-21 08:45
	BusinessManagerTwo     string `gorm:"column:BUSINESS_MANAGER_TWO" json:"business_manager_two" example:"业务负责人2"` //type:string   comment:业务负责人2                              version:2023-07-21 08:45
	OperatingSystem        string `gorm:"column:OPERATING_SYSTEM" json:"operating_system" example:"操作系统"`           //type:string   comment:操作系统                                 version:2023-07-21 08:45
	Remark                 string `gorm:"column:REMARK" json:"remark" example:"备注"`                                 //type:string   comment:备注                                     version:2023-07-21 08:45
	AffiliatedOrganization int64  `gorm:"column:AFFILIATED_ORGANIZATION" json:"affiliated_organization" `           //type:string   comment:所属组织机构                             version:2023-07-21 08:45
	OrganizationName       string `gorm:"-" json:"organization_name" `
	EquipmentRoom          int64  `gorm:"column:EQUIPMENT_ROOM" json:"equipment_room" ` //type:string   comment:所在机房                                 version:2023-07-21 08:45
	RoomName               string `gorm:"-" json:"room_name" `
	OwningCabinet          int64  `gorm:"column:OWNING_CABINET" json:"owning_cabinet" ` //type:string   comment:所在机柜                                 version:2023-07-21 08:45
	CabinetName            string `gorm:"-" json:"cabinet_name" `
	Region                 string `gorm:"column:REGION" json:"region" example:"所在区域"`                             //type:string   comment:所在区域                                 version:2023-07-21 08:45
	CabinetLocation        int64  `gorm:"column:CABINET_LOCATION" json:"cabinet_location" `                       //type:*int     comment:机柜位置                                 version:2023-07-21 08:45
	Abreast                int64  `gorm:"column:ABREAST" json:"abreast" `                                         //type:*int     comment:并排放置(0:否,1:是)                      version:2023-07-21 08:45
	LocationDescription    string `gorm:"column:LOCATION_DESCRIPTION" json:"location_description" example:"位置描述"` //type:string   comment:位置描述                                 version:2023-07-21 08:45
}

// TableName 表名:asset_basic，资产详情。
// 说明:
func (a *AssetBasic) TableName() string {
	return "asset_basic"
}

// 条件查询
func AssetBasicGets(ctx *ctx.Context, query string, limit, offset int) ([]AssetBasic, error) {
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

	var lst []AssetBasic
	err := session.Find(&lst).Error

	return lst, err
}

// 根据map条件查询
func AssetBasicGetsByMap(ctx *ctx.Context, query map[string]interface{}, limit, offset int) ([]AssetBasicDetailsVo, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("ID")
	}

	session = session.Where(query)

	var lst []AssetBasicDetailsVo
	err := session.Model(&AssetBasic{}).Find(&lst).Error

	return lst, err
}

// 按id查询
func AssetBasicGetById(ctx *ctx.Context, id int64) (*AssetBasic, error) {
	var obj *AssetBasic
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func AssetBasicGetsAll(ctx *ctx.Context) ([]AssetBasic, error) {
	var lst []AssetBasic
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加资产详情
func (a *AssetBasic) Add(ctx *ctx.Context) (int64, error) {
	// 这里写AssetBasic的业务逻辑，通过error返回错误

	//启动事务
	tx := DB(ctx).Begin()
	// 实际向库中写入
	err := tx.Create(&a).Error
	if err != nil {
		tx.Rollback()
	}
	//查询设备类型名称
	deviceType, _ := DeviceTypeGetById(ctx, a.DeviceType)
	//插入资产树记录
	var menuTree *AssetTree
	var wTree *AssetTree
	//查询该状态下设备类型是否存在
	err = tx.Where("NAME = ? AND STATUS = ? AND PARENT_ID = ?", deviceType.Name, a.DeviceStatus, int64(0)).Find(&menuTree).Error
	if err != nil {
		tx.Rollback()
	}
	//设备类型不存在，插入设备类型、设备厂商、设备名称
	if menuTree.Id == 0 {
		wTree = BuildObj(deviceType.Name, a.CreatedBy, "type", 0, a.DeviceStatus, a.DeviceType)
		err = tx.Create(wTree).Error
		if err != nil {
			tx.Rollback()
		}
		deviceProducer, _ := DeviceProducerGetById(ctx, a.DeviceProducer)
		wTree = BuildObj(deviceProducer.Alias, a.CreatedBy, "producer", wTree.Id, a.DeviceStatus, a.DeviceProducer)
		err = tx.Create(wTree).Error
		if err != nil {
			tx.Rollback()
		}
		wTree = BuildObj(a.DeviceName, a.CreatedBy, "asset", wTree.Id, a.DeviceStatus, a.Id)
		err = tx.Create(wTree).Error
		if err != nil {
			tx.Rollback()
		}
	} else {
		//设备类型存在，判断设备厂商是否存在
		var menuTree1 *AssetTree
		deviceProducer, _ := DeviceProducerGetById(ctx, a.DeviceProducer)
		err = tx.Where("NAME = ? AND STATUS = ? AND PARENT_ID = ?", deviceProducer.Alias, menuTree.Status, menuTree.Id).Find(&menuTree1).Error

		//设备厂商不存在，插入设备厂商、设备名称
		if menuTree1.Id == 0 {
			wTree = BuildObj(deviceProducer.Alias, a.CreatedBy, "producer", menuTree.Id, a.DeviceStatus, a.DeviceProducer)
			err = tx.Create(wTree).Error
			if err != nil {
				tx.Rollback()
			}
			wTree = BuildObj(a.DeviceName, a.CreatedBy, "asset", wTree.Id, a.DeviceStatus, a.Id)
			err = tx.Create(wTree).Error
			if err != nil {
				tx.Rollback()
			}
		} else {
			//设备厂商存在，插入设备名称
			err = tx.Create(BuildObj(a.DeviceName, a.CreatedBy, "asset", menuTree1.Id, a.DeviceStatus, a.Id)).Error
			if err != nil {
				tx.Rollback()
			}
		}
	}

	tx.Commit()
	return a.Id, err
}

//构建资产树记录
func BuildObj(name, craetedBy, deviceType string, parentId, status, propertyId int64) *AssetTree {
	var assetTree AssetTree
	assetTree.Name = name
	assetTree.ParentId = parentId
	assetTree.Status = status
	assetTree.Type = deviceType
	assetTree.PropertyId = propertyId
	assetTree.CreatedBy = craetedBy
	return &assetTree
}

// 删除资产详情
func (a *AssetBasic) Del(ctx *ctx.Context) error {
	// 这里写AssetBasic的业务逻辑，通过error返回错误

	//删除资产树
	var assetTree AssetTree
	err := DB(ctx).Where("NAME = ?", a.DeviceName).Find(&assetTree).Error
	if err != nil {
		return err
	}
	err = assetTree.Del(ctx)
	if err != nil {
		return err
	}

	// 实际向库中写入
	return DB(ctx).Delete(a).Error
}

// 更新资产详情
func (a *AssetBasic) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写AssetBasic的业务逻辑，通过error返回错误
	var oldAssetBasic *AssetBasic
	var err error
	//更新资产树
	oldAssetBasic, err = AssetBasicGetById(ctx, a.Id)
	if err != nil {
		return err
	}
	if oldAssetBasic.DeviceName != a.DeviceName {
		DB(ctx).Model(&AssetTree{}).Where("NAME = ?", oldAssetBasic.DeviceName).Update("NAME", a.DeviceName)
	}
	if oldAssetBasic.DeviceProducer != a.DeviceProducer {
		DB(ctx).Model(&AssetTree{}).Where("NAME = ?", oldAssetBasic.DeviceProducer).Update("NAME", a.DeviceProducer)
	}
	if oldAssetBasic.DeviceType != a.DeviceType {
		DB(ctx).Model(&AssetTree{}).Where("NAME = ?", oldAssetBasic.DeviceType).Update("NAME", a.DeviceType)
	}

	// 实际向库中写入
	return DB(ctx).Model(a).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func AssetBasicCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetBasic{}).Where(where, args...))
}

// 根据传入条件统计个数
func AssetCountByMap(ctx *ctx.Context, where map[string]interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetBasic{}).Where(where))
}
