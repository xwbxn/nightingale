// Package models  资产变更
// date : 2023-08-04 14:49
// desc : 资产变更
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/toolkits/pkg/logger"
)

// AssetAlter  资产变更。
// 说明:
// 表名:asset_alter
// group: AssetAlter
// version:2023-08-04 14:49
type AssetAlter struct {
	Id               int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键                                          version:2023-08-04 14:49
	AssetId          int64  `gorm:"column:ASSET_ID" json:"asset_id" `                         //type:*int     comment:资产ID                                        version:2023-08-04 14:49
	AlterAt          int64  `gorm:"column:ALTER_AT" json:"alter_at" `                         //type:*int     comment:变更日期                                      version:2023-08-04 15:30
	AlterEventCode   string `gorm:"column:ALTER_EVENT_CODE" json:"alter_event_code" `         //type:string   comment:变更事项编码                                  version:2023-08-04 14:49
	AlterEventKey    string `gorm:"column:ALTER_EVENT_KEY" json:"alter_event_key" `           //type:string   comment:变更事项标签                                  version:2023-08-04 14:49
	BeforeAlter      string `gorm:"column:BEFORE_ALTER" json:"before_alter" `                 //type:string   comment:变更前                                        version:2023-08-04 14:49
	AfterAlter       string `gorm:"column:AFTER_ALTER" json:"after_alter" `                   //type:string   comment:变更后                                        version:2023-08-04 14:49
	AlterSponsor     string `gorm:"column:ALTER_SPONSOR" json:"alter_sponsor" `               //type:string   comment:变更发起人                                    version:2023-08-04 14:49
	AlterStatus      int64  `gorm:"column:ALTER_STATUS" json:"alter_status" `                 //type:*int     comment:确认状态(0:未确认;1:确认)                     version:2023-08-04 14:49
	AlterInstruction string `gorm:"column:ALTER_INSTRUCTION" json:"alter_instruction" `       //type:string   comment:变更说明                                      version:2023-08-04 14:49
	ConfirmOpinion   string `gorm:"column:CONFIRM_OPINION" json:"confirm_opinion" `           //type:string   comment:确认意见                                      version:2023-08-04 14:49
	CreationMode     int64  `gorm:"column:CREATION_MODE" json:"creation_mode" `               //type:*int     comment:创建方式(1:人工录入;2:系统产生;3:信息修改)    version:2023-08-04 14:49
	CreatedBy        string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人                                        version:2023-08-04 14:49
	CreatedAt        int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间                                      version:2023-08-04 14:49
	UpdatedBy        string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人                                        version:2023-08-04 14:49
	UpdatedAt        int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间                                      version:2023-08-04 14:49
}

type AssetAlterVo struct {
	Id                 int64  `gorm:"column:ID;primaryKey" json:"id" `                    //type:*int     comment:主键                                          version:2023-08-04 14:49
	AssetId            int64  `gorm:"column:ASSET_ID" json:"asset_id" `                   //type:*int     comment:资产ID                                        version:2023-08-04 14:49
	AlterAt            int64  `gorm:"column:ALTER_AT" json:"alter_at" `                   //type:*int     comment:变更日期                                      version:2023-08-04 15:30
	AlterEventCode     string `gorm:"column:ALTER_EVENT_CODE" json:"alter_event_code" `   //type:string   comment:变更事项编码                                  version:2023-08-04 14:49
	AlterEventKey      string `gorm:"column:ALTER_EVENT_KEY" json:"alter_event_key" `     //type:string   comment:变更事项标签                                  version:2023-08-04 14:49
	BeforeAlter        string `gorm:"column:BEFORE_ALTER" json:"before_alter" `           //type:string   comment:变更前                                        version:2023-08-04 14:49
	AfterAlter         string `gorm:"column:AFTER_ALTER" json:"after_alter" `             //type:string   comment:变更后                                        version:2023-08-04 14:49
	AlterSponsor       string `gorm:"column:ALTER_SPONSOR" json:"alter_sponsor" `         //type:string   comment:变更发起人                                    version:2023-08-04 14:49
	AlterStatus        int64  `gorm:"column:ALTER_STATUS" json:"alter_status" `           //type:*int     comment:确认状态(0:未确认;1:确认)                     version:2023-08-04 14:49
	AlterInstruction   string `gorm:"column:ALTER_INSTRUCTION" json:"alter_instruction" ` //type:string   comment:变更说明                                      version:2023-08-04 14:49
	ConfirmOpinion     string `gorm:"column:CONFIRM_OPINION" json:"confirm_opinion" `     //type:string   comment:确认意见                                      version:2023-08-04 14:49
	CreationMode       int64  `gorm:"column:CREATION_MODE" json:"creation_mode" `         //type:*int     comment:创建方式(1:人工录入;2:系统产生;3:信息修改)    version:2023-08-04 14:49
	ManagementIp       string `gorm:"-" json:"management_ip" `                            //type:string   comment:管理IP                                   version:2023-07-21 08:45
	SerialNumber       string `gorm:"-" json:"serial_number" `                            //type:string   comment:序列号                                   version:2023-07-21 08:45
	DeviceTypeName     string `gorm:"-" json:"device_type_name" `
	DeviceProducerName string `gorm:"-" json:"device_producer_name" `
	DeviceModelName    string `gorm:"-" json:"device_model_name" `
	RoomName           string `gorm:"-" json:"room_name" `
	CabinetName        string `gorm:"-" json:"cabinet_name" `
	UNumber            int64  `gorm:"-" json:"u_number" `
}

// TableName 表名:asset_alter，资产变更。
// 说明:
func (a *AssetAlter) TableName() string {
	return "asset_alter"
}

// 条件查询
func AssetAlterGets(ctx *ctx.Context, query string, limit, offset int) ([]AssetAlter, error) {
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

	var lst []AssetAlter
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func AssetAlterGetById(ctx *ctx.Context, id int64) (*AssetAlter, error) {
	var obj *AssetAlter
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按资产id查询
func AssetAlterGetByAssetId(ctx *ctx.Context, assetId int64, limit, offset int) ([]AssetAlterVo, error) {
	logger.Debug("---------------------")
	logger.Debug(limit)
	logger.Debug(offset)
	logger.Debug("+---------------------")
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	var obj []AssetAlterVo
	err := session.Debug().Model(&AssetAlter{}).Where("ASSET_ID = ?", assetId).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func AssetAlterGetsAll(ctx *ctx.Context) ([]AssetAlter, error) {
	var lst []AssetAlter
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加资产变更
func (a *AssetAlter) Add(ctx *ctx.Context) error {
	// 这里写AssetAlter的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 删除资产变更
func (a *AssetAlter) Del(ctx *ctx.Context) error {
	// 这里写AssetAlter的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(a).Error
}

// 更新资产变更
func (a *AssetAlter) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写AssetAlter的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(a).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func AssetAlterCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetAlter{}).Where(where, args...))
}
