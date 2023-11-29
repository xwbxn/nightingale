// Package models  接口管理
// date : 2023-10-20 16:47
// desc : 接口管理
package models

import (
	"errors"
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"github.com/ccfos/nightingale/v6/pkg/prom"
	"github.com/prometheus/common/model"
	"github.com/toolkits/pkg/logger"
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

type Serie struct {
	SeriesName string                   `json:"seriesName"`
	DataPoint  []map[string]interface{} `json:"data"`
}

type Series struct {
	Series []Serie `json:"series"`
}

func (as *ApiService) Execute(ctx *ctx.Context, api prom.API) (interface{}, error) {
	series := Series{
		Series: make([]Serie, 0),
	}
	if as.Type == "sql" {
		var dataPoint []map[string]interface{}
		err := DB(ctx).Raw(as.Script).Find(&dataPoint).Error
		if err != nil {
			return nil, err
		}

		s := Serie{
			SeriesName: as.Name,
			DataPoint:  dataPoint,
		}
		series.Series = append(series.Series, s)
	}
	if as.Type == "promql" {
		r := prom.Range{
			Start: time.Now().Add(-30 * time.Minute),
			End:   time.Now(),
			Step:  1 * time.Minute,
		}
		value, warnings, err := api.QueryRange(ctx.Ctx, as.Script, r)
		if len(warnings) > 0 {
			logger.Error(err)
			return nil, err
		}
		items, ok := value.(model.Matrix)
		if !ok {
			return nil, errors.New("prom查询结果无法解析")
		}
		for _, item := range items {
			s := Serie{
				SeriesName: as.Name,
				DataPoint:  make([]map[string]interface{}, 0),
			}
			for _, v := range item.Values {
				s.DataPoint = append(s.DataPoint, map[string]interface{}{
					"name":  v.Timestamp.Time().Format(time.Kitchen),
					"value": float64(v.Value),
				})
			}
			series.Series = append(series.Series, s)
		}
	}
	return series, nil
}
