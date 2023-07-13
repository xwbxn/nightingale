package router

import (
	"strings"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func (rt *Router) groupsList(c *gin.Context) {

	stime, etime := getTimeRange(c)
	dsIds := queryDatasourceIds(c)
	busiGroupId := ginx.QueryInt64(c, "bgid", 0)
	query := ginx.QueryStr(c, "query", "")
	limit := ginx.QueryInt(c, "limit", 20)
	prod := ginx.QueryStr(c, "prods", "")
	name := ginx.QueryStr(c, "name", "")
	parent_id := ginx.QueryInt64(c, "parent_id", 0)
	path := ginx.QueryStr(c, "path", "")
	if prod == "" {
		prod = ginx.QueryStr(c, "rule_prods", "")
	}

	prods := []string{}
	if prod != "" {
		prods = strings.Split(prod, ",")
	}

	total, err := models.GroupsTotal(rt.Ctx, prods, busiGroupId, stime, etime, dsIds, query, name, parent_id, path)
	ginx.Dangerous(err)

	list, err := models.GroupsGets(rt.Ctx, prods, busiGroupId, stime, etime, dsIds, query, name, parent_id, path, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  list,
		"total": total,
	}, nil)
}

func (rt *Router) groupsGet(c *gin.Context) {
	eid := ginx.UrlParamInt64(c, "eid")
	event, err := models.GroupsGetById(rt.Ctx, eid)
	ginx.Dangerous(err)

	if event == nil {
		ginx.Bomb(404, "No such alert event")
	}

	ginx.NewRender(c).Data(event, err)
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
		err = m.UpdateAll(rt.Ctx, eid, f.Name, f.ParentId, f.Path)
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
