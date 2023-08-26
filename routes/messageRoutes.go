// @Title  messageRoutes
// @Description  程序的留言板管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    MessageRoutes
// @description   给gin引擎挂上留言板相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func MessageRoutes(r *gin.Engine) *gin.Engine {

	// TODO 留言板管理的路由分组
	messageRoutes := r.Group("/message")

	// TODO 创建留言板controller
	messageController := controller.NewMessageController()

	// TODO 创建留言
	messageRoutes.POST("/create/:id", middleware.AuthMiddleware(), messageController.Create)

	// TODO 删除留言
	messageRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), messageController.Delete)

	// TODO 查看留言列表
	messageRoutes.GET("/list/:id", messageController.PageList)

	return r
}
