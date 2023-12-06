package operationlog

import (
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/aop"
	"github.com/ccfos/nightingale/v6/pkg/ctx"

	"github.com/toolkits/pkg/concurrent/semaphore"
	"github.com/toolkits/pkg/logger"
)

type Consumer struct {
	ctx *ctx.Context
}

// 创建一个 Consumer 实例
func NewConsumer(ctx *ctx.Context) *Consumer {
	return &Consumer{
		ctx: ctx,
	}
}

func (e *Consumer) LoopConsume() {
	sema := semaphore.NewSemaphore(1)
	duration := time.Duration(1000) * time.Millisecond
	for {
		events := aop.OperationlogQueue.PopBackBy(100)
		if len(events) == 0 {
			time.Sleep(duration)
			continue
		}
		e.consume(events, sema)
	}
}

func (e *Consumer) consume(events []interface{}, sema *semaphore.Semaphore) {
	for i := range events {
		if events[i] == nil {
			continue
		}

		// logger.Debug((events[0]))
		// for _, val := range events {
		// 	// logger.Debug(val.(type))
		// 	m, ok := events[i].(map[string]interface{})
		// 	if !ok {
		// 		logger.Error("类型断言失败")
		// 	}
		// 	logger.Debug(m)
		// }
		event, ok := events[i].(map[string]interface{})
		if !ok {
			logger.Error("类型断言失败")
		}
		log := &models.OperationLog{
			Type:        event["type"].(string),
			Object:      event["object"].(string),
			Description: event["description"].(string),
			User:        event["user_name"].(string),
			OperTime:    event["oper_time"].(int64),
			OperUrl:     event["oper_url"].(string),
			OperParam:   event["oper_param"].(string),
			JsonResult:  "",
			Status:      event["status"].(int64),
			ErrorMsg:    event["error_msg"].(string),
			CreatedBy:   event["user_name"].(string),
			// CreatedAt:  int64(event["created_at"].(float64)),
			// UpdatedBy:  "",
			// UpdatedAt:  int64(event["updated_at"].(float64)),
		}
		sema.Acquire()
		go func(event *models.OperationLog) {
			defer sema.Release()
			e.consumeOne(log)
		}(log)

		// event := events[i].(*models.OperationLog)
		// sema.Acquire()
		// go func(event *models.OperationLog) {
		// 	defer sema.Release()
		// 	e.consumeOne(event)
		// }(event)
	}
}

func (e *Consumer) consumeOne(event *models.OperationLog) {

	e.persist(event)

}

func (e *Consumer) persist(event *models.OperationLog) {

	logger.Debug("操作日志数据开始入库:", *event)
	err := event.Add(e.ctx)
	if err != nil {
		logger.Errorf("event%+v persist err:%v", event, err)
	}
}
