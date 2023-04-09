// @Title  recordRoutes
// @Description  程序的提交管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    RecordRoutes
// @description   给gin引擎挂上提交相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func RecordRoutes(r *gin.Engine) *gin.Engine {

	// TODO 提交管理的路由分组
	recordRoutes := r.Group("/record")

	// TODO 创建提交controller
	recordController := controller.NewRecordController()

	// TODO 创建提交
	recordRoutes.POST("/create", middleware.AuthMiddleware(), recordController.Create)

	// TODO 查看id指定提交状态
	recordRoutes.GET("/show/:id", middleware.AuthMiddleware(), recordController.Show)

	// TODO 查看某类特定提交列表
	recordRoutes.GET("/list", middleware.AuthMiddleware(), recordController.PageList)

	// TODO 订阅提交列表
	recordRoutes.GET("/publish/list", recordController.PublishPageList)

	// TODO 查看提交的测试通过情况
	recordRoutes.GET("/list/case/:id", recordController.CaseList)

	// TODO 查看某个测试的情况
	recordRoutes.GET("/case/:id", recordController.Case)

	return r
}
