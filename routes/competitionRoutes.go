// @Title  competitionRoutes
// @Description  程序的竞赛管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    CompetitionRoutes
// @description   给gin引擎挂上竞赛相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CompetitionRoutes(r *gin.Engine) *gin.Engine {

	// TODO 竞赛管理的路由分组
	competitionRoutes := r.Group("/competition")

	// TODO 创建竞赛controller
	competitionController := controller.NewCompetitionController()

	// TODO 创建竞赛
	competitionRoutes.POST("/create/:id", middleware.AuthMiddleware(), competitionController.Create)

	// TODO 查看竞赛
	competitionRoutes.GET("/show/:id", competitionController.Show)

	// TODO 更新竞赛
	competitionRoutes.PUT("/update/:id", middleware.AuthMiddleware(), competitionController.Update)

	// TODO 删除竞赛
	competitionRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), competitionController.Delete)

	// TODO 查看竞赛列表
	competitionRoutes.GET("/list/:id", competitionController.PageList)

	// TODO 查看指定用户的竞赛列表
	competitionRoutes.GET("/user/list/:id", competitionController.UserList)

	// TODO 查看指定竞赛用户打题数量与罚时

	// TODO 查看竞赛排行

	return r
}
