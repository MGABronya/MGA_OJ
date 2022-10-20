// @Title  threadRoutes
// @Description  程序的题解回复管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    ThreadRoutes
// @description   给gin引擎挂上题解回复相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ThreadRoutes(r *gin.Engine) *gin.Engine {

	// TODO 回复管理的路由分组
	threadRoutes := r.Group("/thread")

	// TODO 创建回复controller
	threadController := controller.NewThreadController()

	// TODO 创建回复
	threadRoutes.POST("/create/:id", middleware.AuthMiddleware(), threadController.Create)

	// TODO 查看回复
	threadRoutes.GET("/show/:id", threadController.Show)

	// TODO 更新回复
	threadRoutes.PUT("/update/:id", middleware.AuthMiddleware(), threadController.Update)

	// TODO 删除回复
	threadRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), threadController.Delete)

	// TODO 查看回复列表
	threadRoutes.GET("/list/:id", threadController.PageList)

	return r
}
