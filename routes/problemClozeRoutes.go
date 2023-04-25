// @Title  problemClozeRoutes
// @Description  程序的填空题管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    problemClozeRoutes
// @description   给gin引擎挂上填空题相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ProblemClozeRoutes(r *gin.Engine) *gin.Engine {

	// TODO 题目管理的路由分组
	problemClozeRoutes := r.Group("/problem/Cloze")

	// TODO 创建题目controller
	problemClozeController := controller.NewProblemClozeController()

	// TODO 创建题目
	problemClozeRoutes.POST("/create/:id", middleware.AuthMiddleware(), problemClozeController.Create)

	// TODO 查看题目
	problemClozeRoutes.GET("/show/:id", middleware.AuthMiddleware(), problemClozeController.Show)

	// TODO 修改题目
	problemClozeRoutes.GET("/update/:id", middleware.AuthMiddleware(), problemClozeController.Update)

	// TODO 删除题目
	problemClozeRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), problemClozeController.Delete)

	// TODO 查看题目列表
	problemClozeRoutes.GET("/list/:id", middleware.AuthMiddleware(), problemClozeController.PageList)

	// TODO 提交测试
	problemClozeRoutes.POST("/submit/:id", middleware.AuthMiddleware(), problemClozeController.Submit)

	// TODO 查看提交
	problemClozeRoutes.GET("/submit/show/:id", middleware.AuthMiddleware(), problemClozeController.ShowSubmit)

	// TODO 查看提交列表
	problemClozeRoutes.GET("/submit/list/:user_id/:problem_id", middleware.AuthMiddleware(), problemClozeController.SubmitList)

	return r
}
