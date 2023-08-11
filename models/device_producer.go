// Package models  设备厂商
// date : 2023-07-08 14:43
// desc : 设备厂商
package models

import (
	"log"
	"reflect"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

const DOWNLOADNUM = 6

var defaultHeight = 25.0

// DeviceProducer  设备厂商。
// 说明:
// 表名:device_producer
// group: DeviceProducer
// version:2023-07-08 14:43
type DeviceProducer struct {
	Id               int64  `gorm:"column:ID;primaryKey" json:"id" `                                 //type:*int     comment:主键            version:2023-07-08 14:43
	Alias            string `gorm:"column:ALIAS" json:"alias" cn:"厂商简称"`                             //type:string   comment:厂商简称        version:2023-07-08 14:43
	ChineseName      string `gorm:"column:CHINESE_NAME" json:"chinese_name" cn:"中文名称"`               //type:string   comment:中文名称        version:2023-07-08 14:43
	CompanyName      string `gorm:"column:COMPANY_NAME" json:"company_name" cn:"公司全称"`               //type:string   comment:公司全称        version:2023-07-08 14:43
	Official         string `gorm:"column:OFFICIAL" json:"official" cn:"官方站点"`                       //type:string   comment:官方站点        version:2023-07-08 14:43
	IsDomestic       int64  `gorm:"column:IS_DOMESTIC" json:"is_domestic" cn:"是否国产"`                 //type:*int     comment:是否国产        version:2023-07-08 14:43
	IsDisplayChinese int64  `gorm:"column:IS_DISPLAY_CHINESE" json:"is_display_chinese" cn:"是否显示中文"` //type:*int     comment:是否显示中文    version:2023-07-08 14:43
	CreatedBy        string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"`        //type:string   comment:创建人          version:2023-07-08 14:43
	CreatedAt        int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"`        //type:*int     comment:创建时间        version:2023-07-08 14:43
	UpdatedBy        string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"`        //type:string   comment:更新人          version:2023-07-08 14:43
	UpdatedAt        int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"`        //type:*int     comment:更新时间        version:2023-07-08 14:43
}

// TableName 表名:device_producer，设备厂商。
// 说明:
func (d *DeviceProducer) TableName() string {
	return "device_producer"
}

// 条件查询
func DeviceProducerGets(ctx *ctx.Context, query string, limit, offset int) ([]DeviceProducer, error) {
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

	var lst []DeviceProducer
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func DeviceProducerGetById(ctx *ctx.Context, id int64) (*DeviceProducer, error) {
	var obj *DeviceProducer
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func DeviceProducerGetsAll(ctx *ctx.Context) ([]DeviceProducer, error) {
	var lst []DeviceProducer
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加设备厂商
func (d *DeviceProducer) Add(ctx *ctx.Context) error {
	// 这里写DeviceProducer的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(d).Error
}

// 删除设备厂商
func (d *DeviceProducer) Del(ctx *ctx.Context) error {
	// 这里写DeviceProducer的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(d).Error
}

// 更新设备厂商
func (d *DeviceProducer) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写DeviceProducer的业务逻辑，通过error返回错误
	// s := GetChinese(*d)
	// 实际向库中写入
	return DB(ctx).Model(d).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func DeviceProducerCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&DeviceProducer{}).Where(where, args...))

}

//读取excel
func ReadExcel(ctx *ctx.Context, xlsx *excelize.File, createdBy string) (int, error) {

	sheets := xlsx.GetSheetMap()
	sheet1 := sheets[1]
	rows := xlsx.GetRows(sheet1)
	maps := make(map[int]string)
	count := 0

	for i, row := range rows {
		if i == 0 { //取得第一行的所有数据---execel表头
			for index, colCell := range row {
				maps[index] = colCell
			}
			// fmt.Println("列信息", cols)

		} else {
			data := &DeviceProducer{}
			dType := reflect.ValueOf(new(DeviceProducer)).Elem()
			g := reflect.ValueOf(data).Elem()
			for j, colCell := range row {

				for k := 0; k < dType.NumField(); k++ {
					fieldInfo := dType.Type().Field(k)
					if fieldInfo.Tag.Get("cn") == maps[j] {
						switch fieldType := fieldInfo.Type.Kind(); fieldType {
						case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
							s1, _ := strconv.Atoi(colCell)
							g.FieldByName(fieldInfo.Name).SetInt(int64(s1))
						case reflect.String:
							g.FieldByName(fieldInfo.Name).SetString(colCell)
						case reflect.Bool:
							g.FieldByName(fieldInfo.Name).SetBool(colCell == "true")
						default:
							log.Printf("field type %s not support yet", fieldType)
						}
					}
				}
			}
			data.CreatedBy = createdBy
			errAdd := data.Add(ctx)
			if errAdd != nil {
				return count, errAdd
			}
			count++
		}
	}
	return count, nil
}

//根据条件写入excel(暂时不用)
// func WriterExcel(ctx *ctx.Context, query string) ([]map[string]string, []map[string]interface{}, error) {
// 	session := DB(ctx)

// 	// 这里使用列名的硬编码构造查询参数, 避免从前台传入造成注入风险
// 	if query != "" {
// 		q := "%" + query + "%"
// 		session = session.Where("id like ?", q)
// 	}

// 	var lst []DeviceProducer
// 	err := session.Find(&lst).Error
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	dType := reflect.ValueOf(new(DeviceProducer)).Elem()

// 	var firstrow [2][DOWNLOADNUM]string
// 	count := 0
// 	for k := 0; k < dType.NumField(); k++ {
// 		fieldInfo := dType.Type().Field(k)
// 		if fieldInfo.Tag.Get("cn") != "" {
// 			firstrow[1][count] = fieldInfo.Tag.Get("cn")
// 			firstrow[0][count] = fieldInfo.Name
// 			count++
// 		}

// 	}
// 	//定义首行标题
// 	dataKey := make([]map[string]string, 0)

// 	for i := 0; i < count; i++ {
// 		dataKey = append(dataKey, map[string]string{
// 			"key":    firstrow[0][i],
// 			"title":  firstrow[1][i],
// 			"width":  "20",
// 			"is_num": "0",
// 		})
// 	}

// 	//填充数据
// 	data := make([]map[string]interface{}, 0)
// 	if len(lst) > 0 {
// 		for _, v := range lst {
// 			var isDomestic string
// 			if v.IsDomestic == 0 {
// 				isDomestic = "否"
// 			} else {
// 				isDomestic = "是"
// 			}
// 			var isDisplayChinese string
// 			if v.IsDisplayChinese == 0 {
// 				isDisplayChinese = "否"
// 			} else {
// 				isDisplayChinese = "是"
// 			}

// 			data = append(data, map[string]interface{}{
// 				"Alias":            v.Alias,
// 				"ChineseName":      v.ChineseName,
// 				"CompanyName":      v.CompanyName,
// 				"Official":         v.Official,
// 				"IsDomestic":       isDomestic,
// 				"IsDisplayChinese": isDisplayChinese,
// 			})
// 		}
// 	}
// 	return dataKey, data, nil

// }
