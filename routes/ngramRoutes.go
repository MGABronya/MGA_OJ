// @Title  ngramRoutes
// @Description  程序的文本相似度相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"

	"github.com/gin-gonic/gin"
)

// @title    ngramRoutes
// @description   给gin引擎挂上文本相似度相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func NgramRoutes(r *gin.Engine) *gin.Engine {

	// TODO 文本相似度的路由分组
	ngramRoutes := r.Group("/ngram")

	// TODO 创建文本相似度controller
	ngramController := controller.NewNgramController()

	// TODO 计算文本相似度
	ngramRoutes.POST("/similarity", ngramController.ComputeSimilarity)

	// TODO 计算矩阵图连通块
	ngramRoutes.POST("/judge/:judge", ngramController.JudgeSimilarity)

	return r
}
