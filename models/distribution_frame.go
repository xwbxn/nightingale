// Package models  配线架信息
// date : 2023-07-16 10:14
// desc : 配线架信息
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// DistributionFrame  配线架信息。
// 说明:
// 表名:distribution_frame
// group: DistributionFrame
// version:2023-07-16 10:14
type DistributionFrame struct {
	Id              int64  `gorm:"column:ID;primaryKey" json:"id" `                                                                                               //type:*int     comment:主键          version:2023-07-16 10:14
	RoomId          int64  `gorm:"column:ROOM_ID" json:"room_id" cn:"所属机房" validate:"required" source:"type=table,table=computer_room,field=room_name"`           //type:*int     comment:所属机房                             version:2023-07-25 17:23
	CabinetId       int64  `gorm:"column:CABINET_ID" json:"cabinet_id" cn:"机柜编号" validate:"required" source:"type=table,table=device_cabinet,field=cabinet_code"` //type:*int     comment:所属机柜                             version:2023-07-29 15:08
	DisFrame        string `gorm:"column:DIS_FRAME" json:"dis_frame" cn:"配线架编号" validate:"required"`                                                              //type:string   comment:配线架编号    version:2023-07-16 10:14
	DisName         string `gorm:"column:DIS_NAME" json:"dis_name" cn:"配线架名称"`                                                                                    //type:string   comment:配线架名称    version:2023-07-16 10:14
	ProducerId      int64  `gorm:"column:PRODUCER_ID" json:"producer_id" cn:"厂商" source:"type=table,table=device_producer,field=alias"`                           //type:*int     comment:厂商                                 version:2023-07-25 17:25
	Model           int64  `gorm:"column:MODEL" json:"model" cn:"型号" source:"type=table,table=device_model,field=model"`                                          //type:string   comment:型号          version:2023-07-16 10:14
	Specification   string `gorm:"column:SPECIFICATION" json:"specification" cn:"规格"`                                                                             //type:string   comment:规格          version:2023-07-16 10:14
	DisType         int64  `gorm:"column:DIS_TYPE" json:"dis_type" cn:"配线架类型" source:"type=option,value=[双绞线;光纤配线架]"`                                             //type:*int     comment:配线架类型(0:双绞线;1:光纤配线架)    version:2023-07-16 11:02
	TotalPortNum    int64  `gorm:"column:TOTAL_PORT_NUM" json:"total_port_num" cn:"总端口数" validate:"required,gte=1"`                                               //type:*int     comment:总端口数      version:2023-07-16 10:14
	UsedPortNum     int64  `gorm:"column:USED_PORT_NUM" json:"used_port_num" validate:"omitempty,gte=0"`                                                          //type:*int     comment:已用端口数    version:2023-07-16 10:14
	PortPrefix      string `gorm:"column:PORT_PREFIX" json:"port_prefix" cn:"端口前缀"`                                                                               //type:string   comment:端口前缀      version:2023-07-16 10:14
	CabinetLocation int64  `gorm:"column:CABINET_LOCATION" json:"cabinet_location" cn:"所在机柜位置(U)" validate:"omitempty,gte=0"`                                     //type:*int     comment:机柜位置(U)   version:2023-07-16 10:14
	PurchaseAt      int64  `gorm:"column:PURCHASE_AT" json:"purchase_at" cn:"采购日期" source:"type=date,value=2006-01-02"`                                           //type:*int     comment:采购日期      version:2023-07-16 10:14
	DisPicture      string `gorm:"column:DIS_PICTURE" json:"dis_picture" `                                                                                        //type:string   comment:配线架图片    version:2023-07-16 10:14
	Use             string `gorm:"column:USE" json:"use" cn:"用途"`                                                                                                 //type:string   comment:用途          version:2023-07-16 10:14
	DutyPersonOne   string `gorm:"column:DUTY_PERSON_ONE" json:"duty_person_one" cn:"责任人1" validate:"required"`                                                   //type:string   comment:责任人1       version:2023-07-16 10:14
	DutyPersonTwo   string `gorm:"column:DUTY_PERSON_TWO" json:"duty_person_two" cn:"责任人2"`                                                                       //type:string   comment:责任人2       version:2023-07-16 10:14
	Unumber         int64  `gorm:"column:UNUMBER" json:"unumber" cn:"U数" validate:"omitempty,gte=0" source:"type=option,value=[1U;2U]"`                           //type:*int     comment:U数           version:2023-07-16 10:14
	CreatedBy       string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                      //type:string   comment:创建人        version:2023-07-16 10:14
	CreatedAt       int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                      //type:*int     comment:创建时间      version:2023-07-16 10:14
	UpdatedBy       string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                      //type:string   comment:更新人        version:2023-07-16 10:14
	UpdatedAt       int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                      //type:*int     comment:更新时间      version:2023-07-16 10:14
}

// TableName 表名:distribution_frame，配线架信息。
// 说明:
func (d *DistributionFrame) TableName() string {
	return "distribution_frame"
}

// 条件查询
func DistributionFrameGets(ctx *ctx.Context, query string, limit, offset int) ([]DistributionFrame, error) {
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

	var lst []DistributionFrame
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DistributionFrameGetById(ctx *ctx.Context, id int64) (*DistributionFrame, error) {
	var obj *DistributionFrame
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DistributionFrameGetsAll(ctx *ctx.Context) ([]DistributionFrame, error) {
	var lst []DistributionFrame
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加配线架信息
func (d *DistributionFrame) Add(ctx *ctx.Context) error {
	// 这里写DistributionFrame的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除配线架信息
func (d *DistributionFrame) Del(ctx *ctx.Context) error {
	// 这里写DistributionFrame的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新配线架信息
func (d *DistributionFrame) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DistributionFrame的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DistributionFrameCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DistributionFrame{}).Where(where, args...))
}
