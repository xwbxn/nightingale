package router

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/ccfos/nightingale/v6/pkg/txt"
	xmltool "github.com/ccfos/nightingale/v6/pkg/xml"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func parseAggrRules(c *gin.Context) []*models.AggrRule {
	aggrRules := strings.Split(ginx.QueryStr(c, "rule", ""), "::") // e.g. field:group_name::field:severity::tagkey:ident

	if len(aggrRules) == 0 {
		ginx.Bomb(http.StatusBadRequest, "rule empty")
	}

	rules := make([]*models.AggrRule, len(aggrRules))
	for i := 0; i < len(aggrRules); i++ {
		pair := strings.Split(aggrRules[i], ":")
		if len(pair) != 2 {
			ginx.Bomb(http.StatusBadRequest, "rule invalid")
		}

		if !(pair[0] == "field" || pair[0] == "tagkey") {
			ginx.Bomb(http.StatusBadRequest, "rule invalid")
		}

		rules[i] = &models.AggrRule{
			Type:  pair[0],
			Value: pair[1],
		}
	}

	return rules
}

func (rt *Router) alertCurEventsCard(c *gin.Context) {
	stime, etime := getTimeRange(c)
	severity := ginx.QueryInt(c, "severity", -1)
	query := ginx.QueryStr(c, "query", "")
	busiGroupId := ginx.QueryInt64(c, "bgid", 0)
	dsIds := queryDatasourceIds(c)
	rules := parseAggrRules(c)

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

	// 最多获取50000个，获取太多也没啥意义
	list, err := models.AlertCurEventGets(rt.Ctx, prods, busiGroupId, stime, etime, severity, dsIds, cates, query, 50000, 0)
	ginx.Dangerous(err)

	cardmap := make(map[string]*AlertCard)
	for _, event := range list {
		title := event.GenCardTitle(rules)
		if _, has := cardmap[title]; has {
			cardmap[title].Total++
			cardmap[title].EventIds = append(cardmap[title].EventIds, event.Id)
			if event.Severity < cardmap[title].Severity {
				cardmap[title].Severity = event.Severity
			}
		} else {
			cardmap[title] = &AlertCard{
				Total:    1,
				EventIds: []int64{event.Id},
				Title:    title,
				Severity: event.Severity,
			}
		}
	}

	titles := make([]string, 0, len(cardmap))
	for title := range cardmap {
		titles = append(titles, title)
	}

	sort.Strings(titles)

	cards := make([]*AlertCard, len(titles))
	for i := 0; i < len(titles); i++ {
		cards[i] = cardmap[titles[i]]
	}

	sort.SliceStable(cards, func(i, j int) bool {
		if cards[i].Severity != cards[j].Severity {
			return cards[i].Severity < cards[j].Severity
		}
		return cards[i].Total > cards[j].Total
	})

	ginx.NewRender(c).Data(cards, nil)
}

type AlertCard struct {
	Title    string  `json:"title"`
	Total    int     `json:"total"`
	EventIds []int64 `json:"event_ids"`
	Severity int     `json:"severity"`
}

func (rt *Router) alertCurEventsCardDetails(c *gin.Context) {
	var f idsForm
	ginx.BindJSON(c, &f)

	list, err := models.AlertCurEventGetByIds(rt.Ctx, f.Ids)
	if err == nil {
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
	}

	ginx.NewRender(c).Data(list, err)
}

// alertCurEventsGetByRid
func (rt *Router) alertCurEventsGetByRid(c *gin.Context) {
	rid := ginx.QueryInt64(c, "rid")
	dsId := ginx.QueryInt64(c, "dsid")
	ginx.NewRender(c).Data(models.AlertCurEventGetByRuleIdAndDsId(rt.Ctx, rid, dsId))
}

// 列表方式，拉取活跃告警
func (rt *Router) alertCurEventsList(c *gin.Context) {
	stime, etime := getTimeRange(c)
	severity := ginx.QueryInt(c, "severity", -1)
	query := ginx.QueryStr(c, "query", "")
	limit := ginx.QueryInt(c, "limit", 20)
	busiGroupId := ginx.QueryInt64(c, "bgid", 0)
	dsIds := queryDatasourceIds(c)

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

	total, err := models.AlertCurEventTotal(rt.Ctx, prods, busiGroupId, stime, etime, severity, dsIds, cates, query)
	ginx.Dangerous(err)

	list, err := models.AlertCurEventGets(rt.Ctx, prods, busiGroupId, stime, etime, severity, dsIds, cates, query, limit, ginx.Offset(c, limit))
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

func (rt *Router) alertCurEventDel(c *gin.Context) {
	var f idsForm
	ginx.BindJSON(c, &f)
	f.Verify()

	set := make(map[int64]struct{})

	for i := 0; i < len(f.Ids); i++ {
		event, err := models.AlertCurEventGetById(rt.Ctx, f.Ids[i])
		ginx.Dangerous(err)

		if _, has := set[event.GroupId]; !has {
			rt.bgrwCheck(c, event.GroupId)
			set[event.GroupId] = struct{}{}
		}
	}

	ginx.NewRender(c).Message(models.AlertCurEventDel(rt.Ctx, f.Ids))
}

func (rt *Router) alertCurEventGet(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	event, err := models.AlertCurEventGetById(rt.Ctx, eid)
	ginx.Dangerous(err)

	if event == nil {
		ginx.Bomb(404, "No such active event")
	}

	if event.AssetId != 0 {
		asset, err := rt.assetCache.Get(event.AssetId)
		ginx.Dangerous(err)
		event.AssetName = asset.Name
		event.AssetIp = asset.Ip
	}

	ginx.NewRender(c).Data(event, nil)
}

// @Summary      当前告警过滤器
// @Description  当前告警过滤器
// @Tags         历史告警和当前告警-西航
// @Accept       json
// @Produce      json
// @Param        type query   int     false  "类型"
// @Param        start query   int     false  "开始时间"
// @Param        end query   int     false  "结束时间"
// @Param        query query   string     false  "搜索框"
// @Param        group query   int     false  "业务组id"
// @Param        limit query   int     false  "条数"
// @Param        page query   int     false  "页码"
// @Success      200  {array}  models.AlertHisEvent
// @Router       /api/n9e/alert-cur-events/list/xh [get]
// @Security     ApiKeyAuth
func (rt *Router) alertCurEventsListXH(c *gin.Context) {
	// stime, etime := getTimeRange(c)

	fType := ginx.QueryInt64(c, "type", -1)
	start := ginx.QueryInt64(c, "start", -1)
	end := ginx.QueryInt64(c, "end", -1)
	query := ginx.QueryStr(c, "query", "")
	group := ginx.QueryInt64(c, "group", -1)
	limit := ginx.QueryInt(c, "limit", 20)
	page := ginx.QueryInt(c, "page", 1)
	ids := make([]int64, 0)

	if query != "" {
		assets := rt.assetCache.GetAll()
		for _, asset := range assets {
			if fType == -1 {
				if strings.Contains(asset.Name, query) || strings.Contains(asset.Type, query) || strings.Contains(asset.Ip, query) {
					ids = append(ids, asset.Id)
				}
			} else if fType == 2 {
				if strings.Contains(asset.Type, query) {
					ids = append(ids, asset.Id)
				}
			}
		}
	}

	total, err := models.AlertCurEventXHTotal(rt.Ctx, fType, start, end, group, ids, query)
	ginx.Dangerous(err)

	list, err := models.AlertCurEventXHGets(rt.Ctx, fType, start, end, group, ids, query, limit, (page-1)*limit)
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

// @Summary      批量删除当前告警
// @Description  批量删除当前告警
// @Tags         历史告警和当前告警-西航
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200
// @Router       /api/n9e/alert-cur-events/batch-del [post]
// @Security     ApiKeyAuth
func (rt *Router) alertCurEventBatchDel(c *gin.Context) {
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
	tx := models.DB(rt.Ctx).Begin()

	err := models.AlertCurEventDelByIdsTx(tx, ids)
	ginx.Dangerous(err)
	err = models.AlertHisEventDelByIdsTx(tx, ids)

	ginx.NewRender(c).Message(err)
}

// @Summary      EXCEL导出当前告警
// @Description  EXCEL导出当前告警
// @Tags         历史告警和当前告警-西航
// @Accept       json
// @Produce      json
// @Param        ftype query   int     false  "类型"
// @Param        type query   int     false  "类型"
// @Param        start query   int     false  "开始时间"
// @Param        end query   int     false  "结束时间"
// @Param        query query   string     false  "搜索框"
// @Param        group query   int     false  "业务组id"
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200
// @Router       /api/n9e/alert-cur-events/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) exportCurEventXH(c *gin.Context) {

	fileType := ginx.QueryInt64(c, "ftype", -1)

	var f map[string]interface{}
	ginx.BindJSON(c, &f)
	var list []models.AlertCurEvent
	var err error

	idsTemp, idsOk := f["ids"]
	ids := make([]int64, 0)
	// var err error
	if idsOk {
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, int64(val.(float64)))
		}
		lstP, err := models.AlertCurEventGetByIds(rt.Ctx, ids)
		ginx.Dangerous(err)
		for _, val := range lstP {
			list = append(list, *val)
		}
	} else {
		fType := ginx.QueryInt64(c, "type", -1)
		start := ginx.QueryInt64(c, "start", -1)
		end := ginx.QueryInt64(c, "end", -1)
		query := ginx.QueryStr(c, "query", "")
		group := ginx.QueryInt64(c, "group", -1)

		if query != "" {
			assets := rt.assetCache.GetAll()
			for _, asset := range assets {
				if fType == -1 {
					if strings.Contains(asset.Name, query) || strings.Contains(asset.Type, query) || strings.Contains(asset.Ip, query) {
						ids = append(ids, asset.Id)
					}
				} else if fType == 2 {
					if strings.Contains(asset.Type, query) {
						ids = append(ids, asset.Id)
					}
				}
			}
		}

		list, err = models.AlertCurEventXHGets(rt.Ctx, fType, start, end, group, ids, query, -1, -1)
		ginx.Dangerous(err)

	}
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
	alterlst := make([]models.AlterImport, 0)
	for _, val := range list {
		alertRule, err := models.AlertRuleGetById(rt.Ctx, val.RuleId)
		ginx.Dangerous(err)
		promSql := ""
		if alertRule != nil {
			promSql = alertRule.RuleConfigCn
		}
		alterlst = append(alterlst, models.AlterImport{
			RuleName:         val.RuleName,
			AssetName:        val.AssetName,
			AssetIp:          val.AssetIp,
			Severity:         strconv.Itoa(val.Severity),
			TriggerTime:      strconv.FormatInt(val.TriggerTime, 10),
			TriggerValue:     val.TriggerValue,
			PromSql:          promSql,
			PromEvalInterval: val.PromEvalInterval,
			PromForDuration:  val.PromForDuration,
		})
	}
	datas := make([]interface{}, 0)
	if len(list) > 0 {
		for _, v := range alterlst {
			datas = append(datas, v)

		}
	}
	if fileType == 1 {
		excels.NewMyExcel("当前告警信息").ExportDataInfo(datas, "cn", rt.Ctx, c)
	} else if fileType == 2 {
		for index, alter := range alterlst {
			if alter.Severity == "1" {
				alterlst[index].Severity = "紧急"
			} else if alter.Severity == "2" {
				alterlst[index].Severity = "一般"
			} else if alter.Severity == "3" {
				alterlst[index].Severity = "事件"
			}
			trigger_time, _ := strconv.ParseInt(alter.TriggerTime, 10, 64)
			alterlst[index].TriggerTime = time.Unix(trigger_time, 0).Format("2006-01-02 15:04:05")
		}

		dataXml := models.AlterXml{
			AlterData: alterlst,
		}
		xmltool.ExportXml(c, dataXml)
	} else if fileType == 3 {
		dataTxt := make([]string, 0)
		str := "规则标题\t资产名称\t资产IP\t告警级别\t指标\t触发时间\t触发时值\t执行频率/s\t持续时长/s\n"
		dataTxt = append(dataTxt, str)
		for _, alter := range alterlst {
			severity := ""
			if alter.Severity == "1" {
				severity = "紧急"
			} else if alter.Severity == "2" {
				severity = "一般"
			} else if alter.Severity == "3" {
				severity = "事件"
			}
			// time.Parse("2006-01-02 15:04:05", alter.TriggerTime)
			// time.Unix(alter.TriggerTime, 0).Format("2006-01-02 15:04:05")
			trigger_time, _ := strconv.ParseInt(alter.TriggerTime, 10, 64)
			str = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\n",
				alter.RuleName, alter.AssetName, alter.AssetIp, severity, alter.PromSql,
				time.Unix(trigger_time, 0).Format("2006-01-02 15:04:05"), alter.TriggerValue,
				alter.PromEvalInterval, alter.PromForDuration)
			dataTxt = append(dataTxt, str)
		}
		txt.ExportTxt(c, dataTxt)
	}

}
