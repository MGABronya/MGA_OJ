// @Title  emailRoutes
// @Description  程序的邮件相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    emailRoutes
// @description   给gin引擎挂上邮件相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func EmailRoutes(r *gin.Engine) *gin.Engine {

	// TODO 邮件的路由分组
	emailRoutes := r.Group("/email")

	// TODO 创建邮件controller
	emailController := controller.NewEmailController()

	// TODO 发送邮件
	emailRoutes.POST("/send/:id", middleware.AuthMiddleware(), emailController.Send)

	// TODO 接收邮件
	emailRoutes.POST("/receive", middleware.AuthMiddleware(), emailController.Receive)

	return r
}
