package router

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/toolkits/pkg/logger"
)

func TestInject(t *testing.T) {
	FileInfo, _ := ioutil.ReadDir("etc/client/")
	// logger.Debug(dir)
	logger.Debug(len(FileInfo))
	// var array []string
	for index := range FileInfo {
		fmt.Println("------------------------------")
		name := FileInfo[index].Name()
		fmt.Println(name)
		size := FileInfo[index].Size()
		fmt.Println(size)
		mTime := FileInfo[index].ModTime()
		fmt.Println(mTime)
		// partName := strings.Split(name, "-")
		// archName := strings.Split(partName[len(partName)-1], ".")[0]
		// if os == partName[2] && arch == archName {
		// 	array = append(array, partName[1])
		// }
	}
}
