// Package models  设备型号
// date : 2023-07-08 14:55
// desc : 设备型号
package models

import (
	"reflect"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// DeviceModel  设备型号。
// 说明:
// 表名:device_model
// group: DeviceModel
// version:2023-07-08 14:55
type DeviceModel struct {
	Id                 int64   `gorm:"column:ID;primaryKey" json:"id" `                                                                               //type:*int       comment:主键                       version:2023-07-08 14:55
	Name               string  `gorm:"column:NAME" json:"name" cn:"型号名称"`                                                                             //type:string     comment:型号名称                   version:2023-07-08 14:55
	Subtype            string  `gorm:"column:SUBTYPE" json:"subtype"  cn:"子类型"`                                                                       //type:string     comment:子类型                     version:2023-07-08 14:55
	CreatedBy          string  `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                      //type:string     comment:创建人                     version:2023-07-08 14:55
	CreatedAt          int64   `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                      //type:*int       comment:创建时间                   version:2023-07-08 14:55
	UpdatedBy          string  `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                      //type:string     comment:更新人                     version:2023-07-08 14:55
	UpdatedAt          int64   `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                      //type:*int       comment:更新时间                   version:2023-07-08 14:55
	ProducerId         int64   `gorm:"column:PRODUCER_ID" json:"producer_id" cn:"厂商"   source:"type=table,table=device_producer,field=chinese_name" ` //type:*int       comment:厂商-ID;来源设备厂商信息   version:2023-07-08 14:55
	Model              string  `gorm:"column:MODEL" json:"model"  cn:"型号"`                                                                            //type:string     comment:型号                       version:2023-07-08 14:55
	Series             string  `gorm:"column:SERIES" json:"series"  cn:"系列"`                                                                          //type:string     comment:系列                       version:2023-07-08 14:55
	UNumber            int64   `gorm:"column:U_NUMBER" json:"u_number"  cn:"U数" source:"type=range,value=[1-32;38]" `                                 //type:*int       comment:U数                        version:2023-07-08 14:55
	OutlineStructure   string  `gorm:"column:OUTLINE_STRUCTURE" json:"outline_structure"  cn:"外形结构"`                                                  //type:string     comment:外形结构                   version:2023-07-08 14:55
	Specifications     string  `gorm:"column:SPECIFICATIONS" json:"specifications"  cn:"规格"`                                                          //type:string     comment:规格                       version:2023-07-08 14:55
	MaximumMemory      float64 `gorm:"column:MAXIMUM_MEMORY" json:"maximum_memory"  cn:"最大内存" `                                                       //type:*float64   comment:最大内存（M）              version:2023-07-08 14:55
	WorkingConsumption float64 `gorm:"column:WORKING_CONSUMPTION" json:"working_consumption"  cn:"工作功耗"`                                              //type:*float64   comment:工作功耗(W)                version:2023-07-08 14:55
	RatedConsumption   float64 `gorm:"column:RATED_CONSUMPTION" json:"rated_consumption"  cn:"额定功耗"`                                                  //type:*float64   comment:额定功耗(W)                version:2023-07-08 14:55
	PeakConsumption    float64 `gorm:"column:PEAK_CONSUMPTION" json:"peak_consumption"  cn:"峰值功耗"`                                                    //type:*float64   comment:峰值功耗(W)                version:2023-07-08 14:55
	Weight             float64 `gorm:"column:WEIGHT" json:"weight"  cn:"设备重量"`                                                                        //type:*float64   comment:设备重量(kg)               version:2023-07-08 14:55
	Enlistment         int64   `gorm:"column:ENLISTMENT" json:"enlistment"   cn:"服役期限" source:"type=option,value=[男;女]"  `                            //type:*int       comment:服役期限(月)               version:2023-07-08 14:55
	OutBandVersion     string  `gorm:"column:OUT_BAND_VERSION" json:"out_band_version" cn:"带外版本"`                                                     //type:string     comment:带外版本                   version:2023-07-08 14:55
	Describe           string  `gorm:"column:DESCRIBE" json:"describe" cn:"描述"`                                                                       //type:string     comment:描述                       version:2023-07-08 14:55
}

func GetChinese() map[string]string {
	var d DeviceModel
	dType := reflect.TypeOf(d)
	var mapLit map[string]string
	for i := 0; i < dType.NumField(); i++ {
		fieldType := dType.Field(i)
		if fieldType.Tag.Get("cn") != "" {
			mapLit[fieldType.Tag.Get("cn")] = fieldType.Name + ";" + fieldType.Type.Name()
		}
	}
	return mapLit
}

// TableName 表名:device_model，设备型号。
// 说明:
func (d *DeviceModel) TableName() string {
	return "device_model"
}

// 条件查询
func DeviceModelGets(ctx *ctx.Context, query string, limit, offset int) ([]DeviceModel, error) {
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

	var lst []DeviceModel
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DeviceModelGetById(ctx *ctx.Context, id int64) (*DeviceModel, error) {
	var obj *DeviceModel
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DeviceModelGetsAll(ctx *ctx.Context) ([]DeviceModel, error) {
	var lst []DeviceModel
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加设备型号
func (d *DeviceModel) Add(ctx *ctx.Context) error {
	// 这里写DeviceModel的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除设备型号
func (d *DeviceModel) Del(ctx *ctx.Context) error {
	// 这里写DeviceModel的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新设备型号
func (d *DeviceModel) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DeviceModel的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DeviceModelCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DeviceModel{}).Where(where, args...))
}
