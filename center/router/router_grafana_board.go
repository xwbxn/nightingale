package router

import (
	"net/http"
	"time"

	"github.com/didi/nightingale/v5/src/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

type grafanaBoardForm struct {
	Name       string `json:"name"`
	Ident      string `json:"ident"`
	Tags       string `json:"tags"`
	Configs    string `json:"configs"`
	Public     int    `json:"public"`
	GrafanaId  int64  `json:"grafana_id"`
	GrafanaUrl string `json:"grafana_url"`
}

func grafanaBoardAdd(c *gin.Context) {
	var f grafanaBoardForm
	ginx.BindJSON(c, &f)

	me := c.MustGet("user").(*models.User)

	board := &models.GrafanaBoard{
		GroupId:    ginx.UrlParamInt64(c, "id"),
		Name:       f.Name,
		Ident:      f.Ident,
		Tags:       f.Tags,
		Configs:    f.Configs,
		GrafanaId:  f.GrafanaId,
		GrafanaUrl: f.GrafanaUrl,
		CreateBy:   me.Username,
		UpdateBy:   me.Username,
	}

	err := board.Add()
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(board, nil)
}

func grafanaBoardGet(c *gin.Context) {
	bid := ginx.UrlParamStr(c, "bid")
	board, err := models.GrafanaBoardGet("id = ? or ident = ?", bid, bid)
	ginx.Dangerous(err)

	if board == nil {
		ginx.Bomb(http.StatusNotFound, "No such dashboard")
	}

	if board.Public == 0 {
		auth()(c)
		user()(c)

		bgroCheck(c, board.GroupId)
	}

	ginx.NewRender(c).Data(board, nil)
}

func grafanaBoardPureGet(c *gin.Context) {
	board, err := models.GrafanaBoardGetByID(ginx.UrlParamInt64(c, "bid"))
	ginx.Dangerous(err)

	if board == nil {
		ginx.Bomb(http.StatusNotFound, "No such dashboard")
	}

	ginx.NewRender(c).Data(board, nil)
}

// bgrwCheck
func grafanaBoardDel(c *gin.Context) {
	var f idsForm
	ginx.BindJSON(c, &f)
	f.Verify()

	for i := 0; i < len(f.Ids); i++ {
		bid := f.Ids[i]

		board, err := models.GrafanaBoardGet("id = ?", bid)
		ginx.Dangerous(err)

		if board == nil {
			continue
		}

		// check permission
		bgrwCheck(c, board.GroupId)

		ginx.Dangerous(board.Del())
	}

	ginx.NewRender(c).Message(nil)
}

func GrafanaBoard(id int64) *models.GrafanaBoard {
	obj, err := models.GrafanaBoardGet("id=?", id)
	ginx.Dangerous(err)

	if obj == nil {
		ginx.Bomb(http.StatusNotFound, "No such dashboard")
	}

	return obj
}

// bgrwCheck
func grafanaBoardPut(c *gin.Context) {
	var f grafanaBoardForm
	ginx.BindJSON(c, &f)

	me := c.MustGet("user").(*models.User)
	bo := GrafanaBoard(ginx.UrlParamInt64(c, "bid"))

	// check permission
	bgrwCheck(c, bo.GroupId)

	can, err := bo.GrafanaBoardCanRenameIdent(f.Ident)
	ginx.Dangerous(err)

	if !can {
		ginx.Bomb(http.StatusOK, "Ident duplicate")
	}

	bo.Name = f.Name
	bo.Ident = f.Ident
	bo.Tags = f.Tags
	bo.GrafanaId = f.GrafanaId
	bo.GrafanaUrl = f.GrafanaUrl
	bo.UpdateBy = me.Username
	bo.UpdateAt = time.Now().Unix()

	err = bo.Update("name", "ident", "tags", "grafana_id", "grafana_url", "update_by", "update_at")
	ginx.NewRender(c).Data(bo, err)
}

func grafanaBoardGets(c *gin.Context) {
	bgid := ginx.UrlParamInt64(c, "id")
	query := ginx.QueryStr(c, "query", "")

	boards, err := models.GrafanaBoardGets(bgid, query)
	ginx.NewRender(c).Data(boards, err)
}
