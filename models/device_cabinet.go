// Package models  机柜信息
// date : 2023-07-11 15:14
// desc : 机柜信息
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// DeviceCabinet  机柜信息。
// 说明:
// 表名:device_cabinet
// group: DeviceCabinet
// version:2023-07-11 15:14
type DeviceCabinet struct {
	Id                     int64   `gorm:"column:ID;primaryKey" json:"id" `                                                                                           //type:*int       comment:主键                                            version:2023-07-11 15:14
	CabinetId              string  `gorm:"column:CABINET_ID" json:"cabinet_id" cn:"机柜ID"`                                                                             //type:string     comment:机柜ID                                          version:2023-07-11 15:14
	EquipmentRoom          string  `gorm:"column:EQUIPMENT_ROOM" json:"equipment_room" cn:"所在机房"`                                                                     //type:string     comment:所在机房                                        version:2023-07-11 15:14
	CabinetCode            string  `gorm:"column:CABINET_CODE" json:"cabinet_code" cn:"机柜编号" validate:"required"`                                                     //type:string     comment:机柜编号                                        version:2023-07-11 15:14
	CabinetName            string  `gorm:"column:CABINET_NAME" json:"cabinet_name" cn:"机柜名称"`                                                                         //type:string     comment:机柜名称                                        version:2023-07-11 15:14
	CabinetCompanyName     string  `gorm:"column:CABINET_COMPANY_NAME" json:"cabinet_company_name" cn:"厂商"`                                                           //type:string     comment:厂商                                            version:2023-07-11 15:14
	CabinetModel           string  `gorm:"column:CABINET_MODEL" json:"cabinet_model" cn:"型号"`                                                                         //type:string     comment:型号                                            version:2023-07-11 15:14
	CabinetPicture         string  `gorm:"column:CABINET_PICTURE" json:"cabinet_picture" `                                                                            //type:string     comment:机柜图片                                        version:2023-07-11 15:14
	Unumber                int64   `gorm:"column:UNUMBER" json:"unumber" cn:"规格(U数)" validate:"required,gte=0,lte=99"`                                                //type:*int       comment:规格(U数)                                       version:2023-07-11 15:14
	RowNumber              int64   `gorm:"column:ROW_NUMBER" json:"row_number" cn:"所在行" validate:"required,gte=1,lte=9999"`                                           //type:*int       comment:所在行                                          version:2023-07-11 15:14
	RowName                string  `gorm:"column:ROW_NAME" json:"row_name" cn:"所在行名称"`                                                                                //type:string     comment:所在行名称                                      version:2023-07-11 15:14
	ColumnNumber           int64   `gorm:"column:COLUMN_NUMBER" json:"column_number" cn:"所在列" validate:"required,gte=1,lte=9999"`                                     //type:*int       comment:所在列                                          version:2023-07-11 15:14
	ColumnName             int64   `gorm:"column:COLUMN_NAME" json:"column_name" cn:"所在列名称"`                                                                          //type:*int       comment:所在列名称                                      version:2023-07-11 15:14
	MainPowerSupply        string  `gorm:"column:MAIN_POWER_SUPPLY" json:"main_power_supply" cn:"主要供电来源"`                                                             //type:string     comment:主要供电来源                                    version:2023-07-11 15:14
	StandbyPowerSupply     string  `gorm:"column:STANDBY_POWER_SUPPLY" json:"standby_power_supply" cn:"备用供电来源"`                                                       //type:string     comment:临时供电来源                                    version:2023-07-11 15:14
	PowerSupplyMode        string  `gorm:"column:POWER_SUPPLY_MODE" json:"power_supply_mode" `                                                                        //type:string     comment:供电方式                                        version:2023-07-11 15:14
	PowerConsumption       int64   `gorm:"column:POWER_CONSUMPTION" json:"power_consumption" cn:"电源功耗" validate:"required,gte=0"`                                     //type:*int       comment:电源功耗                                        version:2023-07-11 15:14
	RatedVoltage           int64   `gorm:"column:RATED_VOLTAGE" json:"rated_voltage" cn:"额定电压" validate:"omitempty,gte=0"`                                            //type:*int       comment:额定电压                                        version:2023-07-11 15:14
	RatedCurrent           int64   `gorm:"column:RATED_CURRENT" json:"rated_current" cn:"额定电流" validate:"omitempty,gte=0"`                                            //type:*int       comment:额定电流                                        version:2023-07-11 15:14
	Use                    string  `gorm:"column:USE" json:"use" cn:"用途"`                                                                                             //type:string     comment:用途                                            version:2023-07-11 15:14
	CabinetType            int64   `gorm:"column:CABINET_TYPE" json:"cabinet_type" cn:"机柜类型" source:"type=option,value=[大一体机机柜;普通机柜;屏蔽机柜]"`                           //type:*int       comment:机柜类型;1:大一体机机柜;2:普通机柜;3:屏蔽机柜   version:2023-07-11 15:14
	ReservedCabinet        int64   `gorm:"column:RESERVED_CABINET" json:"reserved_cabinet" cn:"预留机柜" validate:"omitempty,oneof=0 1" source:"type=option,value=[是;否]"` //type:*int       comment:预留机柜                                        version:2023-07-11 15:14
	UnavailableSpace       int64   `gorm:"column:UNAVAILABLE_SPACE" json:"unavailable_space" cn:"不可用空间" validate:"omitempty,gte=0,ltefield=Unumber"`                  //type:*int       comment:不可用空间                                      version:2023-07-11 15:14
	DutyPersonOne          string  `gorm:"column:DUTY_PERSON_ONE" json:"duty_person_one" cn:"责任人1"`                                                                   //type:string     comment:责任人1                                         version:2023-07-11 15:14
	DutyPersonTwo          string  `gorm:"column:DUTY_PERSON_TWO" json:"duty_person_two" cn:"责任人2"`                                                                   //type:string     comment:责任人2                                         version:2023-07-11 15:14
	CabinetBearingCapacity float64 `gorm:"column:CABINET_BEARING_CAPACITY" json:"cabinet_bearing_capacity" cn:"机柜承重" validate:"omitempty,gte=0"`                      //type:*float64   comment:机柜承重                                        version:2023-07-11 15:14
	CabinetArea            float64 `gorm:"column:CABINET_AREA" json:"cabinet_area" validate:"omitempty,gte=0"`                                                        //type:*float64   comment:机柜面积                                        version:2023-07-11 15:14
	ServicePartition       string  `gorm:"column:SERVICE_PARTITION" json:"service_partition" `                                                                        //type:string     comment:业务分区                                        version:2023-07-11 15:14
	PowerPlugNumber        int64   `gorm:"column:POWER_PLUG_NUMBER" json:"power_plug_number" cn:"电源插头数量" validate:"omitempty,gte=0"`                                  //type:*int       comment:电源插头数量                                    version:2023-07-11 15:14
	CreatedBy              string  `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                  //type:string     comment:创建人                                          version:2023-07-11 15:14
	CreatedAt              int64   `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                  //type:*int       comment:创建时间                                        version:2023-07-11 15:14
	UpdatedBy              string  `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                  //type:string     comment:更新人                                          version:2023-07-11 15:14
	UpdatedAt              int64   `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                  //type:*int       comment:更新时间                                        version:2023-07-11 15:14
}

// TableName 表名:device_cabinet，机柜信息。
// 说明:
func (d *DeviceCabinet) TableName() string {
	return "device_cabinet"
}

// 条件查询
func DeviceCabinetGets(ctx *ctx.Context, query string, limit, offset int) ([]DeviceCabinet, error) {
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

	var lst []DeviceCabinet
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DeviceCabinetGetById(ctx *ctx.Context, id int64) (*DeviceCabinet, error) {
	var obj *DeviceCabinet
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DeviceCabinetGetsAll(ctx *ctx.Context) ([]DeviceCabinet, error) {
	var lst []DeviceCabinet
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加机柜信息
func (d *DeviceCabinet) Add(ctx *ctx.Context) error {
	// 这里写DeviceCabinet的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除机柜信息
func (d *DeviceCabinet) Del(ctx *ctx.Context) error {
	// 这里写DeviceCabinet的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新机柜信息
func (d *DeviceCabinet) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DeviceCabinet的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DeviceCabinetCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DeviceCabinet{}).Where(where, args...))
}
