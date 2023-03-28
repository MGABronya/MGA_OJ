// @Title  specialRoutes
// @Description  程序的特判管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    SpecialRoutes
// @description   给gin引擎挂上特判相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func SpecialRoutes(r *gin.Engine) *gin.Engine {

	// TODO 特判管理的路由分组
	specialJudgeRoutes := r.Group("/specialjudge")

	// TODO 创建特判controller
	specialJudgeController := controller.NewSpecialJudgeController()

	// TODO 创建特判程序
	specialJudgeRoutes.POST("/create", middleware.AuthMiddleware(), specialJudgeController.Create)

	// TODO 更新特判程序
	specialJudgeRoutes.PUT("/update/:id", middleware.AuthMiddleware(), specialJudgeController.Update)

	return r
}
