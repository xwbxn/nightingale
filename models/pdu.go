// Package models  PDU
// date : 2023-07-16 10:13
// desc : PDU
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// Pdu  PDU。
// 说明:
// 表名:pdu
// group: Pdu
// version:2023-07-16 10:13
type Pdu struct {
	Id            int64   `gorm:"column:ID;primaryKey" json:"id" `                                                                                                              //type:*int       comment:主键                  version:2023-07-16 10:13
	AssetsCode    string  `gorm:"column:ASSETS_CODE" json:"assets_code" cn:"资产编号" validate:"required"`                                                                          //type:string     comment:资产编号              version:2023-07-16 10:13
	Name          string  `gorm:"column:NAME" json:"name" cn:"名称"`                                                                                                              //type:string     comment:名称                  version:2023-07-16 10:13
	Brand         int64   `gorm:"column:BRAND" json:"brand" cn:"品牌" source:"type=table,table=device_producer,property=id;producer_type,field=alias,val=component_brand"`        //type:string     comment:品牌                  version:2023-07-16 10:13
	Model         int64   `gorm:"column:MODEL" json:"model" cn:"型号" source:"type=table,table=device_model,property=id,field=model" `                                            //type:string     comment:型号                  version:2023-07-16 10:13
	Standard      int64   `gorm:"column:STANDARD" json:"standard" cn:"标准" validate:"required" source:"type=option,value=[新国标;国标]"`                                              //type:*int       comment:标准(1:新国标;2:国标)    version:2023-07-16 11:05
	JackNum       int64   `gorm:"column:JACK_NUM" json:"jack_num" cn:"插孔数" validate:"required,gte=1"`                                                                           //type:*int       comment:插孔数                version:2023-07-16 10:13
	LimitVoltage  float64 `gorm:"column:LIMIT_VOLTAGE" json:"limit_voltage" cn:"限制电压(V)" validate:"omitempty,gt=0"`                                                             //type:*float64   comment:限制电压(V)           version:2023-07-16 10:13
	MaxElectric   float64 `gorm:"column:MAX_ELECTRIC" json:"max_electric" cn:"最大耐冲击电压(KA)" validate:"omitempty,gt=0"`                                                           //type:*float64   comment:最大耐冲击电压(KA)    version:2023-07-16 10:13
	Use           string  `gorm:"column:USE" json:"use" cn:"用途"`                                                                                                                //type:string     comment:用途                  version:2023-07-16 10:13
	PurchaseAt    int64   `gorm:"column:PURCHASE_AT" json:"purchase_at" cn:"采购日期" source:"type=date,value=2006-01-02 15:04:05"`                                                 //type:*int       comment:采购日期              version:2023-07-16 10:13
	Power         float64 `gorm:"column:POWER" json:"power" cn:"功率" validate:"omitempty,gt=0"`                                                                                  //type:*float64   comment:功率                  version:2023-07-16 10:13
	UnitPrice     float64 `gorm:"column:UNIT_PRICE" json:"unit_price" cn:"单价" validate:"omitempty,gt=0"`                                                                        //type:*float64   comment:单价                  version:2023-07-16 10:13
	BelongRoom    int64   `gorm:"column:BELONG_ROOM" json:"belong_room" cn:"所在机房" validate:"required" source:"type=table,table=computer_room,property=id,field=room_name" `     //type:*int       comment:所在机房                 version:2023-07-29 15:22
	CabinetId     int64   `gorm:"column:CABINET_ID" json:"cabinet_id" cn:"所在机柜编号" validate:"required" source:"type=table,table=device_cabinet,property=id,field=cabinet_code" ` //type:*int       comment:所在机柜编号             version:2023-07-29 15:22
	DutyPersonOne string  `gorm:"column:DUTY_PERSON_ONE" json:"duty_person_one" cn:"责任人1" validate:"required"`                                                                  //type:string     comment:责任人1               version:2023-07-16 10:13
	DutyPersonTwo string  `gorm:"column:DUTY_PERSON_TWO" json:"duty_person_two" cn:"责任人2"`                                                                                      //type:string     comment:责任人2               version:2023-07-16 10:13
	CreatedBy     string  `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                                     //type:string     comment:创建人                version:2023-07-16 10:13
	CreatedAt     int64   `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                                     //type:*int       comment:创建时间              version:2023-07-16 10:13
	UpdatedBy     string  `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                                     //type:string     comment:更新人                version:2023-07-16 10:13
	UpdatedAt     int64   `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                                     //type:*int       comment:更新时间              version:2023-07-16 10:13
}

// TableName 表名:pdu，PDU。
// 说明:
func (p *Pdu) TableName() string {
	return "pdu"
}

// 条件查询
func PduGets(ctx *ctx.Context, query string, limit, offset int) ([]Pdu, error) {
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

	var lst []Pdu
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func PduGetById(ctx *ctx.Context, id int64) (*Pdu, error) {
	var obj *Pdu
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func PduGetsAll(ctx *ctx.Context) ([]Pdu, error) {
	var lst []Pdu
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加PDU
func (p *Pdu) Add(ctx *ctx.Context) error {
	// 这里写Pdu的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(p).Error
}

// 删除PDU
func (p *Pdu) Del(ctx *ctx.Context) error {
	// 这里写Pdu的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(p).Error
}

// 更新PDU
func (p *Pdu) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写Pdu的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(p).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func PduCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&Pdu{}).Where(where, args...))
}
