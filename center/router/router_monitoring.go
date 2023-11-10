// Package models  监控
// date : 2023-10-08 16:45
// desc : 监控
package router

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/prom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取监控
// @Description  根据主键获取监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.Monitoring
// @Router       /api/n9e/xh/monitoring/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) monitoringGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	monitoring, err := models.MonitoringGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if monitoring == nil {
		ginx.Bomb(404, "No such monitoring")
	}

	ginx.NewRender(c).Data(monitoring, nil)
}

// @Summary      查询监控
// @Description  根据条件查询监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        query query   string  false  "搜索栏"
// @Param        assetType query   string  false  "资产类型"
// @Param        datasource query   int  false  "数据源类型"
// @Success      200  {array}  models.Monitoring
// @Router       /api/n9e/xh/monitoring/filter [get]
// @Security     ApiKeyAuth
func (rt *Router) monitoringGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	assetType := ginx.QueryStr(c, "assetType", "")
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")
	datasource := ginx.QueryInt(c, "datasource", -1)

	m := make(map[string]interface{})
	// if dataSource != -1 {
	// 	m["datasource.id"] = dataSource
	// }
	// if assetType != ""{
	// 	m["assets.type"] = assetType
	// }
	total, err := models.MonitoringMapCount(rt.Ctx, m, query, assetType, datasource)
	ginx.Dangerous(err)
	lst, err := models.MonitoringMapGets(rt.Ctx, m, query, limit, (page-1)*limit, assetType, datasource)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      查询所有监控
// @Description  查询所有监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.Monitoring
// @Router       /api/n9e/xh/monitoring/ [get]
// @Security     ApiKeyAuth
func (rt *Router) monitoringAllGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.MonitoringCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.MonitoringAllGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建监控
// @Description  创建监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        body  body   models.Monitoring true "add monitoring"
// @Success      200
// @Router       /api/n9e/xh/monitoring/ [post]
// @Security     ApiKeyAuth
func (rt *Router) monitoringAdd(c *gin.Context) {
	var f models.Monitoring
	ginx.BindJSON(c, &f)
	me := c.MustGet("user").(*models.User)

	var monitor = models.Monitoring{
		AssetId:        f.AssetId,
		MonitoringName: f.MonitoringName,
		DatasourceId:   f.DatasourceId,
		MonitoringSql:  f.MonitoringSql,
		Status:         f.Status,
		TargetId:       f.TargetId,
		Remark:         f.Remark,
		CreatedBy:      me.Username,
		CreatedAt:      time.Now().Unix(),
		UpdatedBy:      me.Username,
		UpdatedAt:      time.Now().Unix(),
	}

	// 更新模型
	err := monitor.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新监控
// @Description  更新监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        body  body   models.Monitoring true "update monitoring"
// @Success      200
// @Router       /api/n9e/xh/monitoring/ [put]
// @Security     ApiKeyAuth
func (rt *Router) monitoringPut(c *gin.Context) {
	var f models.Monitoring
	ginx.BindJSON(c, &f)

	old, err := models.MonitoringGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "monitoring not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除监控
// @Description  根据主键删除监控
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/xh/monitoring/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) monitoringDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	monitoring, err := models.MonitoringGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if monitoring == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(monitoring.Del(rt.Ctx))
}

// @Summary      监控开关
// @Description  根据主键更新监控状态
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param		 status	query	int		true "新状态"
// @Param        body	body	[]int64 true "add ids"
// @Success      200
// @Router       /api/n9e/xh/monitoring/status/{id} [post]
// @Security     ApiKeyAuth
func (rt *Router) monitoringStatusUpdate(c *gin.Context) {
	status := ginx.QueryInt64(c, "status", 1)
	var ids []int64
	ginx.BindJSON(c, &ids)

	ginx.NewRender(c).Message(models.MonitoringUpdateStatus(rt.Ctx, ids, status))
}

// @Summary      获取监控数据
// @Description  获取监控数据
// @Tags         监控
// @Accept       json
// @Produce      json
// @Param        start query   int     false  "开始时间"
// @Param        end query   int     false  "结束时间"
// @Param        body  body   map[string]interface{} true "update query"
// @Success      200
// @Router       /api/n9e/xh/monitoring/data [post]
// @Security     ApiKeyAuth
func (rt *Router) monitoringData(c *gin.Context) {

	startT := ginx.QueryInt64(c, "start", -1)
	endT := ginx.QueryInt64(c, "end", -1)
	if startT == -1 || endT == -1 {
		ginx.Bomb(http.StatusOK, "参数为空")
	}
	if startT >= endT {
		ginx.Bomb(http.StatusOK, "时间区间错误")
	}

	end := time.Unix(endT, 0)
	start := time.Unix(startT, 0)

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
	lst, err := models.MonitoringGetByBatchId(rt.Ctx, ids)
	ginx.Dangerous(err)

	datasourceMap := make(map[int64]interface{}, 0)
	dsMon := make(map[int64][]string, 0)
	dat := make(map[int64]interface{}, 0)
	for _, monitoring := range lst {
		_, datasourceOk := datasourceMap[monitoring.DatasourceId]
		if !datasourceOk {
			datasource, err := models.DatasourceGet(rt.Ctx, monitoring.DatasourceId)
			ginx.Dangerous(err)
			datasourceMap[monitoring.DatasourceId] = nil

			err = datasource.DB2FE()
			ginx.Dangerous(err)

			client := &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: datasource.HTTPJson.TLS.SkipTlsVerify,
					},
				},
			}

			fullURL := datasource.HTTPJson.Url + "/api/v1/label/__name__/values"
			req, err := http.NewRequest("GET", fullURL, nil)
			if err != nil {
				logger.Errorf("Error creating request: %v", err)
				ginx.Bomb(http.StatusOK, fmt.Errorf("request url:%s failed", fullURL).Error())
			}
			if datasource.AuthJson.BasicAuth && datasource.AuthJson.BasicAuthUser != "" {
				req.SetBasicAuth(datasource.AuthJson.BasicAuthUser, datasource.AuthJson.BasicAuthPassword)
			}

			for k, v := range datasource.HTTPJson.Headers {
				req.Header.Set(k, v)
			}
			resp, err := client.Do(req)
			if err != nil {
				logger.Errorf("Error making request: %v\n", err)
				ginx.Bomb(http.StatusOK, fmt.Errorf("request url:%s failed", fullURL).Error())
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				logger.Errorf("Error making request: %v\n", resp.StatusCode)
				ginx.Bomb(http.StatusOK, fmt.Errorf("request url:%s failed code:%d", fullURL, resp.StatusCode).Error())
			}
			// 读取响应Body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				ginx.Bomb(http.StatusOK, "获取数据失败")
				return
			}
			var result map[string]interface{}
			err = json.Unmarshal(body, &result)
			if err != nil {
				logger.Errorf("解析JSON发生错误:", err)
				ginx.Bomb(http.StatusOK, "获取数据失败")
				return
			}
			dataT, dataOk := result["data"]
			data := make([]string, 0)
			if dataOk {
				for _, val := range dataT.([]interface{}) {
					data = append(data, val.(string))
				}
			}
			dsMon[datasource.Id] = data
		}
	}
	for _, monitoring := range lst {
		sql := monitoring.MonitoringSql
		for _, val := range dsMon[monitoring.DatasourceId] {
			if strings.Contains(sql, val) {
				sqlLst := strings.Split(sql, val)
				sql = ""
				for index, sqlPart := range sqlLst {
					if index == 0 {
						sql += sqlPart
					} else {
						if sqlPart[0:1] == "{" {
							sql += val + "{asset_id='" + strconv.FormatInt(monitoring.AssetId, 10) + "'," + sqlPart[1:strings.Count(sqlPart, "")-1]
						} else {
							sql += val + "{asset_id='" + strconv.FormatInt(monitoring.AssetId, 10) + "'}" + sqlPart
						}
					}
				}
			}
		}
		var r prom.Range
		if endT-startT <= 60*60*24 {
			r = prom.Range{Start: start, End: end, Step: prom.DefaultStep * 2}
		} else if endT-startT > 60*60*24 && endT-startT <= 60*60*24*7 {
			r = prom.Range{Start: start, End: end, Step: prom.DefaultStep * 2 * 60}
		} else if endT-startT > 60*60*24*7 {
			r = prom.Range{Start: start, End: end, Step: prom.DefaultStep * 2 * 60 * 60}
		}
		value, warnings, err := rt.PromClients.GetCli(monitoring.DatasourceId).QueryRange(context.Background(), sql, r)
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
				dataVal[val.Timestamp.Unix()] = float64(val.Value)
			}
		}
		dat[monitoring.Id] = dataVal
	}
	ginx.NewRender(c).Data(dat, err)
}
