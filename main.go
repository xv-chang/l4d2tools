package main

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"l4d2Manage/utils"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type WgetModel struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ProxyModel struct {
	Url         string            `json:"url"`
	Method      string            `json:"method"`
	ContentType string            `json:"contentType"`
	Data        map[string]string `json:"data"`
}

type DonwloadFile struct {
	Name string
	Url  string
}

type WaitQueue struct {
	list []*DonwloadFile
}

func (q *WaitQueue) Enqueue(val *DonwloadFile) {
	q.list = append(q.list, val)
}

func (q *WaitQueue) Dequeue() *DonwloadFile {
	val := q.list[0]
	q.list = q.list[1:]
	return val
}

func (q *WaitQueue) Len() int {
	return len(q.list)
}

var queue *WaitQueue
var downloadStatus int
var conf *utils.WebConfig
var mines = map[string]string{
	".js":  "application/x-javascript",
	".css": "text/css",
	".mp4": "video/mp4",
	".jpg": "image/jpeg",
	".png": "image/png",
}

func main() {

	conf = utils.LoadConfig("config.json")
	list := make([]*DonwloadFile, 0)
	queue = &WaitQueue{list}
	downloadStatus = 0
	r := gin.Default()
	r.Use(MineType())
	r.Use(Cors())
	r.StaticFile("/", "./app/index.html")
	r.StaticFile("/favicon.ico", "./app/favicon.ico")
	r.StaticFS("/assets", http.Dir("./app/assets"))
	r.Use(Auth())
	r.POST("/proxy", func(c *gin.Context) {
		var model ProxyModel
		err := c.BindJSON(&model)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "msg": err.Error()})
			return
		}
		resp, _ := utils.Fetch(model.Url, model.Method, model.ContentType, model.Data)
		c.String(200, resp)
	})

	r.GET("/maps", func(c *gin.Context) {

	})
	r.POST("/maps/download", func(c *gin.Context) {
		var model WgetModel
		err := c.BindJSON(&model)
		if err != nil {
			c.JSON(500, gin.H{"ok": false, "msg": err.Error()})
			return
		}
		fileName := path.Join(conf.AddonsDir, model.Name)
		if utils.PathExists(fileName) {
			c.JSON(500, gin.H{"ok": false, "msg": "服务器上文件已存在" + fileName})
			return
		}
		queue.Enqueue(&DonwloadFile{Name: fileName, Url: model.Url})
		c.JSON(200, gin.H{"ok": true, "msg": "已添加到下载列表"})
	})

	r.GET("/downloads", func(c *gin.Context) {
		c.JSON(200, queue.list)
	})
	go runDownload()
	r.Run(conf.WebUrl)
}

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Auth")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		// 处理请求
		c.Next()
	}
}

func MineType() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		ext := path.Ext(url)
		contentType := mines[ext]
		if len(contentType) != 0 {
			c.Header("Content-Type", contentType)
		}
		c.Next()
	}
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Auth")
		if auth != "" {
			if auth != conf.Auth {
				c.AbortWithStatus(401)
				return
			}
		} else {
			c.Next()
		}
	}
}

func download(url string, filename string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func runDownload() {
	for {
		if downloadStatus == 0 && queue.Len() > 0 {
			item := queue.Dequeue()
			downloadStatus = 1
			log.Infoln("开始下载%v", item.Name)
			err := download(item.Url, item.Name)
			log.Infoln("下载完成%v", item.Name)
			utils.Check(err)
			if strings.HasSuffix(item.Name, ".zip") {
				log.Infoln("开始解压%v", item.Name)
				unzipVPK(item.Name, conf.AddonsDir)
				go removeFile(item.Name)
			}
			downloadStatus = 0
		}
		time.Sleep(3 * time.Second)
	}
}

func unzipVPK(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	var decodeName string
	for _, f := range zipReader.File {
		if f.Flags == 0 {
			//如果标致位是0  则是默认的本地编码   默认为gbk
			i := bytes.NewReader([]byte(f.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			decodeName = string(content)
		} else {
			//如果标志为是 1 << 11也就是 2048  则是utf-8编码
			decodeName = f.Name
		}
		if strings.HasSuffix(decodeName, ".vpk") {
			newName := strings.Replace(decodeName, "/", "_", -1)
			fpath := filepath.Join(destDir, newName)
			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func removeFile(file string) {
	log.Infoln("删除zip文件%v", file)
	time.Sleep(3 * time.Second)
	err := os.Remove(file)
	utils.Check(err)
}
