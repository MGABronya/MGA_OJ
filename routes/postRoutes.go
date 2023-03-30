// @Title  postRoutes
// @Description  程序的题解管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    PostRoutes
// @description   给gin引擎挂上题解相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func PostRoutes(r *gin.Engine) *gin.Engine {

	// TODO 讨论管理的路由分组
	postRoutes := r.Group("/post")

	// TODO 创建讨论controller
	postController := controller.NewPostController()

	// TODO 创建讨论
	postRoutes.POST("/create/:id", middleware.AuthMiddleware(), postController.Create)

	// TODO 查看讨论
	postRoutes.GET("/show/:id", postController.Show)

	// TODO 更新讨论
	postRoutes.PUT("/update/:id", middleware.AuthMiddleware(), postController.Update)

	// TODO 删除讨论
	postRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), postController.Delete)

	// TODO 查看讨论列表
	postRoutes.GET("/list/:id", postController.PageList)

	// TODO 查看指定用户的讨论列表
	postRoutes.GET("/user/list/:id", postController.UserList)

	// TODO 点赞、点踩讨论
	postRoutes.POST("/like/:id", middleware.AuthMiddleware(), postController.Like)

	// TODO 取消点赞、点踩状态
	postRoutes.DELETE("/cancle/like/:id", middleware.AuthMiddleware(), postController.CancelLike)

	// TODO 查看点赞、点踩数量
	postRoutes.GET("/like/number/:id", postController.LikeNumber)

	// TODO 查看点赞、点踩列表
	postRoutes.GET("/like/list/:id", postController.LikeList)

	// TODO 查看用户当前点赞状态
	postRoutes.GET("/like/show/:id", middleware.AuthMiddleware(), postController.LikeShow)

	// TODO 查看用户点赞、点踩列表
	postRoutes.GET("/likes", middleware.AuthMiddleware(), postController.Likes)

	// TODO 收藏
	postRoutes.POST("/collect/:id", middleware.AuthMiddleware(), postController.Collect)

	// TODO 取消收藏
	postRoutes.DELETE("/cancel/collect/:id", middleware.AuthMiddleware(), postController.CancelCollect)

	// TODO 查看收藏状态
	postRoutes.GET("/collect/show/:id", middleware.AuthMiddleware(), postController.CollectShow)

	// TODO 查看收藏列表
	postRoutes.GET("/collect/list/:id", postController.CollectList)

	// TODO 查看收藏数量
	postRoutes.GET("/collect/number/:id", postController.CollectNumber)

	// TODO 查看用户收藏夹
	postRoutes.GET("/collects", middleware.AuthMiddleware(), postController.Collects)

	// TODO 游览题解
	postRoutes.POST("/visit/:id", middleware.AuthMiddleware(), postController.Visit)

	// TODO 游览数量
	postRoutes.GET("/visit/number/:id", postController.VisitNumber)

	// TODO 游览列表
	postRoutes.GET("/visit/list/:id", postController.VisitList)

	// TODO 指定用户的游览历史
	postRoutes.GET("/visits/:id", postController.Visits)

	// TODO 创建题解标签
	postRoutes.POST("/label/:id/:label", middleware.AuthMiddleware(), postController.LabelCreate)

	// TODO 删除题解标签
	postRoutes.DELETE("/label/:id/:label", middleware.AuthMiddleware(), postController.LabelDelete)

	// TODO 查看题解标签
	postRoutes.GET("/label/:id", postController.LabelShow)

	// TODO 按文本搜索题解
	postRoutes.GET("/search/:id/:text", postController.Search)

	// TODO 按标签搜索题解
	postRoutes.GET("/search/label/:id", postController.SearchLabel)

	// TODO 按文本和标签交集搜索题解
	postRoutes.GET("/search/with/label/:id/:text", postController.SearchWithLabel)

	// TODO 获取题解热度排行
	postRoutes.GET("/hot/rank/:id", postController.HotRanking)

	return r
}
