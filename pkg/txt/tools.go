package txt

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

func ExportTxt(c *gin.Context, txts []string) {
	// 将结构体数组导出为txt文件
	file, err := os.Create("data.txt")
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}
	defer file.Close()
	// str := "名称\t类型\tIP\t厂商\t资产位置\t状态\n"
	// _, err = io.WriteString(file, str)
	// if err != nil {
	// 	fmt.Println("写入文件失败:", err)
	// 	return
	// }
	for _, data := range txts {
		// status := ""
		// if asset.Status == 0 {
		// 	status = "下线"
		// } else if asset.Status == 1 {
		// 	status = "正常"
		// }
		// str := fmt.Sprintf("%s\t\t%s\t%s\t%s\t%s\t%s\n", asset.Name, asset.Type, asset.Ip, asset.Manufacturers, asset.Position, status)
		_, err := io.WriteString(file, data)
		if err != nil {
			fmt.Println("写入文件失败:", err)
			return
		}
	}
	//设置文件类型
	c.Header("Content-Type", "text/plain")
	//设置文件名称
	rand.Seed(time.Now().UnixNano())
	name := "资产信息" + strconv.FormatInt(rand.Int63n(time.Now().Unix()), 10)
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(name))
	c.File(file.Name())
	ginx.NewRender(c)
}
