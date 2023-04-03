// @Title  rejudgeRoutes
// @Description  程序的重判管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    RejudgeRoutes
// @description   给gin引擎挂上重判相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func RejudgeRoutes(r *gin.Engine) *gin.Engine {

	// TODO 重判管理的路由分组
	rejudgeRoutes := r.Group("/rejudge")

	// TODO 创建重判controller
	rejudgeController := controller.NewRejudgeController()

	// TODO 进行重判
	rejudgeRoutes.PUT("/do", middleware.AuthMiddleware(), rejudgeController.Do)

	// TODO 对某场比赛结果进行清空
	rejudgeRoutes.DELETE("/competiton/delete/:id", middleware.AuthMiddleware(), rejudgeController.CompetitionDelete)

	// TODO 对某场比赛结果重新进行分数统计
	rejudgeRoutes.PUT("/competiton/score/:id", middleware.AuthMiddleware(), rejudgeController.CompetitionScore)

	return r
}
