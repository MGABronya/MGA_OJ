// @Title  badgeRoutes
// @Description  程序的徽章管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    BadgeRoutes
// @description   给gin引擎挂上徽章相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func BadgeRoutes(r *gin.Engine) *gin.Engine {

	// TODO 徽章管理的路由分组
	badgeRoutes := r.Group("/badge")

	// TODO 创建徽章controller
	badgeController := controller.NewBadgeController()

	// TODO 创建徽章
	badgeRoutes.POST("/create", middleware.AuthMiddleware(), badgeController.Create)

	// TODO 更新徽章
	badgeRoutes.PUT("/udpate/:id", middleware.AuthMiddleware(), badgeController.Update)

	// TODO 查看徽章
	badgeRoutes.GET("/show/:id", badgeController.Show)

	// TODO 删除徽章
	badgeRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), badgeController.Show)

	// TODO 查看徽章列表
	badgeRoutes.GET("/list", badgeController.PageList)

	// TODO 查看用户徽章
	badgeRoutes.GET("/user/show/:user/:badge", badgeController.UserShow)

	// TODO 查看用户徽章列表
	badgeRoutes.GET("/user/list/:id", badgeController.UserList)

	// TODO 查看变量列表
	badgeRoutes.GET("/behavior/list", badgeController.BehaviorList)

	// TODO 查看变量描述
	badgeRoutes.GET("/behavior/description/:id", badgeController.BehaviorShow)

	// TODO 用户连接
	badgeRoutes.GET("/publish", middleware.AuthMiddleware(), badgeController.Publish)

	// TODO 查看某用户的某行为统计
	badgeRoutes.GET("/evaluate/expression/:user/:expression", badgeController.EvaluateExpression)

	return r
}
