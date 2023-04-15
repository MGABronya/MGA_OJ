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

// @title    CompetitionRandomGroupRoutes
// @description   给gin引擎挂上竞赛相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CompetitionRandomGroupRoutes(r *gin.Engine) *gin.Engine {

	// TODO 竞赛管理的路由分组
	competitionRandomGroupRoutes := r.Group("/competition/random/group")

	// TODO 创建竞赛controller
	competitionRandomGroupController := controller.NewCompetitionRandomGroupController()

	// TODO 报名比赛
	competitionRandomGroupRoutes.POST("/enter", middleware.AuthMiddleware(), competitionRandomGroupController.Enter)

	// TODO 查看报名状态
	competitionRandomGroupRoutes.GET("/enter/condition", middleware.AuthMiddleware(), competitionRandomGroupController.EnterCondition)

	// TODO 取消报名
	competitionRandomGroupRoutes.DELETE("/cancel/enter", middleware.AuthMiddleware(), competitionRandomGroupController.CancelEnter)

	// TODO 查看报名列表
	competitionRandomGroupRoutes.GET("/enter/list", competitionRandomGroupController.EnterPage)

	// TODO 实时报告情况
	competitionRandomGroupRoutes.GET("/enter/publish", competitionRandomGroupController.EnterPublish)

	return r
}
