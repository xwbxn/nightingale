package router

import (
	"github.com/didi/nightingale/v5/src/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

type BusiOverview struct {
	GroupId    int64
	GroupName  string
	GroupLabel string
	Targets    int64
	Emergency  int64
	Warning    int64
	Notice     int64
}

func overviewGet(c *gin.Context) {
	me := c.MustGet("user").(*models.User)

	busiGroups, err := me.BusiGroups(-1, "")
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

		bov.Emergency, _ = models.AlertCurEventTotal("", group.Id, 0, 0, 1, []string{}, []string{}, "")
		bov.Emergency, _ = models.AlertCurEventTotal("", group.Id, 0, 0, 2, []string{}, []string{}, "")
		bov.Emergency, _ = models.AlertCurEventTotal("", group.Id, 0, 0, 3, []string{}, []string{}, "")
		bov.Targets, _ = models.TargetTotal(group.Id, []string{}, "")

		ret = append(ret, bov)
	}

	ginx.NewRender(c).Data(ret, err)
}
