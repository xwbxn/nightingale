// Package models  人员信息
// date : 2023-08-25 13:56
// desc : 人员信息
package router

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
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
// @Param        user_group_id  query   int64     false  "用户组id"
// @Param        role  query   string     false  "角色"
// @Param        status  query   int64     false  "状态"
// @Param        query  query   string     false  "搜索栏"
// @Success      200  {array}  models.User
// @Success      200
// @Router       /api/n9e/users/xh [get]
// @Security     ApiKeyAuth
func (rt *Router) userGetsXH(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	useGroupId := ginx.QueryInt64(c, "user_group_id", -1)
	role := ginx.QueryStr(c, "role", "")
	status := ginx.QueryInt64(c, "status", -1)
	query := ginx.QueryStr(c, "query", "")
	userIds := make([]int64, 0)
	var err error

	if useGroupId != -1 {
		userIds, err = models.MemberIds(rt.Ctx, useGroupId)
		ginx.Dangerous(err)
	}

	total, err := models.UserCountMap(rt.Ctx, role, query, useGroupId, status, userIds)
	ginx.Dangerous(err)

	list, err := models.UserMap(rt.Ctx, role, query, useGroupId, status, userIds, limit, (page-1)*limit)
	ginx.Dangerous(err)

	user := c.MustGet("user").(*models.User)

	ginx.NewRender(c).Data(gin.H{
		"list":  list,
		"total": total,
		"admin": user.IsAdmin(),
	}, nil)
}

func (rt *Router) userGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.UserTotal(rt.Ctx, query)
	ginx.Dangerous(err)

	list, err := models.UserGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	user := c.MustGet("user").(*models.User)

	ginx.NewRender(c).Data(gin.H{
		"list":  list,
		"total": total,
		"admin": user.IsAdmin(),
	}, nil)
}

// @Summary      用户过滤器查询
// @Description  用户过滤器查询
// @Tags         人员信息
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "条数"
// @Param        type query   string     false  "类型"
// @Param        role  query   string     false  "角色"
// @Param        status  query   int64     false  "状态"
// @Param        query  query   string     false  "搜索栏"
// @Success      200  {array}  models.User
// @Success      200
// @Router       /api/n9e/users [get]
// @Security     ApiKeyAuth
func (rt *Router) userFilterGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	typeF := ginx.QueryStr(c, "type", "")
	role := ginx.QueryStr(c, "role", "")
	status := ginx.QueryInt64(c, "status", -1)
	query := ginx.QueryStr(c, "query", "")

	var err error

	total, err := models.UserFilterCountMap(rt.Ctx, role, query, status, typeF)
	ginx.Dangerous(err)

	list, err := models.UserFilterMap(rt.Ctx, role, query, status, limit, (page-1)*limit, typeF)
	ginx.Dangerous(err)

	user := c.MustGet("user").(*models.User)

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
	GroupName      []int64      `json:"group_name"`
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

	ginx.NewRender(c).Message(u.Add(rt.Ctx, f.GroupName))
}

func (rt *Router) userProfileGet(c *gin.Context) {
	user, err := models.UserGetXH(rt.Ctx, ginx.UrlParamInt64(c, "id"))
	ginx.NewRender(c).Data(user, err)
}

type userProfileForm struct {
	Nickname       string       `json:"nickname"`
	Phone          string       `json:"phone"`
	Email          string       `json:"email"`
	Roles          []string     `json:"roles"`
	Contacts       ormx.JSONObj `json:"contacts"`
	Status         int64        `json:"status"`          //用户状态（1：启用；0：禁用）
	OrganizationId int64        `json:"organization_id"` //组织id
	GroupName      []int64      `json:"group_name"`
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

	ginx.NewRender(c).Message(target.UpdateAllFields(rt.Ctx, f.GroupName))
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
	if id == 1 {
		ginx.Bomb(http.StatusOK, "管理员账号不能被删除")
	}
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
	for _, id := range ids {
		if id == 1 {
			ginx.Bomb(http.StatusOK, "管理员账号不能被删除")
		}
	}

	err := models.UserBatchDel(rt.Ctx, ids)

	ginx.NewRender(c).Message(err)
}

// @Summary      导入用户模板
// @Description  导入用户模板
// @Tags         人员信息
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.UserImport
// @Router       /api/n9e/xh/users/templet [post]
// @Security     ApiKeyAuth
func (rt *Router) templeUserXH(c *gin.Context) {

	datas := make([]interface{}, 0)
	datas = append(datas, models.UserImport{})
	excels.NewMyExcel("用户信息").ExportTempletToWeb(datas, nil, "cn", "source", 1, rt.Ctx, c)
}

// @Summary      EXCEL导入用户信息
// @Description  EXCEL导入用户信息
// @Tags         人员信息
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "file"
// @Success      200
// @Router       /api/n9e/xh/users/import-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) importUserXH(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		ginx.Bomb(http.StatusOK, "上传文件出错")
	}
	//go gin 校验文件是否是excel文件
	//读excel流
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		ginx.Bomb(http.StatusOK, "文件类型错误，请上传EXCEL文件（.xlsx,.xls）")
	}
	// 判断文件中是否有有效的Sheet
	if len(xlsx.GetSheetMap()) == 0 {
		ginx.Bomb(http.StatusOK, "文件不存在sheet")
	}

	//解析excel的数据
	userImports, _, lxRrr := excels.ReadExce[models.UserImport](xlsx, rt.Ctx)
	if lxRrr != nil {
		ginx.Bomb(http.StatusOK, "解析excel文件失败")
		return
	}
	logger.Debug(userImports)
	me := c.MustGet("user").(*models.User)
	// contacts := make(map[string]string)
	contacts, _ := json.Marshal(map[string]string{})
	tx := models.DB(rt.Ctx).Begin()
	for index, entity := range userImports {

		if entity.Password != entity.IsPassword {
			ginx.Bomb(http.StatusOK, "第"+strconv.Itoa(index+1)+"行数据，密码不一致")
		}

		password, err := models.CryptoPass(rt.Ctx, entity.Password)
		ginx.Dangerous(err)

		roles := ""
		if entity.Role1 != "" {
			roles += entity.Role1
		}
		if entity.Role2 != "" && entity.Role1 != entity.Role2 {
			roles += " "
			roles += entity.Role2
		}
		if entity.Role3 != "" && entity.Role1 != entity.Role3 && entity.Role2 != entity.Role3 {
			roles += " "
			roles += entity.Role3
		}
		entity.Contacts.UnmarshalJSON(contacts)

		u := models.User{
			Username:       entity.Username,
			Password:       password,
			Nickname:       entity.Nickname,
			Phone:          entity.Phone,
			Email:          entity.Email,
			Portrait:       "",
			Roles:          roles,
			Status:         0,
			OrganizationId: 0,
			Contacts:       entity.Contacts,
			CreateBy:       me.Username,
			UpdateBy:       me.Username,
		}

		// 校验
		olduser, err := models.UserGetByUsername(rt.Ctx, u.Username)
		if err != nil {
			tx.Rollback()
		}
		ginx.Dangerous(err)

		if olduser != nil {
			ginx.Bomb(http.StatusOK, "第"+strconv.Itoa(index)+"数据，用户名已存在")
		}
		err = u.AddTx(rt.Ctx, tx)
		ginx.Dangerous(err)
	}
	tx.Commit()
	ginx.NewRender(c).Data(err, nil)
}
