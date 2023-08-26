// @Title  problemMCQsRoutes
// @Description  程序的选择题管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    problemMCQsRoutes
// @description   给gin引擎挂上选择题相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ProblemMCQsRoutes(r *gin.Engine) *gin.Engine {

	// TODO 题目管理的路由分组
	problemMCQsRoutes := r.Group("/problem/MCQs")

	// TODO 创建题目controller
	problemMCQsController := controller.NewProblemMCQsController()

	// TODO 创建题目
	problemMCQsRoutes.POST("/create/:id", middleware.AuthMiddleware(), problemMCQsController.Create)

	// TODO 查看题目
	problemMCQsRoutes.GET("/show/:id", middleware.AuthMiddleware(), problemMCQsController.Show)

	// TODO 修改题目
	problemMCQsRoutes.GET("/update/:id", middleware.AuthMiddleware(), problemMCQsController.Update)

	// TODO 删除题目
	problemMCQsRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), problemMCQsController.Delete)

	// TODO 查看题目列表
	problemMCQsRoutes.GET("/list/:id", middleware.AuthMiddleware(), problemMCQsController.PageList)

	// TODO 提交测试
	problemMCQsRoutes.POST("/submit/:id", middleware.AuthMiddleware(), problemMCQsController.Submit)

	// TODO 查看提交
	problemMCQsRoutes.GET("/submit/show/:id", middleware.AuthMiddleware(), problemMCQsController.ShowSubmit)

	// TODO 查看提交列表
	problemMCQsRoutes.GET("/submit/list/:user_id/:problem_id", middleware.AuthMiddleware(), problemMCQsController.SubmitList)

	return r
}
