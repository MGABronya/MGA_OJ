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

	return r
}
