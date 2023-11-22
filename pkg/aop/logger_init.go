package aop

import (

	// "ginDemo/config"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/toolkits/pkg/logger"
)

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

		key := c.Keys

		accepted := c.Accepted

		url := c.Request.URL

		body := c.Request.PostForm

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// remoteIP, _ := c.Request.Response.Location()

		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
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
