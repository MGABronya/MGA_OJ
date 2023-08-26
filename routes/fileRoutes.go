// @Title  fileRoutes
// @Description  程序的文件管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"

	"github.com/gin-gonic/gin"
)

// @title    FileRoutes
// @description   给gin引擎挂上文件相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func FileRoutes(r *gin.Engine) *gin.Engine {

	// TODO 文件管理的路由分组
	fileRoutes := r.Group("/file")

	// TODO 创建文件controller
	fileController := controller.NewFileController()

	// TODO 上传文件
	fileRoutes.POST("/upload", fileController.Upload)

	// TODO 下载文件
	fileRoutes.GET("/download/:id", fileController.Download)

	return r
}
