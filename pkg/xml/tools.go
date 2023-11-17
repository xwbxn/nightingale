package xmltool

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

// func exportTemplet()

func ExportXml[T any](c *gin.Context, xmls T, fileName string) {
	// 创建用于写入XML文件的临时文件
	file, err := os.CreateTemp("", "data.xml")
	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "导出失败")
		return
	}
	defer os.Remove(file.Name())
	// 创建XML编码器，将数据写入文件
	encoder := xml.NewEncoder(file)

	encoder.Indent("", "  ")

	procInst := xml.ProcInst{
		Target: "xml",
		Inst:   []byte("version=\"1.0\" encoding=\"UTF-8\"\n"),
	}
	if err := encoder.EncodeToken(procInst); err != nil {
		ginx.Bomb(http.StatusBadRequest, "导出失败")
	}
	err = encoder.Encode(xmls)

	if err != nil {
		ginx.Bomb(http.StatusBadRequest, "导出失败")
	}
	encoder.Flush()
	file.Close()
	fmt.Println("XML导出成功！")
	//设置文件类型
	c.Header("Content-Type", "application/xml;charset=utf8")
	//设置文件名称
	rand.Seed(time.Now().UnixNano())
	name := fileName + strconv.FormatInt(rand.Int63n(time.Now().Unix()), 10)
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(name))
	c.File(file.Name())
	ginx.NewRender(c)
}
