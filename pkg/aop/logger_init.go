package aop

import (

	// "ginDemo/config"

	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toolkits/pkg/logger"
	"gorm.io/gorm"
)

var log_obj map[string]string
var log_post map[string]string
var log_put map[string]string
var log_delete map[string]string

type OperationLog struct {
	Id          int64          `gorm:"column:ID;primaryKey" json:"id" `                          //type:BIGINT       comment:日志主键      version:2023-10-21 09:14
	Type        string         `gorm:"column:TYPE" json:"type" `                                 //type:string       comment:类型          version:2023-10-21 09:10
	Object      string         `gorm:"column:OBJECT" json:"object" `                             //type:string       comment:对象          version:2023-10-21 09:10
	Description string         `gorm:"column:DESCRIPTION" json:"description" `                   //type:string       comment:描述          version:2023-10-21 09:10
	User        string         `gorm:"column:USER" json:"user" `                                 //type:string       comment:用户          version:2023-10-21 09:10
	OperTime    int64          `gorm:"column:OPER_TIME" json:"oper_time" `                       //type:*int         comment:执行时间      version:2023-10-21 09:10
	OperUrl     string         `gorm:"column:OPER_URL" json:"oper_url" `                         //type:string       comment:请求URL       version:2023-10-21 09:10
	OperParam   string         `gorm:"column:OPER_PARAM" json:"oper_param" `                     //type:string       comment:请求参数      version:2023-10-21 09:10
	JsonResult  string         `gorm:"column:JSON_RESULT" json:"json_result" `                   //type:string       comment:返回参数      version:2023-10-21 09:10
	Status      int64          `gorm:"column:STATUS" json:"status" `                             //type:*int         comment:操作状态码    version:2023-10-21 09:10
	ErrorMsg    string         `gorm:"column:ERROR_MSG" json:"error_msg" `                       //type:string       comment:错误消息      version:2023-10-21 09:10
	CreatedBy   string         `gorm:"column:CREATED_BY" json:"created_by" swaggerignore:"true"` //type:string       comment:创建人        version:2023-10-21 09:10
	CreatedAt   int64          `gorm:"column:CREATED_AT" json:"created_at" swaggerignore:"true"` //type:*int         comment:创建时间      version:2023-10-21 09:10
	UpdatedBy   string         `gorm:"column:UPDATED_BY" json:"updated_by" swaggerignore:"true"` //type:string       comment:更新人        version:2023-10-21 09:10
	UpdatedAt   int64          `gorm:"column:UPDATED_AT" json:"updated_at" swaggerignore:"true"` //type:*int         comment:更新时间      version:2023-10-21 09:10
	DeletedAt   gorm.DeletedAt `gorm:"column:DELETED_AT" json:"deleted_at" swaggerignore:"true"` //type:*time.Time   comment:删除时间      version:2023-10-21 09:10
}

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	logger.Debug("开始记录日志了")

	// logFilePath := config.Log_FILE_PATH
	// logFileName := config.LOG_FILE_NAME

	//日志文件
	// fileName := path.Join(logFilePath, logFileName)

	//写入文件
	// src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	// if err != nil {
	// 	fmt.Println("err", err)
	// }

	//实例化
	logger := logrus.New()

	//设置输出
	// logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{TimestampFormat: "2006-01-02 15:04:05"})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		params := c.Params
		// 将 gin.Params 转化为字符串
		paramsStr := ""
		for _, param := range params {
			paramsStr += param.Key + ":" + param.Value + " "
		}

		key := c.Keys

		accepted := c.Accepted

		url := c.Request.URL

		body := c.Request.PostForm

		userName := ""
		user, ok := key["username"]
		if ok {
			userName = user.(string)
		}

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		err := c.Errors.String()

		// remoteIP, _ := c.Request.Response.Location()

		log_obj = map[string]string{
			"xh/assets": "资产", "xh/monitoring": "监控", "alert": "告警", "user-config": "用户配置",
			"user": "用户信息", "api-service": "接口管理", "/xh/license": "许可管理", "login": "登录", "operation-log": "操作日志",
		}

		log_post = map[string]string{
			"xh/assets": "创建资产", "xh/assets/batch-del": "批量删除资产", "xh/assets/batch-update": "批量修改资产",
			"xh/asset/export-xls": "导出资产", "xh/asset/import-xls": "导入资产", "xh/asset/templet": "导出资产模板",
			"xh/monitoring": "创建监控", "xh/monitoring/data": "获取监控数据", "xh/monitoring/status": "监控开关", "xh/monitoring/batch-del": "批量删除监控",
			"alert-cur-events/batch-del": "批量删除当前告警", "alert-his-events/batch-del": "批量删除历史告警",
			"alert-events/export-xls": "导出告警", "user-config/picture": "选择图片上传",
			"users/batch-del": "批量删除用户", "users/update-property": "批量修改用户状态/组织",
			"users/import-xls": "导入用户信息", "users/templet": "导入用户模板", "users": "创建用户",
			"api-service": "创建接口管理", "xh/license/add-file": "上传证书", "xh/license-config": "创建许可管理", "xh/license/export-xls": "批量导出证书",
			"login": "登录", "/operation-log/export-xls": "导出操作日志数据",
		}
		log_put = map[string]string{
			"xh/assets": "更新资产", "xh/assets-expansion": "更新资产扩展", "xh/assets/batch-update": "批量修改资产",
			"xh/monitoring": "更新监控", "user-config": "更新用户配置", "user-config/logo": "更新用户图标",
			"xh/license/update": "更新证书", "xh/license-config": "更新许可配置", "api-service": "更新接口",
		}
		log_delete = map[string]string{
			"xh/monitoring": "删除监控", "user-config": "删除用户配置", "api-service": "删除接口",
		}

		operType := GetOperType(reqMethod, reqUri)
		operObj := GetOperObj(reqUri, userName)
		operDes := GetOperDes(reqMethod, reqUri, clientIP)
		// operuser := key["username"].(string)
		// var operLog = OperationLog{
		// 	Type:        operType,
		// 	Object:      operObj,
		// 	Description: operDes,
		// 	// User:        operuser,
		// 	OperTime:   startTime.Unix(),
		// 	OperUrl:    reqUri,
		// 	OperParam:  "",
		// 	JsonResult: "",
		// 	Status:     int64(statusCode),
		// 	ErrorMsg:   "",
		// 	CreatedBy:  "root",
		// 	CreatedAt:  time.Now().Unix(),
		// 	UpdatedBy:  "",
		// 	UpdatedAt:  time.Now().Unix(),
		// }
		parts := strings.Split(reqUri, "?")
		data := make(map[string]interface{})
		data["type"] = operType
		data["object"] = operObj
		data["user_name"] = userName
		data["description"] = operDes
		data["oper_time"] = startTime.Unix()
		data["oper_url"] = parts[0]
		data["oper_param"] = paramsStr
		data["json_result"] = ""
		data["req_method"] = reqMethod
		data["status"] = int64(statusCode)
		data["error_msg"] = err

		if strings.Contains(reqUri, "/api/n9e") && reqMethod != "GET" {

			if operType != "" {
				if !OperationlogQueue.PushFront(data) {
					logger.Warningf("event_push_queue: queue is full, event:%+v", data)
				}
				logger.Debug("日志进队列了")
			}

		}
		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
			"type":         operType,
			"params":       params,
			"key":          key,
			"accepted":     accepted,
			"url":          url,
			"body":         body,
			// "remoteIP":     remoteIP,
		}).Info()
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetOperType(m string, r string) string {
	var str strings.Builder
	order_log_post := orderedMap(log_post)
	order_log_put := orderedMap(log_put)
	order_log_delete := orderedMap(log_delete)
	if m == "PUT" {
		for _, k := range order_log_put {
			if strings.Contains(r, k) {
				str.Reset()
				str.WriteString(log_put[k])
			}
		}
	} else if m == "DELETE" {
		for _, k := range order_log_delete {
			if strings.Contains(r, k) {
				str.Reset()
				str.WriteString(log_delete[k])
			}
		}
	} else {
		for _, k := range order_log_post {
			if strings.Contains(r, k) {
				str.Reset()
				str.WriteString(log_post[k])
			}
		}
	}
	return str.String()
}

func GetOperObj(r string, user string) string {
	var str strings.Builder
	for k, v := range log_obj {
		if strings.Contains(r, k) {
			if k == "login" {
				str.WriteString(user)
				break
			} else {
				str.Reset()
				str.WriteString(v)
				str.WriteString("模块")
				break
			}
		}
	}
	return str.String()
}

func GetOperDes(m string, r string, ip string) string {
	var str strings.Builder
	opertype := GetOperType(m, r)
	str.WriteString(opertype)
	str.WriteString(",操作IP: ")
	str.WriteString(ip)
	return str.String()
}

func orderedMap(order map[string]string) []string {
	keys := make([]string, 0, len(order))
	for k := range order {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// import (

// 	// "net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/sirupsen/logrus"
// )

// func InitLogrus() *logrus.Logger {

// 	var log = logrus.New() // 创建一个log示例

// 	log.Formatter = &logrus.JSONFormatter{} // 设置为json格式的日志
// 	// file, err := os.OpenFile("./gin_log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644) // 创建一个log日志文件
// 	// if err != nil {
// 	// 	fmt.Println("创建文件/打开文件失败！")
// 	// 	return err
// 	// }
// 	// log.Out = file               // 设置log的默认文件输出
// 	// gin.SetMode(gin.ReleaseMode) // 发布版本
// 	gin.DefaultWriter = log.Out  // gin框架自己记录的日志也会输出
// 	log.Level = logrus.InfoLevel // 设置日志级别
// 	return log
// }

// func main() {
// 	err := initLogrus()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	r := gin.Default()
// 	r.GET("/logrus", func(c *gin.Context) {
// 		//log日志信息的写入
// 		log.WithFields(logrus.Fields{
// 			"url":    c.Request.RequestURI, //自定义显示的字段
// 			"method": c.Request.Method,
// 			"params": c.Query("name"),
// 			"IP":     c.ClientIP(),
// 		}).Info()
// 		resData := struct {
// 			Code int         `json:"code"`
// 			Msg  string      `json:"msg"`
// 			Data interface{} `json:"data"`
// 		}{http.StatusOK, "响应成功", "OK"}
// 		c.JSON(http.StatusOK, resData)
// 	})
// 	r.Run(":9090")
// }

//todo:文档地址：https://github.com/sirupsen/logrus
