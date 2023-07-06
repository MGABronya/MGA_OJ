// @Title  realNameRoutes
// @Description  程序的实名管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    realNameRoutes
// @description   给gin引擎挂上实名相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func RealNameRoutes(r *gin.Engine) *gin.Engine {

	// TODO 实名管理的路由分组
	realNameRoutes := r.Group("/real/name")

	// TODO 创建实名controller
	realNameController := controller.NewRealNameController()

	// TODO 创建实名
	realNameRoutes.POST("/create", middleware.AuthMiddleware(), realNameController.Create)

	// TODO 查看实名
	realNameRoutes.GET("/show/:id", realNameController.Show)

	// TODO 修改实名
	realNameRoutes.PUT("/update", middleware.AuthMiddleware(), realNameController.Update)

	// TODO 解除实名
	realNameRoutes.DELETE("/delete", middleware.AuthMiddleware(), realNameController.Delete)

	// TODO 查看实名列表
	realNameRoutes.DELETE("/list", middleware.AuthMiddleware(), realNameController.PageList)

	// TODO 上传实名
	realNameRoutes.PUT("/upload", middleware.AuthMiddleware(), realNameController.Upload)

	return r
}
