// @Title  userRoutes
// @Description  程序的用户管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"

	"github.com/gin-gonic/gin"
)

// @title    UserRoutes
// @description   给gin引擎挂上用户相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func UserRoutes(r *gin.Engine) *gin.Engine {

	// TODO 问题管理的路由分组
	problemRoutes := r.Group("/user")

	// TODO 创建用户controller
	userController := controller.NewUserController()

	// TODO 验证码获取
	problemRoutes.GET("/verify/:id", userController.VerifyEmail)

	// TODO 用户注册
	problemRoutes.POST("/regist", userController.Register)

	// TODO 用户登录
	problemRoutes.POST("/login", userController.Login)

	return r
}
