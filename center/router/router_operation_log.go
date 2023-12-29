// Package models  操作日志
// date : 2023-10-21 09:13
// desc : 操作日志
package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	models "github.com/ccfos/nightingale/v6/models"
	excels "github.com/ccfos/nightingale/v6/pkg/excel"
	"github.com/ccfos/nightingale/v6/pkg/txt"
	xmltool "github.com/ccfos/nightingale/v6/pkg/xml"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// @Summary      获取操作日志
// @Description  根据主键获取操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.OperationLog
// @Router       /api/n9e/operation-log/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) operationLogGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	operationLog, err := models.OperationLogGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if operationLog == nil {
		ginx.Bomb(404, "No such operation_log")
	}

	ginx.NewRender(c).Data(operationLog, nil)
}

// @Summary      查询操作日志
// @Description  根据条件查询操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.OperationLog
// @Router       /api/n9e/operation-log/ [get]
// @Security     ApiKeyAuth
func (rt *Router) operationLogGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.OperationLogCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.OperationLogGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      过滤器查询操作日志
// @Description  过滤器条件查询操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        page query   int     false  "页码"
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Param        start query   int  false  "开始时间"
// @Param        end query   int  false  "结束时间"
// @Param        filterType query string  false  "类型"
// @Param        modelType query string  false  "模块类型"
// @Success      200  {array}  models.OperationLog
// @Router       /api/n9e/operation-log/filter [get]
// @Security     ApiKeyAuth
func (rt *Router) operationLogFilterGets(c *gin.Context) {
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")
	start := ginx.QueryInt64(c, "start", 0)
	end := ginx.QueryInt64(c, "end", time.Now().Unix())
	filterType := ginx.QueryStr(c, "filterType", "")
	modelTye := ginx.QueryStr(c, "modelType", "")

	if end == -1 {
		end = time.Now().Unix()
	}
	if end < start {
		ginx.Bomb(http.StatusOK, "时间区间错误")
	}

	total, err := models.FilterLogCount(rt.Ctx, query, start, end, filterType, modelTye)
	ginx.Dangerous(err)
	lst, err := models.FilterLogGets(rt.Ctx, query, (page-1)*limit, limit, start, end, filterType, modelTye)
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建操作日志
// @Description  创建操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        body  body   models.OperationLog true "add operationLog"
// @Success      200
// @Router       /api/n9e/operation-log/ [post]
// @Security     ApiKeyAuth
func (rt *Router) operationLogAdd(c *gin.Context) {
	var f models.OperationLog
	ginx.BindJSON(c, &f)

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.CreatedBy = me.Username

	// 更新模型
	err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(err)
}

// @Summary      更新操作日志
// @Description  更新操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        body  body   models.OperationLog true "update operationLog"
// @Success      200
// @Router       /api/n9e/operation-log/ [put]
// @Security     ApiKeyAuth
func (rt *Router) operationLogPut(c *gin.Context) {
	var f models.OperationLog
	ginx.BindJSON(c, &f)

	old, err := models.OperationLogGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "operation_log not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除操作日志
// @Description  根据主键删除操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/operation-log/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) operationLogDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	operationLog, err := models.OperationLogGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if operationLog == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(operationLog.Del(rt.Ctx))
}

type OperationLogXml struct {
	OperationLog []*models.OperationLogImport `xml:"operation_log"`
}

// @Summary      EXCEL导出操作日志
// @Description  EXCEL导出操作日志
// @Tags         操作日志
// @Accept       json
// @Produce      json
// @Param        ftype query   int     false  "类型"
// @Param        query query   string  false  "查询条件"
// @Param        start query   int  false  "开始时间"
// @Param        end query   int  false  "结束时间"
// @Param        filterType query string  false  "类型"
// @Param        modelType query string  false  "模块类型"
// @Param        body  body   map[string]interface{} true "add query"
// @Success      200
// @Router       /api/n9e/operation-log/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) exportOperationLog(c *gin.Context) {

	fType := ginx.QueryInt64(c, "ftype", -1)

	var f map[string]interface{}
	ginx.BindJSON(c, &f)
	var lst []*models.OperationLogImport
	var err error

	idsTemp, idsOk := f["ids"]
	ids := make([]int64, 0)
	// var err error
	if idsOk {
		for _, val := range idsTemp.([]interface{}) {
			ids = append(ids, int64(val.(float64)))
		}
		lst, err = models.OperationLogGetByIds(rt.Ctx, ids)
		ginx.Dangerous(err)
		for key := range lst {
			otime, err := strconv.ParseInt(lst[key].OperTime, 10, 64)
			ginx.Dangerous(err)
			lst[key].OperTime = time.Unix(otime, 0).Format("2006-01-02 15:04:05")
		}

	} else {
		query := ginx.QueryStr(c, "query", "")
		start := ginx.QueryInt64(c, "start", 0)
		end := ginx.QueryInt64(c, "end", time.Now().Unix())
		filterType := ginx.QueryStr(c, "filterType", "")
		modelTye := ginx.QueryStr(c, "modelType", "")

		if end == -1 {
			end = time.Now().Unix()
		}
		if end < start {
			ginx.Bomb(http.StatusOK, "时间区间错误")
		}

		logLst, err := models.FilterLogGets(rt.Ctx, query, -1, -1, start, end, filterType, modelTye)
		ginx.Dangerous(err)

		for _, val := range logLst {
			lst = append(lst, &models.OperationLogImport{
				Id:          val.Id,
				Type:        val.Type,
				Object:      val.Object,
				Description: val.Description,
				User:        val.User,
				OperTime:    time.Unix(val.OperTime, 0).Format("2006-01-02 15:04:05"),
			})
		}
	}

	datas := make([]interface{}, 0)
	if len(lst) > 0 {
		for _, v := range lst {
			datas = append(datas, *v)

		}
	}
	if fType == 1 {
		if len(lst) == 0 {
			datas = append(datas, models.OperationLogImport{})
			excels.NewMyExcel("操作日志").ExportTempletToWeb(datas, nil, "cn", "source", 0, rt.Ctx, c)
		} else {
			excels.NewMyExcel("操作日志").ExportDataInfo(datas, "cn", rt.Ctx, c)
		}
	} else if fType == 2 {

		dataXml := OperationLogXml{
			OperationLog: lst,
		}
		xmltool.ExportXml(c, dataXml, "操作日志")
	} else if fType == 3 {
		dataTxt := make([]string, 0)
		str := "日志编号\t操作类型\t系统模块\t描述\t操作人员\t操作时间\n"
		dataTxt = append(dataTxt, str)
		for _, log := range lst {
			str = fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s\n", log.Id, log.Type, log.Object, log.Description, log.User, log.OperTime)
			dataTxt = append(dataTxt, str)
		}
		txt.ExportTxt(c, dataTxt, "操作日志")
	}

}
