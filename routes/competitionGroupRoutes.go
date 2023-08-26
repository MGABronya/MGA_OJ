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

// @title    CompetitionGroupRoutes
// @description   给gin引擎挂上竞赛相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CompetitionGroupRoutes(r *gin.Engine) *gin.Engine {

	// TODO 竞赛管理的路由分组
	competitionGroupRoutes := r.Group("/competition/group")

	// TODO 创建竞赛controller
	competitionGroupController := controller.NewCompetitionGroupController()

	// TODO 报名比赛
	competitionGroupRoutes.POST("/enter/:competition_id/:group_id", middleware.AuthMiddleware(), competitionGroupController.Enter)

	// TODO 查看报名状态
	competitionGroupRoutes.GET("/enter/condition/:id", middleware.AuthMiddleware(), competitionGroupController.EnterCondition)

	// TODO 取消报名
	competitionGroupRoutes.DELETE("/cancel/enter/:competition_id/:group_id", middleware.AuthMiddleware(), competitionGroupController.CancelEnter)

	// TODO 查看报名列表
	competitionGroupRoutes.GET("/enter/list/:id", middleware.AuthMiddleware(), competitionGroupController.EnterPage)

	// TODO 比赛提交
	competitionGroupRoutes.POST("/submit/:id", middleware.AuthMiddleware(), competitionGroupController.Submit)

	// TODO 查看提交内容
	competitionGroupRoutes.GET("/show/record/:id", middleware.AuthMiddleware(), competitionGroupController.ShowRecord)

	// TODO 获取多篇提交
	competitionGroupRoutes.GET("/search/list/:id", middleware.AuthMiddleware(), competitionGroupController.SearchList)

	// TODO 订阅提交列表
	competitionGroupRoutes.GET("/publish/list/:id", competitionGroupController.PublishPageList)

	// TODO 查看一篇提交的测试通过情况
	competitionGroupRoutes.GET("/case/list/:id", competitionGroupController.CaseList)

	// TODO 查看一篇测试的通过情况
	competitionGroupRoutes.GET("/case/:id", competitionGroupController.Case)

	// TODO 黑客一篇提交
	competitionGroupRoutes.POST("/hack/:id", middleware.AuthMiddleware(), competitionGroupController.Hack)

	// TODO 计算比赛分数
	competitionGroupRoutes.POST("/score/:id", middleware.AuthMiddleware(), competitionGroupController.CompetitionScore)

	return r
}
