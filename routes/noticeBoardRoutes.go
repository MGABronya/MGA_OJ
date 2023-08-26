// @Title  noticeRoutes
// @Description  程序的公告栏管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    noticeBoardRoutes
// @description   给gin引擎挂上公告栏相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func NoticeBoardRoutes(r *gin.Engine) *gin.Engine {

	// TODO 通告公告栏的路由分组
	noticeBoardRoutes := r.Group("/notice/board")

	// TODO 创建公告controller
	noticeBoardController := controller.NewNoticeBoardController()

	// TODO 创建公告
	noticeBoardRoutes.POST("/create", middleware.AuthMiddleware(), noticeBoardController.Create)

	// TODO 查看公告
	noticeBoardRoutes.GET("/show/:id", noticeBoardController.Show)

	// TODO 更新公告
	noticeBoardRoutes.PUT("/update/:id", middleware.AuthMiddleware(), noticeBoardController.Update)

	// TODO 删除公告
	noticeBoardRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), noticeBoardController.Delete)

	// TODO 公告栏列表
	noticeBoardRoutes.GET("/list", noticeBoardController.PageList)

	return r
}
