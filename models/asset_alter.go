// Package models  资产变更
// date : 2023-08-04 14:49
// desc : 资产变更
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
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
	ConfirmBy        string `gorm:"column:CONFIRM_BY" json:"confirm_by" `                     //type:string   comment:确认人                                        version:2023-08-11 14:23
	ConfirmOpinion   string `gorm:"column:CONFIRM_OPINION" json:"confirm_opinion" `           //type:string   comment:确认意见                                      version:2023-08-04 14:49
	CreationMode     int64  `gorm:"column:CREATION_MODE" json:"creation_mode" `               //type:*int     comment:创建方式(1:人工录入;2:系统产生;3:信息修改)    version:2023-08-04 14:49
	CreatedBy        string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人                                        version:2023-08-04 14:49
	CreatedAt        int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间                                      version:2023-08-04 14:49
	UpdatedBy        string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人                                        version:2023-08-04 14:49
	UpdatedAt        int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间                                      version:2023-08-04 14:49
}

type AssetAlterVo struct {
	Id                 int64  `gorm:"column:ID;primaryKey" json:"id" `                                                                                                                 //type:*int     comment:主键                                          version:2023-08-04 14:49
	AssetId            int64  `gorm:"column:ASSET_ID" json:"asset_id" `                                                                                                                //type:*int     comment:资产ID                                        version:2023-08-04 14:49
	AlterAt            int64  `gorm:"column:ALTER_AT" json:"alter_at" cn:"变更日期" source:"type=date,value=2006-01-02"`                                                                   //type:*int     comment:变更日期                                      version:2023-08-04 15:30
	AlterEventCode     string `gorm:"column:ALTER_EVENT_CODE" json:"alter_event_code" val:""`                                                                                          //type:string   comment:变更事项编码                                  version:2023-08-04 14:49
	AlterEventKey      string `gorm:"column:ALTER_EVENT_KEY" json:"alter_event_key" cn:"变更事项" source:"type=table,table=dict_data,property=type_code;dict_key,field=dict_value" val:""` //type:string   comment:变更事项标签                                  version:2023-08-04 14:49
	BeforeAlter        string `gorm:"column:BEFORE_ALTER" json:"before_alter" cn:"变更前"`                                                                                                //type:string   comment:变更前                                        version:2023-08-04 14:49
	AfterAlter         string `gorm:"column:AFTER_ALTER" json:"after_alter" cn:"变更后"`                                                                                                  //type:string   comment:变更后                                        version:2023-08-04 14:49
	AlterSponsor       string `gorm:"column:ALTER_SPONSOR" json:"alter_sponsor" cn:"变更人"`                                                                                              //type:string   comment:变更发起人                                    version:2023-08-04 14:49
	AlterStatus        int64  `gorm:"column:ALTER_STATUS" json:"alter_status" cn:"变更状态" source:"type=option,value=[未确认;确认]"`                                                           //type:*int     comment:确认状态(0:未确认;1:确认)                     version:2023-08-04 14:49
	AlterInstruction   string `gorm:"column:ALTER_INSTRUCTION" json:"alter_instruction" cn:"变更说明"`                                                                                     //type:string   comment:变更说明                                      version:2023-08-04 14:49
	ConfirmBy          string `gorm:"column:CONFIRM_BY" json:"confirm_by" cn:"确认人"`                                                                                                    //type:string   comment:确认人                                        version:2023-08-11 14:23
	ConfirmOpinion     string `gorm:"column:CONFIRM_OPINION" json:"confirm_opinion" `                                                                                                  //type:string   comment:确认意见                                      version:2023-08-04 14:49
	CreationMode       int64  `gorm:"column:CREATION_MODE" json:"creation_mode" cn:"创建方式" source:"type=option,value=[人工录入;系统产生;信息修改]"`                                                 //type:*int     comment:创建方式(1:人工录入;2:系统产生;3:信息修改)    version:2023-08-04 14:49
	ManagementIp       string `gorm:"-" json:"management_ip" `                                                                                                                         //type:string   comment:管理IP                                   version:2023-07-21 08:45
	SerialNumber       string `gorm:"-" json:"serial_number" `                                                                                                                         //type:string   comment:序列号                                   version:2023-07-21 08:45
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
func AssetAlterGetByMap(ctx *ctx.Context, where map[string]interface{}, limit, offset int) ([]AssetAlterVo, error) {

	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	var obj []AssetAlterVo
	var err error
	start, startOk := where["start"]
	end, endOk := where["end"]
	if startOk && endOk {
		delete(where, "start")
		delete(where, "end")
		err = session.Model(&AssetAlter{}).Where(where).Where("ALTER_AT >= ? AND ALTER_AT <= ?", start, end).Find(&obj).Error
	} else if !startOk && endOk {
		delete(where, "end")
		err = session.Model(&AssetAlter{}).Where(where).Where("ALTER_AT <= ?", end).Find(&obj).Error
	} else if startOk && !endOk {
		delete(where, "start")
		err = session.Model(&AssetAlter{}).Where(where).Where("ALTER_AT >= ?", start).Find(&obj).Error
	} else {
		err = session.Model(&AssetAlter{}).Where(where).Find(&obj).Error
	}
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

// 批量删除资产扩展
func AssetAlterBatchDel(tx *gorm.DB, ids []int64) error {
	//删除资产扩展
	err := tx.Where("ASSET_ID IN ?", ids).Delete(&AssetAlter{}).Error
	if err != nil {
		tx.Rollback()
	}
	return err
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

// 根据assetId统计个数
func AssetAlterCountByMap(ctx *ctx.Context, where map[string]interface{}) (num int64, err error) {
	start, startOk := where["start"]
	end, endOk := where["end"]
	if startOk && endOk {
		delete(where, "start")
		delete(where, "end")
		num, err = Count(DB(ctx).Model(&AssetAlter{}).Where(where).Where("ALTER_AT >= ? AND ALTER_AT <= ?", start, end))
		where["start"] = start
		where["end"] = end
		return num, err
	} else if !startOk && endOk {
		delete(where, "end")
		num, err = Count(DB(ctx).Model(&AssetAlter{}).Where(where).Where("ALTER_AT <= ?", end))
		where["end"] = end
		return num, err
	} else if startOk && !endOk {
		delete(where, "start")
		num, err = Count(DB(ctx).Model(&AssetAlter{}).Where(where).Where("ALTER_AT >= ?", start))
		where["start"] = start
		return num, err
	} else {
		return Count(DB(ctx).Model(&AssetAlter{}).Where(where))
	}

}
