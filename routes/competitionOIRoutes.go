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

// @title    CompetitionOIRoutes
// @description   给gin引擎挂上竞赛相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CompetitionOIRoutes(r *gin.Engine) *gin.Engine {

	// TODO 竞赛管理的路由分组
	competitionOIRoutes := r.Group("/competition/OI")

	// TODO 创建竞赛controller
	competitionOIController := controller.NewCompetitionOIController()

	// TODO 报名比赛
	competitionOIRoutes.POST("/enter/:id", middleware.AuthMiddleware(), competitionOIController.Enter)

	// TODO 查看报名状态
	competitionOIRoutes.GET("/enter/condition/:id", middleware.AuthMiddleware(), competitionOIController.EnterCondition)

	// TODO 取消报名
	competitionOIRoutes.DELETE("/cancel/enter/:id", middleware.AuthMiddleware(), competitionOIController.CancelEnter)

	// TODO 查看报名列表
	competitionOIRoutes.GET("/enter/list/:id", competitionOIController.EnterPage)

	// TODO 比赛提交
	competitionOIRoutes.POST("/submit/:id", middleware.AuthMiddleware(), competitionOIController.Submit)

	// TODO 查看提交内容
	competitionOIRoutes.GET("/show/record/:id", middleware.AuthMiddleware(), competitionOIController.ShowRecord)

	// TODO 获取多篇提交
	competitionOIRoutes.GET("/search/list/:id", middleware.AuthMiddleware(), competitionOIController.SearchList)

	// TODO 订阅提交列表
	competitionOIRoutes.GET("/publish/list/:id", competitionOIController.PublishPageList)

	// TODO 查看一篇提交的测试通过情况
	competitionOIRoutes.GET("/case/list/:id", competitionOIController.CaseList)

	// TODO 查看一篇测试的通过情况
	competitionOIRoutes.GET("/case/:id", competitionOIController.Case)

	return r
}
