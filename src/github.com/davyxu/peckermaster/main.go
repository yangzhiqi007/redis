package main

import (
	"fmt"
	"github.com/davyxu/cellmesh/service"
	"github.com/davyxu/golog"
	"github.com/davyxu/peckermaster/login"
	"github.com/davyxu/peckermaster/model"
	"github.com/davyxu/peckermaster/procmgr"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

var (
	flagDevMode   = service.CommandLine.Bool("dev", false, "enable dev mode, not host page, only api")
	flagPageDir   = service.CommandLine.String("pagedir", "../page", "dir contains pages")
	flagWebListen = service.CommandLine.String("listen", ":9096", "override address in sd")
	flagDB        = service.CommandLine.String("db", "../db/DB.json", "db store data")
)

var log = golog.New("main")

func startWebServer() error {

	e := gin.Default()

	model.DBFileName = *flagDB

	if err := model.DB.Load(); err != nil {
		log.Errorf("load config failed, %s", err.Error())
		return err
	}

	if *flagDevMode {
		// npm+webpack, 只做api服务,允许跨域
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true

		// http header自定义字段需要在这里允许
		config.AddAllowHeaders("Access-Token")

		e.Use(cors.New(config))

	} else {

		// 正常做web服务器, 不运行跨域访问
		e.Use(gzip.Gzip(gzip.DefaultCompression))
		e.LoadHTMLFiles(filepath.Join(*flagPageDir, "index.html"))
		e.Static("/static", filepath.Join(*flagPageDir, "static"))
		e.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"title": "Main website",
			})
		})
	}

	model.InitAdminPassword()
	login.InitRouter(&e.RouterGroup)
	procmgr.InitRouter(&e.RouterGroup)

	go func() {
		err := e.Run(*flagWebListen)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	return nil
}

func main() {

	service.Init("peckermaster")

	if startWebServer() != nil {
		os.Exit(1)
	}

	service.WaitExitSignal()
}
