package router

import (
	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func (rt *Router) organizationGets(c *gin.Context) {
	list, err := models.OrganizationList(rt.Ctx)
	ginx.Dangerous(err)
	ginx.NewRender(c).Data(list, nil)
}

func (rt *Router) organizationGet(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "id")
	event, err := models.OrganizationGetById(rt.Ctx, eid)
	ginx.Dangerous(err)

	if event == nil {
		ginx.Bomb(404, "No such organizaion")
	}

	ginx.NewRender(c).Data(event, err)
}

func (rt *Router) organizationPut(c *gin.Context) {
	var f models.Organization
	ginx.BindJSON(c, &f)

	m, err := models.OrganizationGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)

	if m == nil {
		ginx.Bomb(404, "No such organization")
	} else {
		err = m.Update(rt.Ctx, f)
		ginx.Dangerous(err)
	}
	ginx.NewRender(c).Message(err)
}

func (rt *Router) organizationDel(c *gin.Context) {
	orgId := ginx.UrlParamInt64(c, "id")
	org, err := models.OrganizationGetById(rt.Ctx, orgId)
	ginx.Dangerous(err)
	if org == nil {
		ginx.Bomb(404, "No such organization")
	}

	childrenCount, err := models.OrganizationCount(rt.Ctx, "parent_id = ?", org.Id)
	if childrenCount > 0 {
		ginx.Bomb(404, "This organization hava suborganization")
	}

	assetsCount, err := models.AssetCount(rt.Ctx, "organization_id = ?", org.Id)
	if assetsCount > 0 {
		ginx.Bomb(404, "This organization hava assets")
	}

	var list = []int64{orgId}

	ginx.NewRender(c).Message(models.OrganizationDels(rt.Ctx, list))

}

func (rt *Router) organizationAdd(c *gin.Context) {
	var f models.Organization
	ginx.BindJSON(c, &f)

	u := models.Organization{
		Name:     f.Name,
		ParentId: f.ParentId,
		Path:     f.Path,
	}

	ginx.NewRender(c).Message(u.Add(rt.Ctx))
}

// @Summary      根据Ids获取组织树
// @Description  根据Ids获取组织树
// @Tags         组织树管理
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "add ids"
// @Success      200 {array}  models.Organization
// @Router       /api/n9e/organization/name/ [post]
// @Security     ApiKeyAuth
func (rt *Router) organizationGetsByIds(c *gin.Context) {
	var f []int64
	ginx.BindJSON(c, &f)

	organizations, err := models.OrganizationGetsByIds(rt.Ctx, f)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(organizations, err)
}
