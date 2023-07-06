// @Title  replyRoutes
// @Description  程序的讨论回复管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    ReplyRoutes
// @description   给gin引擎挂上讨论回复相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ReplyRoutes(r *gin.Engine) *gin.Engine {

	// TODO 回复管理的路由分组
	replyRoutes := r.Group("/reply")

	// TODO 创建回复controller
	replyController := controller.NewReplyController()

	// TODO 创建回复
	replyRoutes.POST("/create/:id", middleware.AuthMiddleware(), replyController.Create)

	// TODO 查看回复
	replyRoutes.GET("/show/:id", replyController.Show)

	// TODO 更新回复
	replyRoutes.PUT("/update/:id", middleware.AuthMiddleware(), replyController.Update)

	// TODO 删除回复
	replyRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), replyController.Delete)

	// TODO 查看回复列表
	replyRoutes.GET("/list/:id", replyController.PageList)

	// TODO 查看指定用户的回复列表
	replyRoutes.GET("/user/list/:id", replyController.UserList)

	// TODO 点赞、点踩回复
	replyRoutes.POST("/like/:id", middleware.AuthMiddleware(), replyController.Like)

	// TODO 取消点赞、点踩状态
	replyRoutes.DELETE("/cancel/like/:id", middleware.AuthMiddleware(), replyController.CancelLike)

	// TODO 查看点赞、点踩数量
	replyRoutes.GET("/like/number/:id", replyController.LikeNumber)

	// TODO 查看点赞、点踩列表
	replyRoutes.GET("/like/list/:id", replyController.LikeList)

	// TODO 查看用户当前点赞状态
	replyRoutes.GET("/like/show/:id", middleware.AuthMiddleware(), replyController.LikeShow)

	// TODO 查看用户点赞、点踩列表
	replyRoutes.GET("/likes/:id", replyController.Likes)

	// TODO 获取讨论回复热度排行
	replyRoutes.GET("/hot/rank/:id", replyController.HotRanking)

	return r
}
