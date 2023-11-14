package prom

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ctx"
	json "github.com/json-iterator/go"
	"github.com/prometheus/common/model"

	// "github.com/pkg/errors"
	// "github.com/prometheus/common/model"
	// "github.com/prometheus/prometheus/pkg/labels"
	// "github.com/prometheus/prometheus/promql/parser"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

/* 声明全局变量 */
var dsMon map[int64][]string

func BuildPromSql(ctx *ctx.Context, metrics, assetId string, datasource int64) string {
	if len(dsMon) == 0 {
		getAllMetrics(ctx)
	}
	var sql string
	for _, val := range dsMon[datasource] {
		if strings.Contains(metrics, val) {
			sqlLst := strings.Split(metrics, val)
			sql = ""
			for index, sqlPart := range sqlLst {
				if index == 0 {
					sql += sqlPart
				} else {
					if sqlPart[0:1] == "{" {
						sql += val + "{asset_id='" + assetId + "'," + sqlPart[1:strings.Count(sqlPart, "")-1]
					} else {
						sql += val + "{asset_id='" + assetId + "'}" + sqlPart
					}
				}
			}
		}
	}
	return sql
}

func GetValue(value model.Value) float64 {

	var res float64 = 0

	if value.Type() == model.ValVector {
		items, ok := value.(model.Vector)
		if !ok {
			return res
		}

		for _, item := range items {
			if math.IsNaN(float64(item.Value)) {
				continue
			}
			floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(item.Value)), 64)
			res = floatVal
		}
	} else if value.Type() == model.ValMatrix {
		items, ok := value.(model.Matrix)
		if !ok {
			return res
		}

		for _, item := range items {
			if len(item.Values) == 0 {
				return res
			}
			if math.IsNaN(float64(item.Values[0].Value)) {
				continue
			}

			floatVal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(item.Values[0].Value)), 64)
			res = floatVal
		}
	}
	return res
}

func getAllMetrics(ctx *ctx.Context) {

	datasources, err := models.DatasourceGetAll(ctx)
	ginx.Dangerous(err)

	for _, datasource := range datasources {
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

// func injectLabels(expr parser.Expr, match labels.MatchType, name, value string) {
// 	switch e := expr.(type) {
// 	case *parser.AggregateExpr:
// 		injectLabels(e.Expr, match, name, value)
// 	case *parser.Call:
// 		for _, v := range e.Args {
// 			injectLabels(v, match, name, value)
// 		}
// 	case *parser.ParenExpr:
// 		injectLabels(e.Expr, match, name, value)
// 	case *parser.UnaryExpr:
// 		injectLabels(e.Expr, match, name, value)
// 	case *parser.BinaryExpr:
// 		injectLabels(e.LHS, match, name, value)
// 		injectLabels(e.RHS, match, name, value)
// 	case *parser.VectorSelector:
// 		l := genMetricLabel(match, name, value)
// 		e.LabelMatchers = append(e.LabelMatchers, l)
// 		return
// 	case *parser.MatrixSelector:
// 		injectLabels(e.VectorSelector, match, name, value)
// 	case *parser.SubqueryExpr:
// 		injectLabels(e.Expr, match, name, value)
// 	case *parser.NumberLiteral, *parser.StringLiteral:
// 		return
// 	default:
// 		panic(errors.Errorf("unhandled expression of type: %T", expr))
// 	}
// 	return
// }

// func genMetricLabel(match labels.MatchType, name, value string) *labels.Matcher {
// 	m, err := labels.NewMatcher(match, name, value)
// 	if nil != err {
// 		return nil
// 	}

// 	return m
// }

// func main() {
// 	ql := `rate(http_requests_total[5m])[30m:1m]`
// 	expr, _ := parser.ParseExpr(ql)
// 	injectLabels(expr, labels.MatchEqual, "asset_id", "123")
// 	fmt.Println(expr.String())
// }
