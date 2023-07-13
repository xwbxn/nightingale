package router

import (
	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func (rt *Router) groupsList(c *gin.Context) {
	list, err := models.GroupsList(rt.Ctx)
	ginx.Dangerous(err)
	// var apilist []GroupTreeModel
	ginx.NewRender(c).Data(list, nil)
}

func (rt *Router) groupsGet(c *gin.Context) {
	// eid := ginx.UrlParamInt64(c, "eid")
	// event, err := models.GroupsGetById(rt.Ctx, eid)
	// ginx.Dangerous(err)

	// if event == nil {
	// 	ginx.Bomb(404, "No such alert event")
	// }

	// ginx.NewRender(c).Data(event, err)
}

func (rt *Router) updateGroups(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	m, err := models.GroupsGetById(rt.Ctx, eid)
	ginx.Dangerous(err)

	if m == nil {
		ginx.Bomb(404, "No such group")
	} else {
		var f models.Groups
		ginx.BindJSON(c, &f)
		// err = m.UpdateAll(rt.Ctx, eid, f.Name, f.ParentId, f.Path)
		ginx.Dangerous(err)
	}
	ginx.NewRender(c).Message(err)
}

func (rt *Router) groupsDel(c *gin.Context) {
	var f idsForm
	ginx.BindJSON(c, &f)
	f.Verify()

	// param(busiGroupId) for protect
	ginx.NewRender(c).Message(models.GroupsDels(rt.Ctx, f.Ids, ginx.UrlParamInt64(c, "id")))
}
