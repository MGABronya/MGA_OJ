// @Title  postRoutes
// @Description  程序的题解管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    PostRoutes
// @description   给gin引擎挂上题解相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func PostRoutes(r *gin.Engine) *gin.Engine {

	// TODO 讨论管理的路由分组
	postRoutes := r.Group("/post")

	// TODO 创建讨论controller
	postController := controller.NewPostController()

	// TODO 创建讨论
	postRoutes.POST("/create/:id", middleware.AuthMiddleware(), postController.Create)

	// TODO 查看讨论
	postRoutes.GET("/show/:id", postController.Show)

	// TODO 更新讨论
	postRoutes.PUT("/update/:id", middleware.AuthMiddleware(), postController.Update)

	// TODO 删除讨论
	postRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), postController.Delete)

	// TODO 查看讨论列表
	postRoutes.GET("/list/:id", postController.PageList)

	return r
}
