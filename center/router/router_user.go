// Package models  人员信息
// date : 2023-08-25 13:56
// desc : 人员信息
package router

import (
	"net/http"
	"strings"

	"github.com/ccfos/nightingale/v6/models"
	"github.com/ccfos/nightingale/v6/pkg/ormx"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

func (rt *Router) userFindAll(c *gin.Context) {
	list, err := models.UserGetAll(rt.Ctx)
	ginx.NewRender(c).Data(list, err)
}

// @Summary      过滤器查询
// @Description  过滤器查询
// @Tags         人员信息
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        organization_id  query   int64     false  "组织id"
// @Param        status  query   int64     false  "状态"
// @Param        query  query   string     false  "搜索栏"
// @Success      200  {array}  models.User
// @Success      200
// @Router       /api/n9e/users/update-property/ [post]
// @Security     ApiKeyAuth
func (rt *Router) userGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	organizationId := ginx.QueryInt64(c, "organization_id", -1)
	status := ginx.QueryInt64(c, "status", -1)
	query := ginx.QueryStr(c, "query", "")

	m := make(map[string]interface{})
	if organizationId != -1 {
		m["organization_id"] = organizationId
	}
	if status != -1 {
		m["status"] = status
	}

	total, err := models.UserCountMap(rt.Ctx, m, query)
	ginx.Dangerous(err)

	list, err := models.UserMap(rt.Ctx, m, query, limit, (page-1)*limit)
	ginx.Dangerous(err)

	user := c.MustGet("user").(*models.User)
	logger.Debug(list)
	logger.Debug(total)
	logger.Debug(user.IsAdmin())

	ginx.NewRender(c).Data(gin.H{
		"list":  list,
		"total": total,
		"admin": user.IsAdmin(),
	}, nil)
}

type userAddForm struct {
	Username       string       `json:"username" binding:"required"`
	Password       string       `json:"password" binding:"required"`
	Nickname       string       `json:"nickname"`
	Phone          string       `json:"phone"`
	Email          string       `json:"email"`
	Portrait       string       `json:"portrait"`
	Roles          []string     `json:"roles" binding:"required"`
	Contacts       ormx.JSONObj `json:"contacts"`
	Status         int64        `json:"status"`          //用户状态（1：启用；0：禁用）
	OrganizationId int64        `json:"organization_id"` //组织id
}

func (rt *Router) userAddPost(c *gin.Context) {
	var f userAddForm
	ginx.BindJSON(c, &f)

	password, err := models.CryptoPass(rt.Ctx, f.Password)
	ginx.Dangerous(err)

	if len(f.Roles) == 0 {
		ginx.Bomb(http.StatusBadRequest, "roles empty")
	}

	user := c.MustGet("user").(*models.User)

	u := models.User{
		Username:       f.Username,
		Password:       password,
		Nickname:       f.Nickname,
		Phone:          f.Phone,
		Email:          f.Email,
		Portrait:       f.Portrait,
		Roles:          strings.Join(f.Roles, " "),
		Status:         0,
		OrganizationId: f.OrganizationId,
		Contacts:       f.Contacts,
		CreateBy:       user.Username,
		UpdateBy:       user.Username,
	}

	ginx.NewRender(c).Message(u.Add(rt.Ctx))
}

func (rt *Router) userProfileGet(c *gin.Context) {
	user := User(rt.Ctx, ginx.UrlParamInt64(c, "id"))
	ginx.NewRender(c).Data(user, nil)
}

type userProfileForm struct {
	Nickname       string       `json:"nickname"`
	Phone          string       `json:"phone"`
	Email          string       `json:"email"`
	Roles          []string     `json:"roles"`
	Contacts       ormx.JSONObj `json:"contacts"`
	Status         int64        `json:"status"`          //用户状态（1：启用；0：禁用）
	OrganizationId int64        `json:"organization_id"` //组织id
}

func (rt *Router) userProfilePut(c *gin.Context) {
	var f userProfileForm
	ginx.BindJSON(c, &f)

	if len(f.Roles) == 0 {
		ginx.Bomb(http.StatusBadRequest, "roles empty")
	}

	target := User(rt.Ctx, ginx.UrlParamInt64(c, "id"))
	target.Nickname = f.Nickname
	target.Phone = f.Phone
	target.Email = f.Email
	target.Roles = strings.Join(f.Roles, " ")
	target.Contacts = f.Contacts
	target.Status = f.Status
	target.OrganizationId = f.OrganizationId
	target.UpdateBy = c.MustGet("username").(string)

	ginx.NewRender(c).Message(target.UpdateAllFields(rt.Ctx))
}

type userPasswordForm struct {
	Password string `json:"password" binding:"required"`
}

func (rt *Router) userPasswordPut(c *gin.Context) {
	var f userPasswordForm
	ginx.BindJSON(c, &f)

	target := User(rt.Ctx, ginx.UrlParamInt64(c, "id"))

	cryptoPass, err := models.CryptoPass(rt.Ctx, f.Password)
	ginx.Dangerous(err)

	ginx.NewRender(c).Message(target.UpdatePassword(rt.Ctx, cryptoPass, c.MustGet("username").(string)))
}

func (rt *Router) userDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	target, err := models.UserGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if target == nil {
		ginx.NewRender(c).Message(nil)
		return
	}

	ginx.NewRender(c).Message(target.Del(rt.Ctx))
}

// @Summary      查询人员列表
// @Description  查询人员列表
// @Tags         人员信息
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.userNameVo
// @Router       /api/n9e/user/getNames/ [get]
// @Security     ApiKeyAuth
func (rt *Router) userNameGets(c *gin.Context) {

	target, err := models.UserNameGets(rt.Ctx)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(target, nil)
}

// @Summary      批量修改用户状态/组织
// @Description  批量修改用户状态/组织
// @Tags         人员信息
// @Accept       json
// @Produce      json
// @Param        type  query   string     true  "status/organization"
// @Param        property  query   int     true  "新状态/新组织id"
// @Param        body  body   []int64 true "add ids"
// @Success      200
// @Router       /api/n9e/users/update-property/ [post]
// @Security     ApiKeyAuth
func (rt *Router) userPropertyUpdate(c *gin.Context) {

	propertyType := ginx.QueryStr(c, "type", "")
	if propertyType == "" {
		ginx.Bomb(http.StatusOK, "类型错误!")
	}
	newproperty := ginx.QueryInt64(c, "property", -1)
	if newproperty == -1 {
		ginx.Bomb(http.StatusOK, "参数错误!")
	}

	var ids []int64
	ginx.BindJSON(c, &ids)
	logger.Debug(ids)

	var err error
	if propertyType == "status" {
		err = models.UpdateBatch(rt.Ctx, ids, map[string]interface{}{"status": newproperty})
	} else if propertyType == "organization" {
		err = models.UpdateBatch(rt.Ctx, ids, map[string]interface{}{"organization_id": newproperty})
	}

	ginx.NewRender(c).Message(err)
}

// @Summary      批量删除用户
// @Description  批量删除用户(后续关联业务后，可能需要添加校验)
// @Tags         人员信息
// @Accept       json
// @Produce      json
// @Param        body  body   []int64 true "add ids"
// @Success      200
// @Router       /api/n9e/users/batch-del/ [post]
// @Security     ApiKeyAuth
func (rt *Router) userDels(c *gin.Context) {

	var ids []int64
	ginx.BindJSON(c, &ids)

	err := models.UpdateBatchDel(rt.Ctx, ids)

	ginx.NewRender(c).Message(err)
}
