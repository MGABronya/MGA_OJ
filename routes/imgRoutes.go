// @Title  imgRoutes
// @Description  程序的图片管理相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"

	"github.com/gin-gonic/gin"
)

// @title    ImgRoutes
// @description   给gin引擎挂上图片相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ImgRoutes(r *gin.Engine) *gin.Engine {

	// TODO 图片管理的路由分组
	imgRoutes := r.Group("/img")

	// TODO 创建图片controller
	imgController := controller.NewImgController()

	// TODO 上传图片
	imgRoutes.POST("/upload", imgController.Upload)

	return r
}
