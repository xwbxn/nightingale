package router

import (
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/models"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func getTimeRange(c *gin.Context) (stime, etime int64) {
	stime = ginx.QueryInt64(c, "stime", 0)
	etime = ginx.QueryInt64(c, "etime", 0)
	hours := ginx.QueryInt64(c, "hours", 0)
	now := time.Now().Unix()
	if hours != 0 {
		stime = now - 3600*hours
		etime = now + 3600*24
	}

	if stime != 0 && etime == 0 {
		etime = now + 3600*24
	}
	return
}

func (rt *Router) alertHisEventsList(c *gin.Context) {
	stime, etime := getTimeRange(c)

	severity := ginx.QueryInt(c, "severity", -1)
	recovered := ginx.QueryInt(c, "is_recovered", -1)
	query := ginx.QueryStr(c, "query", "")
	limit := ginx.QueryInt(c, "limit", 20)
	busiGroupId := ginx.QueryInt64(c, "bgid", 0)
	dsIds := queryDatasourceIds(c)
	status := ginx.QueryInt(c, "status", 0) //添加status参数

	prod := ginx.QueryStr(c, "prods", "")
	if prod == "" {
		prod = ginx.QueryStr(c, "rule_prods", "")
	}

	prods := []string{}
	if prod != "" {
		prods = strings.Split(prod, ",")
	}

	cate := ginx.QueryStr(c, "cate", "$all")
	cates := []string{}
	if cate != "$all" {
		cates = strings.Split(cate, ",")
	}

	total, err := models.AlertHisEventTotal(rt.Ctx, prods, busiGroupId, stime, etime, severity, recovered, dsIds, cates, query, status)
	ginx.Dangerous(err)

	list, err := models.AlertHisEventGets(rt.Ctx, prods, busiGroupId, stime, etime, severity, recovered, dsIds, cates, query, limit, ginx.Offset(c, limit), status)
	ginx.Dangerous(err)

	cache := make(map[int64]*models.UserGroup)
	for i := 0; i < len(list); i++ {
		list[i].FillNotifyGroups(rt.Ctx, cache)
	}

	ginx.NewRender(c).Data(gin.H{
		"list":  list,
		"total": total,
	}, nil)
}

func (rt *Router) alertHisEventGet(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	event, err := models.AlertHisEventGetById(rt.Ctx, eid)
	ginx.Dangerous(err)

	if event == nil {
		ginx.Bomb(404, "No such alert event")
	}

	ginx.NewRender(c).Data(event, err)
}

/**
人工处理，已恢复
**/
func (rt *Router) alertHisEventSolve(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	event, err := models.AlertHisEventGetById(rt.Ctx, eid)
	ginx.Dangerous(err)

	if event == nil {
		ginx.Bomb(404, "No such alert event")
	} else {
		if event.IsRecovered == 0 {
			ginx.Dangerous("指标尚未恢复")
		} else {
			me := c.MustGet("user").(*models.User)
			var f bodyModel
			ginx.BindJSON(c, &f)
			err = event.UpdateStatus(rt.Ctx, eid, 1, f.Remark, me.Username)
			ginx.Dangerous(err)
		}
	}

	ginx.NewRender(c).Message(err)
}

type bodyModel struct {
	Remark string `json:"remark"`
}

/**
人工处理，已关闭
**/
func (rt *Router) alertHisEventClose(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	event, err := models.AlertHisEventGetById(rt.Ctx, eid)
	ginx.Dangerous(err)

	if event == nil {
		ginx.Bomb(404, "No such alert event")
	} else {
		me := c.MustGet("user").(*models.User)
		var f bodyModel
		ginx.BindJSON(c, &f)
		err = event.UpdateStatus(rt.Ctx, eid, 2, f.Remark, me.Username)

	}

	ginx.NewRender(c).Data(event, err)
}
