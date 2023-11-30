package router

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
	"github.com/toolkits/pkg/logger"
)

type SystemLog struct {
	Name       string `json:"name"`
	Size       string `json:"size"`
	UpdateTime int64  `json:"update_time"`
}

// @Summary      获取系统日志
// @Description  获取系统日志
// @Tags         系统日志
// @Accept       json
// @Produce      json
// @Param        filter    query    string  false  "筛选框(“file_name”：文件名称)"
// @Param        query    query    string  false  "搜索框"
// @Param        start    query    int64  false  "开始时间"
// @Param        end    query    int64  false  "结束时间"
// @Param        page    query    int  false  "页码"
// @Param        limit    query    int  false  "条数"
// @Success      200 {object} SystemLog
// @Router       /api/n9e/xh/sys-log/filter [get]
// @Security     ApiKeyAuth
func (rt *Router) systemLogGets(c *gin.Context) {
	filter := ginx.QueryStr(c, "filter", "")
	query := ginx.QueryStr(c, "query", "")
	start := ginx.QueryInt64(c, "start", -1)
	end := ginx.QueryInt64(c, "end", -1)
	page := ginx.QueryInt(c, "page", 1)
	limit := ginx.QueryInt(c, "limit", 20)
	if filter == "" && query != "" {
		ginx.Bomb(http.StatusOK, "参数错误")
	}
	FileInfo, _ := ioutil.ReadDir("logs/")
	logger.Debug(len(FileInfo))
	systemLogs := make([]SystemLog, 0)

	for index := range FileInfo {
		sizeT := float64(FileInfo[index].Size())
		size := ""

		if sizeT <= math.Pow(10, 3) {
			sizeF, err := strconv.ParseFloat(fmt.Sprintf("%.2f", sizeT), 64)
			ginx.Dangerous(err)
			size = strconv.FormatFloat(sizeF, 'f', -1, 64) + "b"
		} else if sizeT <= math.Pow(10, 6) {
			sizeF, err := strconv.ParseFloat(fmt.Sprintf("%.2f", sizeT/math.Pow(10, 3)), 64)
			ginx.Dangerous(err)
			size = strconv.FormatFloat(sizeF, 'f', -1, 64) + "kb"
		} else if sizeT <= math.Pow(10, 9) {
			sizeF, err := strconv.ParseFloat(fmt.Sprintf("%.2f", sizeT/math.Pow(10, 6)), 64)
			ginx.Dangerous(err)
			size = strconv.FormatFloat(sizeF, 'f', -1, 64) + "mb"
		} else if sizeT <= math.Pow(10, 12) {
			sizeF, err := strconv.ParseFloat(fmt.Sprintf("%.2f", sizeT/math.Pow(10, 9)), 64)
			ginx.Dangerous(err)
			size = strconv.FormatFloat(sizeF, 'f', -1, 64) + "gb"
		}
		if filter == "file_name" && strings.Contains(FileInfo[index].Name(), query) {
			if start != -1 && end != -1 {
				if FileInfo[index].ModTime().Unix() >= start && FileInfo[index].ModTime().Unix() <= end {
					systemLogs = append(systemLogs, SystemLog{
						Name:       FileInfo[index].Name(),
						Size:       size,
						UpdateTime: FileInfo[index].ModTime().Unix(),
					})
				}
			} else if start != -1 && end == -1 {
				if FileInfo[index].ModTime().Unix() >= start {
					systemLogs = append(systemLogs, SystemLog{
						Name:       FileInfo[index].Name(),
						Size:       size,
						UpdateTime: FileInfo[index].ModTime().Unix(),
					})
				}
			} else if start == -1 && end != -1 {
				if FileInfo[index].ModTime().Unix() <= end {
					systemLogs = append(systemLogs, SystemLog{
						Name:       FileInfo[index].Name(),
						Size:       size,
						UpdateTime: FileInfo[index].ModTime().Unix()})
				}
			} else {
				systemLogs = append(systemLogs, SystemLog{
					Name:       FileInfo[index].Name(),
					Size:       size,
					UpdateTime: FileInfo[index].ModTime().Unix()})
			}
		} else if filter == "" {
			if start != -1 && end != -1 {
				if FileInfo[index].ModTime().Unix() >= start && FileInfo[index].ModTime().Unix() <= end {
					systemLogs = append(systemLogs, SystemLog{
						Name:       FileInfo[index].Name(),
						Size:       size,
						UpdateTime: FileInfo[index].ModTime().Unix(),
					})
				}
			} else if start != -1 && end == -1 {
				if FileInfo[index].ModTime().Unix() >= start {
					systemLogs = append(systemLogs, SystemLog{
						Name:       FileInfo[index].Name(),
						Size:       size,
						UpdateTime: FileInfo[index].ModTime().Unix(),
					})
				}
			} else if start == -1 && end != -1 {
				if FileInfo[index].ModTime().Unix() <= end {
					systemLogs = append(systemLogs, SystemLog{
						Name:       FileInfo[index].Name(),
						Size:       size,
						UpdateTime: FileInfo[index].ModTime().Unix()})
				}
			} else {
				systemLogs = append(systemLogs, SystemLog{
					Name:       FileInfo[index].Name(),
					Size:       size,
					UpdateTime: FileInfo[index].ModTime().Unix()})
			}
		}
	}
	sort.Slice(systemLogs, func(i, j int) bool {
		return systemLogs[i].UpdateTime < systemLogs[j].UpdateTime
	})
	if (page-1)*limit > len(systemLogs) {
		ginx.Bomb(http.StatusOK, "参数错误")
	}
	endData := page * limit
	if page*limit > len(systemLogs) {
		endData = len(systemLogs)
	}
	ginx.NewRender(c).Data(gin.H{
		"list":  systemLogs[(page-1)*limit : endData],
		"total": len(systemLogs),
	}, nil)
}

// @Summary      下载系统日志
// @Description  下载系统日志
// @Tags         系统日志
// @Accept       json
// @Produce      json
// @Param        filter    query    string  false  "筛选框(“file_name”：文件名称)"
// @Param        query    query    string  false  "搜索框"
// @Param        start    query    int64  false  "开始时间"
// @Param        end    query    int64  false  "结束时间"
// @Param        body  body   map[string]interface{} false "{“name”:[“ss.log”, “ee.log”]}"
// @Success      200
// @Router       /api/n9e/xh/sys-log/export-xls [post]
// @Security     ApiKeyAuth
func (rt *Router) exportSystemLogXH(c *gin.Context) {
	filter := ginx.QueryStr(c, "filter", "")
	query := ginx.QueryStr(c, "query", "")
	start := ginx.QueryInt64(c, "start", -1)
	end := ginx.QueryInt64(c, "end", -1)
	if filter == "" && query != "" {
		ginx.Bomb(http.StatusOK, "参数错误")
	}
	var f map[string]interface{}
	ginx.BindJSON(c, &f)

	namesT, nameOk := f["names"]
	names := make([]string, 0)
	if nameOk {
		for _, val := range namesT.([]interface{}) {
			names = append(names, val.(string))
		}
	} else {
		FileInfo, _ := ioutil.ReadDir("logs/")
		for index := range FileInfo {
			if filter == "file_name" && strings.Contains(FileInfo[index].Name(), query) {
				// names = append(names, FileInfo[index].Name())
				if start != -1 && end != -1 {
					if FileInfo[index].ModTime().Unix() >= start && FileInfo[index].ModTime().Unix() <= end {
						names = append(names, FileInfo[index].Name())
					}
				} else if start != -1 && end == -1 {
					if FileInfo[index].ModTime().Unix() >= start {
						names = append(names, FileInfo[index].Name())
					}
				} else if start == -1 && end != -1 {
					if FileInfo[index].ModTime().Unix() <= end {
						names = append(names, FileInfo[index].Name())
					}
				} else {
					names = append(names, FileInfo[index].Name())
				}
			} else if filter == "" {
				if start != -1 && end != -1 {
					if FileInfo[index].ModTime().Unix() >= start && FileInfo[index].ModTime().Unix() <= end {
						names = append(names, FileInfo[index].Name())
					}
				} else if start != -1 && end == -1 {
					if FileInfo[index].ModTime().Unix() >= start {
						names = append(names, FileInfo[index].Name())
					}
				} else if start == -1 && end != -1 {
					if FileInfo[index].ModTime().Unix() <= end {
						names = append(names, FileInfo[index].Name())
					}
				} else {
					names = append(names, FileInfo[index].Name())
				}
			}
		}

	}

	// 创建临时文件来存储压缩包

	zipFile, err := os.CreateTemp("", "download*.zip")
	ginx.Dangerous(err)
	defer os.Remove(zipFile.Name())

	// 创建 Zip Writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if len(names) != 0 {
		for _, name := range names {
			file, err := os.Open("logs/" + name)
			ginx.Dangerous(err)
			defer file.Close()
			// 创建 Zip 文件
			zipEntry, err := zipWriter.Create(filepath.Base("logs/" + name))
			ginx.Dangerous(err)
			// 将文件内容复制到 Zip 文件中
			_, err = io.Copy(zipEntry, file)
			ginx.Dangerous(err)

		}
	}
	// 关闭zip文件写入器
	err = zipWriter.Close()
	ginx.Dangerous(err)
	// 设置响应头，告诉浏览器返回的是一个压缩包文件
	c.Header("Content-Type", "application/octet-stream")
	// 设置文件下载的名称
	c.Header("Content-Disposition", "attachment; filename=log.zip")
	// 将临时文件内容写入响应体
	c.File(zipFile.Name())
}
