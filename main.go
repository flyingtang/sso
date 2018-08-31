package main

import (
	"sso-v2/cmds"
	"sso-v2/configs"
	"sso-v2/controllers"
	"sso-v2/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 初始化日志格式
func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
}

func main() {
	// 解析命令行参数
	cmd := cmds.ParseCmd()

	// 解析全局配置参数
	config := configs.NewConfig(cmd.Env)

	// 初始化数据库
	db, _ := models.NewMysql()
	defer db.Close()

	// gin 部分
	cmd.SetMode() // 根据参数启动对应的环境
	router := gin.Default()
	config.SetSession(router)
	//  设置静态资源目录挂载
	router.Static("/static", "./views")
	router.LoadHTMLGlob("views/**/*")
	// 路由
	router.GET("/signup",controllers.GetSignup)
	router.POST("/signup", controllers.Signup)
	router.GET("/login", controllers.GetLogin)
	router.POST("/login", controllers.Login)
	// router.POST("/login", controllers.Login)

	router.GET("/authorize", controllers.Authorize)
	addr := config.GetHttpAddrPort()
	router.Run(addr)
}
