// Package models  接口管理
// date : 2023-10-20 16:47
// desc : 接口管理
package models

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/prom"
	"github.com/prometheus/common/model"
	"github.com/toolkits/pkg/logger"
	"github.com/xwb1989/sqlparser"
)

// ApiService  接口管理。
// 说明:
// 表名:api_service
// group: ApiService
// version:2023-10-20 16:47
type ApiService struct {
	Id           int64  `gorm:"column:ID;primaryKey" json:"id" `                          //type:*int     comment:主键                           version:2023-10-20 16:47
	CreatedBy    string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人                         version:2023-10-20 16:47
	CreatedAt    int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间                       version:2023-10-20 16:47
	UpdatedBy    string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人                         version:2023-10-20 16:47
	UpdatedAt    int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间                       version:2023-10-20 16:47
	DeletedAt    int64  `gorm:"column:DELETED_AT" json:"deleted_at" `                     //type:*int     comment:删除时间                       version:2023-10-20 16:47
	Name         string `gorm:"column:NAME" json:"name" `                                 //type:string   comment:名称                           version:2023-10-20 16:47
	Type         string `gorm:"column:TYPE" json:"type" `                                 //type:string   comment:类型;sql or promql             version:2023-10-20 16:47
	DatasourceId int64  `gorm:"column:DATASOURCE_ID" json:"datasource_id" `               //type:*int     comment:数据源;promql 需要指定数据源   version:2023-10-20 16:47
	Url          string `gorm:"column:URL" json:"url" `                                   //type:string   comment:URL                            version:2023-10-20 16:47
	Script       string `gorm:"column:SCRIPT" json:"script" `                             //type:string   comment:执行脚本                       version:2023-10-20 16:47
	ValueField   string `gorm:"column:VALUE_FIELD" json:"value_field"`
}

// TableName 表名:api_service，接口管理。
// 说明:
func (a *ApiService) TableName() string {
	return "api_service"
}

// 条件查询
func ApiServiceGets(ctx *ctx.Context, query string, limit, offset int) ([]ApiService, error) {
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

	var lst []ApiService
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func ApiServiceGetById(ctx *ctx.Context, id int64) (*ApiService, error) {
	var obj *ApiService
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func ApiServiceGetsAll(ctx *ctx.Context) ([]ApiService, error) {
	var lst []ApiService
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加接口管理
func (a *ApiService) Add(ctx *ctx.Context) error {
	// 这里写ApiService的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(a).Error
}

// 删除接口管理
func (a *ApiService) Del(ctx *ctx.Context) error {
	// 这里写ApiService的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(a).Error
}

// 更新接口管理
func (a *ApiService) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写ApiService的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(a).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func ApiServiceCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&ApiService{}).Where(where, args...))
}

// 按url查询
func ApiServiceGetByUrl(ctx *ctx.Context, url string) (*ApiService, error) {
	var obj *ApiService
	err := DB(ctx).Take(&obj, &ApiService{Url: url}).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (as *ApiService) IsDangerous() bool {
	if as.Type == "promql" {
		return false
	}
	if as.Type == "sql" {
		stmt, err := sqlparser.Parse(as.Script)
		if err != nil {
			return false
		}
		switch stmt.(type) {
		case *sqlparser.Select:
			return false
		}
	}
	return true
}

type ApiSerie struct {
	SeriesName string                   `json:"seriesName"`
	DataPoint  []map[string]interface{} `json:"data"`
}

type ApiSeries struct {
	Series []ApiSerie `json:"series"`
}

// 执行api接口
// 返回promtheus api格式
func (as *ApiService) Execute(ctx *ctx.Context, api prom.API) (interface{}, error) {
	data := []model.Samples{}
	samples := model.Samples{}

	if as.Type == "sql" { //sql模式
		var dataPoints []map[string]interface{}
		err := DB(ctx).Raw(as.Script).Find(&dataPoints).Error
		if err != nil {
			return nil, err
		}

		now := model.Now()
		for _, row := range dataPoints {
			sample := &model.Sample{
				Metric:    make(model.Metric),
				Timestamp: now,
			}
			for name, value := range row {
				var strVal string
				var floatVal float64

				switch value.(type) {
				case string:
					strVal = value.(string)
					floatVal, _ = strconv.ParseFloat(strVal, 64)
				case int64:
					floatVal = float64(value.(int64))
					strVal = fmt.Sprintf("%f", floatVal)
				case float64:
					floatVal = value.(float64)
					strVal = fmt.Sprintf("%f", floatVal)
				default:
					logger.Errorf("execute apiservice convert v error: %v", value)
					return nil, errors.New("execute apiservice convert v error")
				}

				sample.Metric[model.LabelName(name)] = model.LabelValue(strVal)
				if name == as.ValueField {
					sample.Value = model.SampleValue(floatVal)
				}

			}
			samples = append(samples, sample)
		}
		data = append(data, samples)
	}
	return data, nil
}
