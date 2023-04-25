// @Title  examRoutes
// @Description  程序的测试管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    examRoutes
// @description   给gin引擎挂上测试相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ExamRoutes(r *gin.Engine) *gin.Engine {

	// TODO 测试管理的路由分组
	examRoutes := r.Group("/exam")

	// TODO 创建测试controller
	examController := controller.NewExamController()

	// TODO 创建测试
	examRoutes.POST("/create/:id", middleware.AuthMiddleware(), examController.Create)

	// TODO 查看测试
	examRoutes.GET("/show/:id", middleware.AuthMiddleware(), examController.Show)

	// TODO 修改测试
	examRoutes.GET("/update/:id", middleware.AuthMiddleware(), examController.Update)

	// TODO 删除测试
	examRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), examController.Delete)

	// TODO 查看测试列表
	examRoutes.GET("/list/:id", middleware.AuthMiddleware(), examController.PageList)

	// TODO 查看用户分数
	examRoutes.GET("/score/show/:user_id/:exam_id", middleware.AuthMiddleware(), examController.ScoreShow)

	// TODO 修改用户分数
	examRoutes.PUT("/score/update/:user_id/:exam_id", middleware.AuthMiddleware(), examController.Update)

	// TODO 查看分数列表
	examRoutes.GET("/score/list/:id", middleware.AuthMiddleware(), examController.ScoreList)

	return r
}
