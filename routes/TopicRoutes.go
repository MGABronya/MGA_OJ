// @Title  topicRoutes
// @Description  程序的主题管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    TopicRoutes
// @description   给gin引擎挂上主题相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func TopicRoutes(r *gin.Engine) *gin.Engine {

	// TODO 主题管理的路由分组
	topicRoutes := r.Group("/topic")

	// TODO 创建主题controller
	topicController := controller.NewTopicController()

	// TODO 创建主题
	topicRoutes.POST("/create", middleware.AuthMiddleware(), topicController.Create)

	// TODO 查看主题
	topicRoutes.GET("/show/:id", topicController.Show)

	// TODO 更新主题
	topicRoutes.PUT("/update/:id", middleware.AuthMiddleware(), topicController.Update)

	// TODO 删除主题
	topicRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), topicController.Delete)

	// TODO 查看主题列表
	topicRoutes.GET("/list", topicController.PageList)

	// TODO 查看某一用户的主题列表
	topicRoutes.GET("/user/list/:id", topicController.UserList)

	// TODO 查看某一主题的题目列表
	topicRoutes.GET("/problem/list/:id", topicController.ProblemList)

	// TODO 点赞、点踩主题
	topicRoutes.POST("/like/:id", middleware.AuthMiddleware(), topicController.Like)

	// TODO 取消点赞、点踩状态
	topicRoutes.DELETE("/cancle/like/:id", middleware.AuthMiddleware(), topicController.CancelLike)

	// TODO 查看点赞、点踩数量
	topicRoutes.GET("/like/number/:id", topicController.LikeNumber)

	// TODO 查看点赞、点踩列表
	topicRoutes.GET("/like/list/:id", topicController.LikeList)

	// TODO 查看用户当前点赞状态
	topicRoutes.GET("/like/show/:id", middleware.AuthMiddleware(), topicController.LikeShow)

	// TODO 查看用户点赞、点踩列表
	topicRoutes.GET("/likes", middleware.AuthMiddleware(), topicController.Likes)

	// TODO 收藏
	topicRoutes.POST("/collect/:id", middleware.AuthMiddleware(), topicController.Collect)

	// TODO 取消收藏
	topicRoutes.DELETE("/cancel/collect/:id", middleware.AuthMiddleware(), topicController.CancelCollect)

	// TODO 查看收藏状态
	topicRoutes.GET("/collect/show/:id", middleware.AuthMiddleware(), topicController.CollectShow)

	// TODO 查看收藏列表
	topicRoutes.GET("/collect/list/:id", topicController.CollectList)

	// TODO 查看收藏数量
	topicRoutes.GET("/collect/number/:id", topicController.CollectNumber)

	// TODO 查看用户收藏夹
	topicRoutes.GET("/collects", middleware.AuthMiddleware(), topicController.Collects)

	// TODO 游览主题
	topicRoutes.POST("/visit/:id", middleware.AuthMiddleware(), topicController.Visit)

	// TODO 游览数量
	topicRoutes.GET("/visit/number/:id", topicController.VisitNumber)

	// TODO 游览列表
	topicRoutes.GET("/visit/list/:id", topicController.VisitList)

	// TODO 游览历史
	topicRoutes.GET("/visits", middleware.AuthMiddleware(), topicController.Visits)

	// TODO 创建主题标签
	topicRoutes.POST("/label/:id/:label", middleware.AuthMiddleware(), topicController.LabelCreate)

	// TODO 删除主题标签
	topicRoutes.DELETE("/label/:id/:label", middleware.AuthMiddleware(), topicController.LabelDelete)

	// TODO 查看主题标签
	topicRoutes.GET("/label/:id", topicController.LabelShow)

	return r
}
