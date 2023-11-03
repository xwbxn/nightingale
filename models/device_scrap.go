// Package models  设备报废
// date : 2023-9-08 14:42
// desc : 设备报废
package models

import (
	"strconv"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// DeviceScrap  设备报废。
// 说明:
// 表名:device_scrap
// group: DeviceScrap
// version:2023-9-08 14:42
type DeviceScrap struct {
	Id                    int64          `gorm:"column:ID;primaryKey" json:"id" `                                                                                                                    //type:*int     comment:主键            version:2023-9-08 16:02
	AssetId               int64          `gorm:"column:ASSET_ID" json:"asset_id" `                                                                                                                   //type:*int     comment:资产ID          version:2023-9-08 17:20
	ScrapAt               int64          `gorm:"column:SCRAP_AT" json:"scrap_at" cn:"报废日期" source:"type=date,value=2006-01-02"`                                                                      //type:*int     comment:报废时间        version:2023-9-08 16:02
	SerialNumber          string         `gorm:"column:SERIAL_NUMBER" json:"serial_number" cn:"序列号" `                                                                                                //type:string   comment:序列号          version:2023-9-08 16:02
	DeviceName            string         `gorm:"column:DEVICE_NAME" json:"device_name" cn:"设备名称" `                                                                                                   //type:string   comment:设备名称        version:2023-9-08 16:02
	OldManagementIp       string         `gorm:"column:OLD_MANAGEMENT_IP" json:"old_management_ip" cn:"原管理IP" `                                                                                      //type:string   comment:管理IP          version:2023-9-08 16:02
	DeviceProducer        int64          `gorm:"column:DEVICE_PRODUCER" json:"device_producer" cn:"厂商" source:"type=table,table=device_producer,property=id;producer_type,field=alias,val=producer"` //type:*int     comment:厂商            version:2023-9-08 16:02
	DeviceModel           int64          `gorm:"column:DEVICE_MODEL" json:"device_model" cn:"型号" source:"type=table,table=device_model,property=id,field=name"`                                      //type:*int     comment:型号            version:2023-9-08 16:02
	DeviceType            int64          `gorm:"column:DEVICE_TYPE" json:"device_type" cn:"设备类型" source:"type=table,table=device_type,property=id;types,field=name,val=1"`                           //type:*int     comment:设备类型        version:2023-9-08 16:02
	AssetCode             string         `gorm:"column:ASSET_CODE" json:"asset_code" `                                                                                                               //type:string   comment:资产编号        version:2023-9-08 16:02
	OldDeviceManager      string         `gorm:"column:OLD_DEVICE_MANAGER" json:"old_device_manager" `                                                                                               //type:string   comment:原责任人        version:2023-9-08 16:02
	OldBelongOrganization int64          `gorm:"column:OLD_BELONG_ORGANIZATION" json:"old_belong_organization" cn:"所属部门" source:"type=table,table=organization,property=id,field=name"`              //type:*int     comment:所属组织机构    version:2023-9-08 16:02
	Remark                string         `gorm:"column:REMARK" json:"remark" cn:"报废说明" `                                                                                                             //type:string   comment:报废说明        version:2023-9-08 16:02
	OldDatacenter         int64          `gorm:"column:OLD_DATACENTER" json:"old_datacenter" cn:"原数据中心" source:"type=table,table=datacenter,property=id,field=datacenter_name"`                      //type:*int     comment:原数据中心      version:2023-9-08 16:02
	OldLocation           string         `gorm:"column:OLD_LOCATION" json:"old_location" cn:"原设备位置" `                                                                                                //type:string   comment:原所在位置      version:2023-9-08 16:02
	CreatedBy             string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                                           //type:string   comment:创建人          version:2023-9-08 16:02
	CreatedAt             int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                                           //type:*int     comment:创建时间        version:2023-9-08 16:02
	UpdatedBy             string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                                           //type:string   comment:更新人          version:2023-9-08 16:02
	UpdatedAt             int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                                           //type:*int     comment:更新时间        version:2023-9-08 16:02
	DeletedAt             gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`                                                                                           //type:*int       comment:删除时间        version:2023-9-08 16:39
}

// TableName 表名:device_scrap，设备报废。
// 说明:
func (d *DeviceScrap) TableName() string {
	return "device_scrap"
}

// 条件查询
func DeviceScrapGets(ctx *ctx.Context, query string, dataRange, datacenter, room int64, limit, offset int) ([]DeviceScrap, error) {
	session := DB(ctx)
	// 分页
	if limit > -1 {
		session = session.Limit(limit).Offset(offset).Order("id")
	}

	// 这里使用列名的硬编码构造查询参数, 避免从前台传入造成注入风险
	if datacenter != -1 {
		session.Where("OLD_DATACENTER = ?", datacenter)
	}
	if room != -1 {
		roomstr := strconv.FormatInt(room, 10) + "%"
		session.Where("OLD_LOCATION LIKE ?", roomstr)
	}
	if dataRange != -1 {
		session.Where("scrap_at >= ?", dataRange)
	}
	if query != "" {
		query = "%" + query + "%"
		session.Where("OLD_MANAGEMENT_IP LIKE ? OR DEVICE_NAME LIKE ? OR SERIAL_NUMBER LIKE ? OR DEVICE_PRODUCE LIKE ? OR DEVICE_MODEL LIKE ? OR ASSET_CODE LIKE ?", query, query, query, query, query, query)
	}

	var lst []DeviceScrap
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DeviceScrapGetById(ctx *ctx.Context, id int64) (*DeviceScrap, error) {
	var obj *DeviceScrap
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按id批量查询
func DeviceScrapBatchGetById(ctx *ctx.Context, ids []int64) ([]DeviceScrap, error) {
	var obj []DeviceScrap
	err := DB(ctx).Model(&DeviceScrap{}).Where("ID IN ?", ids).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DeviceScrapGetsAll(ctx *ctx.Context) ([]DeviceScrap, error) {
	var lst []DeviceScrap
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加设备报废
func (d *DeviceScrap) Add(ctx *ctx.Context) error {
	// 这里写DeviceScrap的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 增加设备报废(事务)
func (d *DeviceScrap) AddTx(tx *gorm.DB) error {
	// 这里写DeviceScrap的业务逻辑，通过error返回错误

	// 实际向库中写入
	return tx.Create(d).Error
}

// 删除设备报废
func (d *DeviceScrap) Del(ctx *ctx.Context) error {
	// 这里写DeviceScrap的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 删除资产详情
func DeviceScrapBatchDel(tx *gorm.DB, ids []int64) error {
	// 这里写DeviceScrap的业务逻辑，通过error返回错误

	//删除资产详情
	err := tx.Where("ID IN ?", ids).Delete(&DeviceScrap{}).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 更新设备报废
func (d *DeviceScrap) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DeviceScrap的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DeviceScrapCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DeviceScrap{}).Where(where, args...))
}

// 根据传入条件统计个数
func DeviceScrapByMap(ctx *ctx.Context, where map[string]interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DeviceScrap{}).Where(where))
}

// 根据过滤器统计个数
func DeviceScrapFindCount(ctx *ctx.Context, query string, dataRange, datacenter, room int64) (num int64, err error) {
	session := DB(ctx)
	if datacenter != -1 {
		session.Where("OLD_DATACENTER = ?", datacenter)
	}
	if room != -1 {
		roomstr := strconv.FormatInt(room, 10) + "%"
		session.Where("OLD_LOCATION LIKE ?", roomstr)
	}
	if dataRange != -1 {
		session.Where("scrap_at >= ?", dataRange)
	}
	if query != "" {
		query = "%" + query + "%"
		session.Where("OLD_MANAGEMENT_IP LIKE ? OR DEVICE_NAME LIKE ? OR SERIAL_NUMBER LIKE ? OR DEVICE_PRODUCE LIKE ? ORR DEVICE_MODEL LIKE ? OR ASSET_CODE LIKE ? OR", query, query, query, query, query, query)
	}
	err = session.Model(&DeviceScrap{}).Count(&num).Error
	return num, err
}
