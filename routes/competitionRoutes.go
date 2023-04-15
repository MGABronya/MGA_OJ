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
	competitionRoutes.GET("/list", competitionController.PageList)

	// TODO 新建密码
	competitionRoutes.POST("/passwd/create/:id", middleware.AuthMiddleware(), competitionController.CreatePasswd)

	// TODO 删除密码
	competitionRoutes.DELETE("/passwd/delete/:id", middleware.AuthMiddleware(), competitionController.DeletePasswd)

	// TODO 创建比赛标签
	competitionRoutes.POST("/label/:id/:label", middleware.AuthMiddleware(), competitionController.LabelCreate)

	// TODO 删除比赛标签
	competitionRoutes.DELETE("/label/:id/:label", middleware.AuthMiddleware(), competitionController.LabelDelete)

	// TODO 查看比赛标签
	competitionRoutes.GET("/label/:id", competitionController.LabelShow)

	// TODO 按文本搜索比赛
	competitionRoutes.GET("/search/:text", competitionController.Search)

	// TODO 按标签搜索比赛
	competitionRoutes.GET("/search/label", competitionController.SearchLabel)

	// TODO 按文本和标签交集搜索比赛
	competitionRoutes.GET("/search/with/label/:text", competitionController.SearchWithLabel)

	// TODO 重判
	competitionRoutes.PUT("/rejudge/:id", middleware.AuthMiddleware(), competitionController.Rejudge)

	// TODO 竞赛结果清除
	competitionRoutes.DELETE("/data/delete/:id", middleware.AuthMiddleware(), competitionController.CompetitionDataDelete)

	// TODO 查看指定竞赛成员排名
	competitionRoutes.GET("/member/rank/:competition/:member", competitionController.RankMember)

	// TODO 查看指定竞赛成员罚时情况
	competitionRoutes.GET("/member/show/:competition/:member", competitionController.MemberShow)

	// TODO 查看竞赛排行
	competitionRoutes.GET("/rank/list/:id", competitionController.RankList)

	// TODO 进行滚榜
	competitionRoutes.GET("/rolling/list/:id", competitionController.RollingList)

	return r
}
