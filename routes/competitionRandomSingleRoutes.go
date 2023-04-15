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

// @title    CompetitionRandomSingleRoutes
// @description   给gin引擎挂上竞赛相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CompetitionRandomSingleRoutes(r *gin.Engine) *gin.Engine {

	// TODO 竞赛管理的路由分组
	competitionRandomSingleRoutes := r.Group("/competition/random/single")

	// TODO 创建竞赛controller
	competitionRandomSingleController := controller.NewCompetitionRandomSingleController()

	// TODO 报名比赛
	competitionRandomSingleRoutes.POST("/enter", middleware.AuthMiddleware(), competitionRandomSingleController.Enter)

	// TODO 查看报名状态
	competitionRandomSingleRoutes.GET("/enter/condition", middleware.AuthMiddleware(), competitionRandomSingleController.EnterCondition)

	// TODO 取消报名
	competitionRandomSingleRoutes.DELETE("/cancel/enter", middleware.AuthMiddleware(), competitionRandomSingleController.CancelEnter)

	// TODO 查看报名列表
	competitionRandomSingleRoutes.GET("/enter/list", competitionRandomSingleController.EnterPage)

	// TODO 实时报告情况
	competitionRandomSingleRoutes.GET("/enter/publish", competitionRandomSingleController.EnterPublish)

	return r
}
