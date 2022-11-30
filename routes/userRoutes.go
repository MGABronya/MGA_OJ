// @Title  userRoutes
// @Description  程序的用户管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    UserRoutes
// @description   给gin引擎挂上用户相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func UserRoutes(r *gin.Engine) *gin.Engine {

	// TODO 用户管理的路由分组
	userRoutes := r.Group("/user")

	// TODO 创建用户controller
	userController := controller.NewUserController()

	// TODO 验证码获取
	userRoutes.GET("/verify/:id", userController.VerifyEmail)

	// TODO 用户注册
	userRoutes.POST("/regist", userController.Register)

	// TODO 用户登录
	userRoutes.POST("/login", userController.Login)

	// TODO 找回密码
	userRoutes.PUT("/security", userController.Security)

	// TODO 更新密码
	userRoutes.PUT("/update/password", middleware.AuthMiddleware(), userController.UpdatePass)

	// TODO 返回当前登录的用户
	userRoutes.GET("/info", middleware.AuthMiddleware(), userController.Info)

	// TODO 修改用户信息
	userRoutes.PUT("/update", middleware.AuthMiddleware(), userController.Update)

	// TODO 修改用户等级
	userRoutes.PUT("/update/level/:id/:level", middleware.AuthMiddleware(), userController.UpdateLevel)

	// TODO 修改用户头像
	userRoutes.PUT("/update/icon", middleware.AuthMiddleware(), userController.UpdateIcon)

	return r
}
