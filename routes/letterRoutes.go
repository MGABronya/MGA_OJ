// @Title  letterRoutes
// @Description  程序的私信管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    LetterRoutes
// @description   给gin引擎挂上私信相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func LetterRoutes(r *gin.Engine) *gin.Engine {

	// TODO 私信管理的路由分组
	letterRoutes := r.Group("/letter")

	// TODO 创建私信controller
	letterController := controller.NewLetterController()

	// TODO 创建私信
	letterRoutes.POST("/send/:id", middleware.AuthMiddleware(), letterController.Send)

	// TODO 查看获取多篇连接
	letterRoutes.GET("/link/list", middleware.AuthMiddleware(), letterController.LinkList)

	// TODO 列出聊天列表
	letterRoutes.GET("/list/:id", middleware.AuthMiddleware(), letterController.ChatList)

	// TODO 移除指定连接
	letterRoutes.DELETE("/remove/link/:id", middleware.AuthMiddleware(), letterController.RemoveLink)

	// TODO 已读
	letterRoutes.PUT("/read/:id", middleware.AuthMiddleware(), letterController.Read)

	// TODO 订阅私信
	letterRoutes.GET("/receive/:id", middleware.WsAuthMiddleware(), letterController.Receive)

	// TODO 订阅连接
	letterRoutes.GET("/receivelink", middleware.WsAuthMiddleware(), letterController.ReceiveLink)

	// TODO 用户私信拉黑某用户
	letterRoutes.POST("/block/:id", middleware.AuthMiddleware(), letterController.Block)

	// TODO 移除某用户私信的黑名单
	letterRoutes.DELETE("/remove/black/:id", middleware.AuthMiddleware(), letterController.RemoveBlack)

	// TODO 查看私信黑名单
	letterRoutes.GET("/black/list", middleware.AuthMiddleware(), letterController.BlackList)

	return r
}
