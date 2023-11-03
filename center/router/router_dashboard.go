package router

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/prom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

type ChartPure struct {
	Configs string `json:"configs"`
	Weight  int    `json:"weight"`
}

type ChartGroupPure struct {
	Name   string      `json:"name"`
	Weight int         `json:"weight"`
	Charts []ChartPure `json:"charts"`
}

type DashboardPure struct {
	Name        string           `json:"name"`
	Tags        string           `json:"tags"`
	Configs     string           `json:"configs"`
	ChartGroups []ChartGroupPure `json:"chart_groups"`
}

type MetricJson struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

/***
  对外提供的资产信息列表
  author:guoxp
*/

type AssetJson struct {
	Id       int64               `json:"id"`
	Name     string              `json:"name"`
	Ip       string              `json:"ip"`
	Label    string              `json:"label"`
	Status   int64               `json:"status"`
	UpdateAt int64               `json:"update_at"`
	Category string              `json:"category"`
	Type     string              `json:"type"`
	Metrics  []map[string]string `json:"metrics"`
	GroupId  int64               `json:"group_id"`
	Tags     []string            `json:"tags"`
	Url      string              `json:"url"`
	Severity int64               `json:"severity"`
}

func (rt *Router) getDashboardAssetsByFE(c *gin.Context) {

	category := ginx.QueryStr(c, "category", "")
	atype := ginx.QueryStr(c, "atype", "")
	groupId := ginx.QueryInt64(c, "group_id", -1)

	var data []*AssetJson
	lst := rt.assetCache.GetAll()
	for _, item := range lst {
		ar, _ := rt.assetCache.GetType(item.Type)

		if category != "" && ar.Category != category {
			continue
		}
		if atype != "" && item.Type != atype {
			continue
		}
		if groupId > -1 && item.OrganizationId != groupId { //这里使用orgid作为group返回查询条件
			continue
		}

		metrics := []map[string]string{}
		for _, m := range item.Metrics {
			metrics = append(metrics, m)
		}

		//将名称中含有的IP拿出来
		ip := ""
		name := item.Name
		compileRegex := regexp.MustCompile("（(.*?)）") // 中文括号，例如：华南地区（广州） -> 广州
		matchArr := compileRegex.FindStringSubmatch(item.Name)

		if len(matchArr) == 0 {
			compileRegex = regexp.MustCompile("\\((.*?)\\)") // 兼容英文括号并取消括号的转义，例如：华南地区 (广州) -> 广州。
			matchArr = compileRegex.FindStringSubmatch(item.Name)
		}
		if len(matchArr) != 0 {
			ipTemp := matchArr[len(matchArr)-1]
			if len(strings.Split(ipTemp, ".")) == 4 {
				ip = ipTemp
				nameTemp := strings.Split(item.Name, matchArr[0])
				name = nameTemp[0]
			}
		}
		data = append(data, &AssetJson{
			Id:       item.Id,
			Name:     name,
			Ip:       ip,
			Label:    item.Label,
			Status:   item.Health,
			UpdateAt: item.HealthAt,
			Category: ar.Category,
			Type:     item.Type,
			Metrics:  metrics,
			GroupId:  item.OrganizationId, //这里使用orgid作为group返回查询条件
			Tags:     item.TagsJSON,
		})
	}
	//ws.SetMessage(1, data) //socket推送内容

	ginx.NewRender(c).Data(data, nil)
}

func (rt *Router) getOrganizationTreeByFE(c *gin.Context) {
	list, err := models.OrganizationTreeGetsFE(rt.Ctx)
	ginx.Dangerous(err)
	ginx.NewRender(c).Data(list, nil)
}

// @Summary      告警列表接口前端接口返回
// @Description  告警列表接口前端接口返回
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Success      200  {array}  []models.FeAlert
// @Router       /api/n9e/dashboard/alert-cur-events/ceshi [get]
// @Security     ApiKeyAuth
func (rt *Router) getAlertListByFE(c *gin.Context) {
	list, err := models.AlertFeList(rt.Ctx)

	//从资产缓存中更新orgid
	for _, item := range list {
		asset, has := rt.assetCache.Get(int64(item.AssetId))
		if has {
			item.AssetName = asset.Name
			item.OrganizeId = int(asset.OrganizationId)
		}
	}

	ginx.NewRender(c).Data(list, err)
}

func (rt *Router) getDashboardAssetStatistics(c *gin.Context) {
	key := ginx.QueryStr(c, "key", "type")

	data := map[string]int64{}
	lst := rt.assetCache.GetAll()

	if key == "type" {
		sort.Slice(lst, func(i, j int) bool {
			return lst[i].Type > lst[j].Type
		})

		for _, item := range lst {
			_, has := data[item.Type]
			if !has {
				data[item.Type] = 0
			}
			data[item.Type] += 1
		}
	}

	ginx.NewRender(c).Data(data, nil)
}

// @Summary      首页数据统计接口
// @Description  首页数据统计接口
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api/n9e/dashboard/count/ceshi [get]
// @Security     ApiKeyAuth
func (rt *Router) getDashboardDataCount(c *gin.Context) {
	logger.Debug("kaishitongji")

	m := make(map[string]interface{})

	//安全运行时长
	dictData, err := models.DictDataGetByMap(rt.Ctx, map[string]interface{}{"type_code": "system_initialize", "dict_key": "uptime"})
	ginx.Dangerous(err)
	start, err := strconv.ParseInt(dictData[0].DictValue, 10, 64)
	ginx.Dangerous(err)
	runtime := time.Now().Unix() - start
	m["runtime"] = runtime

	//告警统计
	alarm := make(map[string]interface{})
	current, err := models.AlertCurCount(rt.Ctx)
	ginx.Dangerous(err)
	unprocessed, err := models.AlertHisCount(rt.Ctx)
	ginx.Dangerous(err)
	alarm["current"] = current
	alarm["unprocessed"] = unprocessed
	//历史告警数
	hisAll, err := models.AlertHisCountMap(rt.Ctx, map[string]interface{}{})
	ginx.Dangerous(err)
	alarm["history"] = hisAll
	alarm["confirm"] = hisAll - unprocessed
	m["alarm"] = alarm

	//资产分类
	assetType := make(map[string]int64)
	//TODO 排序规则未给出
	// deviceSort := make(map[string]int64, 0)
	// appAndCloudSort := make(map[string]int64, 0)
	var core int64
	lst := rt.assetCache.GetAll()
	for _, val := range lst {
		if val.Type == "主机(exporter)" || val.Type == "主机" {
			for _, m := range val.Metrics {

				if m["name"] == "CPU核数" {
					flo, err := strconv.ParseFloat(m["value"], 64)
					ginx.Dangerous(err)
					core += int64(flo)
				}
			}
		}
		num, Ok := assetType[val.Type]
		if Ok {
			assetType[val.Type] = num + 1
		} else {
			assetType[val.Type] = 1
		}
	}
	deviceNum := 0
	for key := range assetType {
		ar, _ := rt.assetCache.GetType(key)
		if ar.Category == "网络设备" {
			deviceNum++
		}
	}
	device := make([]interface{}, deviceNum)
	appAndCloud := make([]interface{}, len(assetType)-deviceNum)
	var dNum, aNum int64
	// mom := -1
	for key, val := range assetType {
		ar, _ := rt.assetCache.GetType(key)
		//统计环比
		mom, err := models.AssetMom(rt.Ctx, key)
		ginx.Dangerous(err)
		if ar.Category == "网络设备" {
			device[dNum] = map[string]interface{}{"name": key, "count": val, "mom": mom, "sort": dNum}
			dNum++

			//排序
			// device[deviceSort[key]] = map[string]interface{}{"name":key,"count":val}
		} else if ar.Category == "业务资产" {
			appAndCloud[aNum] = map[string]interface{}{"name": key, "count": val, "mom": mom, "sort": aNum}
			aNum++
			//排序
			// appAndCloud[deviceSort[key]] = map[string]interface{}{"name":key,"count":val}
		}
	}
	if len(device) == 0 {
		device = append(device)
	}
	if len(appAndCloud) == 0 {
		appAndCloud = append(appAndCloud)
	}
	m["device"] = device
	m["app_and_cloud"] = appAndCloud

	//TODO
	overview := make(map[string]int64)
	overview["hardware"] = 99
	overview["network"] = 99
	overview["cores"] = core
	overview["application"] = int64(len(assetType) - deviceNum)
	m["overview"] = overview

	ginx.NewRender(c).Data(m, err)
}

// @Summary      资产详情接口
// @Description  资产详情接口
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "add []int64"
// @Success      200  {array}  AssetJson
// @Router       /api/n9e/dashboard/asset/details/ceshi [post]
// @Security     ApiKeyAuth
func (rt *Router) AssetDetails(c *gin.Context) {
	// assetId := ginx.QueryInt64(c, "id", -1)

	f := make([]int64, 0)
	ginx.BindJSON(c, &f)
	logger.Debug(f)

	curAlert, err := models.AlertFeList(rt.Ctx)
	ginx.Dangerous(err)

	var data []*AssetJson
	for _, assetId := range f {
		asset, Ok := rt.assetCache.Get(assetId)
		if !Ok {
			ginx.Bomb(http.StatusOK, "该资产不存在!")
		}
		ar, _ := rt.assetCache.GetType(asset.Type)

		metrics := []map[string]string{}
		for _, m := range asset.Metrics {
			if m["name"] == "入口流量" || m["name"] == "出口流量" || m["name"] == "吞吐量" {
				netData, err := strconv.ParseFloat(strings.Split(m["value"], ".")[0], 64)
				logger.Debug(netData)
				ginx.Dangerous(err)
				if netData < math.Pow(10, 3) {
					m["value"] = strings.Split(m["value"], ".")[0] + "Kb/s"
				} else if netData >= math.Pow(10, 3) && netData < math.Pow(10, 6) {
					netDT, err := strconv.ParseFloat(fmt.Sprintf("%.2f", netData/math.Pow(10, 3)), 64)
					ginx.Dangerous(err)
					logger.Debug(netDT)
					m["value"] = strconv.FormatFloat(netDT, 'f', -1, 64) + "Mb/s"
				} else if netData >= math.Pow(10, 6) && netData < math.Pow(10, 9) {
					netDT, err := strconv.ParseFloat(fmt.Sprintf("%.2f", netData/math.Pow(10, 6)), 64)
					ginx.Dangerous(err)
					logger.Debug(netDT)
					m["value"] = strconv.FormatFloat(netDT, 'f', -1, 64) + "Gb/s"
					logger.Debug(m["value"])
				}
			}
			metrics = append(metrics, m)
		}

		//将名称中含有的IP拿出来
		ip := ""
		name := asset.Name
		compileRegex := regexp.MustCompile("（(.*?)）") // 中文括号，例如：华南地区（广州） -> 广州
		matchArr := compileRegex.FindStringSubmatch(asset.Name)

		if len(matchArr) == 0 {
			compileRegex = regexp.MustCompile("\\((.*?)\\)") // 兼容英文括号并取消括号的转义，例如：华南地区 (广州) -> 广州。
			matchArr = compileRegex.FindStringSubmatch(asset.Name)
		}
		if len(matchArr) != 0 {
			ipTemp := matchArr[len(matchArr)-1]
			if len(strings.Split(ipTemp, ".")) == 4 {
				ip = ipTemp
				nameTemp := strings.Split(asset.Name, matchArr[0])
				name = nameTemp[0]
			}
		}

		var severity int64
		severity = 3
		for _, alertVal := range curAlert {
			if assetId == int64(alertVal.AssetId) {
				if severity < int64(alertVal.Severity) {
					severity = int64(alertVal.Severity)
				}

			}
		}

		data = append(data, &AssetJson{
			Id:       asset.Id,
			Name:     name,
			Ip:       ip,
			Label:    asset.Label,
			Status:   asset.Health,
			UpdateAt: asset.HealthAt,
			Category: ar.Category,
			Type:     asset.Type,
			Metrics:  metrics,
			GroupId:  asset.OrganizationId, //这里使用orgid作为group返回查询条件
			Tags:     asset.TagsJSON,
			Url:      ar.Dashboard,
			Severity: severity,
		})
	}
	ginx.NewRender(c).Data(data, nil)
}

//告警详情接口
// @Summary      告警详情接口
// @Description  告警详情接口
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Param        id    query    int  true  "资产id"
// @Success      200  {array}  models.FeAlert
// @Router       /api/n9e/dashboard/alarm/details [get]
// @Security     ApiKeyAuth
func (rt *Router) AlarmDetails(c *gin.Context) {
	assetId := ginx.QueryStr(c, "id", "")
	if assetId == "" {
		ginx.Bomb(http.StatusOK, "该资产不存在!")
	}
	list, err := models.AlertFeListByAssetId(rt.Ctx, assetId)

	//从资产缓存中更新orgid
	for index, item := range list {
		asset, has := rt.assetCache.Get(int64(item.AssetId))
		if has {
			list[index].Url = "/alert-cur-events/" + strconv.FormatInt(item.Id, 10)
			list[index].OrganizeId = int(asset.OrganizationId)
		}
	}

	ginx.NewRender(c).Data(list, err)
}

// @Summary      历史告警搜索记录查询
// @Description  历史告警搜索记录查询
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Success      200  {array}  []string
// @Router       /api/n9e/dashboard/his-alarms/his-query/ceshi [get]
// @Security     ApiKeyAuth
func (rt *Router) AlarmHisQueryGet(c *gin.Context) {

	jsonData, err := rt.Redis.Get(rt.Ctx.Ctx, "hisQuery").Bytes()
	if err != nil {
		if err.Error() == "redis: nil" {
			ginx.Bomb(http.StatusOK, "无历史搜索记录")
		} else {
			ginx.Dangerous(err)
		}
	}

	// 将JSON字符串转换为结构体数组
	var str []string
	err = json.Unmarshal(jsonData, &str)
	ginx.Dangerous(err)

	logger.Debug(str)

	ginx.NewRender(c).Data(str, err)
}

//告警详情接口
// @Summary      历史告警搜索删除
// @Description  历史告警搜索删除
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Param        record    query    string  true  "搜索记录"
// @Success      200  {array}  models.FeAlert
// @Router       /api/n9e/dashboard/his-alarms/his-query/ceshi [delete]
// @Security     ApiKeyAuth
func (rt *Router) AlarmHisDel(c *gin.Context) {

	record := ginx.QueryStr(c, "record", "")
	if record == "" {
		ginx.Bomb(http.StatusOK, "参数错误")
	}

	jsonData, err := rt.Redis.Get(rt.Ctx.Ctx, "hisQuery").Bytes()
	if err != nil {
		if err.Error() == "redis: nil" {
			ginx.Bomb(http.StatusOK, "无历史搜索记录,不可删除")
		} else {
			ginx.Dangerous(err)
		}
	}

	// 将JSON字符串转换为数组
	var str []string
	err = json.Unmarshal(jsonData, &str)
	ginx.Dangerous(err)

	var new []string
	for _, val := range str {
		if val != record {
			new = append(new, val)
		}
	}

	// 将数组转换为JSON字符串
	data, err := json.Marshal(new)
	ginx.Dangerous(err)
	err = rt.Redis.Set(rt.Ctx.Ctx, "hisQuery", data, 0).Err()

	logger.Debug(str)

	ginx.NewRender(c).Message(err)
}

// @Summary      历史告警过滤条件查询
// @Description  历史告警过滤条件查询
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.FeAlert
// @Router       /api/n9e/dashboard/his-alarm/filter/ceshi [get]
// @Security     ApiKeyAuth
func (rt *Router) AlarmHisFilter(c *gin.Context) {

	enevtType := make(map[int64]string)
	rule, err := models.RuleNameGet(rt.Ctx)
	ginx.Dangerous(err)
	for _, val := range rule {
		enevtType[val.RuleId] = val.RuleName
	}

	severity := make(map[int64]string)
	severity[0] = "紧急"
	severity[1] = "一般"
	severity[2] = "事件"

	dataRange := make(map[int64]string)
	dataRange[300] = "5分钟内"
	dataRange[600] = "10分钟内"
	dataRange[900] = "15分钟内"
	dataRange[1800] = "30分钟内"
	dataRange[3600] = "1小时内"
	dataRange[7200] = "2小时内"
	dataRange[21600] = "6小时内"
	dataRange[43200] = "12小时内"
	dataRange[86400] = "24小时内"

	group := make(map[int64]string)
	groups, err := models.GroupNameGet(rt.Ctx)
	ginx.Dangerous(err)
	for _, val := range groups {
		group[val.GroupId] = val.GroupName
	}

	ginx.NewRender(c).Data(gin.H{
		"enevt_type": enevtType,
		"severity":   severity,
		"data-range": dataRange,
		"group":      group,
	}, err)
}

//告警详情接口
// @Summary      历史告警搜索记录
// @Description  历史告警搜索记录
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Param        enevt_type    query    int64  false  "事件类型"
// @Param        severity    query    int64  false  "告警级别"
// @Param        data-range    query    int64  false  "时间范围"
// @Param        group    query    int64  false  "业务组"
// @Param        query    query    string  false  "搜索框"
// @Param        limit    query    int  false  "条数"
// @Param        page    query    int  false  "页码"
// @Success      200  {array}  models.FeAlert
// @Router       /api/n9e/dashboard/his-alarms/ceshi [get]
// @Security     ApiKeyAuth
func (rt *Router) AlarmHisGet(c *gin.Context) {

	enevtType := ginx.QueryInt64(c, "enevt_type", -1)
	severity := ginx.QueryInt64(c, "severity", -1)
	dataRange := ginx.QueryInt64(c, "data-range", -1)
	group := ginx.QueryInt64(c, "group", -1)
	query := ginx.QueryStr(c, "query", "")
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)

	where := make(map[string]interface{})
	if enevtType != -1 {
		where["rule_id"] = enevtType
	}
	if severity != -1 {
		where["severity"] = severity
	}
	if group != -1 {
		where["group_id"] = group
	}

	if query != "" {
		jsonData, err := rt.Redis.Get(rt.Ctx.Ctx, "hisQuery").Bytes()
		if err != nil {
			if err.Error() != "redis: nil" {
				ginx.Dangerous(err)
			}
		}
		var str []string
		if len(jsonData) != 0 {
			// 将JSON字符串转换为数组
			err = json.Unmarshal(jsonData, &str)
			ginx.Dangerous(err)
		}

		var new []string
		new = append(new, query)
		for _, val := range str {
			if val != query {
				new = append(new, val)
			}
			if len(new) == 10 {
				break
			}
		}
		// 将数组转换为JSON字符串
		data, err := json.Marshal(new)
		ginx.Dangerous(err)
		err = rt.Redis.Set(rt.Ctx.Ctx, "hisQuery", data, 0).Err()
		ginx.Dangerous(err)
	}
	ids, err := models.AssetIdByNameTypeIp(rt.Ctx, query)
	ginx.Dangerous(err)

	total, err := models.AlertHisCountFilter(rt.Ctx, where, dataRange, query, ids)
	ginx.Dangerous(err)

	lst, err := models.AlertHisFilter(rt.Ctx, where, dataRange, query, ids, limit, (page-1)*limit)
	ginx.Dangerous(err)

	//从资产缓存中更新orgid
	for index, item := range lst {
		asset, has := rt.assetCache.Get(int64(item.AssetId))
		if has {
			lst[index].Url = "/alert-his-events/" + strconv.FormatInt(item.Id, 10)
			lst[index].AssetName = asset.Name
			lst[index].OrganizeId = int(asset.OrganizationId)
		}
	}

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, err)
}

// @Summary      列表数据测试
// @Description  列表数据测试
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/n9e/dashboard/data/ceshi [get]
// @Security     ApiKeyAuth
func (rt *Router) DataBoardPerson(c *gin.Context) {

	// me := c.MustGet("user").(*models.User)

	asset, Ok := rt.assetCache.Get(111)
	if !Ok {
		ginx.Bomb(http.StatusOK, "该资产不存在!")
	}
	end := time.Now()
	start := end.Add(-time.Hour * 24)
	r := prom.Range{Start: start, End: end, Step: prom.DefaultStep * 2 * 60}
	str := "cpu_usage_active{asset_id=" + "'" + strconv.FormatInt(asset.Id, 10) + "'" + "}"
	value, warnings, err := rt.PromClients.GetCli(1).QueryRange(context.Background(), str, r)
	ginx.Dangerous(err)

	if len(warnings) > 0 {
		logger.Error(err)
		return
	}
	logger.Debug(value)

	data := make(map[int64]float64)

	items, ok := value.(model.Matrix)
	if !ok {
		return
	}

	for _, item := range items {
		if len(item.Values) == 0 {
			return
		}
		for _, val := range item.Values {
			if math.IsNaN(float64(val.Value)) {
				break
			}
			_, timeSOk := data[val.Timestamp.Unix()]
			if timeSOk {
				continue
			}
			data[val.Timestamp.Unix()] = float64(val.Value)
		}
	}
	logger.Debug(data)
}

// @Summary      筛选资产
// @Description  筛选资产
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Param        filter    query    string  false  "资产类型"
// @Param        query    query    string  false  "搜索框"
// @Success      200
// @Router       /api/n9e/dashboard/user/list/ceshi [get]
// @Security     ApiKeyAuth
func (rt *Router) DataBoardAssetList(c *gin.Context) {

	filter := ginx.QueryStr(c, "filter", "")
	query := ginx.QueryStr(c, "query", "")

	var data []*AssetJson

	ansAssets := make([]models.Asset, 0)
	assets := rt.assetCache.GetAll()
	for _, val := range assets {
		if val.Type != "网络交换机" && val.Type != "主机" {
			continue
		}
		if query == "" {
			if filter == "" {
				ansAssets = append(ansAssets, *val)
			} else {
				if (filter == "网络设备" && val.Type == "网络交换机") || (filter == "ECS" && (val.Type == "主机")) {
					ansAssets = append(ansAssets, *val)
				}
			}
		} else {
			ids, err := models.OrgIdByName(rt.Ctx, query)
			ginx.Dangerous(err)
			if filter == "" {
				if strings.Contains(val.Name, query) || strings.Contains(val.Type, query) || models.IsContain(ids, val.OrganizationId) {
					ansAssets = append(ansAssets, *val)
				}
			} else {
				if (filter == "网络设备" && val.Type == "网络交换机") || (filter == "ECS" && (val.Type == "主机")) {
					if strings.Contains(val.Name, query) || strings.Contains(val.Type, query) || models.IsContain(ids, val.OrganizationId) {
						ansAssets = append(ansAssets, *val)
					}
				}
			}
		}
	}

	for _, asset := range ansAssets {
		ar, _ := rt.assetCache.GetType(asset.Type)
		//将名称中含有的IP拿出来
		ip := ""
		name := asset.Name
		compileRegex := regexp.MustCompile("（(.*?)）") // 中文括号，例如：华南地区（广州） -> 广州
		matchArr := compileRegex.FindStringSubmatch(asset.Name)

		if len(matchArr) == 0 {
			compileRegex = regexp.MustCompile("\\((.*?)\\)") // 兼容英文括号并取消括号的转义，例如：华南地区 (广州) -> 广州。
			matchArr = compileRegex.FindStringSubmatch(asset.Name)
		}
		if len(matchArr) != 0 {
			ipTemp := matchArr[len(matchArr)-1]
			if len(strings.Split(ipTemp, ".")) == 4 {
				ip = ipTemp
				nameTemp := strings.Split(asset.Name, matchArr[0])
				name = nameTemp[0]
			}
		}

		data = append(data, &AssetJson{
			Id:       asset.Id,
			Name:     name,
			Ip:       ip,
			Label:    asset.Label,
			Status:   asset.Health,
			UpdateAt: asset.HealthAt,
			Category: ar.Category,
			Type:     asset.Type,
			GroupId:  asset.OrganizationId, //这里使用orgid作为group返回查询条件
			Tags:     asset.TagsJSON,
		})
	}

	ginx.NewRender(c).Data(data, nil)
}

// @Summary      获取pageName列表
// @Description  获取pageName列表
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Success      200 {array} string
// @Router       /api/n9e/dashboard/user/page-name/get [get]
// @Security     ApiKeyAuth
func (rt *Router) DataBoardPageNameGets(c *gin.Context) {
	me := c.MustGet("user").(*models.User)

	lst, err := models.DashBoardUserPageNameByUser(rt.Ctx, me.Id)
	ginx.NewRender(c).Data(lst, err)
}

// @Summary      删除pageName
// @Description  删除pageName
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Param        page_name query   string     false  "页签"
// @Success      200 {array} string
// @Router       /api/n9e/dashboard/user/page-name/del [delete]
// @Security     ApiKeyAuth
func (rt *Router) DataBoardPageNameDel(c *gin.Context) {
	me := c.MustGet("user").(*models.User)
	pageName := ginx.QueryStr(c, "page_name", "")

	tx := models.DB(rt.Ctx).Begin()

	err := models.DeleteByUserAndPageName(tx, me.Id, pageName)
	tx.Commit()
	ginx.NewRender(c).Message(err)
}

// @Summary      创建/删除用户看板
// @Description  创建/删除用户看板
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Param        page_name query   string     false  "页签"
// @Param        body  body   []models.DashboardUser true "add DashboardUser"
// @Success      200
// @Router       /api/n9e/dashboard/user/add [post]
// @Security     ApiKeyAuth
func (rt *Router) DataBoardAsset(c *gin.Context) {
	pageName := ginx.QueryStr(c, "page_name", "")

	f := make([]models.DashboardUser, 0)
	ginx.BindJSON(c, &f)

	if len(f) == 0 {
		ginx.Bomb(http.StatusOK, "参数为空")
	}

	me := c.MustGet("user").(*models.User)

	for index := range f {
		f[index].UserId = me.Id
		f[index].PageName = pageName
		f[index].CreatedBy = me.Username
	}

	tx := models.DB(rt.Ctx).Begin()

	err := models.DeleteByUserAndPageName(tx, me.Id, pageName)
	ginx.Dangerous(err)

	err = models.AddDashBoardUser(tx, f)
	tx.Commit()
	ginx.NewRender(c).Message(err)
}

// @Summary      数据面板
// @Description  数据面板
// @Tags         大屏展示
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        page_name query   string     false  "页签名称"
// @Param        start query   int     false  "开始时间"
// @Param        end query   int     false  "结束时间"
// @Success      200
// @Router       /api/n9e/dashboard/user/data [get]
// @Security     ApiKeyAuth
func (rt *Router) DataBoardAssetsGet(c *gin.Context) {
	logger.Debug("9999999999999999999999999")

	models.AssetTypeGetsAll()

	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 4)
	pageName := ginx.QueryStr(c, "page_name", "")
	startT := ginx.QueryInt64(c, "start", -1)
	endT := ginx.QueryInt64(c, "end", -1)
	if pageName == "" || startT == -1 || endT == -1 {
		ginx.Bomb(http.StatusOK, "参数为空")
	}
	if startT >= endT {
		ginx.Bomb(http.StatusOK, "时间区间错误")
	}

	data := make([]interface{}, 0)
	me := c.MustGet("user").(*models.User)
	total, err := models.DashBoardUserCountByUserAndPageName(rt.Ctx, me.Id, pageName)
	ginx.Dangerous(err)
	lst, err := models.DashBoardUserByUserAndPageName(rt.Ctx, me.Id, pageName, limit, (page-1)*limit)
	ginx.Dangerous(err)

	for _, val := range lst {
		m := make(map[string]interface{})
		asset, has := rt.assetCache.Get(val.AssetsId)
		if !has {
			ginx.Bomb(http.StatusOK, "该资产不存在")
		}
		m["sort"] = val.Sort
		m["page_name"] = val.PageName

		ar, _ := rt.assetCache.GetType(asset.Type)
		//将名称中含有的IP拿出来
		ip := ""
		name := asset.Name
		compileRegex := regexp.MustCompile("（(.*?)）") // 中文括号，例如：华南地区（广州） -> 广州
		matchArr := compileRegex.FindStringSubmatch(asset.Name)

		if len(matchArr) == 0 {
			compileRegex = regexp.MustCompile("\\((.*?)\\)") // 兼容英文括号并取消括号的转义，例如：华南地区 (广州) -> 广州。
			matchArr = compileRegex.FindStringSubmatch(asset.Name)
		}
		if len(matchArr) != 0 {
			ipTemp := matchArr[len(matchArr)-1]
			if len(strings.Split(ipTemp, ".")) == 4 {
				ip = ipTemp
				nameTemp := strings.Split(asset.Name, matchArr[0])
				name = nameTemp[0]
			}
		}

		m["id"] = asset.Id
		m["name"] = name
		m["ip"] = ip
		m["status"] = asset.Health
		m["label"] = asset.Label
		m["update_at"] = asset.HealthAt
		m["category"] = ar.Category
		m["type"] = asset.Type
		m["group_id"] = asset.OrganizationId //这里使用orgid作为group返回查询条件
		m["tags"] = asset.TagsJSON

		// end := time.Now()
		// start := end.Add(-time.Hour * 24)
		end := time.Unix(endT, 0)
		start := time.Unix(startT, 0)
		r := prom.Range{Start: start, End: end, Step: prom.DefaultStep * 2}
		// str := "cpu_usage_active{asset_id=" + "'" + strconv.FormatInt(asset.Id, 10) + "'" + "}"

		query := make(map[string]string, 0)
		if val.Type == "主机" {
			query["cpu_usage_active"] = "cpu_usage_active{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}"
			query["net_bits_recv"] = "net_bits_recv{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}"
			query["net_bits_sent"] = "net_bits_sent{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}"
			query["disk_used_percent"] = "sum by(asset_id) (disk_used{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}) / sum by(asset_id) (disk_total{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'})"
			query["mem_used_percent"] = "mem_used_percent{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}"
			query["diskio_read"] = "sum by(asset_id) (rate(diskio_read_bytes{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}[1m]))"
			query["diskio_write"] = "sum by(asset_id) (rate(diskio_write_bytes{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}[1m]))"
			query["netstat_tcp_established"] = "netstat_tcp_established{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}"
			query["system_load1"] = "system_load1{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}"
		} else {
			query["switch_in"] = "sum by(asset_id) (rate(switch_legacy_if_in{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}[5m]))"
			query["switch_out"] = "sum by(asset_id) (rate(switch_legacy_if_out{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}[5m]))"
			query["cpu_usage_active"] = "switch_legacy_cpu_util{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}"
			query["mem_used_percent"] = "switch_legacy_mem_util{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'}"
			query["drop_package"] = "sum by(asset_id) (switch_legacy_if_in_discards{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'})[1m]"
			query["switch_in_percent"] = "sum by(asset_id) (switch_legacy_if_in_speed_percent{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'})[1m]"
			query["switch_out_percent"] = "sum by(asset_id) (switch_legacy_if_out_speed_percent{asset_id='" + strconv.FormatInt(asset.Id, 10) + "'})[1m]"
		}

		for key, str := range query {
			value, warnings, err := rt.PromClients.GetCli(1).QueryRange(context.Background(), str, r)
			ginx.Dangerous(err)

			if len(warnings) > 0 {
				logger.Error(err)
				return
			}
			dataVal := make(map[int64]float64)

			items, ok := value.(model.Matrix)
			if !ok {
				return
			}

			for _, item := range items {
				if len(item.Values) == 0 {
					return
				}
				for _, val := range item.Values {
					if math.IsNaN(float64(val.Value)) {
						break
					}
					_, timeSOk := dataVal[val.Timestamp.Unix()]
					if timeSOk {
						continue
					}
					if key == "disk_used_percent" {
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(val.Value)*100), 64)
						dataVal[val.Timestamp.Unix()] = floatVal
					} else {
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(val.Value)), 64)
						dataVal[val.Timestamp.Unix()] = floatVal
					}
				}
			}
			m[key] = dataVal
		}

		//当前值
		mCur := make(map[string]float64)
		for key, str := range query {
			value, warnings, err := rt.PromClients.GetCli(1).Query(context.Background(), str, time.Now())
			ginx.Dangerous(err)

			if len(warnings) > 0 {
				logger.Error(err)
				return
			}
			// dataVal := make(map[int64]float64)
			logger.Debug("_____________________________")
			logger.Debug(key)
			logger.Debug(value.Type())
			logger.Debug(value)
			if value.Type() == model.ValVector {
				items, ok := value.(model.Vector)
				if !ok {
					return
				}

				for _, item := range items {
					if math.IsNaN(float64(item.Value)) {
						continue
					}
					if key == "disk_used_percent" {
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(item.Value)*100), 64)
						mCur[key] = floatVal
					} else {
						floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(item.Value)), 64)
						mCur[key] = floatVal
					}
				}
			} else if value.Type() == model.ValMatrix {
				items, ok := value.(model.Matrix)
				if !ok {
					return
				}

				for _, item := range items {
					if len(item.Values) == 0 {
						return
					}
					if math.IsNaN(float64(item.Values[0].Value)) {
						continue
					}

					floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(item.Values[0].Value)), 64)
					mCur[key] = floatVal
				}
			}
			m["cur"] = mCur
		}
		data = append(data, m)
	}
	ginx.NewRender(c).Data(gin.H{
		"list":  data,
		"total": total,
	}, nil)
}
