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

// @title    CompetitionStandardSingleRoutes
// @description   给gin引擎挂上竞赛相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CompetitionStandardSingleRoutes(r *gin.Engine) *gin.Engine {

	// TODO 竞赛管理的路由分组
	competitionStandardSingleRoutes := r.Group("/competition/standard/single")

	// TODO 创建竞赛controller
	competitionStandardSingleController := controller.NewCompetitionStandardSingleController()

	// TODO 报名比赛
	competitionStandardSingleRoutes.POST("/enter/:id", middleware.AuthMiddleware(), competitionStandardSingleController.Enter)

	// TODO 查看报名状态
	competitionStandardSingleRoutes.GET("/enter/condition/:id", middleware.AuthMiddleware(), competitionStandardSingleController.EnterCondition)

	// TODO 取消报名
	competitionStandardSingleRoutes.DELETE("/cancel/enter/:id", middleware.AuthMiddleware(), competitionStandardSingleController.CancelEnter)

	// TODO 查看报名列表
	competitionStandardSingleRoutes.GET("/enter/list/:id", middleware.AuthMiddleware(), competitionStandardSingleController.EnterPage)

	return r
}
