// @Title  levelRoutes
// @Description  程序的权限管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    LetterRoutes
// @description   给gin引擎挂上权限相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func LevelRoutes(r *gin.Engine) *gin.Engine {

	// TODO 权限管理的路由分组
	levelRoutes := r.Group("/level")

	// TODO 创建权限controller
	levelController := controller.NewLevelController()

	// TODO 修改某用户权限
	levelRoutes.POST("/update/:id/:level", middleware.AuthMiddleware(), levelController.Update)

	return r
}
