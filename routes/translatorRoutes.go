// @Title  translatorRoutes
// @Description  程序的翻译相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"

	"github.com/gin-gonic/gin"
)

// @title    translatorRoutes
// @description   给gin引擎挂上翻译相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func TranslatorRoutes(r *gin.Engine) *gin.Engine {

	// TODO 翻译管理的路由分组
	translatorRoutes := r.Group("/translator")

	// TODO 翻译controller
	translatorController := controller.NewTranslatorController()

	// TODO 翻译
	translatorRoutes.POST("/translate", translatorController.Translate)

	return r
}
