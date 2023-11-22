package router

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/models"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/i18n"
)

// Return all, front-end search and paging
func (rt *Router) alertRuleGets(c *gin.Context) {
	busiGroupId := ginx.UrlParamInt64(c, "id")
	ars, err := models.AlertRuleGets(rt.Ctx, busiGroupId)
	if err == nil {
		cache := make(map[int64]*models.UserGroup)
		for i := 0; i < len(ars); i++ {
			ars[i].FillNotifyGroups(rt.Ctx, cache)
			ars[i].FillSeverities()
			if ars[i].AssetId != 0 {
				asset, err := rt.assetCache.Get(ars[i].AssetId)
				ginx.Dangerous(err)
				ars[i].AssetName = asset.Name
				ars[i].AssetIp = asset.Ip
			}
		}
	}
	ginx.NewRender(c).Data(ars, err)
}

// @Summary      告警规则过滤器
// @Description  告警规则过滤器
// @Tags         告警规则
// @Accept       json
// @Produce      json
// @Param        busiGroupId query   int     false  "busiGroupId"
// @Param        filter query   int     false  "告警等级（1：告警级别；2：资产类型；3：资产id）"
// @Param        query query   string     false  "搜索框"
// @Param        limit query   int     false  "条数"
// @Param        page query   int     false  "页码"
// @Success      200  {array}  models.AlertHisEvent
// @Router       /api/n9e/busi-group/{id}/alert-rules/xh [get]
// @Security     ApiKeyAuth
func (rt *Router) alertRuleGetsXH(c *gin.Context) {
	busiGroupId := ginx.UrlParamInt64(c, "id")
	filter := ginx.QueryStr(c, "filter", "")
	query := ginx.QueryStr(c, "query", "")
	limit := ginx.QueryInt(c, "limit", 20)
	page := ginx.QueryInt(c, "page", 1)
	ids := make([]int64, 0)
	if filter != "" && query == "" {
		ginx.Bomb(http.StatusOK, "参数错误")
	}

	if query != "" {
		assets := rt.assetCache.GetAll()
		if filter == "" {
			for _, asset := range assets {
				if strings.Contains(asset.Name, query) || strings.Contains(asset.Type, query) || strings.Contains(asset.Ip, query) {
					ids = append(ids, asset.Id)
				}
			}
		} else if filter == "2" {
			for _, asset := range assets {
				if strings.Contains(asset.Type, query) {
					ids = append(ids, asset.Id)
				}
			}
		}
	}
	total, err := models.AlertRuleGetsTotal(rt.Ctx, busiGroupId, filter, query, ids)
	ginx.Dangerous(err)

	ars, err := models.AlertRuleGetsFilter(rt.Ctx, busiGroupId, filter, query, ids, limit, (page - 1*limit))
	ginx.Dangerous(err)
	if err == nil {
		cache := make(map[int64]*models.UserGroup)
		for i := 0; i < len(ars); i++ {
			ars[i].FillNotifyGroups(rt.Ctx, cache)
			ars[i].FillSeverities()
			if ars[i].AssetId != 0 {
				asset, err := rt.assetCache.Get(ars[i].AssetId)
				ginx.Dangerous(err)
				ars[i].AssetName = asset.Name
				ars[i].AssetIp = asset.Ip
			}
		}
	}
	ginx.NewRender(c).Data(gin.H{
		"list":  ars,
		"total": total,
	}, err)
}

func (rt *Router) alertRulesGetByService(c *gin.Context) {
	prods := []string{}
	prodStr := ginx.QueryStr(c, "prods", "")
	if prodStr != "" {
		prods = strings.Split(ginx.QueryStr(c, "prods", ""), ",")
	}

	query := ginx.QueryStr(c, "query", "")
	algorithm := ginx.QueryStr(c, "algorithm", "")
	cluster := ginx.QueryStr(c, "cluster", "")
	cate := ginx.QueryStr(c, "cate", "$all")
	cates := []string{}
	if cate != "$all" {
		cates = strings.Split(cate, ",")
	}

	disabled := ginx.QueryInt(c, "disabled", -1)
	ars, err := models.AlertRulesGetsBy(rt.Ctx, prods, query, algorithm, cluster, cates, disabled)
	if err == nil {
		cache := make(map[int64]*models.UserGroup)
		for i := 0; i < len(ars); i++ {
			ars[i].FillNotifyGroups(rt.Ctx, cache)
		}
	}
	ginx.NewRender(c).Data(ars, err)
}

// single or import
func (rt *Router) alertRuleAddByFE(c *gin.Context) {
	username := c.MustGet("username").(string)

	var lst []models.AlertRule
	ginx.BindJSON(c, &lst)

	count := len(lst)
	if count == 0 {
		ginx.Bomb(http.StatusBadRequest, "input json is empty")
	}

	bgid := ginx.UrlParamInt64(c, "id")
	reterr := rt.alertRuleAdd(lst, username, bgid, c.GetHeader("X-Language"))

	ginx.NewRender(c).Data(reterr, nil)
}

func (rt *Router) alertRuleAddByImport(c *gin.Context) {
	username := c.MustGet("username").(string)

	var lst []models.AlertRule
	ginx.BindJSON(c, &lst)

	count := len(lst)
	if count == 0 {
		ginx.Bomb(http.StatusBadRequest, "input json is empty")
	}

	bgid := ginx.UrlParamInt64(c, "id")
	reterr := rt.alertRuleAdd(lst, username, bgid, c.GetHeader("X-Language"))

	ginx.NewRender(c).Data(reterr, nil)
}

func (rt *Router) alertRuleAddByService(c *gin.Context) {
	var lst []models.AlertRule
	ginx.BindJSON(c, &lst)

	count := len(lst)
	if count == 0 {
		ginx.Bomb(http.StatusBadRequest, "input json is empty")
	}
	reterr := rt.alertRuleAddForService(lst, "")
	ginx.NewRender(c).Data(reterr, nil)
}

func (rt *Router) alertRuleAddForService(lst []models.AlertRule, username string) map[string]string {
	count := len(lst)
	// alert rule name -> error string
	reterr := make(map[string]string)
	for i := 0; i < count; i++ {
		lst[i].Id = 0
		if username != "" {
			lst[i].CreateBy = username
			lst[i].UpdateBy = username
		}

		if err := lst[i].FE2DB(rt.Ctx); err != nil {
			reterr[lst[i].Name] = err.Error()
			continue
		}

		if err := lst[i].Add(rt.Ctx); err != nil {
			reterr[lst[i].Name] = err.Error()
		} else {
			reterr[lst[i].Name] = ""
		}
	}
	return reterr
}

func (rt *Router) alertRuleAdd(lst []models.AlertRule, username string, bgid int64, lang string) map[string]string {
	count := len(lst)
	// alert rule name -> error string
	reterr := make(map[string]string)
	for i := 0; i < count; i++ {
		lst[i].Id = 0
		lst[i].GroupId = bgid
		if username != "" {
			lst[i].CreateBy = username
			lst[i].UpdateBy = username
		}

		if err := lst[i].FE2DB(rt.Ctx); err != nil {
			reterr[lst[i].Name] = i18n.Sprintf(lang, err.Error())
			continue
		}

		if err := lst[i].Add(rt.Ctx); err != nil {
			reterr[lst[i].Name] = i18n.Sprintf(lang, err.Error())
		} else {
			reterr[lst[i].Name] = ""
		}
	}
	return reterr
}

func (rt *Router) alertRuleDel(c *gin.Context) {
	var f idsForm
	ginx.BindJSON(c, &f)
	f.Verify()

	// param(busiGroupId) for protect
	ginx.NewRender(c).Message(models.AlertRuleDels(rt.Ctx, f.Ids, ginx.UrlParamInt64(c, "id")))
}

func (rt *Router) alertRuleDelByService(c *gin.Context) {
	var f idsForm
	ginx.BindJSON(c, &f)
	f.Verify()
	ginx.NewRender(c).Message(models.AlertRuleDels(rt.Ctx, f.Ids))
}

func (rt *Router) alertRulePutByFE(c *gin.Context) {
	var f models.AlertRule
	ginx.BindJSON(c, &f)

	arid := ginx.UrlParamInt64(c, "arid")
	ar, err := models.AlertRuleGetById(rt.Ctx, arid)
	ginx.Dangerous(err)

	if ar == nil {
		ginx.NewRender(c, http.StatusNotFound).Message("No such AlertRule")
		return
	}

	rt.bgrwCheck(c, ar.GroupId)

	f.UpdateBy = c.MustGet("username").(string)
	ginx.NewRender(c).Message(ar.Update(rt.Ctx, f))
}

func (rt *Router) alertRulePutByService(c *gin.Context) {
	var f models.AlertRule
	ginx.BindJSON(c, &f)

	arid := ginx.UrlParamInt64(c, "arid")
	ar, err := models.AlertRuleGetById(rt.Ctx, arid)
	ginx.Dangerous(err)

	if ar == nil {
		ginx.NewRender(c, http.StatusNotFound).Message("No such AlertRule")
		return
	}
	ginx.NewRender(c).Message(ar.Update(rt.Ctx, f))
}

type alertRuleFieldForm struct {
	Ids    []int64                `json:"ids"`
	Fields map[string]interface{} `json:"fields"`
	Action string                 `json:"action"`
}

// update one field: cluster note severity disabled prom_eval_interval prom_for_duration notify_channels notify_groups notify_recovered notify_repeat_step callbacks runbook_url append_tags
func (rt *Router) alertRulePutFields(c *gin.Context) {
	var f alertRuleFieldForm
	ginx.BindJSON(c, &f)

	if len(f.Fields) == 0 {
		ginx.Bomb(http.StatusBadRequest, "fields empty")
	}

	f.Fields["update_by"] = c.MustGet("username").(string)
	f.Fields["update_at"] = time.Now().Unix()

	for i := 0; i < len(f.Ids); i++ {
		ar, err := models.AlertRuleGetById(rt.Ctx, f.Ids[i])
		ginx.Dangerous(err)

		if ar == nil {
			continue
		}

		if f.Action == "callback_add" {
			// 增加一个 callback 地址
			if callbacks, has := f.Fields["callbacks"]; has {
				callback := callbacks.(string)
				if !strings.Contains(ar.Callbacks, callback) {
					ginx.Dangerous(ar.UpdateFieldsMap(rt.Ctx, map[string]interface{}{"callbacks": ar.Callbacks + " " + callback}))
					continue
				}
			}
		}

		if f.Action == "callback_del" {
			// 删除一个 callback 地址
			if callbacks, has := f.Fields["callbacks"]; has {
				callback := callbacks.(string)
				ginx.Dangerous(ar.UpdateFieldsMap(rt.Ctx, map[string]interface{}{"callbacks": strings.ReplaceAll(ar.Callbacks, callback, "")}))
				continue
			}
		}

		for k, v := range f.Fields {
			ginx.Dangerous(ar.UpdateColumn(rt.Ctx, k, v))
		}
	}

	ginx.NewRender(c).Message(nil)
}

func (rt *Router) alertRuleGet(c *gin.Context) {
	arid := ginx.UrlParamInt64(c, "arid")

	ar, err := models.AlertRuleGetById(rt.Ctx, arid)
	ginx.Dangerous(err)

	if ar == nil {
		ginx.NewRender(c, http.StatusNotFound).Message("No such AlertRule")
		return
	}

	err = ar.FillNotifyGroups(rt.Ctx, make(map[int64]*models.UserGroup))
	ginx.Dangerous(err)

	if ar.AssetId != 0 {
		asset, err := rt.assetCache.Get(ar.AssetId)
		ginx.Dangerous(err)
		ar.AssetName = asset.Name
		ar.AssetIp = asset.Ip
	}

	ginx.NewRender(c).Data(ar, err)
}

func (rt *Router) alertRuleGetsByIds(c *gin.Context) {

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

	ars, err := models.AlertRuleGetByIds(rt.Ctx, ids)
	ginx.Dangerous(err)

	for index := range ars {
		err = ars[index].FillNotifyGroups(rt.Ctx, make(map[int64]*models.UserGroup))
		ginx.Dangerous(err)

		if ars[index].AssetId != 0 {
			asset, err := rt.assetCache.Get(ars[index].AssetId)
			ginx.Dangerous(err)
			ars[index].AssetName = asset.Name
			ars[index].AssetIp = asset.Ip
		}
	}

	ginx.NewRender(c).Data(ars, err)
}

// pre validation before save rule
func (rt *Router) alertRuleValidation(c *gin.Context) {
	var f models.AlertRule //new
	ginx.BindJSON(c, &f)

	arid := ginx.UrlParamInt64(c, "arid")
	ar, err := models.AlertRuleGetById(rt.Ctx, arid)
	ginx.Dangerous(err)

	if ar == nil {
		ginx.NewRender(c, http.StatusNotFound).Message("No such AlertRule")
		return
	}

	rt.bgrwCheck(c, ar.GroupId)

	if len(f.NotifyChannelsJSON) > 0 && len(f.NotifyGroupsJSON) > 0 { //Validation NotifyChannels
		ngids := make([]int64, 0, len(f.NotifyChannelsJSON))
		for i := range f.NotifyGroupsJSON {
			id, _ := strconv.ParseInt(f.NotifyGroupsJSON[i], 10, 64)
			ngids = append(ngids, id)
		}
		userGroups := rt.UserGroupCache.GetByUserGroupIds(ngids)
		uids := make([]int64, 0)
		for i := range userGroups {
			uids = append(uids, userGroups[i].UserIds...)
		}
		users := rt.UserCache.GetByUserIds(uids)
		//If any users have a certain notify channel's token, it will be okay. Otherwise, this notify channel is absent of tokens.
		ancs := make([]string, 0, len(f.NotifyChannelsJSON)) //absent Notify Channels
		for i := range f.NotifyChannelsJSON {
			flag := true
			for ui := range users {
				if _, b := users[ui].ExtractToken(f.NotifyChannelsJSON[i]); b {
					flag = false
					break
				}
			}
			if flag {
				ancs = append(ancs, f.NotifyChannelsJSON[i])
			}
		}

		if len(ancs) > 0 {
			ginx.NewRender(c).Message(i18n.Sprintf(c.GetHeader("X-Language"), "All users are missing notify channel configurations. Please check for missing tokens (each channel should be configured with at least one user). %s", ancs))
			return
		}

	}

	ginx.NewRender(c).Message("")
}
