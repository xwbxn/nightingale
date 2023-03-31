package router

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/ginx"
)

type InstallVar struct {
	Token      string
	ServerHost string
	ServerPort int64
}

func (rt *Router) CategrefGetStart(c *gin.Context) {
	busigroup := ginx.QueryStr(c, "busigroup")

	// 生成随机token，5分钟有效，在下载探针时验证
	rand.Seed(time.Now().UnixNano())
	rand_arr := make([]byte, 32)
	rand.Read(rand_arr)
	token := hex.EncodeToString(rand_arr)

	err := rt.Redis.Set(c.Request.Context(), token, busigroup, time.Minute*5).Err()
	if err != nil {
		ginx.Bomb(http.StatusInternalServerError, err.Error())
	}

	serverHost := rt.getServerHost(c)
	serverPort := rt.getServerPort()
	cmd := fmt.Sprintf("curl http://%s:%d/categraf/install?token=%s", serverHost, serverPort, token)
	ginx.NewRender(c).Data(map[string]string{"url": cmd}, nil)
}

func (rt *Router) CategrafInstall(c *gin.Context) {
	token := ginx.QueryStr(c, "token")

	serverHost := rt.getServerHost(c)
	serverPort := rt.getServerPort()
	vars := InstallVar{
		Token:      token,
		ServerHost: serverHost,
		ServerPort: serverPort,
	}

	template_file := path.Join("etc", "categraf", "install.tpl")
	tmpl := template.Must(template.ParseFiles(template_file))

	var body bytes.Buffer
	tmpl.Execute(&body, vars)
	c.Writer.Write(body.Bytes())
}

type CategrafConfig struct {
	Busigroup      string
	ServerHost     string
	WriterUrl      string
	WriterUser     string
	WriterPass     string
	HttpEnable     string
	IbexEnable     string
	IbexServer     string
	ProviderServer string
}

func (rt *Router) CategrafDownload(c *gin.Context) {
	token := ginx.QueryStr(c, "token")

	busigroup, err := rt.Redis.Get(c.Request.Context(), token).Result()
	if err != nil {
		ginx.Bomb(http.StatusInternalServerError, err.Error())
	}

	tmpdir, err := os.MkdirTemp("", "")
	if err != nil {
		ginx.Bomb(http.StatusInternalServerError, err.Error())
	}

	os.Mkdir(tmpdir+"/categraf", 755)
	homedir, _ := os.Getwd()
	src := path.Join(homedir, "etc", "categraf", "dist")
	dst := path.Join(tmpdir, "categraf")
	cmd := exec.Command("cp", "-r", src, dst)
	err = CopyDir(src, dst)
	if err != nil {
		ginx.Bomb(http.StatusInternalServerError, err.Error())
	}

	serverHost := rt.getServerHost(c)
	writerUrl := rt.Center.Categraf.WriterUrl
	if writerUrl == "" {
		writerUrl = fmt.Sprintf("http://%s:17000/prometheus/v1/write", serverHost)
	}

	categrafConfig := CategrafConfig{
		Busigroup:      busigroup,
		WriterUrl:      writerUrl,
		WriterUser:     rt.Center.Categraf.WriterUser,
		WriterPass:     rt.Center.Categraf.WriterPass,
		IbexEnable:     "false",
		IbexServer:     serverHost + ":20090",
		HttpEnable:     "false",
		ProviderServer: serverHost + ":20000",
	}

	cfg_file_name := path.Join(tmpdir, "categraf", "conf", "config.toml")
	cfg_templ := template.Must(template.ParseFiles(cfg_file_name))
	cfg_file, err := os.OpenFile(cfg_file_name, os.O_RDWR, 0644)
	if err != nil {
		ginx.Bomb(http.StatusInternalServerError, err.Error())
	}
	cfg_templ.Execute(cfg_file, categrafConfig)
	cfg_file.Close()

	os.Chdir(tmpdir)
	cmd = exec.Command("tar", "-czf", "categraf.tar.gz", "categraf")
	err = cmd.Run()
	if err != nil {
		ginx.Bomb(http.StatusInternalServerError, err.Error())
	}

	downfile := path.Join(tmpdir, "categraf.tar.gz")
	c.FileAttachment(downfile, "categraf.tar.gz")
}

// 如果配置未指定host, 则从请求的host头中取，用于生成链接
func (rt *Router) getServerHost(c *gin.Context) string {
	serverHost := rt.Center.Categraf.ServerHost
	if serverHost == "" {
		serverHost = strings.Split(c.Request.Host, ":")[0]
	}
	return serverHost
}

// 如果配置未指定port, 则设置为默认值8080
func (rt *Router) getServerPort() int64 {
	serverPort := rt.Center.Categraf.ServerPort
	if serverPort == 0 {
		serverPort = 8080
	}
	return serverPort
}

/**
 * 拷贝文件夹,同时拷贝文件夹中的文件
 * @param srcPath 需要拷贝的文件夹路径: D:/test
 * @param destPath 拷贝到的位置: D:/backup/
 */
func CopyDir(srcPath string, destPath string) error {
	//检测目录正确性
	if srcInfo, err := os.Stat(srcPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("srcPath不是一个正确的目录！")
			fmt.Println(e.Error())
			return e
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !destInfo.IsDir() {
			e := errors.New("destInfo不是一个正确的目录！")
			fmt.Println(e.Error())
			return e
		}
	}
	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			copyFile(path, destNewPath)
		}
		return nil
	})
	if err != nil {
		fmt.Printf(err.Error())
	}
	return err
}

// 生成目录并拷贝文件
func copyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer srcFile.Close()
	//分割path目录
	destSplitPathDirs := strings.Split(dest, "/")
	//检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b, _ := pathExists(destSplitPath)
			if b == false {
				//创建目录
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

// 检测文件夹路径时候存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
