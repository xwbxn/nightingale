// Package models  资产详情
// date : 2023-07-21 08:45
// desc : 资产详情
package models

import (
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/toolkits/pkg/logger"
	"gorm.io/gorm"
)

// AssetBasic  资产详情。
// 说明:
// 表名:asset_basic
// group: AssetBasic
// version:2023-07-21 08:45
type AssetBasic struct {
	Id                     int64          `gorm:"column:ID;primaryKey" json:"id" `                                          //type:*int     comment:主键                                     version:2023-07-21 08:45
	DeviceType             int64          `gorm:"column:DEVICE_TYPE" json:"device_type" `                                   //type:*int   comment:设备类型                                 version:2023-07-21 08:45
	ManagementIp           string         `gorm:"column:MANAGEMENT_IP" json:"management_ip" example:"管理IP"`                 //type:string   comment:管理IP                                   version:2023-07-21 08:45
	DeviceName             string         `gorm:"column:DEVICE_NAME" json:"device_name" example:"设备名称"`                     //type:string   comment:设备名称                                 version:2023-07-21 08:45
	SerialNumber           string         `gorm:"column:SERIAL_NUMBER" json:"serial_number" example:"序列号"`                  //type:string   comment:序列号                                   version:2023-07-21 08:45
	DeviceStatus           int64          `gorm:"column:DEVICE_STATUS" json:"device_status" `                               //type:*int     comment:状态(1:待上线,2:已上线,3:下线,4:报废)    version:2023-07-21 08:45
	ManagedState           int64          `gorm:"column:MANAGED_STATE" json:"managed_state" `                               //type:*int     comment:纳管状态                                 version:2023-07-28 14:36
	DeviceProducer         int64          `gorm:"column:DEVICE_PRODUCER" json:"device_producer" `                           //type:*int     comment:厂商                                     version:2023-07-25 16:30
	DeviceModel            int64          `gorm:"column:DEVICE_MODEL" json:"device_model" `                                 //type:string   comment:型号                                     version:2023-07-21 08:45
	Subtype                string         `gorm:"column:SUBTYPE" json:"subtype" example:"子类型"`                              //type:string   comment:子类型                                   version:2023-07-21 08:45
	OutlineStructure       string         `gorm:"column:OUTLINE_STRUCTURE" json:"outline_structure" example:"外形结构"`         //type:string   comment:外形结构                                 version:2023-07-21 08:45
	Specifications         string         `gorm:"column:SPECIFICATIONS" json:"specifications" example:"规格"`                 //type:string   comment:规格                                     version:2023-07-21 08:45
	UNumber                int64          `gorm:"column:U_NUMBER" json:"u_number" `                                         //type:*int     comment:U数                                      version:2023-07-21 08:45
	UseStorage             string         `gorm:"column:USE_STORAGE" json:"use_storage" example:"使用存储"`                     //type:string   comment:使用存储                                 version:2023-07-21 08:45
	DatacenterId           int64          `gorm:"column:DATACENTER_ID" json:"datacenter_id" `                               //type:*int     comment:数据中心                                 version:2023-08-15 09:23
	RelatedService         string         `gorm:"column:RELATED_SERVICE" json:"related_service" `                           //type:string   comment:关联业务                                 version:2023-08-15 09:23
	ServicePath            string         `gorm:"column:SERVICE_PATH" json:"service_path" `                                 //type:string   comment:业务路径                                 version:2023-08-15 09:23
	DeviceManagerOne       string         `gorm:"column:DEVICE_MANAGER_ONE" json:"device_manager_one" example:"设备负责人1"`     //type:string   comment:设备负责人1                              version:2023-07-21 08:45
	DeviceManagerTwo       string         `gorm:"column:DEVICE_MANAGER_TWO" json:"device_manager_two" example:"设备负责人2"`     //type:string   comment:设备负责人2                              version:2023-07-21 08:45
	BusinessManagerOne     string         `gorm:"column:BUSINESS_MANAGER_ONE" json:"business_manager_one" example:"业务负责人1"` //type:string   comment:业务负责人1                              version:2023-07-21 08:45
	BusinessManagerTwo     string         `gorm:"column:BUSINESS_MANAGER_TWO" json:"business_manager_two" example:"业务负责人2"` //type:string   comment:业务负责人2                              version:2023-07-21 08:45
	OperatingSystem        string         `gorm:"column:OPERATING_SYSTEM" json:"operating_system" example:"操作系统"`           //type:string   comment:操作系统                                 version:2023-07-21 08:45
	Remark                 string         `gorm:"column:REMARK" json:"remark" example:"备注"`                                 //type:string   comment:备注                                     version:2023-07-21 08:45
	AffiliatedOrganization int64          `gorm:"column:AFFILIATED_ORGANIZATION" json:"affiliated_organization" `           //type:string   comment:所属组织机构                             version:2023-07-21 08:45
	EquipmentRoom          int64          `gorm:"column:EQUIPMENT_ROOM" json:"equipment_room" `                             //type:string   comment:所在机房                                 version:2023-07-21 08:45
	OwningCabinet          int64          `gorm:"column:OWNING_CABINET" json:"owning_cabinet" `                             //type:string   comment:所在机柜                                 version:2023-07-21 08:45
	Region                 string         `gorm:"column:REGION" json:"region" example:"所在区域"`                               //type:string   comment:所在区域                                 version:2023-07-21 08:45
	CabinetLocation        int64          `gorm:"column:CABINET_LOCATION" json:"cabinet_location" `                         //type:*int     comment:机柜位置                                 version:2023-07-21 08:45
	Abreast                int64          `gorm:"column:ABREAST" json:"abreast" `                                           //type:*int     comment:并排放置(0:否,1:是)                      version:2023-07-21 08:45
	LocationDescription    string         `gorm:"column:LOCATION_DESCRIPTION" json:"location_description" example:"位置描述"`   //type:string   comment:位置描述                                 version:2023-07-21 08:45
	ExtensionTest          string         `gorm:"column:EXTENSION_TEST" json:"extension_test" `                             //type:string   comment:扩展测试                                 version:2023-08-18 09:14
	CreatedBy              string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                 //type:string   comment:创建人                                   version:2023-07-21 08:45
	CreatedAt              int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                 //type:*int     comment:创建时间                                 version:2023-07-21 08:45
	UpdatedBy              string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                 //type:string   comment:更新人                                   version:2023-07-21 08:45
	UpdatedAt              int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                 //type:*int     comment:更新时间                                 version:2023-07-21 08:45
	DeletedAt              gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`                 //type:*int       comment:删除时间        version:2023-9-08 16:39
}
type AssetBasicExpansionVo struct {
	Id                     int64            `gorm:"column:ID;primaryKey" json:"id" `                                          //type:*int     comment:主键                                     version:2023-07-21 08:45
	DeviceType             int64            `gorm:"column:DEVICE_TYPE" json:"device_type" `                                   //type:*int   comment:设备类型                                 version:2023-07-21 08:45
	ManagementIp           string           `gorm:"column:MANAGEMENT_IP" json:"management_ip" example:"管理IP"`                 //type:string   comment:管理IP                                   version:2023-07-21 08:45
	DeviceName             string           `gorm:"column:DEVICE_NAME" json:"device_name" example:"设备名称"`                     //type:string   comment:设备名称                                 version:2023-07-21 08:45
	SerialNumber           string           `gorm:"column:SERIAL_NUMBER" json:"serial_number" example:"序列号"`                  //type:string   comment:序列号                                   version:2023-07-21 08:45
	DeviceStatus           int64            `gorm:"column:DEVICE_STATUS" json:"device_status" `                               //type:*int     comment:状态(1:待上线,2:已上线,3:下线,4:报废)    version:2023-07-21 08:45
	ManagedState           int64            `gorm:"column:MANAGED_STATE" json:"managed_state" `                               //type:*int     comment:纳管状态                                 version:2023-07-28 14:36
	DeviceProducer         int64            `gorm:"column:DEVICE_PRODUCER" json:"device_producer" `                           //type:*int     comment:厂商                                     version:2023-07-25 16:30
	DeviceModel            int64            `gorm:"column:DEVICE_MODEL" json:"device_model" `                                 //type:string   comment:型号                                     version:2023-07-21 08:45
	Subtype                string           `gorm:"column:SUBTYPE" json:"subtype" example:"子类型"`                              //type:string   comment:子类型                                   version:2023-07-21 08:45
	OutlineStructure       string           `gorm:"column:OUTLINE_STRUCTURE" json:"outline_structure" example:"外形结构"`         //type:string   comment:外形结构                                 version:2023-07-21 08:45
	Specifications         string           `gorm:"column:SPECIFICATIONS" json:"specifications" example:"规格"`                 //type:string   comment:规格                                     version:2023-07-21 08:45
	UNumber                int64            `gorm:"column:U_NUMBER" json:"u_number" `                                         //type:*int     comment:U数                                      version:2023-07-21 08:45
	UseStorage             string           `gorm:"column:USE_STORAGE" json:"use_storage" example:"使用存储"`                     //type:string   comment:使用存储                                 version:2023-07-21 08:45
	DatacenterId           int64            `gorm:"column:DATACENTER_ID" json:"datacenter_id" `                               //type:*int     comment:数据中心                                 version:2023-08-15 09:23
	RelatedService         string           `gorm:"column:RELATED_SERVICE" json:"related_service" `                           //type:string   comment:关联业务                                 version:2023-08-15 09:23
	ServicePath            string           `gorm:"column:SERVICE_PATH" json:"service_path" `                                 //type:string   comment:业务路径                                 version:2023-08-15 09:23
	DeviceManagerOne       string           `gorm:"column:DEVICE_MANAGER_ONE" json:"device_manager_one" example:"设备负责人1"`     //type:string   comment:设备负责人1                              version:2023-07-21 08:45
	DeviceManagerTwo       string           `gorm:"column:DEVICE_MANAGER_TWO" json:"device_manager_two" example:"设备负责人2"`     //type:string   comment:设备负责人2                              version:2023-07-21 08:45
	BusinessManagerOne     string           `gorm:"column:BUSINESS_MANAGER_ONE" json:"business_manager_one" example:"业务负责人1"` //type:string   comment:业务负责人1                              version:2023-07-21 08:45
	BusinessManagerTwo     string           `gorm:"column:BUSINESS_MANAGER_TWO" json:"business_manager_two" example:"业务负责人2"` //type:string   comment:业务负责人2                              version:2023-07-21 08:45
	OperatingSystem        string           `gorm:"column:OPERATING_SYSTEM" json:"operating_system" example:"操作系统"`           //type:string   comment:操作系统                                 version:2023-07-21 08:45
	Remark                 string           `gorm:"column:REMARK" json:"remark" example:"备注"`                                 //type:string   comment:备注                                     version:2023-07-21 08:45
	AffiliatedOrganization int64            `gorm:"column:AFFILIATED_ORGANIZATION" json:"affiliated_organization" `           //type:string   comment:所属组织机构                             version:2023-07-21 08:45
	EquipmentRoom          int64            `gorm:"column:EQUIPMENT_ROOM" json:"equipment_room" `                             //type:string   comment:所在机房                                 version:2023-07-21 08:45
	OwningCabinet          int64            `gorm:"column:OWNING_CABINET" json:"owning_cabinet" `                             //type:string   comment:所在机柜                                 version:2023-07-21 08:45
	Region                 string           `gorm:"column:REGION" json:"region" example:"所在区域"`                               //type:string   comment:所在区域                                 version:2023-07-21 08:45
	CabinetLocation        int64            `gorm:"column:CABINET_LOCATION" json:"cabinet_location" `                         //type:*int     comment:机柜位置                                 version:2023-07-21 08:45
	Abreast                int64            `gorm:"column:ABREAST" json:"abreast" `                                           //type:*int     comment:并排放置(0:否,1:是)                      version:2023-07-21 08:45
	LocationDescription    string           `gorm:"column:LOCATION_DESCRIPTION" json:"location_description" example:"位置描述"`   //type:string   comment:位置描述                                 version:2023-07-21 08:45
	ExtensionTest          string           `gorm:"column:EXTENSION_TEST" json:"extension_test" `                             //type:string   comment:扩展测试                                 version:2023-08-18 09:14
	CreatedBy              string           `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                 //type:string   comment:创建人                                   version:2023-07-21 08:45
	CreatedAt              int64            `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                 //type:*int     comment:创建时间                                 version:2023-07-21 08:45
	UpdatedBy              string           `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                 //type:string   comment:更新人                                   version:2023-07-21 08:45
	UpdatedAt              int64            `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                 //type:*int     comment:更新时间                                 version:2023-07-21 08:45
	DeletedAt              gorm.DeletedAt   `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`                 //type:*int       comment:删除时间        version:2023-9-08 16:39
	BasicExpansion         []AssetExpansion `gorm:"-" json:"basic_expansion" swaggerignore:"true"`
}

type AssetBasicFindVo struct {
	Id             int64  `gorm:"column:ID;primaryKey" json:"id" `                //type:*int     comment:主键                                     version:2023-07-21 08:45
	DeviceType     int64  `gorm:"column:DEVICE_TYPE" json:"device_type" `         //type:*int   comment:设备类型                                 version:2023-07-21 08:45
	DeviceStatus   int64  `gorm:"column:DEVICE_STATUS" json:"device_status" `     //type:*int     comment:状态(1:待上线,2:已上线,3:下线,4:报废)    version:2023-07-21 08:45
	DeviceProducer int64  `gorm:"column:DEVICE_PRODUCER" json:"device_producer" ` //type:*int     comment:厂商                                     version:2023-07-25 16:30
	FinishAt       int64  `gorm:"column:FINISH_AT" json:"finish_at" `             //type:*int     comment:结束日期    version:2023-07-30 10:01
	DeviceManager  string `gorm:"-" json:"device_manager_one" `                   //type:string   comment:设备负责人1                              version:2023-07-21 08:45
	OwningCabinet  int64  `gorm:"column:OWNING_CABINET" json:"owning_cabinet" `   //type:string   comment:所在机柜                                 version:2023-07-21 08:45

}

type AssetBasicStartFindVo struct {
	Id             int64 `gorm:"column:ID;primaryKey" json:"id" `                //type:*int     comment:主键                                     version:2023-07-21 08:45
	DeviceType     int64 `gorm:"column:DEVICE_TYPE" json:"device_type" `         //type:*int   comment:设备类型                                 version:2023-07-21 08:45
	DeviceStatus   int64 `gorm:"column:DEVICE_STATUS" json:"device_status" `     //type:*int     comment:状态(1:待上线,2:已上线,3:下线,4:报废)    version:2023-07-21 08:45
	DeviceProducer int64 `gorm:"column:DEVICE_PRODUCER" json:"device_producer" ` //type:*int     comment:厂商                                     version:2023-07-25 16:30
}

type AssetBasicDetailsVo struct {
	Id                     int64            `gorm:"column:ID;primaryKey" json:"id" `        //type:*int     comment:主键                                     version:2023-07-21 08:45
	DeviceType             int64            `gorm:"column:DEVICE_TYPE" json:"device_type" ` //type:*int   comment:设备类型                                 version:2023-07-21 08:45
	DeviceTypeName         string           `gorm:"-" json:"device_type_name" `
	ManagementIp           string           `gorm:"column:MANAGEMENT_IP" json:"management_ip" example:"管理IP"` //type:string   comment:管理IP                                   version:2023-07-21 08:45
	DeviceName             string           `gorm:"column:DEVICE_NAME" json:"device_name" example:"设备名称"`     //type:string   comment:设备名称                                 version:2023-07-21 08:45
	SerialNumber           string           `gorm:"column:SERIAL_NUMBER" json:"serial_number" example:"序列号"`  //type:string   comment:序列号                                   version:2023-07-21 08:45
	DeviceStatus           int64            `gorm:"column:DEVICE_STATUS" json:"device_status" `               //type:*int     comment:状态(1:待上线,2:已上线,3:下线,4:报废)    version:2023-07-21 08:45
	ManagedState           int64            `gorm:"column:MANAGED_STATE" json:"managed_state" `               //type:*int     comment:纳管状态                                 version:2023-07-28 14:36
	DeviceProducer         int64            `gorm:"column:DEVICE_PRODUCER" json:"device_producer" `           //type:*int     comment:厂商                                     version:2023-07-25 16:30
	DeviceProducerName     string           `gorm:"-" json:"device_producer_name" `
	DeviceModel            int64            `gorm:"column:DEVICE_MODEL" json:"device_model" ` //type:string   comment:型号                                     version:2023-07-21 08:45
	DeviceModelName        string           `gorm:"-" json:"device_model_name" `
	Subtype                string           `gorm:"column:SUBTYPE" json:"subtype" example:"子类型"`                              //type:string   comment:子类型                                   version:2023-07-21 08:45
	OutlineStructure       string           `gorm:"column:OUTLINE_STRUCTURE" json:"outline_structure" example:"外形结构"`         //type:string   comment:外形结构                                 version:2023-07-21 08:45
	Specifications         string           `gorm:"column:SPECIFICATIONS" json:"specifications" example:"规格"`                 //type:string   comment:规格                                     version:2023-07-21 08:45
	UNumber                int64            `gorm:"column:U_NUMBER" json:"u_number" `                                         //type:*int     comment:U数                                      version:2023-07-21 08:45
	UseStorage             string           `gorm:"column:USE_STORAGE" json:"use_storage" example:"使用存储"`                     //type:string   comment:使用存储                                 version:2023-07-21 08:45
	DeviceManagerOne       string           `gorm:"column:DEVICE_MANAGER_ONE" json:"device_manager_one" example:"设备负责人1"`     //type:string   comment:设备负责人1                              version:2023-07-21 08:45
	DeviceManagerTwo       string           `gorm:"column:DEVICE_MANAGER_TWO" json:"device_manager_two" example:"设备负责人2"`     //type:string   comment:设备负责人2                              version:2023-07-21 08:45
	BusinessManagerOne     string           `gorm:"column:BUSINESS_MANAGER_ONE" json:"business_manager_one" example:"业务负责人1"` //type:string   comment:业务负责人1                              version:2023-07-21 08:45
	BusinessManagerTwo     string           `gorm:"column:BUSINESS_MANAGER_TWO" json:"business_manager_two" example:"业务负责人2"` //type:string   comment:业务负责人2                              version:2023-07-21 08:45
	OperatingSystem        string           `gorm:"column:OPERATING_SYSTEM" json:"operating_system" example:"操作系统"`           //type:string   comment:操作系统                                 version:2023-07-21 08:45
	Remark                 string           `gorm:"column:REMARK" json:"remark" example:"备注"`                                 //type:string   comment:备注                                     version:2023-07-21 08:45
	AffiliatedOrganization int64            `gorm:"column:AFFILIATED_ORGANIZATION" json:"affiliated_organization" `           //type:string   comment:所属组织机构                             version:2023-07-21 08:45
	OrganizationName       string           `gorm:"-" json:"organization_name" `
	DatacenterId           int64            `gorm:"column:DATACENTER_ID" json:"datacenter_id" cn:"数据中心" ` //type:*int     comment:数据中心                                 version:2023-08-15 09:23
	EquipmentRoom          int64            `gorm:"column:EQUIPMENT_ROOM" json:"equipment_room" `         //type:string   comment:所在机房                                 version:2023-07-21 08:45
	RoomName               string           `gorm:"-" json:"room_name" `
	OwningCabinet          int64            `gorm:"column:OWNING_CABINET" json:"owning_cabinet" ` //type:string   comment:所在机柜                                 version:2023-07-21 08:45
	CabinetName            string           `gorm:"-" json:"cabinet_name" `
	Region                 string           `gorm:"column:REGION" json:"region" example:"所在区域"`                             //type:string   comment:所在区域                                 version:2023-07-21 08:45
	CabinetLocation        int64            `gorm:"column:CABINET_LOCATION" json:"cabinet_location" `                       //type:*int     comment:机柜位置                                 version:2023-07-21 08:45
	Abreast                int64            `gorm:"column:ABREAST" json:"abreast" `                                         //type:*int     comment:并排放置(0:否,1:是)                      version:2023-07-21 08:45
	LocationDescription    string           `gorm:"column:LOCATION_DESCRIPTION" json:"location_description" example:"位置描述"` //type:string   comment:位置描述                                 version:2023-07-21 08:45
	ExtensionTest          string           `gorm:"column:EXTENSION_TEST" json:"extension_test" `                           //type:string   comment:扩展测试                                 version:2023-08-18 09:14
	BasicExpansion         []AssetExpansion `gorm:"-" json:"basic_expansion" swaggerignore:"true"`
}

type AssetBasicImport struct {
	Id                     int64            `gorm:"column:ID;primaryKey" json:"id" `                                                                                                                               //type:*int     comment:主键                                     version:2023-07-21 08:45
	DeviceType             int64            `gorm:"column:DEVICE_TYPE" json:"device_type" cn:"设备类型" source:"type=table,table=device_type,property=id;types,field=name,val=1"`                                      //type:*int   comment:设备类型                                 version:2023-07-21 08:45
	AssetCode              string           `gorm:"column:ASSET_CODE" json:"asset_code" cn:"资产编号"`                                                                                                                 //type:string   comment:资产编号    version:2023-08-04 09:45
	DeviceName             string           `gorm:"column:DEVICE_NAME" json:"device_name" cn:"设备名称"`                                                                                                               //type:string   comment:设备名称                                 version:2023-07-21 08:45
	SerialNumber           string           `gorm:"column:SERIAL_NUMBER" json:"serial_number" cn:"序列号"`                                                                                                            //type:string   comment:序列号                                   version:2023-07-21 08:45
	ManagementIp           string           `gorm:"column:MANAGEMENT_IP" json:"management_ip" cn:"管理IP"`                                                                                                           //type:string   comment:管理IP                                   version:2023-07-21 08:45
	DeviceStatus           int64            `gorm:"column:DEVICE_STATUS" json:"device_status" `                                                                                                                    //type:*int     comment:状态(1:待上线,2:已上线,3:下线,4:报废)    version:2023-07-21 08:45
	ManagedState           int64            `gorm:"column:MANAGED_STATE" json:"managed_state" `                                                                                                                    //type:*int     comment:纳管状态                                 version:2023-07-28 14:36
	DeviceProducer         int64            `gorm:"column:DEVICE_PRODUCER" json:"device_producer" cn:"厂商" source:"type=table,table=device_producer,property=id;producer_type,field=alias,val=producer" `           //type:*int     comment:厂商                                     version:2023-07-25 16:30
	DeviceModel            string           `gorm:"column:DEVICE_MODEL" json:"device_model" cn:"型号"`                                                                                                               //type:string   comment:型号                                     version:2023-07-21 08:45
	Subtype                string           `gorm:"column:SUBTYPE" json:"subtype" `                                                                                                                                //type:string   comment:子类型                                   version:2023-07-21 08:45
	OutlineStructure       string           `gorm:"column:OUTLINE_STRUCTURE" json:"outline_structure" `                                                                                                            //type:string   comment:外形结构                                 version:2023-07-21 08:45
	Specifications         string           `gorm:"column:SPECIFICATIONS" json:"specifications" `                                                                                                                  //type:string   comment:规格                                     version:2023-07-21 08:45
	UNumber                int64            `gorm:"column:U_NUMBER" json:"u_number" `                                                                                                                              //type:*int     comment:U数                                      version:2023-07-21 08:45
	UseStorage             string           `gorm:"column:USE_STORAGE" json:"use_storage" `                                                                                                                        //type:string   comment:使用存储                                 version:2023-07-21 08:45
	OperatingSystem        string           `gorm:"column:OPERATING_SYSTEM" json:"operating_system" cn:"操作系统" source:"type=table,table=dict_data,property=dict_key;type_code,field=dict_value,val=operate_system"` //type:string   comment:操作系统                                 version:2023-07-21 08:45
	ProductionIp           string           `gorm:"-" json:"production_ip" cn:"生产IP"`
	RelatedService         string           `gorm:"column:RELATED_SERVICE" json:"related_service" cn:"关联业务"`                                                                                                                      //type:string   comment:关联业务                                 version:2023-08-15 09:23
	ServicePath            string           `gorm:"column:SERVICE_PATH" json:"service_path" cn:"业务路径"`                                                                                                                            //type:string   comment:业务路径                                 version:2023-08-15 09:23
	DatacenterId           int64            `gorm:"column:DATACENTER_ID" json:"datacenter_id" cn:"数据中心" source:"type=table,table=datacenter,property=id,field=datacenter_name"`                                                   //type:*int     comment:数据中心                                 version:2023-08-15 09:23
	EquipmentRoom          string           `gorm:"column:EQUIPMENT_ROOM" json:"equipment_room" cn:"所在机房"`                                                                                                                        //type:string   comment:所在机房                                 version:2023-07-21 08:45
	OwningCabinet          string           `gorm:"column:OWNING_CABINET" json:"owning_cabinet" cn:"机柜"`                                                                                                                          //type:string   comment:所在机柜                                 version:2023-07-21 08:45
	CabinetLocation        int64            `gorm:"column:CABINET_LOCATION" json:"cabinet_location" cn:"机柜位置(U)"`                                                                                                                 //type:*int     comment:机柜位置                                 version:2023-07-21 08:45
	CabinetLocationC       string           `gorm:"column:CABINET_LOCATION_C" json:"cabinet_location_c" cn:"机柜位置(C)"`                                                                                                             //type:*int     comment:机柜位置                                 version:2023-07-21 08:45
	Abreast                int64            `gorm:"column:ABREAST" json:"abreast" cn:"并排放置列" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"`                                                                     //type:*int     comment:并排放置(0:否,1:是)                      version:2023-07-21 08:45
	Region                 string           `gorm:"column:REGION" json:"region" cn:"所属分区"`                                                                                                                                        //type:string   comment:所在区域                                 version:2023-07-21 08:45
	DeviceManagerOne       string           `gorm:"column:DEVICE_MANAGER_ONE" json:"device_manager_one" cn:"设备负责人1"`                                                                                                              //type:string   comment:设备负责人1                              version:2023-07-21 08:45
	DeviceManagerTwo       string           `gorm:"column:DEVICE_MANAGER_TWO" json:"device_manager_two" cn:"设备负责人2"`                                                                                                              //type:string   comment:设备负责人2                              version:2023-07-21 08:45
	BusinessManagerOne     string           `gorm:"column:BUSINESS_MANAGER_ONE" json:"business_manager_one" cn:"业务负责人1"`                                                                                                          //type:string   comment:业务负责人1                              version:2023-07-21 08:45
	BusinessManagerTwo     string           `gorm:"column:BUSINESS_MANAGER_TWO" json:"business_manager_two" cn:"业务负责人2"`                                                                                                          //type:string   comment:业务负责人2                              version:2023-07-21 08:45
	Remark                 string           `gorm:"column:REMARK" json:"remark"  `                                                                                                                                                //type:string   comment:备注                                     version:2023-07-21 08:45
	AffiliatedOrganization int64            `gorm:"column:AFFILIATED_ORGANIZATION" json:"affiliated_organization" source:"type=table,table=organization,property=id,field=name"`                                                  //type:string   comment:所属组织机构                             version:2023-07-21 08:45
	LocationDescription    string           `gorm:"column:LOCATION_DESCRIPTION" json:"location_description" `                                                                                                                     //type:string   comment:位置描述                                 version:2023-07-21 08:45
	MaintenanceType        string           `gorm:"column:MAINTENANCE_TYPE" json:"maintenance_type" cn:"维保类型" source:"type=table,table=dict_data,property=dict_key;type_code,field=dict_value,val=maintenance_type" `             //type:*int     comment:维保类型（数据字典）    version:2023-07-30 10:09
	MaintenanceProvider    int64            `gorm:"column:MAINTENANCE_PROVIDER" json:"maintenance_provider" cn:"维保商" source:"type=table,table=device_producer,property=id;producer_type,field=alias,val=third_party_maintenance"` //type:*int     comment:维保商      version:2023-07-30 10:01
	StartAt                int64            `gorm:"column:START_AT" json:"start_at" cn:"维保开始日期" source:"type=date,value=2006-01-02"`                                                                                              //type:*int     comment:开始日期    version:2023-07-30 10:01
	FinishAt               int64            `gorm:"column:FINISH_AT" json:"finish_at" cn:"维保结束日期" source:"type=date,value=2006-01-02"`                                                                                            //type:*int     comment:结束日期    version:2023-07-30 10:01
	BelongDept             int64            `gorm:"column:BELONG_DEPT" json:"belong_dept" cn:"所属部门" source:"type=table,table=organization,property=id,field=name"`                                                                //type:*int     comment:所属部门    version:2023-08-04 09:45
	EquipmentUse           string           `gorm:"column:EQUIPMENT_USE" json:"equipment_use" cn:"设备用途"`                                                                                                                          //type:string   comment:设备用途    version:2023-08-04 09:45
	BasicExpansion         []AssetExpansion `gorm:"-" json:"basic_expansion" swaggerignore:"true" cn:"expansion"`
}

type AssetBasicMainMang struct {
	AssetBasicCopy       AssetBasicExpansionVo `gorm:"-" json:"asset_basic_copy" `
	AssetMaintenanceCopy AssetMaintenanceVo    `gorm:"-" json:"asset_maintenance_copy" `
	AssetManagementCopy  AssetManagement       `gorm:"-" json:"asset_management_copy" `
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

	session = session.Where(query).Where("DEVICE_STATUS != ?", 4)

	var lst []AssetBasicDetailsVo
	err := session.Debug().Model(&AssetBasic{}).Find(&lst).Error

	return lst, err
}

// 根据ids条件查询
func AssetBasicGetsByIds(ctx *ctx.Context, query []int64, limit, offset int) ([]AssetBasicDetailsVo, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("ID")
	}

	session = session.Where("ID IN ?", query)

	var lst []AssetBasicDetailsVo
	err := session.Model(&AssetBasic{}).Find(&lst).Error

	return lst, err
}

// 根据map条件查询(全局过滤器)
func AssetByMap(ctx *ctx.Context, query map[string]interface{}, limit, offset int) ([]AssetBasicDetailsVo, error) {
	session := DB(ctx).Model(&AssetBasic{}).Joins("left join asset_expansion on asset_expansion.ASSET_ID = asset_basic.ID").Joins("left join asset_maintenance on asset_maintenance.ASSET_ID = asset_basic.ID").Joins("left join asset_management on asset_management.ASSET_ID = asset_basic.ID")
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("ID")
	}

	findQuery, ok := query["query"]
	if ok {
		delete(query, "query")
	}
	filter, filterOk := query["filter"]
	if filterOk {
		str := strings.Split(filter.(string), "-")
		delete(query, "filter")
		if str[0] == "alter" {
			session.Model(&AssetBasic{}).Joins("left join asset_alter on asset_alter.ASSET_ID = asset_basic.ID")
			session.Where("asset_alter.ID > ?", 0)
		} else if str[0] == "maint" {
			session.Where("asset_maintenance.FINISH_AT < ?", time.Now().Unix())
		} else if str[0] == "status" {
			query["device_status"] = str[1]
		}

	}
	finishAt, finishOk := query["finish_at"]
	if finishOk {
		delete(query, "finish_at")
		session = session.Where("asset_maintenance.FINISH_AT = ?", finishAt)
	}
	session = session.Where(query)
	_, statusOk := query["DEVICE_STATUS"]
	if !statusOk {
		session = session.Where("asset_basic.DEVICE_STATUS != ?", 4)
	}
	logger.Debug(query)

	var err error
	var ids []int64
	var lst []AssetBasicDetailsVo

	if ok {
		findQuery = "%" + findQuery.(string) + "%"
		err = session.Debug().Model(&AssetBasic{}).Distinct("asset_basic.ID").Where("asset_basic.MANAGEMENT_IP LIKE ? OR asset_basic.DEVICE_NAME LIKE ? OR asset_basic.SERIAL_NUMBER LIKE ? OR asset_basic.DEVICE_PRODUCER LIKE ? OR asset_basic.DEVICE_MODEL LIKE ? OR asset_basic.RELATED_SERVICE LIKE ? OR asset_basic.REMARK LIKE ? OR (asset_expansion.PROPERTY_NAME_CN = ? AND asset_expansion.PROPERTY_VALUE LIKE ?) OR (asset_expansion.PROPERTY_NAME_CN = ? AND asset_expansion.PROPERTY_VALUE LIKE ?) OR (asset_expansion.PROPERTY_CATEGORY = ? AND asset_expansion.PROPERTY_VALUE LIKE ?) OR asset_management.EQUIPMENT_USE LIKE ? OR asset_management.ASSET_CODE LIKE ?", findQuery, findQuery, findQuery, findQuery, findQuery, findQuery, findQuery, "RAID级别", findQuery, "MAC地址", findQuery, "basic_expansion", findQuery, findQuery, findQuery).Find(&ids).Error
	} else {
		err = session.Debug().Model(&AssetBasic{}).Distinct("asset_basic.ID").Find(&ids).Error
	}
	if len(ids) > 0 {
		err = DB(ctx).Model(&AssetBasic{}).Where("ID IN ?", ids).Find(&lst).Error
	}

	return lst, err
}

// 按id查询
func AssetBasicGetById[T any](ctx *ctx.Context, id int64) (*T, error) {
	var obj *T
	err := DB(ctx).Model(&AssetBasic{}).Where("ID = ?", id).Find(&obj).Error
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
func (a *AssetBasic) Add(ctx *ctx.Context, tx *gorm.DB) (int64, error) {
	// 这里写AssetBasic的业务逻辑，通过error返回错误

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
		wTree.ManagementIp = a.ManagementIp
		wTree.SerialNumber = a.SerialNumber
		err = tx.Create(wTree).Error
		if err != nil {
			tx.Rollback()
		}
	} else {
		//设备类型存在，判断设备厂商是否存在
		var menuTree1 *AssetTree
		deviceProducer, _ := DeviceProducerGetById(ctx, a.DeviceProducer)
		err = tx.Where("NAME = ? AND STATUS = ? AND PARENT_ID = ?", deviceProducer.Alias, menuTree.Status, menuTree.Id).Find(&menuTree1).Error
		if err != nil {
			tx.Rollback()
		}
		//设备厂商不存在，插入设备厂商、设备名称
		if menuTree1.Id == 0 {
			wTree = BuildObj(deviceProducer.Alias, a.CreatedBy, "producer", menuTree.Id, a.DeviceStatus, a.DeviceProducer)
			err = tx.Create(wTree).Error
			if err != nil {
				tx.Rollback()
			}
			wTree = BuildObj(a.DeviceName, a.CreatedBy, "asset", wTree.Id, a.DeviceStatus, a.Id)
			wTree.ManagementIp = a.ManagementIp
			wTree.SerialNumber = a.SerialNumber
			err = tx.Create(wTree).Error
			if err != nil {
				tx.Rollback()
			}
		} else {
			//设备厂商存在，插入设备名称
			wTree = BuildObj(a.DeviceName, a.CreatedBy, "asset", menuTree1.Id, a.DeviceStatus, a.Id)
			wTree.ManagementIp = a.ManagementIp
			wTree.SerialNumber = a.SerialNumber
			err = tx.Create(wTree).Error
			if err != nil {
				tx.Rollback()
			}
		}
	}
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
	//开启事务
	tx := DB(ctx).Begin()

	//删除资产树
	var assetTree AssetTree
	err := tx.Where("NAME = ?", a.DeviceName).Find(&assetTree).Error
	if err != nil {
		return err
	}

	err = assetTree.Del(tx)
	if err != nil {
		tx.Rollback()
	}

	// 实际向库中写入
	err = tx.Delete(a).Error
	if err != nil {
		tx.Rollback()
	}

	//删除扩展字段
	err = tx.Delete(&AssetExpansion{}).Where("ASSET_ID = ? AND CONFIG_CATEGORY = form-basic", a.Id).Error
	if err != nil {
		tx.Rollback()
	}

	tx.Commit()
	return err
}

// 更新资产详情
func (a *AssetBasicExpansionVo) Update(ctx *ctx.Context, tx *gorm.DB, name string, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写AssetBasic的业务逻辑，通过error返回错误
	// var oldAssetBasic *AssetBasic
	// var err error
	//更新资产树
	// oldAssetBasic, err = AssetBasicGetById[AssetBasic](ctx, a.Id)
	// if err != nil {
	// 	return err
	// }

	// if oldAssetBasic.DeviceName != a.DeviceName {
	// 	err = tx.Model(&AssetTree{}).Where("NAME = ?", oldAssetBasic.DeviceName).Updates(map[string]interface{}{"NAME": a.DeviceName, "UPDATED_BY": name}).Error
	// 	if err != nil {
	// 		tx.Rollback()
	// 	}
	// }
	// if oldAssetBasic.DeviceProducer != a.DeviceProducer {
	// 	err = tx.Model(&AssetTree{}).Where("NAME = ?", oldAssetBasic.DeviceProducer).Updates(map[string]interface{}{"NAME": a.DeviceProducer, "UPDATED_BY": name}).Error
	// 	if err != nil {
	// 		tx.Rollback()
	// 	}
	// }
	// if oldAssetBasic.DeviceType != a.DeviceType {
	// 	err = tx.Model(&AssetTree{}).Where("NAME = ?", oldAssetBasic.DeviceType).Updates(map[string]interface{}{"NAME": a.DeviceType, "UPDATED_BY": name}).Error
	// 	if err != nil {
	// 		tx.Rollback()
	// 	}
	// }

	// 实际向库中写入
	err := tx.Model(&AssetBasic{}).Where("ID = ?", a.Id).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 根据条件统计个数
func AssetBasicCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetBasic{}).Where(where, args...).Where("DEVICE_STATUS != ?", 4))
}

// 根据传入条件统计个数
func AssetCountByMap(ctx *ctx.Context, where map[string]interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetBasic{}).Where(where).Where("DEVICE_STATUS != ?", 4))
}

// 根据传入条件统计个数
func AssetCountByIds(ctx *ctx.Context, where []int64) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetBasic{}).Where("ID IN ?", where))
}

// 根据map条件查询(全局过滤器)
func AssetCountMap(ctx *ctx.Context, query map[string]interface{}) (num int64, err error) {
	session := DB(ctx).Model(&AssetBasic{}).Joins("left join asset_expansion on asset_expansion.ASSET_ID = asset_basic.ID").Joins("left join asset_maintenance on asset_maintenance.ASSET_ID = asset_basic.ID").Joins("left join asset_management on asset_management.ASSET_ID = asset_basic.ID")

	findQuery, ok := query["query"]
	if ok {
		delete(query, "query")
	}
	filter, filterOk := query["filter"]
	if filterOk {
		str := strings.Split(filter.(string), "-")
		delete(query, "filter")
		if str[0] == "alter" {
			session.Model(&AssetBasic{}).Joins("left join asset_alter on asset_alter.ASSET_ID = asset_basic.ID")
			session.Where("asset_alter.ID > ?", 0)
		} else if str[0] == "maint" {
			session.Where("asset_maintenance.FINISH_AT < ?", time.Now().Unix())
		} else if str[0] == "status" {
			query["device_status"] = str[1]
		}

	}
	finishAt, finishOk := query["finish_at"]
	if finishOk {
		delete(query, "finish_at")
		session = session.Where("asset_maintenance.FINISH_AT = ?", finishAt)
	}
	session = session.Where(query)
	_, statusOk := query["DEVICE_STATUS"]
	if !statusOk {
		session = session.Where("asset_basic.DEVICE_STATUS != ?", 4)
	}

	if ok {
		findQuery = "%" + findQuery.(string) + "%"
		err = session.Debug().Model(&AssetBasic{}).Distinct("asset_basic.ID").Where("asset_basic.MANAGEMENT_IP LIKE ? OR asset_basic.DEVICE_NAME LIKE ? OR asset_basic.SERIAL_NUMBER LIKE ? OR asset_basic.DEVICE_PRODUCER LIKE ? OR asset_basic.DEVICE_MODEL LIKE ? OR asset_basic.RELATED_SERVICE LIKE ? OR asset_basic.REMARK LIKE ? OR (asset_expansion.PROPERTY_NAME_CN = ? AND asset_expansion.PROPERTY_VALUE LIKE ?) OR (asset_expansion.PROPERTY_NAME_CN = ? AND asset_expansion.PROPERTY_VALUE LIKE ?) OR (asset_expansion.PROPERTY_CATEGORY = ? AND asset_expansion.PROPERTY_VALUE LIKE ?) OR asset_management.EQUIPMENT_USE LIKE ? OR asset_management.ASSET_CODE LIKE ?", findQuery, findQuery, findQuery, findQuery, findQuery, findQuery, findQuery, "RAID级别", findQuery, "MAC地址", findQuery, "basic_expansion", findQuery, findQuery, findQuery).Count(&num).Error
	} else {
		err = session.Debug().Model(&AssetBasic{}).Distinct("asset_basic.ID").Count(&num).Error
	}
	return num, err
}

type Result struct {
	Type int64 `gorm:"type" json:"type" `
	Num  int64 `gorm:"num" json:"num" `
}

// 统计已上线/待上线/已下线个数
func AssetStatusCount(ctx *ctx.Context) (res []Result, err error) {
	err = DB(ctx).Model(&AssetBasic{}).Select("DEVICE_STATUS as type, count(DEVICE_STATUS) as num").Group("DEVICE_STATUS").Having("DEVICE_STATUS != 4").Find(&res).Error
	return res, err
}

//判断表是否存在
func HasTableByName(ctx *ctx.Context, name string) bool {
	return DB(ctx).Migrator().HasTable(name)
}

//判断表中是否存在字段
func HasTableFieldByName(ctx *ctx.Context, name, field string) (num int64, err error) {
	return Count(DB(ctx).Where("table_name = ? AND column_name = ?", name, field).Table("information_schema.columns"))
}

//查询数据
func TableGetsByName(ctx *ctx.Context, name string, fields []string) (m []map[string]interface{}, err error) {
	err = DB(ctx).Select(fields).Table(name).Find(&m).Error
	return m, err
}

// 删除资产详情
func AssetBasicBatchDel(tx *gorm.DB, ids []int64) error {
	// 这里写AssetBasic的业务逻辑，通过error返回错误

	//删除资产详情
	err := tx.Where("ID IN ?", ids).Delete(&AssetBasic{}).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

//更新状态
func UpdateStatus(ctx *ctx.Context, assetIds []int64, status int64, name string) error {
	return DB(ctx).Model(&AssetBasic{}).Where("ID IN ?", assetIds).Updates(map[string]interface{}{"DEVICE_STATUS": status, "UPDATED_BY": name}).Error
}

//更新状态
func UpdateTxStatus(tx *gorm.DB, assetIds []int64, status int64, name string) error {
	return tx.Model(&AssetBasic{}).Where("ID IN ?", assetIds).Updates(map[string]interface{}{"DEVICE_STATUS": status, "UPDATED_BY": name}).Error
}

// 更新资产详情(map)
func UpdateBasicMap(ctx *ctx.Context, ids []int64, m map[string]interface{}) error {
	// 这里写AssetManagement的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(&AssetBasic{}).Where("ID IN ?", ids).Updates(m).Error
}

// 更新资产详情(map)
func UpdateBasicTxMap(tx *gorm.DB, id int64, m map[string]interface{}) error {
	// 这里写AssetManagement的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Debug().Model(&AssetBasic{}).Where("ID = ?", id).Updates(m).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

//校验该设备型号是否有资产使用
func AssetBasicCountByModels(ctx *ctx.Context, ids []int64) ([]AssetBasic, error) {
	var obj []AssetBasic
	err := DB(ctx).Model(&AssetBasic{}).Where("DEVICE_MODEL IN ?", ids).Find(&obj).Error
	return obj, err
}
