package router

import (
	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func (rt *Router) organizeList(c *gin.Context) {
	list, err := models.OrganizeList(rt.Ctx)
	ginx.Dangerous(err)
	ginx.NewRender(c).Data(list, nil)
}

func (rt *Router) organizeGet(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	event, err := models.OrganizeGetById(rt.Ctx, eid)
	ginx.Dangerous(err)

	if event == nil {
		ginx.Bomb(404, "No such organize")
	}

	ginx.NewRender(c).Data(event, err)
}

func (rt *Router) updateOrganize(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	m, err := models.OrganizeGetById(rt.Ctx, eid)
	ginx.Dangerous(err)
	var f models.Organize
	ginx.BindJSON(c, &f)
	if m == nil {
		ginx.Bomb(404, "No such organize")
	} else {
		err = m.UpdateAll(rt.Ctx, eid, f.Name, f.ParentId, f.Path)
		ginx.Dangerous(err)
	}
	ginx.NewRender(c).Message(err)
}

func (rt *Router) organizeDel(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	org, err := models.OrganizeGetById(rt.Ctx, eid)
	ginx.Dangerous(err)
	if org == nil {
		ginx.Bomb(404, "No such organize")
	}

	//TODO: 这里不应进行遍历，待优化
	lss, err := models.OrganizeGets(rt.Ctx)
	for i := 0; i < len(lss); i++ {
		if lss[i].ParentId == org.Id {
			ginx.Bomb(404, "This organize hava suborganize")
		}

	}
	lst, err := models.FindAssetByOrg(rt.Ctx, org.Id)
	if len(lst) > 0 {
		ginx.Bomb(404, "This organize hava assets")
	}

	var list = []int64{eid}

	ginx.NewRender(c).Message(models.OrganizeDels(rt.Ctx, list))

}

func (rt *Router) organizeAdd(c *gin.Context) {
	var f models.Organize
	ginx.BindJSON(c, &f)

	u := models.Organize{
		Name:     f.Name,
		ParentId: f.ParentId,
		Path:     f.Path,
	}

	ginx.NewRender(c).Message(u.Add(rt.Ctx))
}

// type IdOrganize struct {
// 	ID int `json:"id"`
// }

// type IdsOrganize struct {
// 	Ids        []int64 `json:"ids"`         //资产id组
// 	OrganizeId int64   `json:"organize_id"` //组织id

// }

// func (rt *Router) updatesOrganize(c *gin.Context) {
// 	var f IdsOrganize
// 	ginx.BindJSON(c, &f)
// 	models.UpdateOrganize(rt.Ctx, f.Ids, f.OrganizeId)
// }
