// Package models  操作日志
// date : 2023-9-17 14:14
// desc : 操作日志
package models

import (
	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// OperationLog  操作日志。
// 说明:
// 表名:operation_log
// group: OperationLog
// version:2023-9-17 14:14
type OperationLog struct {
	OperId        int64  `gorm:"column:OPER_ID;primaryKey" json:"oper_id" `                //type:BIGINT   comment:日志主键                                   version:2023-9-17 14:14
	Title         string `gorm:"column:TITLE" json:"title" `                               //type:string   comment:模块标题                                   version:2023-9-17 14:14
	BusinessType  int64  `gorm:"column:BUSINESS_TYPE" json:"business_type" `               //type:*int     comment:业务类型（0其它 1新增 2修改 3删除）        version:2023-9-17 14:14
	Method        string `gorm:"column:METHOD" json:"method" `                             //type:string   comment:方法名称                                   version:2023-9-17 14:14
	RequestMethod string `gorm:"column:REQUEST_METHOD" json:"request_method" `             //type:string   comment:请求方式                                   version:2023-9-17 14:14
	OperatorType  int64  `gorm:"column:OPERATOR_TYPE" json:"operator_type" `               //type:*int     comment:操作类别（0其它 1后台用户 2手机端用户）    version:2023-9-17 14:14
	OperName      string `gorm:"column:OPER_NAME" json:"oper_name" `                       //type:string   comment:操作人员                                   version:2023-9-17 14:14
	DeptName      string `gorm:"column:DEPT_NAME" json:"dept_name" `                       //type:string   comment:部门名称                                   version:2023-9-17 14:14
	OperUrl       string `gorm:"column:OPER_URL" json:"oper_url" `                         //type:string   comment:请求URL                                    version:2023-9-17 14:14
	OperIp        string `gorm:"column:OPER_IP" json:"oper_ip" `                           //type:string   comment:主机地址                                   version:2023-9-17 14:14
	OperLocation  string `gorm:"column:OPER_LOCATION" json:"oper_location" `               //type:string   comment:操作地点                                   version:2023-9-17 14:14
	OperParam     string `gorm:"column:OPER_PARAM" json:"oper_param" `                     //type:string   comment:请求参数                                   version:2023-9-17 14:14
	JsonResult    string `gorm:"column:JSON_RESULT" json:"json_result" `                   //type:string   comment:返回参数                                   version:2023-9-17 14:14
	Status        int64  `gorm:"column:STATUS" json:"status" `                             //type:*int     comment:操作状态（0正常 1异常）                    version:2023-9-17 14:14
	ErrorMsg      string `gorm:"column:ERROR_MSG" json:"error_msg" `                       //type:string   comment:错误消息                                   version:2023-9-17 14:14
	OperTime      int64  `gorm:"column:OPER_TIME" json:"oper_time" `                       //type:*int     comment:操作时间                                   version:2023-9-17 14:14
	CreatedBy     string `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string   comment:创建人                                     version:2023-9-17 14:14
	CreatedAt     int64  `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int     comment:创建时间                                   version:2023-9-17 14:14
	UpdatedBy     string `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string   comment:更新人                                     version:2023-9-17 14:14
	UpdatedAt     int64  `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int     comment:更新时间                                   version:2023-9-17 14:14
	DeletedAt     string `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:string   comment:删除时间                                   version:2023-9-17 14:14
}

// TableName 表名:operation_log，操作日志。
// 说明:
func (o *OperationLog) TableName() string {
	return "operation_log"
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
