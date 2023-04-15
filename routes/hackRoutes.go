// @Title  hackRoutes
// @Description  程序的黑客管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"

	"github.com/gin-gonic/gin"
)

// @title    hackRoutes
// @description   给gin引擎挂上黑客相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func HackRoutes(r *gin.Engine) *gin.Engine {

	// TODO 黑客管理的路由分组
	hackRoutes := r.Group("/hack")

	// TODO 创建黑客controller
	hackController := controller.NewHackController()

	hackRoutes.GET("/show/:id", hackController.Show)

	hackRoutes.GET("/shownum/:member_id/:competition_id", hackController.ShowNum)

	return r
}
