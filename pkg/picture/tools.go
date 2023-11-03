package picture

import (
	"errors"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

//校验图片格式及大小是否正确
//return
func VerifyPicture(fileHeader *multipart.FileHeader) (fileNameSuffix string, err error) {
	if fileHeader.Size > 1024*1024*5 {
		return "", errors.New("文件超5MB")
	}
	fileName := strings.Split(fileHeader.Filename, ".")
	if fileName[len(fileName)-1] != "bmp" && fileName[len(fileName)-1] != "jpeg" && fileName[len(fileName)-1] != "jpg" && fileName[len(fileName)-1] != "png" {
		return "", errors.New("文件格式错误")
	}
	return fileName[len(fileName)-1], nil
}

//生成图片名称
func GeneratePictureName(fileType, suffix string) (filePath string, err error) {
	path := "etc/picture/"

	_, err = PathExists(path)
	if err != nil {
		return "", err
	}

	name := fileType + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "." + suffix

	return path + name, nil
}

// PathExists 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
	return false, err
}
