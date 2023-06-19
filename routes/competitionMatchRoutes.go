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

// @title    CompetitionMatchRoutes
// @description   给gin引擎挂上竞赛相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CompetitionMatchRoutes(r *gin.Engine) *gin.Engine {

	// TODO 竞赛管理的路由分组
	competitionMatchRoutes := r.Group("/competition/match")

	// TODO 创建竞赛controller
	competitionMatchController := controller.NewCompetitionMatchController()

	// TODO 报名比赛
	competitionMatchRoutes.POST("/enter/:id", middleware.AuthMiddleware(), competitionMatchController.Enter)

	// TODO 查看报名状态
	competitionMatchRoutes.GET("/enter/condition/:id", middleware.AuthMiddleware(), competitionMatchController.EnterCondition)

	// TODO 取消报名
	competitionMatchRoutes.DELETE("/cancel/enter/:id", middleware.AuthMiddleware(), competitionMatchController.CancelEnter)

	// TODO 查看报名列表
	competitionMatchRoutes.GET("/enter/list/:id", middleware.AuthMiddleware(), competitionMatchController.EnterPage)

	return r
}
