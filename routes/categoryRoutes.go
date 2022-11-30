// @Title  categoryRoutes
// @Description  程序的分类管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title   CategoryRoutes
// @description   给gin引擎挂上分类相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CategoryRoutes(r *gin.Engine) *gin.Engine {

	// TODO 分类管理的路由分组
	categoryRoutes := r.Group("/category")

	// TODO 创建分类controller
	categoryController := controller.NewCategoryController()

	// TODO 创建分类
	categoryRoutes.POST("/create", middleware.AuthMiddleware(), categoryController.Create)

	// TODO 查看分类
	categoryRoutes.GET("/show/:id", categoryController.Show)

	// TODO 更新分类
	categoryRoutes.PUT("/update/:id", middleware.AuthMiddleware(), categoryController.Update)

	// TODO 删除分类
	categoryRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), categoryController.Delete)

	// TODO 查看分类列表
	categoryRoutes.GET("/list", categoryController.PageList)

	return r
}
