package router

import (
	"github.com/didi/nightingale/v5/src/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

func overviewGet(c *gin.Context) {
	me := c.MustGet("user").(*models.User)

	busiGroups, err := me.BusiGroups(-1, "")
	if err != nil {
		logger.Errorf("get busiGroup fail: %v", err)
		ginx.Dangerous(err)
	}

	bgids := []int64{}
	for _, group := range busiGroups {
		bgids = append(bgids, group.Id)
	}

	lst, err := models.AlertCurEventOverview(bgids)
	if err != nil {
		ginx.Dangerous(err)
	}

	ginx.NewRender(c).Data(lst, err)
}
