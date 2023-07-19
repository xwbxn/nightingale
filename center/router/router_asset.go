package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/prometheus/common/model"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func (rt *Router) assetsGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	assets, err := models.AssetGetById(rt.Ctx, id)
	ginx.NewRender(c).Data(assets, err)
}

func (rt *Router) assetsGets(c *gin.Context) {
	bgid := ginx.QueryInt64(c, "bgid", -1)
	query := ginx.QueryStr(c, "query", "")
	organizeId := ginx.QueryInt64(c, "organize_id", -1)
	assets, err := models.AssetGets(rt.Ctx, bgid, query, organizeId)
	ginx.NewRender(c).Data(assets, err)
}

type assetsModel struct {
	Id         int64  `json:"id"`
	Version    string `json:"version"`
	Ident      string `json:"ident"`
	GroupId    int64  `json:"group_id"`
	Name       string `json:"name"`
	Label      string `json:"label"`
	Tags       string `json:"tags"`
	Type       string `json:"type"`
	Memo       string `json:"memo"`
	Configs    string `json:"configs"`
	Params     string `json:"params"`
	OrganizeId int64  `json:"organize_id"`
}

func (rt *Router) assetsAdd(c *gin.Context) {
	var f assetsModel
	ginx.BindJSON(c, &f)
	me := c.MustGet("user").(*models.User)

	var assets = models.Asset{
		GroupId:    f.GroupId,
		Name:       f.Name,
		Ident:      f.Ident,
		Label:      f.Label,
		Type:       f.Type,
		Memo:       f.Memo,
		Configs:    f.Configs,
		Params:     f.Params,
		CreateBy:   me.Username,
		CreateAt:   time.Now().Unix(),
		OrganizeId: f.OrganizeId,
	}

	err := assets.Add(rt.Ctx)
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

type assetsForm struct {
	Ids []string `json:"ids" binding:"required"`
}

func (rt *Router) assetDel(c *gin.Context) {
	var f assetsForm
	ginx.BindJSON(c, &f)

	if len(f.Ids) == 0 {
		ginx.Bomb(http.StatusBadRequest, "ids empty")
	}

	ginx.NewRender(c).Message(models.AssetDel(rt.Ctx, f.Ids))
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
	data, err := models.TargetGets(rt.Ctx, bgid, nil, "", 100, 0)
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

type Assetorganize struct {
	Ids []int64 `json:"ids" binding:"required"` //资产id组
	Id  int64   `json:"id"`                     //组织树id
}

func (rt *Router) updatesAssetOrganize(c *gin.Context) {
	var f Assetorganize
	ginx.BindJSON(c, &f)
	ginx.NewRender(c).Message(models.UpdateOrganize(rt.Ctx, f.Ids, f.Id))
}
