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

	"github.com/ccfos/nightingale/v6/models"
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
	Producer        string `json:"producer"`
	Os              string `json:"os"`
	Cpu             int64  `json:"cpu"`
	Memory          int64  `json:"memory"`
	PluginVersion   string `json:"plugin_version"`
	Location        string `json:"location"`
	AssetStatus     string `json:"asset_status"`
	DirectoryId     int64  `json:"directory_id"`
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

	pluginVersion := ""
	pluginVersionT, pluginVersionOk := f["plugin_version"]
	if pluginVersionOk {
		pluginVersion = pluginVersionT.(string)
	}
	location := ""
	locationT, locationOk := f["location"]
	if locationOk {
		location = locationT.(string)
	}
	assetStatus := ""
	assetStatusT, assetStatusOk := f["asset_status"]
	if assetStatusOk {
		assetStatus = assetStatusT.(string)
	}
	memo := ""
	memoT, memoOk := f["memo"]
	if memoOk {
		memo = memoT.(string)
	}

	var assets = models.Asset{
		Name:           f["name"].(string),
		Type:           f["type"].(string),
		Ip:             f["ip"].(string),
		Producer:       f["producer"].(string),
		Os:             f["os"].(string),
		Cpu:            int64(f["cpu"].(float64)),
		Memory:         int64(f["memory"].(float64)),
		PluginVersion:  pluginVersion,
		Location:       location,
		AssetStatus:    assetStatus,
		DirectoryId:    int64(f["directory_id"].(float64)),
		Memo:           memo,
		CreateBy:       me.Username,
		CreateAt:       time.Now().Unix(),
		OrganizationId: int64(f["organization_id"].(float64)),
	}

	tx := models.DB(rt.Ctx).Begin()
	assetId, err := assets.AddTx(tx)
	ginx.Dangerous(err)

	var assetsExpansion []models.AssetsExpansion
	dictDatas, err := models.DictDataGetByMap(rt.Ctx, map[string]interface{}{"type_code": "asset_ext_fields"})
	ginx.Dangerous(err)
	for _, val := range dictDatas {
		value, ok := f[val.DictKey]
		if ok {
			var expansion = models.AssetsExpansion{
				AssetsId:  assetId,
				NameCn:    val.DictValue,
				Name:      val.DictKey,
				Value:     value.(string),
				CreatedBy: me.Username,
			}
			assetsExpansion = append(assetsExpansion, expansion)
		}
	}
	err = models.AssetsExpansionAddTx(tx, assetsExpansion)
	tx.Commit()

	ginx.NewRender(c).Message(err)
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
	oldExpansions, err := models.AssetsExpansionGetsMap(rt.Ctx, map[string]interface{}{"assets_id": int64(f["id"].(float64))})
	ginx.Dangerous(err)
	me := c.MustGet("user").(*models.User)

	if oldAssets == nil {
		ginx.Bomb(http.StatusOK, "assets not found")
	}

	pluginVersion := ""
	pluginVersionT, pluginVersionOk := f["plugin_version"]
	if pluginVersionOk {
		pluginVersion = pluginVersionT.(string)
	}
	location := ""
	locationT, locationOk := f["location"]
	if locationOk {
		location = locationT.(string)
	}
	assetStatus := ""
	assetStatusT, assetStatusOk := f["asset_status"]
	if assetStatusOk {
		assetStatus = assetStatusT.(string)
	}
	memo := ""
	memoT, memoOk := f["memo"]
	if memoOk {
		memo = memoT.(string)
	}

	oldAssets.Name = f["name"].(string)
	oldAssets.Type = f["type"].(string)
	oldAssets.Ip = f["ip"].(string)
	oldAssets.Producer = f["producer"].(string)
	oldAssets.Os = f["os"].(string)
	oldAssets.Cpu = int64(f["cpu"].(float64))
	oldAssets.Memory = int64(f["memory"].(float64))
	oldAssets.PluginVersion = pluginVersion
	oldAssets.Location = location
	oldAssets.AssetStatus = assetStatus
	oldAssets.DirectoryId = int64(f["directory_id"].(float64))
	oldAssets.Memo = memo
	oldAssets.UpdateAt = time.Now().Unix()
	oldAssets.UpdateBy = me.Username

	tx := models.DB(rt.Ctx).Begin()

	err = oldAssets.UpdateTx(tx, "name", "type", "ip", "producer", "os", "cpu", "memory", "plugin_version", "location", "asset_status", "directory_id", "memo", "update_at", "update_by")
	ginx.Dangerous(err)
	for _, val := range oldExpansions {
		if f[val.Name].(string) != val.Value {
			models.AssetsExpansionUpdateTx(tx, map[string]interface{}{"id": val.Id}, map[string]interface{}{"value": f[val.Name], "updated_by": me.Username})
			ginx.Dangerous(err)
		}
	}
	tx.Commit()

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

// @Summary      过滤器-西航
// @Description  过滤器-西航
// @Tags         资产-西航
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200
// @Router       /api/n9e/xh/assets/filter [post]
// @Security     ApiKeyAuth
func (rt *Router) assetGetFilter(c *gin.Context) {
	var f map[string]interface{}
	ginx.BindJSON(c, &f)

	//分页参数
	limitT, limitOk := f["limit"]
	limit := int(limitT.(float64))
	if !limitOk {
		limit = 20
	}
	pageT, pageOk := f["page"]
	page := int(pageT.(float64))
	if !pageOk {
		page = 1
	}

	query := ""
	queryTemp, queryOk := f["query"]
	if !queryOk {
		query = ""
	}

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
	ids := make([]int64, 0)
	directoryId, directoryIdOk := f["directory_id"]
	if directoryIdOk {
		assetsDirTree, err := models.BuildDirTree(rt.Ctx, int64(directoryId.(float64)))
		ginx.Dangerous(err)
		ids, err = models.AssetsDirIds(rt.Ctx, assetsDirTree, ids)
		logger.Debug(ids)
		ginx.Dangerous(err)
	}
	total, err := models.AssetsCountFilter(rt.Ctx, ids, query, queryType)
	ginx.Dangerous(err)

	lst, err := models.AssetsGetsFilter(rt.Ctx, ids, query, queryType, limit, (page-1)*limit)
	ginx.Dangerous(err)

	for _, asset := range lst {
		atype, ok := rt.assetCache.GetType(asset.Type)
		if ok {
			asset.Dashboard = strings.ReplaceAll(atype.Dashboard, "${id}", fmt.Sprintf("%d", asset.Id))
		}
		health, ok := rt.assetCache.Get(asset.Id)
		if ok {
			asset.Health = health.Health
		}
	}

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
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
// @Param        body  body   assetsModel true "add assetsModel"
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
	err := models.AssetDelTx(tx, f.Ids)
	ginx.Dangerous(err)
	err = models.AssetsExpansionDelAssetsIds(tx, f.Ids)
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
