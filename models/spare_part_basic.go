// Package models  备件基础数据
// date : 2023-08-20 15:42
// desc : 备件基础数据
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// SparePartBasic  备件基础数据。
// 说明:
// 表名:spare_part_basic
// group: SparePartBasic
// version:2023-08-20 15:42
type SparePartBasic struct {
	Id                  int64          `gorm:"column:ID;primaryKey" json:"id" `                                                                                                                             //type:*int       comment:主键                   version:2023-08-20 15:42
	ProductId           string         `gorm:"column:PRODUCT_ID" json:"product_id" cn:"商品编号"`                                                                                                               //type:string     comment:商品编号               version:2023-08-20 15:42
	ComponentName       string         `gorm:"column:COMPONENT_NAME" json:"component_name" cn:"部件名称"`                                                                                                       //type:string     comment:部件名称               version:2023-08-20 15:42
	ComponentType       int64          `gorm:"column:COMPONENT_TYPE" json:"component_type" cn:"部件类型" source:"type=table,table=component_type,property=id,field=component_type"`                             //type:*int       comment:部件类型               version:2023-08-20 15:42
	ComponentNum        string         `gorm:"column:COMPONENT_NUM" json:"component_num" cn:"部件号"`                                                                                                          //type:string     comment:部件号                 version:2023-08-20 15:42
	ComponentBrand      string         `gorm:"column:COMPONENT_BRAND" json:"component_brand" cn:"部件品牌" source:"type=table,table=device_producer,property=id;producer_type,field=alias,val=component_brand"` //type:string     comment:部件品牌               version:2023-08-20 15:42
	Specification       string         `gorm:"column:SPECIFICATION" json:"specification" cn:"型号规格"`                                                                                                         //type:string     comment:型号规格               version:2023-08-20 15:42
	ComponentUnit       string         `gorm:"column:COMPONENT_UNIT" json:"component_unit" cn:"部件单位"`                                                                                                       //type:string     comment:部件单位               version:2023-08-20 15:42
	UnitPrice           float64        `gorm:"column:UNIT_PRICE" json:"unit_price" cn:"单价(元)"`                                                                                                              //type:float64    comment:单价(元)               version:2023-08-20 15:42
	DeviceType          int64          `gorm:"column:DEVICE_TYPE" json:"device_type" source:"type=table,table=device_type;types,property=id,field=name,val=2"`                                              //type:*int       comment:设备类型               version:2023-08-21 09:04
	AssetClassification string         `gorm:"column:ASSET_CLASSIFICATION" json:"asset_classification" cn:"资产分类"`                                                                                           //type:string     comment:资产分类               version:2023-08-20 15:42
	BelongLine          string         `gorm:"column:BELONG_LINE" json:"belong_line" cn:"所属条线"`                                                                                                             //type:string     comment:所属条线               version:2023-08-20 15:42
	BelongOrganization  string         `gorm:"column:BELONG_ORGANIZATION" json:"belong_organization" cn:"所属单位"`                                                                                             //type:string     comment:所属单位               version:2023-08-20 15:42
	PurchasingApplicant string         `gorm:"column:PURCHASING_APPLICANT" json:"purchasing_applicant" cn:"采购申请人" `                                                                                         //type:string     comment:采购申请人             version:2023-08-20 15:42
	Supplier            int64          `gorm:"column:SUPPLIER" json:"supplier" cn:"供应商" source:"type=table,table=device_producer,property=id;producer_type,field=alias,val=supplier" `                      //type:*int       comment:供应商                 version:2023-08-20 15:42
	DetailedConfig      string         `gorm:"column:DETAILED_CONFIG" json:"detailed_config" cn:"详细配置"`                                                                                                     //type:string     comment:详细配置               version:2023-08-20 15:42
	Remark              string         `gorm:"column:REMARK" json:"remark" cn:"备注"`                                                                                                                         //type:string     comment:备注                   version:2023-08-20 15:42
	SparePartDetail     int64          `gorm:"column:SPARE_PART_DETAIL" json:"spare_part_detail" cn:"备件明细" validate:"omitempty,oneof=0 1" source:"type=option,value=[否;是]"`                                 //type:*int       comment:备件明细(0:否;1:是)    version:2023-08-20 15:42
	ComponentPicture    string         `gorm:"column:COMPONENT_PICTURE" json:"component_picture" `                                                                                                          //type:string     comment:部件图片               version:2023-08-20 15:42
	CreatedBy           string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`                                                                                                    //type:string     comment:创建人                 version:2023-08-20 15:42
	CreatedAt           int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`                                                                                                    //type:*int       comment:创建时间               version:2023-08-20 15:42
	UpdatedBy           string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`                                                                                                    //type:string     comment:更新人                 version:2023-08-20 15:42
	UpdatedAt           int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`                                                                                                    //type:*int       comment:更新时间               version:2023-08-20 15:42
	DeletedAt           gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"`                                                                                                    //type:*int       comment:删除时间        version:2023-9-08 16:39
}

// TableName 表名:spare_part_basic，备件基础数据。
// 说明:
func (s *SparePartBasic) TableName() string {
	return "spare_part_basic"
}

// 条件查询
func SparePartBasicGets(ctx *ctx.Context, query string, limit, offset int) ([]SparePartBasic, error) {
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

	var lst []SparePartBasic
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func SparePartBasicGetById(ctx *ctx.Context, id int64) (*SparePartBasic, error) {
	var obj *SparePartBasic
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// 按id查询照片路径
func ComponentPictureGetById(ctx *ctx.Context, ids []int64) ([]SparePartBasic, error) {
	var obj []SparePartBasic
	err := DB(ctx).Model(&SparePartBasic{}).Select("COMPONENT_PICTURE").Where("ID IN ?", ids).Find(&obj).Error
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// 查询所有
func SparePartBasicGetsAll(ctx *ctx.Context) ([]SparePartBasic, error) {
	var lst []SparePartBasic
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加备件基础数据
func (s *SparePartBasic) Add(ctx *ctx.Context) error {
	// 这里写SparePartBasic的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(s).Error
}

// 删除备件基础数据
func (s *SparePartBasic) Del(ctx *ctx.Context) error {
	// 这里写SparePartBasic的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(s).Error
}

// 删除备件基础数据
func SparePartBasicBatchDel(ctx *ctx.Context, s []int64) error {
	// 这里写SparePartBasic的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Where("ID IN ?", s).Delete(&SparePartBasic{}).Error
}

// 更新备件基础数据
func (s *SparePartBasic) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写SparePartBasic的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(s).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func SparePartBasicCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&SparePartBasic{}).Where(where, args...))
}
