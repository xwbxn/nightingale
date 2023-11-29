package router

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/ccfos/nightingale/v6/pkg/txt"
	xmltool "github.com/ccfos/nightingale/v6/pkg/xml"
	"github.com/prometheus/common/model"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

func (rt *Router) assetsGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assets, err := models.AssetGetById(rt.Ctx, id)
	ginx.NewRender(c).Data(assets, err)
}

func (rt *Router) assetsGets(c *gin.Context) {
	bgid := ginx.QueryInt64(c, "bgid", -1)
	query := ginx.QueryStr(c, "query", "")
	orgId := ginx.QueryInt64(c, "organization_id", -1)
	assets, err := models.AssetGets(rt.Ctx, bgid, query, orgId)
	ginx.Dangerous(err)
	for _, asset := range assets {
		atype, ok := rt.assetCache.GetType(asset.Type)
		if ok {
			asset.Dashboard = strings.ReplaceAll(atype.Dashboard, "${id}", fmt.Sprintf("%d", asset.Id))
		}
		health, ok := rt.assetCache.Get(asset.Id)
		if ok {
			asset.Health = health.Health
		}
	}
	ginx.NewRender(c).Data(assets, err)
}

type assetsModel struct {
	Id              int64  `json:"id"`
	Version         string `json:"version"`
	Ident           string `json:"ident"`
	GroupId         int64  `json:"group_id"`
	Name            string `json:"name"`
	Label           string `json:"label"`
	Tags            string `json:"tags"`
	Type            string `json:"type"`
	Ip              string `json:"ip"`
	Manufacturers   string `json:"manufacturers"`
	Position        string `json:"position"`
	Memo            string `json:"memo"`
	Configs         string `json:"configs"`
	Params          string `json:"params"`
	OrganizationId  int64  `json:"organization_id"`
	OptionalMetrics string `json:"optional_metrics"`
}

func (rt *Router) assetsAdd(c *gin.Context) {
	var f assetsModel
	ginx.BindJSON(c, &f)
	me := c.MustGet("user").(*models.User)

	var assets = models.Asset{
		GroupId:        f.GroupId,
		Name:           f.Name,
		Ident:          f.Ident,
		Label:          f.Label,
		Type:           f.Type,
		Memo:           f.Memo,
		Configs:        f.Configs,
		Params:         f.Params,
		CreateBy:       me.Username,
		CreateAt:       time.Now().Unix(),
		OrganizationId: f.OrganizationId,
	}

	err := assets.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      创建资产-西航
// @Description  创建资产-西航
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Param        body  body   assetsModel true "add assetsModel"
// @Success      200
// @Router       /api/n9e/xh/assets/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetsAddXH(c *gin.Context) {
	var f map[string]interface{}
	ginx.BindJSON(c, &f)
	me := c.MustGet("user").(*models.User)

	num, err := models.AssetsCountMap(rt.Ctx, map[string]interface{}{"ip": f["ip"]})
	ginx.Dangerous(err)
	if num > 0 {
		ginx.Bomb(http.StatusOK, "IP已存在")
	}

	position := ""
	positionT, positionOk := f["position"]
	if positionOk {
		position = positionT.(string)
	}
	memo := ""
	memoT, memoOk := f["memo"]
	if memoOk {
		memo = memoT.(string)
	}
	manufacturers := ""
	manufacturersT, manufacturersOk := f["manufacturers"]
	if manufacturersOk {
		manufacturers = manufacturersT.(string)
	}

	var group int64 = 0
	groupT, groupOk := f["group_id"]
	if groupOk {
		group = int64(groupT.(float64))
	}

	var assets = models.Asset{
		GroupId:       group,
		Name:          f["name"].(string),
		Type:          f["type"].(string),
		Ip:            f["ip"].(string),
		Manufacturers: manufacturers,
		Position:      position,
		Memo:          memo,
		CreateBy:      me.Username,
		CreateAt:      time.Now().Unix(),
	}

	id, err := assets.AddXH(rt.Ctx)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(id, err)
}

type optionalMetricsForm struct {
	Id              int64             `json:"id"`
	OptionalMetrics []*models.Metrics `json:"optional_metrics"`
}

func (rt *Router) putOptionalMetrics(c *gin.Context) {
	var f optionalMetricsForm
	ginx.BindJSON(c, &f)
	oldAssets, err := models.AssetGet(rt.Ctx, "id=?", f.Id)
	ginx.Dangerous(err)
	me := c.MustGet("user").(*models.User)

	if oldAssets == nil {
		ginx.Bomb(http.StatusOK, "assets not found")
	}

	om, err := json.Marshal(f.OptionalMetrics)
	ginx.Dangerous(err)
	oldAssets.OptionalMetrics = string(om)
	oldAssets.UpdateAt = time.Now().Unix()
	oldAssets.UpdateBy = me.Username

	err = oldAssets.Update(rt.Ctx, "optional_metrics", "update_at", "update_by")
	ginx.NewRender(c).Message(err)
}

func (rt *Router) assetPut(c *gin.Context) {
	var f assetsModel
	ginx.BindJSON(c, &f)
	oldAssets, err := models.AssetGet(rt.Ctx, "id=?", f.Id)
	ginx.Dangerous(err)
	me := c.MustGet("user").(*models.User)

	if oldAssets == nil {
		ginx.Bomb(http.StatusOK, "assets not found")
	}

	oldAssets.Ident = f.Ident
	oldAssets.Label = f.Label
	oldAssets.Name = f.Name
	oldAssets.Configs = f.Configs
	oldAssets.Memo = f.Memo
	oldAssets.Params = f.Params
	oldAssets.UpdateAt = time.Now().Unix()
	oldAssets.UpdateBy = me.Username

	err = oldAssets.Update(rt.Ctx, "name", "params", " ident", "label", "configs", "memo", "update_at", "update_by")
	ginx.NewRender(c).Message(err)
}

// @Summary      更新资产-西航
// @Description  更新资产-西航
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Param        body  body   assetsModel true "add assetsModel"
// @Success      200
// @Router       /api/n9e/xh/assets/ [put]
// @Security     ApiKeyAuth
func (rt *Router) assetPutXH(c *gin.Context) {
	var f map[string]interface{}
	ginx.BindJSON(c, &f)

	oldAssets, err := models.AssetGet(rt.Ctx, "id=?", int64(f["id"].(float64)))
	ginx.Dangerous(err)

	me := c.MustGet("user").(*models.User)

	if oldAssets == nil {
		ginx.Bomb(http.StatusOK, "资产不存在")
	}
	if oldAssets.Ip != f["ip"].(string) {
		num, err := models.AssetsCountMap(rt.Ctx, map[string]interface{}{"ip": f["ip"]})
		ginx.Dangerous(err)
		if num > 0 {
			ginx.Bomb(http.StatusOK, "IP已存在")
		}
	}

	position := ""
	positionT, positionOk := f["position"]
	if positionOk {
		position = positionT.(string)
	}
	memo := ""
	memoT, memoOk := f["memo"]
	if memoOk {
		memo = memoT.(string)
	}
	manufacturers := ""
	manufacturersT, manufacturersOk := f["manufacturers"]
	if manufacturersOk {
		manufacturers = manufacturersT.(string)
	}
	var group int64 = 0
	groupT, groupOk := f["group_id"]
	if groupOk {
		group = int64(groupT.(float64))
	}

	oldAssets.GroupId = group
	oldAssets.Name = f["name"].(string)
	oldAssets.Type = f["type"].(string)
	oldAssets.Ip = f["ip"].(string)
	oldAssets.Manufacturers = manufacturers
	oldAssets.Position = position
	oldAssets.Memo = memo
	oldAssets.UpdateAt = time.Now().Unix()
	oldAssets.UpdateBy = me.Username

	err = oldAssets.Update(rt.Ctx, "group_id", "name", "type", "ip", "manufacturers", "os", "cpu", "memory", "plugin_version", "position", "asset_status", "directory_id", "memo", "update_at", "update_by")
	ginx.Dangerous(err)

	ginx.NewRender(c).Message(err)
}

type assetsForm struct {
	Ids []string `json:"ids" binding:"required"`
}

func (rt *Router) assetDel(c *gin.Context) {
	var f assetsForm
	ginx.BindJSON(c, &f)

	if len(f.Ids) == 0 {
		ginx.Bomb(http.StatusBadRequest, "参数为空")
	}
	tx := models.DB(rt.Ctx).Begin()
	err := models.AssetDelTx(tx, f.Ids)
	ginx.Dangerous(err)
	err = models.AssetsExpansionDelAssetsIds(tx, f.Ids)
	tx.Commit()

	ginx.NewRender(c).Message(err)
}

// @Summary      根据id查询资产-西航
// @Description  根据id查询资产-西航
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Param        asset query   int     true  "资产ID"
// @Success      200
// @Router       /api/n9e/xh/assets/id [get]
// @Security     ApiKeyAuth
func (rt *Router) assetGetById(c *gin.Context) {
	id := ginx.QueryInt64(c, "asset", -1)
	if id == -1 {
		ginx.Bomb(http.StatusOK, "参数错误")
	}

	lst, err := models.AssetsGetsMap(rt.Ctx, map[string]interface{}{"id": id})
	ginx.Dangerous(err)
	if len(lst) == 0 {
		ginx.Bomb(http.StatusOK, "资产不存在")
	}
	for index, asset := range lst {
		atype, ok := rt.assetCache.GetType(asset.Type)
		if ok {
			asset.Dashboard = strings.ReplaceAll(atype.Dashboard, "${id}", fmt.Sprintf("%d", asset.Id))
		}
		health, ok := rt.assetCache.Get(asset.Id)
		if ok {
			asset.Health = health.Health
		}
		assetType, _ := rt.assetCache.GetType(asset.Type)
		assetC, assetCOk := rt.assetCache.Get(asset.Id)
		if assetCOk {
			metrics := make([]map[string]interface{}, 0)

			for _, val := range assetType.Metrics {
				value, valueOk := assetC.Metrics[val.Name]
				if valueOk {
					metrics = append(metrics, map[string]interface{}{"label": val.Metrics, "val": value["value"]})
				}
			}
			lst[index].MetricsList = metrics
		}

		exps, err := models.AssetsExpansionGetsMap(rt.Ctx, map[string]interface{}{"assets_id": asset.Id})
		ginx.Dangerous(err)
		lst[index].Exps = exps
	}

	ginx.NewRender(c).Data(lst[0], nil)
}

// @Summary      过滤器-西航
// @Description  过滤器-西航
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Param        type    query    string  false  "资产类型"
// @Param        filter    query    string  false  "筛选框(“ip”：IP地址；“name”：资产名称；“manufacturers”：厂商；“os”操作系统；“group_id”：业务组；“position”：资产位置)"
// @Param        query    query    string  false  "搜索框"
// @Param        page    query    int  false  "页码"
// @Param        limit    query    int  false  "条数"
// @Success      200
// @Router       /api/n9e/xh/assets/filter [get]
// @Security     ApiKeyAuth
func (rt *Router) assetGetFilter(c *gin.Context) {
	// var f map[string]interface{}
	// ginx.BindJSON(c, &f)

	// //分页参数
	// limitT, limitOk := f["limit"]
	// limit := 20
	// if limitOk {
	// 	limit = int(limitT.(float64))
	// }
	// pageT, pageOk := f["page"]
	// page := 1
	// if pageOk {
	// 	page = int(pageT.(float64))
	// }

	// query := ""
	// queryTemp, queryOk := f["query"]
	// // if !queryOk {
	// // 	query = ""
	// // }

	// queryType := ""
	// filter, filterOk := f["filter"]
	// if filterOk {
	// 	if !queryOk {
	// 		ginx.Bomb(http.StatusBadRequest, "参数为空")
	// 	} else {
	// 		query = "%" + queryTemp.(string) + "%"
	// 	}

	// 	// if filter.(string) == "1" || filter.(string) == "4" {
	// 	// 	queryType = "name"
	// 	// } else if filter.(string) == "2" {
	// 	// 	queryType = "ip"
	// 	// } else if filter.(string) == "3" {
	// 	// 	queryType = "type"
	// 	// }
	// 	if filter.(string) == "group" {
	// 		queryType = "group_id"
	// 	}
	// } else {
	// 	if queryOk {
	// 		query = "%" + queryTemp.(string) + "%"
	// 	}
	// }
	// aType := ""
	// typeT, typeOk := f["type"]
	// if typeOk {
	// 	aType = typeT.(string)
	// }
	// total, err := models.AssetsCountFilter(rt.Ctx, aType, query, queryType)
	// ginx.Dangerous(err)

	// lst, err := models.AssetsGetsFilter(rt.Ctx, aType, query, queryType, limit, (page-1)*limit)
	// ginx.Dangerous(err)

	filter := ginx.QueryStr(c, "filter", "")
	query := ginx.QueryStr(c, "query", "")
	page := ginx.QueryInt(c, "page", 1)
	aType := ginx.QueryStr(c, "type", "")
	limit := ginx.QueryInt(c, "limit", 20)
	if filter == "" && query != "" {
		ginx.Bomb(http.StatusOK, "参数错误")
	}

	total, err := models.AssetsCountFilterNew(rt.Ctx, filter, query, aType)
	ginx.Dangerous(err)

	lst, err := models.AssetsGetsFilterNew(rt.Ctx, filter, query, aType, limit, (page-1)*limit)
	ginx.Dangerous(err)

	for index, asset := range lst {
		atype, ok := rt.assetCache.GetType(asset.Type)
		if ok {
			asset.Dashboard = strings.ReplaceAll(atype.Dashboard, "${id}", fmt.Sprintf("%d", asset.Id))
		}
		health, ok := rt.assetCache.Get(asset.Id)
		if ok {
			asset.Health = health.Health
		}
		assetType, _ := rt.assetCache.GetType(asset.Type)
		assetC, assetCOk := rt.assetCache.Get(asset.Id)
		if assetCOk {
			metrics := make([]map[string]interface{}, 0)

			if len(assetType.Metrics) > 0 {
				for _, val := range assetType.Metrics {
					value, valueOk := assetC.Metrics[val.Name]
					if valueOk {
						metrics = append(metrics, map[string]interface{}{"label": val.Metrics, "val": value["value"]})
					}
				}
			}

			lst[index].MetricsList = metrics
		}

		exps, err := models.AssetsExpansionGetsMap(rt.Ctx, map[string]interface{}{"assets_id": asset.Id})
		ginx.Dangerous(err)
		lst[index].Exps = exps
	}

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      全量资产信息-西航
// @Description  全量资产信息-西航
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/n9e/xh/assets/out [get]
// @Security     ApiKeyAuth
func (rt *Router) assetGetAll(c *gin.Context) {

	lst, err := models.AssetsGetsFilter(rt.Ctx, "", "", "", -1, -1)
	ginx.Dangerous(err)

	for index, asset := range lst {
		atype, ok := rt.assetCache.GetType(asset.Type)
		if ok {
			asset.Dashboard = strings.ReplaceAll(atype.Dashboard, "${id}", fmt.Sprintf("%d", asset.Id))
		}
		health, ok := rt.assetCache.Get(asset.Id)
		if ok {
			asset.Health = health.Health
		}
		assetType, _ := rt.assetCache.GetType(asset.Type)
		assetC, assetCOk := rt.assetCache.Get(asset.Id)
		if assetCOk {
			metrics := make([]map[string]interface{}, 0)

			if len(assetType.Metrics) > 0 {
				for _, val := range assetType.Metrics {
					value, valueOk := assetC.Metrics[val.Name]
					if valueOk {
						metrics = append(metrics, map[string]interface{}{"label": val.Metrics, "val": value["value"]})
					}
				}
			}

			lst[index].MetricsList = metrics
		}

		exps, err := models.AssetsExpansionGetsMap(rt.Ctx, map[string]interface{}{"assets_id": asset.Id})
		ginx.Dangerous(err)
		lst[index].Exps = exps
	}

	ginx.NewRender(c).Data(lst, nil)
}

type UpdateBody struct {
	Ids   []int64     `json:"ids"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// @Summary      批量修改资产属性-西航
// @Description  批量修改资产属性-西航
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Param        body  body   UpdateBody true "add UpdateBody"
// @Success      200
// @Router       /api/n9e/xh/assets/batch-update/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetUpdateXH(c *gin.Context) {
	var f UpdateBody
	ginx.BindJSON(c, &f)

	ginx.NewRender(c).Message(models.UpdateByIds(rt.Ctx, f.Ids, f.Name, f.Value))
}

// @Summary      批量删除资产-西航
// @Description  批量删除资产-西航
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Param        body  body  assetsForm   true "add assetsForm"
// @Success      200
// @Router       /api/n9e/xh/assets/batch-del/ [post]
// @Security     ApiKeyAuth
func (rt *Router) assetDelXH(c *gin.Context) {
	var f assetsForm
	ginx.BindJSON(c, &f)

	if len(f.Ids) == 0 {
		ginx.Bomb(http.StatusBadRequest, "参数为空")
	}
	tx := models.DB(rt.Ctx).Begin()
	//删除资产基本属性
	err := models.AssetDelTx(tx, f.Ids)
	ginx.Dangerous(err)
	//删除资产扩展属性
	err = models.AssetsExpansionDelAssetsIds(tx, f.Ids)
	ginx.Dangerous(err)
	//删除资产所属监控
	err = models.MonitoringDelTxByAssetId(tx, f.Ids)
	ginx.Dangerous(err)
	//删除资产告警规则
	err = models.AlertRuleDelTxByAssetId(tx, f.Ids)
	tx.Commit()

	ginx.NewRender(c).Message(err)
}

func (rt *Router) assetDefaultConfigGet(c *gin.Context) {
	var f map[string]interface{}
	ginx.BindJSON(c, &f)

	assetType := ginx.UrlParamStr(c, "type")
	content, err := models.AssetGenConfig(assetType, f)

	ginx.NewRender(c).Data(map[string]string{"content": content.String()}, err)
}

func (rt *Router) assetIdentGetAll(c *gin.Context) {
	bgid := ginx.QueryInt64(c, "bgid", 0)
	data, err := models.TargetGets(rt.Ctx, []int64{bgid}, nil, "", 0, 100, 0)
	ginx.NewRender(c).Data(data, err)
}

// @Summary      ceshi
// @Description  根据条件查询资产详情/维保/管理
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Success      200 {array}  models.AssetType
// @Router       /api/n9e/assets/types [get]
// @Security     ApiKeyAuth
func (rt *Router) assetGetTypeList(c *gin.Context) {
	data, err := models.AssetTypeGetsAll()

	ginx.NewRender(c).Data(data, err)
}

func (rt *Router) assetGetTags(c *gin.Context) {
	ids := ginx.QueryStr(c, "ids")
	ids = strings.ReplaceAll(ids, ",", " ")
	lst, err := models.AssetGetTags(rt.Ctx, strings.Fields(ids))
	ginx.NewRender(c).Data(lst, err)
}

type assetsTagsForm struct {
	Ids  []string `json:"ids" binding:"required"`
	Tags []string `json:"tags" binding:"required"`
}

func (rt *Router) assetBindTagsByFE(c *gin.Context) {
	var f assetsTagsForm
	ginx.BindJSON(c, &f)

	if len(f.Ids) == 0 {
		ginx.Bomb(http.StatusBadRequest, "ids empty")
	}

	ginx.NewRender(c).Message(rt.assetsBindTags(f))
}

func (rt *Router) assetsBindTags(f assetsTagsForm) error {
	for i := 0; i < len(f.Tags); i++ {
		arr := strings.Split(f.Tags[i], "=")
		if len(arr) != 2 {
			return fmt.Errorf("invalid tag(%s)", f.Tags[i])
		}

		if strings.TrimSpace(arr[0]) == "" || strings.TrimSpace(arr[1]) == "" {
			return fmt.Errorf("invalid tag(%s)", f.Tags[i])
		}

		if strings.IndexByte(arr[0], '.') != -1 {
			return fmt.Errorf("invalid tagkey(%s): cannot contains . ", arr[0])
		}

		if strings.IndexByte(arr[0], '-') != -1 {
			return fmt.Errorf("invalid tagkey(%s): cannot contains -", arr[0])
		}

		if !model.LabelNameRE.MatchString(arr[0]) {
			return fmt.Errorf("invalid tagkey(%s)", arr[0])
		}
	}

	for i := 0; i < len(f.Ids); i++ {
		id, _ := strconv.Atoi(f.Ids[i])
		asset, err := models.AssetGetById(rt.Ctx, int64(id))
		if err != nil {
			return err
		}

		if asset == nil {
			continue
		}

		// 不能有同key的标签，否则附到时序数据上会产生覆盖，让人困惑
		for j := 0; j < len(f.Tags); j++ {
			tagkey := strings.Split(f.Tags[j], "=")[0]
			tagkeyPrefix := tagkey + "="
			if strings.HasPrefix(asset.Tags, tagkeyPrefix) {
				return fmt.Errorf("duplicate tagkey(%s)", tagkey)
			}
		}

		err = asset.AddTags(rt.Ctx, f.Tags)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rt *Router) assetUnbindTagsByFE(c *gin.Context) {
	var f assetsTagsForm
	ginx.BindJSON(c, &f)

	if len(f.Ids) == 0 {
		ginx.Bomb(http.StatusBadRequest, "ids empty")
	}

	ginx.NewRender(c).Message(rt.assetUnbindTags(f))
}

func (rt *Router) assetUnbindTags(f assetsTagsForm) error {
	for i := 0; i < len(f.Ids); i++ {
		id, _ := strconv.Atoi(f.Ids[i])
		asset, err := models.AssetGetById(rt.Ctx, int64(id))
		if err != nil {
			return err
		}

		if asset == nil {
			continue
		}

		err = asset.DelTags(rt.Ctx, f.Tags)
		if err != nil {
			return err
		}
	}
	return nil
}

type assetBgidForm struct {
	Ids  []string `json:"ids" binding:"required"`
	Bgid int64    `json:"bgid"`
}

func (rt *Router) assetUpdateBgid(c *gin.Context) {
	var f assetBgidForm
	ginx.BindJSON(c, &f)

	if len(f.Ids) == 0 {
		ginx.Bomb(http.StatusBadRequest, "ids empty")
	}

	user := c.MustGet("user").(*models.User)
	if user.IsAdmin() {
		ginx.NewRender(c).Message(models.AssetUpdateBgid(rt.Ctx, f.Ids, f.Bgid, false))
		return
	}

	ginx.NewRender(c).Message(models.TargetUpdateBgid(rt.Ctx, f.Ids, f.Bgid, false))
}

type assetNoteForm struct {
	Ids  []string `json:"ids" binding:"required"`
	Note string   `json:"note"`
}

func (rt *Router) assetUpdateNote(c *gin.Context) {
	var f assetNoteForm
	ginx.BindJSON(c, &f)

	if len(f.Ids) == 0 {
		ginx.Bomb(http.StatusBadRequest, "ids empty")
	}

	ginx.NewRender(c).Message(models.AssetUpdateNote(rt.Ctx, f.Ids, f.Note))
}

type AssetOrganizationForm struct {
	Ids []string `json:"ids" binding:"required"` //资产id组
	Id  int64    `json:"id"`                     //组织树id
}

func (rt *Router) assetUpdateOrganization(c *gin.Context) {
	var f AssetOrganizationForm
	ginx.BindJSON(c, &f)
	ginx.NewRender(c).Message(models.AssetUpdateOrganization(rt.Ctx, f.Ids, f.Id))
}

type Person struct {
	Name   string `xml:"name"`
	Age    int    `xml:"age"`
	Gender string `xml:"gender"`
}

//导入xml
// @Summary      导出xml（测试）
// @Description  导出xml
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /api/n9e/xh/assets/xml/ [get]
// @Security     ApiKeyAuth
func (rt *Router) xmlceshi(c *gin.Context) {
	person := Person{
		Name:   "John",
		Age:    25,
		Gender: "Male",
	}
	// c.XML(http.StatusOK, person)
	logger.Debug("111111111111111")
	xmlData, err := xml.MarshalIndent(person, "", "  ")
	if err != nil {
		fmt.Println("XML导出失败：", err)
		return
	}
	file, err := os.Create("person.xml")
	if err != nil {
		fmt.Println("创建文件失败：", err)
		return
	}
	defer file.Close()
	logger.Debug("2222222222222")
	_, err = file.Write(xmlData)
	if err != nil {
		fmt.Println("写入文件失败：", err)
		return
	}
	fmt.Println("XML导出成功！")
	//设置文件类型
	c.Header("Content-Type", "application/xml;charset=utf8")
	//设置文件名称
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape("xml测试"))
	c.Writer.Write(xmlData)
	logger.Debug("333333333333")
	ginx.NewRender(c)
}

// @Summary      导出资产模板
// @Description  导出资产模板
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Asset
// @Router       /api/n9e/xh/asset/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templeAssetXH(c *gin.Context) {

	datas := make([]interface{}, 0)
	datas = append(datas, models.Asset{})
	excels.NewMyExcel("资产信息").ExportTempletToWeb(datas, nil, "cn", "source", 1, rt.Ctx, c)
}

// @Summary      EXCEL导入资产
// @Description  EXCEL导入资产
// @Tags         资产-西航
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/xh/asset/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importAssetXH(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "上传文件出错")
	}
	//读excel流
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "读取excel文件失败")
	}

	//解析excel的数据
	assetImports, _, lxRrr := excels.ReadExce[models.Asset](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusBadRequest, "解析excel文件失败")
		return
	}
	logger.Debug(assetImports)
	me := c.MustGet("user").(*models.User)
	var qty int = 0
	for index, entity := range assetImports {

		//校验
		num, err := models.AssetsCountMap(rt.Ctx, map[string]interface{}{"ip": entity.Ip})
		ginx.Dangerous(err)
		if num > 0 {
			str := "第" + strconv.Itoa(index) + "数据，IP已存在"
			ginx.Bomb(http.StatusOK, str)
		}

		// 循环体
		var f models.Asset = entity
		f.CreateBy = me.Username
		f.AddXH(rt.Ctx)
		qty++
	}
	ginx.NewRender(c).Data(qty, nil)

}

type AssetXml struct {
	AssetData []models.AssetImport `xml:"asset_data"`
}

// @Summary      EXCEL导出资产
// @Description  EXCEL导出资产
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Param        ftype query   int     false  "类型"
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200
// @Router       /api/n9e/xh/asset/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) exportAssetXH(c *gin.Context) {

	fType := ginx.QueryInt64(c, "ftype", -1)

	var f map[string]interface{}
	ginx.BindJSON(c, &f)
	var lst []models.Asset
	var err error

	idsTemp, idsOk := f["ids"]
	ids := make([]int64, 0)
	// var err error
	if idsOk {
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, int64(val.(float64)))
		}
		lst, err = models.AssetGetByIds(rt.Ctx, ids)
		ginx.Dangerous(err)
	} else {
		query := ""
		queryTemp, queryOk := f["query"]

		queryType := ""
		filter, filterOk := f["filter"]
		if filterOk {
			if !queryOk {
				ginx.Bomb(http.StatusBadRequest, "参数为空")
			} else {
				query = "%" + queryTemp.(string) + "%"
			}

			if filter.(string) == "1" || filter.(string) == "4" {
				queryType = "name"
			} else if filter.(string) == "2" {
				queryType = "ip"
			} else if filter.(string) == "3" {
				queryType = "type"
			}
		} else {
			if queryOk {
				query = "%" + queryTemp.(string) + "%"
			}
		}
		aType := ""
		typeT, typeOk := f["type"]
		if typeOk {
			aType = typeT.(string)
		}

		lst, err = models.AssetsGetsFilter(rt.Ctx, aType, query, queryType, -1, -1)
		ginx.Dangerous(err)
	}

	datas := make([]interface{}, 0)
	if len(lst) > 0 {
		for _, v := range lst {
			datas = append(datas, v)

		}
	} else {

	}
	if fType == 1 {
		if len(lst) == 0 {
			datas = append(datas, models.Asset{})
			excels.NewMyExcel("资产信息").ExportTempletToWeb(datas, nil, "cn", "source", 0, rt.Ctx, c)
		} else {
			excels.NewMyExcel("资产信息").ExportDataInfo(datas, "cn", rt.Ctx, c)
		}
	} else if fType == 2 {

		var assetsImport []models.AssetImport
		for _, val := range lst {
			var assetImport models.AssetImport
			assetImport.Ip = val.Ip
			assetImport.Name = val.Name
			assetImport.Manufacturers = val.Manufacturers
			assetImport.Position = val.Position
			if val.Status == 0 {
				assetImport.Status = "下线"
			} else if val.Status == 1 {
				assetImport.Status = "正常"
			}

			assetImport.Type = val.Type
			assetsImport = append(assetsImport, assetImport)
		}
		dataXml := AssetXml{
			AssetData: assetsImport,
		}
		xmltool.ExportXml(c, dataXml, "资产信息")
	} else if fType == 3 {
		dataTxt := make([]string, 0)
		str := "名称\t类型\tIP\t厂商\t资产位置\t状态\n"
		dataTxt = append(dataTxt, str)
		for _, asset := range lst {
			status := ""
			if asset.Status == 0 {
				status = "下线"
			} else if asset.Status == 1 {
				status = "正常"
			}
			str = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\n", asset.Name, asset.Type, asset.Ip, asset.Manufacturers, asset.Position, status)
			dataTxt = append(dataTxt, str)
		}
		txt.ExportTxt(c, dataTxt, "资产信息")
	}

}
