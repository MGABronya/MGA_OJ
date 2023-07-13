// @Title  tagRoutes
// @Description  程序的自动标签管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"

	"github.com/gin-gonic/gin"
)

// @title    TagRoutes
// @description   给gin引擎挂上自动标签相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func TagRoutes(r *gin.Engine) *gin.Engine {

	// TODO 自动标签管理的路由分组
	tagRoutes := r.Group("/tag")

	// TODO 创建自动标签controller
	tagController := controller.NewTagController()

	// TODO 创建自动标签
	tagRoutes.POST("/create/:tag", tagController.Create)

	// TODO 删除自动标签
	tagRoutes.DELETE("/delete/:tag", tagController.Delete)

	// TODO 查看自动标签
	tagRoutes.GET("/show/:tag", tagController.Show)

	// TODO 查看自动标签列表
	tagRoutes.GET("/list", tagController.PageList)

	// TODO 生成自动标签
	tagRoutes.GET("/auto", tagController.Auto)

	return r
}
