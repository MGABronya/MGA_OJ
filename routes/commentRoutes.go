// @Title  commentRoutes
// @Description  程序的讨论管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    CommentRoutes
// @description   给gin引擎挂上讨论相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CommentRoutes(r *gin.Engine) *gin.Engine {

	// TODO 讨论管理的路由分组
	commentRoutes := r.Group("/comment")

	// TODO 创建讨论controller
	commentController := controller.NewCommentController()

	// TODO 创建讨论
	commentRoutes.POST("/create/:id", middleware.AuthMiddleware(), commentController.Create)

	// TODO 查看讨论
	commentRoutes.GET("/show/:id", commentController.Show)

	// TODO 更新讨论
	commentRoutes.PUT("/update/:id", middleware.AuthMiddleware(), commentController.Update)

	// TODO 删除讨论
	commentRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), commentController.Delete)

	// TODO 查看讨论列表
	commentRoutes.GET("/list/:id", commentController.PageList)

	return r
}
