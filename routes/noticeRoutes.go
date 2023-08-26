// @Title  noticeRoutes
// @Description  程序的通告管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    noticeRoutes
// @description   给gin引擎挂上通告相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func NoticeRoutes(r *gin.Engine) *gin.Engine {

	// TODO 通告管理的路由分组
	noticeRoutes := r.Group("/notice")

	// TODO 创建通告controller
	noticeController := controller.NewNoticeController()

	noticeRoutes.POST("/create/:id", middleware.AuthMiddleware(), noticeController.Create)

	noticeRoutes.GET("/publish/:id", noticeController.Publish)

	noticeRoutes.GET("/show/:id", noticeController.Show)

	noticeRoutes.GET("/list/:id", noticeController.PageList)

	return r
}
