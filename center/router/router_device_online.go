// Package models  设备上线下线记录表
// date : 2023-08-27 16:35
// desc : 设备上线下线记录表
package router

import (
	"net/http"

	models "github.com/ccfos/nightingale/v6/models"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

// @Summary      获取设备上线下线记录表
// @Description  根据主键获取设备上线下线记录表
// @Tags         设备上线下线记录表
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200  {object}  models.DeviceOnline
// @Router       /api/n9e/device-online/{id} [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceOnlineGet(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceOnline, err := models.DeviceOnlineGetById(rt.Ctx, id)
	ginx.Dangerous(err)

	if deviceOnline == nil {
		ginx.Bomb(404, "No such device_online")
	}

	ginx.NewRender(c).Data(deviceOnline, nil)
}

// @Summary      查询设备上线下线记录表
// @Description  根据条件查询设备上线下线记录表
// @Tags         设备上线下线记录表
// @Accept       json
// @Produce      json
// @Param        limit query   int     false  "返回条数"
// @Param        query query   string  false  "查询条件"
// @Success      200  {array}  models.DeviceOnline
// @Router       /api/n9e/device-online/ [get]
// @Security     ApiKeyAuth
func (rt *Router) deviceOnlineGets(c *gin.Context) {
	limit := ginx.QueryInt(c, "limit", 20)
	query := ginx.QueryStr(c, "query", "")

	total, err := models.DeviceOnlineCount(rt.Ctx, query)
	ginx.Dangerous(err)
	lst, err := models.DeviceOnlineGets(rt.Ctx, query, limit, ginx.Offset(c, limit))
	ginx.Dangerous(err)

	ginx.NewRender(c).Data(gin.H{
		"list":  lst,
		"total": total,
	}, nil)
}

// @Summary      创建设备上线下线记录表
// @Description  创建设备上线下线记录表
// @Tags         设备上线下线记录表
// @Accept       json
// @Produce      json
// @Param        body  body   map[string]interface{} true "add deviceOnline"
// @Success      200
// @Router       /api/n9e/device-online/ [post]
// @Security     ApiKeyAuth
func (rt *Router) deviceOnlineAdd(c *gin.Context) {
	f := make(map[string]interface{}, 0)
	ginx.BindJSON(c, &f)

	assetBasics := make([]map[string]interface{}, 0)
	var ids []int64
	var onLines []models.DeviceOnline
	des, ok := f["description"]
	var dep string
	if ok {
		dep = des.(string)
	}

	status := int64(f["device_status"].(float64))
	deviceTime := int64(f["device_time"].(float64))
	assetsInterface := f["asset"].([]interface{})
	var directory, clearConfig int64
	resourceFree := make([]string, 0)

	asset := make([]map[string]interface{}, 0)
	for _, val := range assetsInterface {
		asset = append(asset, val.(map[string]interface{}))
	}
	if status == 1 {
		clearConfig = int64(f["clear_config"].(float64))
		directory = int64(f["line_directory"].(float64))
		logger.Debug(clearConfig)
	} else if status == 3 {
		resourceFreeInterface := f["resource_free"].([]interface{})
		directory = int64(f["line_directory"].(float64))

		for _, val := range resourceFreeInterface {
			resourceFree = append(resourceFree, val.(string))
		}
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)

	for _, val := range asset {
		var onLine models.DeviceOnline
		if status == 1 {
			assetBasic := make(map[string]interface{})

			assetBasic["management_ip"] = val["management_ip"]
			assetBasic["updated_by"] = me.Username
			assetBasics = append(assetBasics, assetBasic)
		} else if status == 2 {
			assetBasic := make(map[string]interface{})

			assetBasic["management_ip"] = val["management_ip"]
			assetBasic["equipment_room"] = val["equipment_room"]
			assetBasic["owning_cabinet"] = val["owning_cabinet"]
			assetBasic["u_number"] = val["u_number"]
			assetBasic["device_status"] = status
			assetBasic["cabinet_location"] = val["cabinet_location"]
			assetBasic["related_service"] = val["related_service"]
			assetBasic["updated_by"] = me.Username
			assetBasics = append(assetBasics, assetBasic)
		}

		ids = append(ids, int64(val["id"].(float64)))

		onLine.DeviceStatus = status
		onLine.AssetId = int64(val["id"].(float64))
		onLine.LineAt = deviceTime
		if ok {
			onLine.Description = dep
		}

		if status == 2 {
			onLine.LineDirectory = int64(val["directory"].(float64))
		} else {
			onLine.LineDirectory = directory
		}
		onLine.CreatedBy = me.Username
		onLines = append(onLines, onLine)
	}

	tx := models.DB(rt.Ctx).Begin()

	for index, val := range ids {
		assetTree, err := models.AssetTreeGetByMap(rt.Ctx, map[string]interface{}{"type": "asset", "property_id": val})
		if err != nil {
			tx.Rollback()
			ginx.Dangerous(err)
		}
		if status == 1 || status == 2 {
			err = models.UpdateBasicTxMap(tx, ids[index], assetBasics[index])
			ginx.Dangerous(err)
		} else {
			for _, resVal := range resourceFree {
				err = models.UpdateBasicTxMap(tx, ids[index], map[string]interface{}{"device_status": status, "updated_by": me.Username})
				ginx.Dangerous(err)
				if resVal == "chear_business" {
					err = models.UpdateBasicTxMap(tx, ids[index], map[string]interface{}{"related_service": "", "updated_by": me.Username})
					ginx.Dangerous(err)
					//TODO 接触检测、清除告警信息（暂未开发）
				} else if resVal == "clear_managment_ip" {
					err = models.UpdateBasicTxMap(tx, ids[index], map[string]interface{}{"management_ip": "", "updated_by": me.Username})
					ginx.Dangerous(err)
				} else if resVal == "clear_manage_ip" {
					err = models.UpdateAssetExpansionMap(tx, map[string]interface{}{"asset_id": ids[index], "property_name": "production_ip"}, map[string]interface{}{"property_value": "", "updated_by": me.Username})
					ginx.Dangerous(err)
				} else if resVal == "clear_cabinet" {
					err = models.UpdateBasicTxMap(tx, ids[index], map[string]interface{}{"equipment_room": "", "updated_by": me.Username})
					ginx.Dangerous(err)
				} else if resVal == "clear_room" {
					err = models.UpdateBasicTxMap(tx, ids[index], map[string]interface{}{"owning_cabinet": "", "updated_by": me.Username})
					ginx.Dangerous(err)

				}
			}
		}

		err = models.UpdateTxTree(tx, map[string]interface{}{"id": assetTree.Id}, map[string]interface{}{"status": status, "parent_id": directory, "updated_by": me.Username})
		ginx.Dangerous(err)
	}
	err := models.DeviceOnlineTxBatchAdd(tx, onLines)
	ginx.Dangerous(err)

	tx.Commit()

	// for _, val := range treeIds {
	// 	err := models.AssetTreeDelById(rt.Ctx, val)
	// 	ginx.Dangerous(err)
	// }

	// 更新模型
	// err := f.Add(rt.Ctx)
	ginx.NewRender(c).Message(nil)
}

// @Summary      更新设备上线下线记录表
// @Description  更新设备上线下线记录表
// @Tags         设备上线下线记录表
// @Accept       json
// @Produce      json
// @Param        body  body   models.DeviceOnline true "update deviceOnline"
// @Success      200
// @Router       /api/n9e/device-online/ [put]
// @Security     ApiKeyAuth
func (rt *Router) deviceOnlinePut(c *gin.Context) {
	var f models.DeviceOnline
	ginx.BindJSON(c, &f)

	old, err := models.DeviceOnlineGetById(rt.Ctx, f.Id)
	ginx.Dangerous(err)
	if old == nil {
		ginx.Bomb(http.StatusOK, "device_online not found")
	}

	// 添加审计信息
	me := c.MustGet("user").(*models.User)
	f.UpdatedBy = me.Username

	// 可修改"*"为字段名称，实现更新部分字段功能
	ginx.NewRender(c).Message(old.Update(rt.Ctx, f, "*"))
}

// @Summary      删除设备上线下线记录表
// @Description  根据主键删除设备上线下线记录表
// @Tags         设备上线下线记录表
// @Accept       json
// @Produce      json
// @Param        id    path    string  true  "主键"
// @Success      200
// @Router       /api/n9e/device-online/{id} [delete]
// @Security     ApiKeyAuth
func (rt *Router) deviceOnlineDel(c *gin.Context) {
	id := ginx.UrlParamInt64(c, "id")
	deviceOnline, err := models.DeviceOnlineGetById(rt.Ctx, id)
	// 有错则跳出，无错则继续
	ginx.Dangerous(err)

	if deviceOnline == nil {
		ginx.NewRender(c).Message(nil)
		return
	}
	ginx.NewRender(c).Message(deviceOnline.Del(rt.Ctx))
}
