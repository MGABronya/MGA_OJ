// @Title  remarkRoutes
// @Description  程序的题解回复管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    RemarkRoutes
// @description   给gin引擎挂上题解回复相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func RemarkRoutes(r *gin.Engine) *gin.Engine {

	// TODO 回复管理的路由分组
	remarkRoutes := r.Group("/remark")

	// TODO 创建回复controller
	remarkController := controller.NewRemarkController()

	// TODO 创建回复
	remarkRoutes.POST("/create/:id", middleware.AuthMiddleware(), remarkController.Create)

	// TODO 查看回复
	remarkRoutes.GET("/show/:id", remarkController.Show)

	// TODO 更新回复
	remarkRoutes.PUT("/update/:id", middleware.AuthMiddleware(), remarkController.Update)

	// TODO 删除回复
	remarkRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), remarkController.Delete)

	// TODO 查看回复列表
	remarkRoutes.GET("/list/:id", remarkController.PageList)

	// TODO 查看指定用户的回复列表
	remarkRoutes.GET("/user/list/:id", remarkController.UserList)

	// TODO 点赞、点踩回复
	remarkRoutes.POST("/like/:id", middleware.AuthMiddleware(), remarkController.Like)

	// TODO 取消点赞、点踩状态
	remarkRoutes.DELETE("/cancle/like/:id", middleware.AuthMiddleware(), remarkController.CancelLike)

	// TODO 查看点赞、点踩数量
	remarkRoutes.GET("/like/number/:id", remarkController.LikeNumber)

	// TODO 查看点赞、点踩列表
	remarkRoutes.GET("/like/list/:id", remarkController.LikeList)

	// TODO 查看用户当前点赞状态
	remarkRoutes.GET("/like/show/:id", middleware.AuthMiddleware(), remarkController.LikeShow)

	// TODO 查看指定用户的点赞、点踩列表
	remarkRoutes.GET("/likes/:id", remarkController.Likes)

	// TODO 获取文章回复热度排行
	remarkRoutes.GET("/hot/rank/:id", remarkController.HotRanking)

	return r
}
