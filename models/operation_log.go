// Package models  操作日志
// date : 2023-10-21 09:10
// desc : 操作日志
package models

import (
	"strings"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
	"gorm.io/gorm"
)

// OperationLog  操作日志。
// 说明:
// 表名:operation_log
// group: OperationLog
// version:2023-10-21 09:10
type OperationLog struct {
	Id          int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:BIGINT       comment:日志主键      version:2023-10-21 09:14
	Type        string         `gorm:"column:TYPE" json:"type" `                                 //type:string       comment:类型          version:2023-10-21 09:10
	Object      string         `gorm:"column:OBJECT" json:"object" `                             //type:string       comment:对象          version:2023-10-21 09:10
	Description string         `gorm:"column:DESCRIPTION" json:"description" `                   //type:string       comment:描述          version:2023-10-21 09:10
	User        string         `gorm:"column:USER" json:"user" `                                 //type:string       comment:用户          version:2023-10-21 09:10
	OperTime    int64          `gorm:"column:OPER_TIME" json:"oper_time" `                       //type:*int         comment:执行时间      version:2023-10-21 09:10
	OperUrl     string         `gorm:"column:OPER_URL" json:"oper_url" `                         //type:string       comment:请求URL       version:2023-10-21 09:10
	OperParam   string         `gorm:"column:OPER_PARAM" json:"oper_param" `                     //type:string       comment:请求参数      version:2023-10-21 09:10
	JsonResult  string         `gorm:"column:JSON_RESULT" json:"json_result" `                   //type:string       comment:返回参数      version:2023-10-21 09:10
	Status      int64          `gorm:"column:STATUS" json:"status" `                             //type:*int         comment:操作状态码    version:2023-10-21 09:10
	ErrorMsg    string         `gorm:"column:ERROR_MSG" json:"error_msg" `                       //type:string       comment:错误消息      version:2023-10-21 09:10
	CreatedBy   string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string       comment:创建人        version:2023-10-21 09:10
	CreatedAt   int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int         comment:创建时间      version:2023-10-21 09:10
	UpdatedBy   string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string       comment:更新人        version:2023-10-21 09:10
	UpdatedAt   int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int         comment:更新时间      version:2023-10-21 09:10
	DeletedAt   gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*time.Time   comment:删除时间      version:2023-10-21 09:10
}

// TableName 表名:operation_log，操作日志。
// 说明:
func (o *OperationLog) TableName() string {
	return "operation_log"
}

// 条件查询
func OperationLogGets(ctx *ctx.Context, query string, limit, offset int) ([]OperationLog, error) {
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

	var lst []OperationLog
	err := session.Find(&lst).Error

	return lst, err
}

// 按id查询
func OperationLogGetById(ctx *ctx.Context, id int64) (*OperationLog, error) {
	var obj *OperationLog
	err := DB(ctx).Take(&obj, id).Error
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// 查询所有
func OperationLogGetsAll(ctx *ctx.Context) ([]OperationLog, error) {
	var lst []OperationLog
	err := DB(ctx).Find(&lst).Error

	return lst, err
}

// 增加操作日志
func (o *OperationLog) Add(ctx *ctx.Context) error {
	// 这里写OperationLog的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Create(o).Error
}

// 删除操作日志
func (o *OperationLog) Del(ctx *ctx.Context) error {
	// 这里写OperationLog的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Delete(o).Error
}

// 更新操作日志
func (o *OperationLog) Update(ctx *ctx.Context, updateFrom interface{}, selectField interface{}, selectFields ...interface{}) error {
	// 这里写OperationLog的业务逻辑，通过error返回错误

	// 实际向库中写入
	return DB(ctx).Model(o).Select(selectField, selectFields...).Omit("CREATED_AT", "CREATED_BY").Updates(updateFrom).Error
}

// 根据条件统计个数
func OperationLogCount(ctx *ctx.Context, where string, args ...interface{}) (num int64, err error) {
	return Count(DB(ctx).Model(&OperationLog{}).Where(where, args...))
}

// 过滤器条件统计个数
func FilterLogCount(ctx *ctx.Context, query string, start int64, end int64,
	 filterType string, modelType string)(num int64, err error){
	query = "%" + query + "%"
	var str strings.Builder
	vals := make([]interface{}, 0)

	if filterType == ""{
		
		str.WriteString(" TYPE like ? or ")
		vals = append(vals, query)
		str.WriteString("OBJECT like ? or ")
		vals = append(vals, query)
		str.WriteString("DESCRIPTION like ? or ")
		vals = append(vals, query)
		str.WriteString("USER like ? ")
		vals = append(vals, query)

		err = DB(ctx).Debug().Model(&OperationLog{}).Where(str.String(), vals...).
		Where("`OPER_TIME` BETWEEN ? AND ?", start, end).Count(&num).Error
	}else{
		err = DB(ctx).Debug().Model(&OperationLog{}).Where(filterType + " LIKE ? ", query).
		Where("`OPER_TIME` BETWEEN ? AND ?", start, end).Count(&num).Error
	}
	return num, err
}

// 过滤器条件查询
func FilterLogGets(ctx *ctx.Context, query string, offset int, limit int,
	 start int64, end int64, filterType string, modelType string) (lst []OperationLog, err error){
	query = "%" + query + "%"
	var str strings.Builder
	vals := make([]interface{}, 0)

	if filterType == ""{
		str.WriteString("TYPE like ? or ")
		vals = append(vals, query)
		str.WriteString("OBJECT like ? or ")
		vals = append(vals, query)
		str.WriteString("DESCRIPTION like ? or ")
		vals = append(vals, query)
		str.WriteString("USER like ? ")
		vals = append(vals, query)

		err = DB(ctx).Debug().Model(&OperationLog{}).Where(str.String(), vals...).
		Where("`OPER_TIME` BETWEEN ? AND ?",start,end).Limit(limit).Offset(offset).Find(&lst).Error
	}else{
		err = DB(ctx).Debug().Model(&OperationLog{}).Where(filterType + " LIKE ? ", query).
		Where("`OPER_TIME` BETWEEN ? AND ?",start,end).Limit(limit).Offset(offset).Find(&lst).Error
	}
	return lst, err
}