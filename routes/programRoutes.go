// @Title  programRoutes
// @Description  程序的程序管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    programRoutes
// @description   给gin引擎挂上程序相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ProgramRoutes(r *gin.Engine) *gin.Engine {

	// TODO 通告管理的路由分组
	programRoutes := r.Group("/program")

	// TODO 创建通告controller
	programController := controller.NewProgramController()

	programRoutes.POST("/create", middleware.AuthMiddleware(), programController.Create)

	programRoutes.GET("/show/:id", middleware.AuthMiddleware(), programController.Show)

	programRoutes.PUT("/update/:id", middleware.AuthMiddleware(), programController.Update)

	programRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), programController.Delete)

	programRoutes.GET("/list", middleware.AuthMiddleware(), programController.PageList)

	return r
}
