// @Title  problemNewRoutes
// @Description  程序的赛内题目管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    ProblemNewRoutes
// @description   给gin引擎挂上赛内题目相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ProblemNewRoutes(r *gin.Engine) *gin.Engine {

	// TODO 题目管理的路由分组
	problemNewRoutes := r.Group("/problem/new")

	// TODO 创建题目controller
	problemNewController := controller.NewProblemNewController()

	// TODO 创建题目
	problemNewRoutes.POST("/create", middleware.AuthMiddleware(), problemNewController.Create)

	// TODO 引用题目
	problemNewRoutes.POST("/quote/:competition_id/:problem_id/:score", middleware.AuthMiddleware(), problemNewController.Quote)

	// TODO 重现题目
	problemNewRoutes.POST("/rematch/:competition_id/:problem_id", middleware.AuthMiddleware(), problemNewController.Rematch)

	// TODO 查看题目
	problemNewRoutes.GET("/show/:id", middleware.AuthMiddleware(), problemNewController.Show)

	// TODO 查看题目测试样例数量
	problemNewRoutes.GET("/test/num/:id", problemNewController.TestNum)

	// TODO 更新题目
	problemNewRoutes.PUT("/update/:id", middleware.AuthMiddleware(), problemNewController.Update)

	// TODO 删除题目
	problemNewRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), problemNewController.Delete)

	// TODO 查看题目列表
	problemNewRoutes.GET("/list/:id", problemNewController.PageList)

	return r
}
