// Package models  资产维保
// date : 2023-07-23 09:44
// desc : 资产维保
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/toolkits/pkg/logger"
	"gorm.io/gorm"
	optlock "gorm.io/plugin/optimisticlock"
)

// AssetMaintenance  资产维保。
// 说明:
// 表名:asset_maintenance
// group: AssetMaintenance
// version:2023-07-23 09:44
type AssetMaintenance struct {
	Id                  int64           `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键        version:2023-07-30 10:01
	AssetId             int64           `gorm:"column:ASSET_ID" json:"asset_id" `                         //type:*int     comment:资产ID      version:2023-07-30 10:01
	MaintenanceType     string          `gorm:"column:MAINTENANCE_TYPE" json:"maintenance_type" `         //type:*int     comment:维保类型（数据字典）    version:2023-07-30 10:09
	MaintenanceProvider int64           `gorm:"column:MAINTENANCE_PROVIDER" json:"maintenance_provider" ` //type:*int     comment:维保商      version:2023-07-30 10:01
	StartAt             int64           `gorm:"column:START_AT" json:"start_at" `                         //type:*int     comment:开始日期    version:2023-07-30 10:01
	FinishAt            int64           `gorm:"column:FINISH_AT" json:"finish_at" `                       //type:*int     comment:结束日期    version:2023-07-30 10:01
	MaintenancePeriod   string          `gorm:"column:MAINTENANCE_PERIOD" json:"maintenance_period" `     //type:string   comment:维保期限    version:2023-07-30 10:01
	ProductionAt        int64           `gorm:"column:PRODUCTION_AT" json:"production_at" `               //type:*int     comment:出厂日期    version:2023-07-30 10:01
	Version             optlock.Version `gorm:"column:VERSION" json:"version" swaggerignore:"true"`       //type:*int     comment:版本号      version:2023-07-30 10:01
	CreatedBy           string          `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人      version:2023-07-30 10:01
	CreatedAt           int64           `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间    version:2023-07-30 10:01
	UpdatedBy           string          `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人      version:2023-07-30 10:01
	UpdatedAt           int64           `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间    version:2023-07-30 10:01
	DeletedAt           gorm.DeletedAt  `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*int       comment:删除时间        version:2023-9-08 16:39
}

type AssetMaintenanceVo struct {
	Id                  int64                      `gorm:"column:ID;primaryKey" json:"id" `                                    //type:*int     comment:主键        version:2023-07-30 10:01
	AssetId             int64                      `gorm:"column:ASSET_ID" json:"asset_id" `                                   //type:*int     comment:资产ID      version:2023-07-30 10:01
	MaintenanceType     string                     `gorm:"column:MAINTENANCE_TYPE" json:"maintenance_type" `                   //type:*int     comment:维保类型（数据字典）    version:2023-07-30 10:09
	MaintenanceProvider int64                      `gorm:"column:MAINTENANCE_PROVIDER" json:"maintenance_provider" `           //type:*int     comment:维保商      version:2023-07-30 10:01
	StartAt             int64                      `gorm:"column:START_AT" json:"start_at" `                                   //type:*int     comment:开始日期    version:2023-07-30 10:01
	FinishAt            int64                      `gorm:"column:FINISH_AT" json:"finish_at" `                                 //type:*int     comment:结束日期    version:2023-07-30 10:01
	MaintenancePeriod   string                     `gorm:"column:MAINTENANCE_PERIOD" json:"maintenance_period" example:"维保期限"` //type:string   comment:维保期限    version:2023-07-30 10:01
	ProductionAt        int64                      `gorm:"column:PRODUCTION_AT" json:"production_at" `                         //type:*int     comment:出厂日期    version:2023-07-30 10:01
	ServiceConfig       []MaintenanceServiceConfig `gorm:"-" json:"service_config" `
}

// TableName 表名:asset_maintenance，资产维保。
// 说明:
func (a *AssetMaintenance) TableName() string {
	return "asset_maintenance"
}

// 条件查询
func AssetMaintenanceGets(ctx *ctx.Context, query string, limit, offset int) ([]AssetMaintenance, error) {
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

	var lst []AssetMaintenance
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func AssetMaintenanceGetById(ctx *ctx.Context, id int64) (*AssetMaintenanceVo, error) {
	var obj *AssetMaintenanceVo
	err := DB(ctx).Model(&AssetMaintenance{}).Where("ID = ?", id).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 按assetId查询
func AssetMaintenanceVoGetByAssetId(ctx *ctx.Context, assetId int64) (*AssetMaintenanceVo, error) {
	var obj *AssetMaintenanceVo
	err := DB(ctx).Model(&AssetMaintenance{}).Where("ASSET_ID = ?", assetId).Find(&obj).Error
	if err != nil {
		return nil, err
	}

	if (&AssetMaintenanceVo{} != obj) {
		configLst, err := MaintenanceServiceConfigGetByMaintId(ctx, obj.Id)
		if err != nil {
			return nil, err
		}
		obj.ServiceConfig = configLst
	}

	return obj, nil
}

// 查询所有
func AssetMaintenanceGetsAll(ctx *ctx.Context) ([]AssetMaintenance, error) {
	var lst []AssetMaintenance
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加资产维保
func (a *AssetMaintenance) Add(ctx *ctx.Context) error {
	// 这里写AssetMaintenance的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 增加资产维保
func (a *AssetMaintenance) AddTx(tx *gorm.DB) error {
	// 这里写AssetMaintenance的业务逻辑，通过error返回错误

	// 实际向库中写入
	err := tx.Create(a).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 增加资产维保及配置记录
func (fVo *AssetMaintenanceVo) AddConfig(tx *gorm.DB, name string) error {
	// 这里写AssetMaintenance的业务逻辑，通过error返回错误

	var f AssetMaintenance

	f.AssetId = fVo.AssetId
	f.MaintenanceType = fVo.MaintenanceType
	f.MaintenanceProvider = fVo.MaintenanceProvider
	f.MaintenancePeriod = fVo.MaintenancePeriod
	// 添加审计信息
	f.CreatedBy = name

	// 更新模型
	err := tx.Debug().Create(&f).Error
	if err != nil {
		tx.Rollback()
	}

	//更新服务配置表
	config := fVo.ServiceConfig
	if len(config) != 0 {
		// 添加审计信息及维保ID
		for index := range config {
			config[index].MaintenanceId = f.Id
			config[index].CreatedBy = name
		}
		// 写入维保配置表
		err = tx.Debug().Create(&config).Error
		if err != nil {
			tx.Rollback()
		}
	}
	return err
}

// 删除资产维保
func (a *AssetMaintenanceVo) DelTx(ctx *ctx.Context) error {
	// 这里写AssetMaintenance的业务逻辑，通过error返回错误
	tx := DB(ctx).Begin()

	err := tx.Where("ID = ?", a.Id).Delete(&AssetMaintenance{}).Error
	if err != nil {
		tx.Rollback()
	}
	err = tx.Where("MAINTENANCE_ID = ?", a.Id).Delete(&MaintenanceServiceConfig{}).Error
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return err
}

// 批量删除资产维保
func AssetMaintenanceBatchDel(ctx *ctx.Context, tx *gorm.DB, ids []int64) error {
	//删除资产扩展
	//删除资产维保
	maint := make([]int64, 0)
	for _, val := range ids {
		assetMaint, err := AssetMaintenanceVoGetByAssetId(ctx, val)
		if err != nil {
			tx.Rollback()
		}
		if (&AssetMaintenanceVo{} == assetMaint) {
			continue
		}
		maint = append(maint, (*assetMaint).Id)
	}

	err := tx.Where("ASSET_ID IN ?", ids).Delete(&AssetMaintenance{}).Error
	if err != nil {
		tx.Rollback()
	}
	err = tx.Where("MAINTENANCE_ID IN ?", maint).Delete(&MaintenanceServiceConfig{}).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

// 更新资产维保
func MaintenanceUpdate(ctx *ctx.Context, id int64, name string, updateFrom AssetMaintenanceVo) error {
	// 这里写AssetMaintenance的业务逻辑，通过error返回错误
	tx := DB(ctx).Begin()

	oldConfigs, err := MaintenanceServiceConfigGetByMaintId(ctx, id)
	if err != nil {
		tx.Rollback()
	}
	configs := updateFrom.ServiceConfig
	updateFrom.ServiceConfig = nil

	// 实际向库中写入
	err = tx.Debug().Model(&AssetMaintenance{}).Where("ID = ?", id).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
	if err != nil {
		tx.Rollback()
	}

	logger.Debug(configs)
	for _, val := range configs {
		//新增
		if val.Id == 0 {
			val.MaintenanceId = updateFrom.Id
			err := val.AddTx(tx)
			if err != nil {
				return err
			}
		} else {
			//更新
			for _, oldVal := range oldConfigs {
				if oldVal.Id == val.Id {
					err := oldVal.UpdateTx(tx, val, "*")
					if err != nil {
						return err
					}
					oldVal.Id = 0
					break
				}
			}
		}
	}
	for _, oldVal := range oldConfigs {
		if oldVal.Id == 0 {
			err := oldVal.DelTx(tx)
			if err != nil {
				return err
			}
			break
		}
	}
	tx.Commit()
	return err
}

// 根据条件统计个数
func AssetMaintenanceCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetMaintenance{}).Where(where, args...))
}

// 根据条件统计个数
func AssetMaintenanceCountMap(ctx *ctx.Context, where map[string]interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&AssetMaintenance{}).Where(where))
}

// 删除资产维保
func (a *AssetMaintenance) BatchDel(ctx *ctx.Context) error {
	// 这里写AssetMaintenance的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(a).Error
}
