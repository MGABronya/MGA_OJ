// @Title  chatRoutes
// @Description  程序的群聊管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    ChatRoutes
// @description   给gin引擎挂上群聊相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ChatRoutes(r *gin.Engine) *gin.Engine {

	// TODO 私信管理的路由分组
	chatRoutes := r.Group("/chat")

	// TODO 创建群聊controller
	chatController := controller.NewChatController()

	// TODO 创建群聊消息
	chatRoutes.POST("/send/:id", middleware.AuthMiddleware(), chatController.Send)

	// TODO 查看获取多篇连接
	chatRoutes.GET("/link/list", middleware.AuthMiddleware(), chatController.LinkList)

	// TODO 列出聊天列表
	chatRoutes.GET("/list/:id", middleware.AuthMiddleware(), chatController.ChatList)

	// TODO 移除指定连接
	chatRoutes.DELETE("/remove/link/:id", middleware.AuthMiddleware(), chatController.RemoveLink)

	// TODO 建立实时接收
	chatRoutes.GET("/receive/:id", middleware.AuthMiddleware(), chatController.Receive)

	// TODO 建立连接实时接收
	chatRoutes.GET("/receivelink", middleware.AuthMiddleware(), chatController.ReceiveLink)

	return r
}
