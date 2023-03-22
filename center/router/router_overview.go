package router

import (
	"github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

type BusiOverview struct {
	GroupId    int64  `json:"id"`
	GroupName  string `json:"name"`
	GroupLabel string `json:"label"`
	Targets    int64  `json:"targets"`
	Emergency  int64  `json:"emergency"`
	Warning    int64  `json:"warning"`
	Notice     int64  `json:"notice"`
}

func (rt *Router) overviewGet(c *gin.Context) {
	me := c.MustGet("user").(*models.User)

	busiGroups, err := me.BusiGroups(rt.Ctx, -1, "")
	if err != nil {
		logger.Errorf("get busiGroup fail: %v", err)
		ginx.Dangerous(err)
	}

	ret := []BusiOverview{}
	for _, group := range busiGroups {
		bov := BusiOverview{}
		bov.GroupId = group.Id
		bov.GroupName = group.Name
		bov.GroupLabel = group.LabelValue

		bov.Emergency, _ = models.AlertCurEventTotal(rt.Ctx, []string{}, group.Id, 0, 0, 1, []int64{}, []string{}, "")
		bov.Warning, _ = models.AlertCurEventTotal(rt.Ctx, []string{}, group.Id, 0, 0, 2, []int64{}, []string{}, "")
		bov.Notice, _ = models.AlertCurEventTotal(rt.Ctx, []string{}, group.Id, 0, 0, 3, []int64{}, []string{}, "")
		bov.Targets, _ = models.TargetTotal(rt.Ctx, group.Id, []int64{}, "")

		ret = append(ret, bov)
	}

	ginx.NewRender(c).Data(ret, err)
}
