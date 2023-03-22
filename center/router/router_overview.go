package router

import (
	"github.com/didi/nightingale/v5/src/models"
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
		bov.Warning, _ = models.AlertCurEventTotal("", group.Id, 0, 0, 2, []string{}, []string{}, "")
		bov.Notice, _ = models.AlertCurEventTotal("", group.Id, 0, 0, 3, []string{}, []string{}, "")
		bov.Targets, _ = models.TargetTotal(group.Id, []string{}, "")

		ret = append(ret, bov)
	}

	ginx.NewRender(c).Data(ret, err)
}
