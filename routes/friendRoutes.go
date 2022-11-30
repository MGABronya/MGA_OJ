// @Title  friendRoutes
// @Description  程序的好友管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    friendRoutes
// @description   给gin引擎挂上好友相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func FriendRoutes(r *gin.Engine) *gin.Engine {

	// TODO 好友管理的路由分组
	friendRoutes := r.Group("/friend")

	// TODO 创建好友controller
	friendController := controller.NewFriendController()

	// TODO 用户申请添加某个好友
	friendRoutes.POST("/apply/:id", middleware.AuthMiddleware(), friendController.Apply)

	// TODO 用户查看发出的好友申请
	friendRoutes.GET("/applying/list", middleware.AuthMiddleware(), friendController.ApplyingList)

	// TODO 好友查看接收到的好友申请
	friendRoutes.GET("/applied/list", middleware.AuthMiddleware(), friendController.AppliedList)

	// TODO 用户通过好友申请
	friendRoutes.POST("/consent/:id", middleware.AuthMiddleware(), friendController.Consent)

	// TODO 用户拒绝申请
	friendRoutes.PUT("/refuse/:id", middleware.AuthMiddleware(), friendController.Refuse)

	// TODO 用户拉黑某用户
	friendRoutes.POST("/block/:id", middleware.AuthMiddleware(), friendController.Block)

	// TODO 移除某用户的黑名单
	friendRoutes.DELETE("/remove/black/:id", middleware.AuthMiddleware(), friendController.RemoveBlack)

	// TODO 查看黑名单
	friendRoutes.GET("/black/list", middleware.AuthMiddleware(), friendController.BlackList)

	// TODO 用户删除某个好友
	friendRoutes.DELETE("/quit/:id", middleware.AuthMiddleware(), friendController.Quit)

	return r
}
