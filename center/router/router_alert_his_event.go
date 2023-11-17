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
		if list[i].AssetId != 0 {
			asset, err := rt.assetCache.Get(list[i].AssetId)
			ginx.Dangerous(err)
			list[i].AssetName = asset.Name
			list[i].AssetIp = asset.Ip
		}
	}

	ginx.NewRender(c).Data(gin.H{
		"list":  list,
		"total": total,
	}, nil)
}

// @Summary      历史告警过滤器
// @Description  历史告警过滤器
// @Tags         历史告警和当前告警-西航
// @Accept       json
// @Produce      json
// @Param        severity query   int     false  "告警等级"
// @Param        query query   string     false  "搜索框"
// @Param        group query   int     false  "业务组id"
// @Param        limit query   int     false  "条数"
// @Param        page query   int     false  "页码"
// @Success      200  {array}  models.AlertHisEvent
// @Router       /api/n9e/alert-his-events/list/xh [get]
// @Security     ApiKeyAuth
// func (rt *Router) alertHisEventsListXH(c *gin.Context) {
// 	// stime, etime := getTimeRange(c)

// 	group := ginx.QueryInt64(c, "group", -1)
// 	query := ginx.QueryStr(c, "query", "")
// 	severity := ginx.QueryInt64(c, "severity", -1)
// 	limit := ginx.QueryInt(c, "limit", 20)
// 	page := ginx.QueryInt(c, "page", 1)
// 	ids := make([]int64, 0)

// 	if query != "" {
// 		assets := rt.assetCache.GetAll()
// 		for _, asset := range assets {
// 			// if fType == -1 {
// 			if strings.Contains(asset.Name, query) || strings.Contains(asset.Type, query) || strings.Contains(asset.Ip, query) {
// 				ids = append(ids, asset.Id)
// 			}
// 			// } else if fType == 2 {
// 			// 	if strings.Contains(asset.Type, query) {
// 			// 		ids = append(ids, asset.Id)
// 			// 	}
// 			// }
// 		}
// 	}

// 	total, err := models.AlertHisEventXHTotal(rt.Ctx, severity, group, ids, query)
// 	ginx.Dangerous(err)

// 	list, err := models.AlertHisEventXHGets(rt.Ctx, severity, group, ids, query, limit, (page-1)*limit)
// 	ginx.Dangerous(err)

// 	cache := make(map[int64]*models.UserGroup)
// 	for i := 0; i < len(list); i++ {
// 		list[i].FillNotifyGroups(rt.Ctx, cache)
// 		if list[i].AssetId != 0 {
// 			asset, err := rt.assetCache.Get(list[i].AssetId)
// 			ginx.Dangerous(err)
// 			list[i].AssetName = asset.Name
// 			list[i].AssetIp = asset.Ip
// 		}
// 	}

// 	ginx.NewRender(c).Data(gin.H{
// 		"list":  list,
// 		"total": total,
// 	}, nil)
// }

func (rt *Router) alertHisEventGet(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	event, err := models.AlertHisEventGetById(rt.Ctx, eid)
	ginx.Dangerous(err)

	if event == nil {
		ginx.Bomb(404, "No such alert event")
	}

	if event.AssetId != 0 {
		asset, err := rt.assetCache.Get(event.AssetId)
		ginx.Dangerous(err)
		event.AssetName = asset.Name
		event.AssetIp = asset.Ip
	}

	ginx.NewRender(c).Data(event, err)
}

/**
人工处理，已恢复
**/
func (rt *Router) alertHisEventSolve(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	event, err := models.AlertHisEventGetById(rt.Ctx, eid)
	me := c.MustGet("user").(*models.User)
	var f bodyModel
	ginx.BindJSON(c, &f)
	ginx.Dangerous(err)

	if event == nil {
		ginx.Bomb(404, "No such alert event")
	} else {
		if event.IsRecovered == 0 {
			ginx.Dangerous("指标尚未恢复")
		} else {
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
	me := c.MustGet("user").(*models.User)
	ginx.Dangerous(err)
	var f bodyModel
	ginx.BindJSON(c, &f)

	if event == nil {
		ginx.Bomb(404, "No such alert event")
	} else {
		err = event.UpdateStatus(rt.Ctx, eid, 2, f.Remark, me.Username)
	}
	ginx.NewRender(c).Data(event, err)
}

// @Summary      批量删除历史告警
// @Description  批量删除历史告警
// @Tags         历史告警和当前告警-西航
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200
// @Router       /api/n9e/alert-his-events/batch-del [post]
// @Security     ApiKeyAuth
func (rt *Router) alertHisEventBatchDel(c *gin.Context) {
	var f map[string]interface{}
	ginx.BindJSON(c, &f)

	idsTemp, idsOk := f["ids"]
	ids := make([]int64, 0)
	// var err error
	if idsOk {
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, int64(val.(float64)))
		}
	}

	err := models.AlertHisEventDelByIds(rt.Ctx, ids)

	ginx.NewRender(c).Message(err)
}

// @Summary      EXCEL导出历史告警
// @Description  EXCEL导出历史告警
// @Tags         历史告警和当前告警-西航
// @Accept       json
// @Produce      json
// @Param        severity query   int     false  "告警等级"
// @Param        query query   string     false  "搜索框"
// @Param        group query   int     false  "业务组id"
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200
// @Router       /api/n9e/alert-his-events/export-xls [post]
// @Security     ApiKeyAuth
// func (rt *Router) exportHisEventXH(c *gin.Context) {

// 	fileType := ginx.QueryInt64(c, "ftype", -1)

// 	var f map[string]interface{}
// 	ginx.BindJSON(c, &f)
// 	var list []models.AlertHisEvent
// 	var err error

// 	idsTemp, idsOk := f["ids"]
// 	ids := make([]int64, 0)
// 	// var err error
// 	if idsOk {
// 		for _, val := range idsTemp.([]interface{}) {
// 			ids = append(ids, int64(val.(float64)))
// 		}
// 		list, err = models.AlertHisEventGetByIds(rt.Ctx, ids)
// 		ginx.Dangerous(err)
// 	} else {
// 		group := ginx.QueryInt64(c, "group", -1)
// 		query := ginx.QueryStr(c, "query", "")
// 		severity := ginx.QueryInt64(c, "severity", -1)
// 		limit := ginx.QueryInt(c, "limit", 20)
// 		page := ginx.QueryInt(c, "page", 1)
// 		ids := make([]int64, 0)

// 		if query != "" {
// 			assets := rt.assetCache.GetAll()
// 			for _, asset := range assets {
// 				if strings.Contains(asset.Name, query) || strings.Contains(asset.Type, query) || strings.Contains(asset.Ip, query) {
// 					ids = append(ids, asset.Id)
// 				}
// 			}
// 		}

// 		list, err = models.AlertEventXHGets[models.AlertHisEvent](rt.Ctx, 1, severity, group, ids, query, limit, (page-1)*limit)
// 		ginx.Dangerous(err)

// 	}
// 	cache := make(map[int64]*models.UserGroup)
// 	for i := 0; i < len(list); i++ {
// 		list[i].FillNotifyGroups(rt.Ctx, cache)
// 		if list[i].AssetId != 0 {
// 			asset, err := rt.assetCache.Get(list[i].AssetId)
// 			ginx.Dangerous(err)
// 			list[i].AssetName = asset.Name
// 			list[i].AssetIp = asset.Ip
// 		}
// 	}
// 	alterlst := make([]models.AlterImport, 0)
// 	for _, val := range list {
// 		alertRule, err := models.AlertRuleGetById(rt.Ctx, val.RuleId)

// 		promSql := ""
// 		if alertRule != nil {
// 			promSql = alertRule.RuleConfigCn
// 		}
// 		ginx.Dangerous(err)
// 		alterlst = append(alterlst, models.AlterImport{
// 			RuleName:         val.RuleName,
// 			AssetName:        val.AssetName,
// 			AssetIp:          val.AssetIp,
// 			Severity:         strconv.Itoa(val.Severity),
// 			TriggerTime:      strconv.FormatInt(val.TriggerTime, 10),
// 			TriggerValue:     val.TriggerValue,
// 			PromSql:          promSql,
// 			PromEvalInterval: val.PromEvalInterval,
// 			PromForDuration:  val.PromForDuration,
// 		})
// 	}
// 	datas := make([]interface{}, 0)
// 	if len(list) > 0 {
// 		for _, v := range alterlst {
// 			datas = append(datas, v)

// 		}
// 	}
// 	if fileType == 1 {
// 		excels.NewMyExcel("当前告警信息").ExportDataInfo(datas, "cn", rt.Ctx, c)
// 	} else if fileType == 2 {
// 		for index, alter := range alterlst {
// 			if alter.Severity == "1" {
// 				alterlst[index].Severity = "紧急"
// 			} else if alter.Severity == "2" {
// 				alterlst[index].Severity = "一般"
// 			} else if alter.Severity == "3" {
// 				alterlst[index].Severity = "事件"
// 			}
// 			trigger_time, _ := strconv.ParseInt(alter.TriggerTime, 10, 64)
// 			alterlst[index].TriggerTime = time.Unix(trigger_time, 0).Format("2006-01-02 15:04:05")
// 		}

// 		dataXml := models.AlterXml{
// 			AlterData: alterlst,
// 		}
// 		xmltool.ExportXml(c, dataXml, "历史告警")
// 	} else if fileType == 3 {
// 		dataTxt := make([]string, 0)
// 		str := "规则标题\t资产名称\t资产IP\t告警级别\t指标\t触发时间\t触发时值\t执行频率/s\t持续时长/s\n"
// 		dataTxt = append(dataTxt, str)
// 		for _, alter := range alterlst {
// 			severity := ""
// 			if alter.Severity == "1" {
// 				severity = "紧急"
// 			} else if alter.Severity == "2" {
// 				severity = "一般"
// 			} else if alter.Severity == "3" {
// 				severity = "事件"
// 			}
// 			// time.Parse("2006-01-02 15:04:05", alter.TriggerTime)
// 			// time.Unix(alter.TriggerTime, 0).Format("2006-01-02 15:04:05")
// 			trigger_time, _ := strconv.ParseInt(alter.TriggerTime, 10, 64)
// 			str = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\n",
// 				alter.RuleName, alter.AssetName, alter.AssetIp, severity, alter.PromSql,
// 				time.Unix(trigger_time, 0).Format("2006-01-02 15:04:05"), alter.TriggerValue,
// 				alter.PromEvalInterval, alter.PromForDuration)
// 			dataTxt = append(dataTxt, str)
// 		}
// 		txt.ExportTxt(c, dataTxt, "历史告警")
// 	}

// }
