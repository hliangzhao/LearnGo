package main

/*
RESTful (REpresentational State Transfer) API

不符合restful风格的API设计：
	/getAllEmployees
	/getAllExternalEmployees
	/createEmployee
	/updateEmployee

符合restful风格的API设计：
（1）名词应当使用复数。
	GET /employees                                      // 获取所有员工
	GET /employees/56                                   // 获取编号为56的员工
	GET /employees?state=external                       // 获取所有状态为外勤的员工
	GET /employees?state=external&maturity=senior       // 获取所有状态为外勤的高级员工
	PUT /employees/56 ...                               // 更新编号为56的员工
	POST /employees                                     // 创建一个新员工
	DELETE /employees                                   // 删除所有员工
	DELETE /employees/56                                // 删除编号为56的员工

（2）用动词区分资源请求和非资源请求。
GET /translate?from=de_ED&to=en_US&text=Hallo
*/

import (
	`github.com/gin-gonic/gin`
	`github.com/hliangzhao/LearnGo/11-gin/controllers`
	`github.com/hliangzhao/LearnGo/11-gin/middlewares`
	`log`
)

func main() {
	server := gin.Default()
	server.Use(middlewares.MyAuth())

	// serve一个静态资源文件夹
	server.Static("/resources", "./resources")
	// serve一个静态资源文件
	server.StaticFile("/csapp", "./resources/CSAPP-1-1.mp4")

	// 将controller和url绑定
	server.GET("/ping", func(context *gin.Context) {
		context.String(200, "%s", "pong")
	})
	videoController := controllers.NewVideoController()
	// 将共同url前缀抽取出来组成group
	videoGroup := server.Group("/videos")

	// 使用中间件
	videoGroup.Use(middlewares.MyLogger())

	videoGroup.GET("/", videoController.GetAll)
	videoGroup.PUT("/:id", videoController.Update)
	videoGroup.POST("/", videoController.Create)
	videoGroup.DELETE("/:id", videoController.Delete)
	// server.GET("/videos", videoController.GetAll)
	// server.PUT("/videos/:id", videoController.Update)
	// server.POST("/videos", videoController.Create)
	// server.DELETE("/videos/:id", videoController.Delete)

	// 启动服务器
	log.Fatalln(server.Run("localhost:8080"))
}
