// @Title  testRoutes
// @Description  程序的测试管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"

	"github.com/gin-gonic/gin"
)

// @title    TestRoutes
// @description   给gin引擎挂上表单相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func TestRoutes(r *gin.Engine) *gin.Engine {

	// TODO 测试管理的路由分组
	testRoutes := r.Group("/test")

	// TODO 创建测试controller
	testController := controller.NewTestController()

	// TODO 创建测试
	testRoutes.POST("/create", testController.Create)

	return r
}
