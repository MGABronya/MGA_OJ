// @Title  problemFileRoutes
// @Description  程序的选择题管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    problemFileRoutes
// @description   给gin引擎挂上文件题相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ProblemFileRoutes(r *gin.Engine) *gin.Engine {

	// TODO 题目管理的路由分组
	problemFileRoutes := r.Group("/problem/File")

	// TODO 创建题目controller
	problemFileController := controller.NewProblemFileController()

	// TODO 创建题目
	problemFileRoutes.POST("/create/:id", middleware.AuthMiddleware(), problemFileController.Create)

	// TODO 查看题目
	problemFileRoutes.GET("/show/:id", middleware.AuthMiddleware(), problemFileController.Show)

	// TODO 修改题目
	problemFileRoutes.GET("/update/:id", middleware.AuthMiddleware(), problemFileController.Update)

	// TODO 删除题目
	problemFileRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), problemFileController.Delete)

	// TODO 查看题目列表
	problemFileRoutes.GET("/list/:id", middleware.AuthMiddleware(), problemFileController.PageList)

	// TODO 提交测试
	problemFileRoutes.POST("/submit/:id", middleware.AuthMiddleware(), problemFileController.Submit)

	// TODO 查看提交
	problemFileRoutes.GET("/submit/show/:id", middleware.AuthMiddleware(), problemFileController.ShowSubmit)

	// TODO 查看提交列表
	problemFileRoutes.GET("/submit/list/:user_id/:problem_id", middleware.AuthMiddleware(), problemFileController.SubmitList)

	return r
}
