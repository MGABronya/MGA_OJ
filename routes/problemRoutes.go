// @Title  problemRoutes
// @Description  程序的题目管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    ProblemRoutes
// @description   给gin引擎挂上题目相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ProblemRoutes(r *gin.Engine) *gin.Engine {

	// TODO 问题管理的路由分组
	problemRoutes := r.Group("/problem")

	// TODO 创建题目controller
	problemController := controller.NewProblemController()

	// TODO 创建题目
	problemRoutes.POST("/create", middleware.AuthMiddleware(), problemController.Create)

	// TODO 查看题目
	problemRoutes.GET("/show/:id", problemController.Show)

	// TODO 更新题目
	problemRoutes.PUT("/update/:id", middleware.AuthMiddleware(), problemController.Update)

	// TODO 删除题目
	problemRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), problemController.Delete)

	// TODO 查看问题列表
	problemRoutes.GET("/list", problemController.PageList)

	return r
}
