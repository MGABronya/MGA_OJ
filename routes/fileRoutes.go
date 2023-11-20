// @Title  fileRoutes
// @Description  程序的文件管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

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
	fileRoutes.POST("/upload/:path", middleware.AuthMiddleware(), fileController.Upload)

	// TODO 下载文件
	fileRoutes.GET("/download/:id", middleware.AuthMiddleware(), fileController.Download)

	// TODO 解压某文件
	fileRoutes.PUT("/unzip", middleware.AuthMiddleware(), fileController.Unzip)

	// TODO 查看目录
	fileRoutes.GET("/path/:id", middleware.AuthMiddleware(), fileController.ShowPath)

	// TODO 创建目录
	fileRoutes.PUT("/mkdir/:id", middleware.AuthMiddleware(), fileController.MkDir)

	// TODO 复制
	fileRoutes.PUT("/cp", middleware.AuthMiddleware(), fileController.CP)

	// TODO 删除
	fileRoutes.DELETE("/rm/:id", middleware.AuthMiddleware(), fileController.RM)

	// TODO 重命名
	fileRoutes.PUT("/rename", middleware.AuthMiddleware(), fileController.Rename)

	// TODO 复制目录
	fileRoutes.PUT("/all/cp", middleware.AuthMiddleware(), fileController.CPAll)

	// TODO 删除目录
	fileRoutes.DELETE("/all/rm/:id", middleware.AuthMiddleware(), fileController.RMAll)

	return r
}
