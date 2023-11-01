// @Title  heartRoutes
// @Description  程序的心跳相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"

	"github.com/gin-gonic/gin"
)

// @title    HeartRoutes
// @description   给gin引擎挂上心跳相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func HeartRoutes(r *gin.Engine) *gin.Engine {

	// TODO 心跳管理的路由分组
	heartRoutes := r.Group("/heart")

	// TODO 创建心跳controller
	heartController := controller.NewHeartController()

	// TODO 查看指定时间段的心跳情况
	heartRoutes.GET("/show/:id/:start/:end", heartController.Show)

	// TODO 订阅心跳长连接
	heartRoutes.GET("/publish/:id", heartController.Publish)

	// TODO 查看近10s内的心跳忙碌占比
	heartRoutes.GET("/percentage", heartController.Percentage)

	return r
}
