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

// @title    CompetitionSingleRoutes
// @description   给gin引擎挂上竞赛相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CompetitionSingleRoutes(r *gin.Engine) *gin.Engine {

	// TODO 竞赛管理的路由分组
	competitionSingleRoutes := r.Group("/competition/single")

	// TODO 创建竞赛controller
	competitionSingleController := controller.NewCompetitionSingleController()

	// TODO 报名比赛
	competitionSingleRoutes.POST("/enter/:id", middleware.AuthMiddleware(), competitionSingleController.Enter)

	// TODO 查看报名状态
	competitionSingleRoutes.GET("/enter/condition/:id", middleware.AuthMiddleware(), competitionSingleController.EnterCondition)

	// TODO 取消报名
	competitionSingleRoutes.DELETE("/cancel/enter/:id", middleware.AuthMiddleware(), competitionSingleController.CancelEnter)

	// TODO 查看报名列表
	competitionSingleRoutes.GET("/enter/list/:id", competitionSingleController.EnterPage)

	// TODO 比赛提交
	competitionSingleRoutes.POST("/submit/:id", middleware.AuthMiddleware(), competitionSingleController.Submit)

	// TODO 查看提交内容
	competitionSingleRoutes.GET("/show/record/:id", middleware.AuthMiddleware(), competitionSingleController.ShowRecord)

	// TODO 获取多篇提交
	competitionSingleRoutes.GET("/search/list/:id", middleware.AuthMiddleware(), competitionSingleController.SearchList)

	// TODO 订阅提交列表
	competitionSingleRoutes.GET("/publish/list/:id", competitionSingleController.PublishPageList)

	// TODO 查看一篇提交的测试通过情况
	competitionSingleRoutes.GET("/case/list/:id", competitionSingleController.CaseList)

	// TODO 查看一篇测试的通过情况
	competitionSingleRoutes.GET("/case/:id", competitionSingleController.Case)

	// TODO 黑客一篇提交
	competitionSingleRoutes.POST("/hack/:id", middleware.AuthMiddleware(), competitionSingleController.Hack)

	// TODO 计算比赛分数
	competitionSingleRoutes.POST("/score/:id", middleware.AuthMiddleware(), competitionSingleController.CompetitionScore)

	return r
}
