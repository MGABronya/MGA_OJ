// @Title  competitionRoutes
// @Description  程序的竞赛管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    CompetitionStandardGroupRoutes
// @description   给gin引擎挂上竞赛相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CompetitionStandardGroupRoutes(r *gin.Engine) *gin.Engine {

	// TODO 竞赛管理的路由分组
	competitionStandardGroupRoutes := r.Group("/competition/standard/group")

	// TODO 创建竞赛controller
	competitionStandardGroupController := controller.NewCompetitionStandardGroupController()

	// TODO 报名比赛
	competitionStandardGroupRoutes.POST("/enter/:id", middleware.AuthMiddleware(), competitionStandardGroupController.Enter)

	// TODO 查看报名状态
	competitionStandardGroupRoutes.GET("/enter/condition/:id", middleware.AuthMiddleware(), competitionStandardGroupController.EnterCondition)

	// TODO 取消报名
	competitionStandardGroupRoutes.DELETE("/cancel/enter/:group_id/:competition_id", middleware.AuthMiddleware(), competitionStandardGroupController.CancelEnter)

	// TODO 查看报名列表
	competitionStandardGroupRoutes.GET("/enter/list/:id", middleware.AuthMiddleware(), competitionStandardGroupController.EnterPage)

	return r
}
