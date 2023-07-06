// @Title  articleRoutes
// @Description  程序的文章管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    ArticleRoutes
// @description   给gin引擎挂上文章相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ArticleRoutes(r *gin.Engine) *gin.Engine {

	// TODO 文章管理的路由分组
	articleRoutes := r.Group("/article")

	// TODO 创建文章controller
	articleController := controller.NewArticleController()

	// TODO 创建文章
	articleRoutes.POST("/create", middleware.AuthMiddleware(), articleController.Create)

	// TODO 查看文章
	articleRoutes.GET("/show/:id", articleController.Show)

	// TODO 更新文章
	articleRoutes.PUT("/update/:id", middleware.AuthMiddleware(), articleController.Update)

	// TODO 删除文章
	articleRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), articleController.Delete)

	// TODO 查看文章列表
	articleRoutes.GET("/list", articleController.PageList)

	// TODO 查看指定用户的文章列表
	articleRoutes.GET("/user/list/:id", articleController.UserList)

	// TODO 查看指定分类的文章列表
	articleRoutes.GET("/category/list/:id", articleController.CategoryList)

	// TODO 点赞或点踩文章
	articleRoutes.POST("/like/:id", middleware.AuthMiddleware(), articleController.Like)

	// TODO 取消点赞、点踩状态
	articleRoutes.DELETE("/cancel/like/:id", middleware.AuthMiddleware(), articleController.CancelLike)

	// TODO 查看点赞、点踩数量
	articleRoutes.GET("/like/number/:id", articleController.LikeNumber)

	// TODO 查看点赞、点踩列表
	articleRoutes.GET("/like/list/:id", articleController.LikeList)

	// TODO 查看用户点赞、点踩列表
	articleRoutes.GET("/likes/:id", articleController.Likes)

	// TODO 查看用户当前点赞状态
	articleRoutes.GET("/like/show/:id", middleware.AuthMiddleware(), articleController.LikeShow)

	// TODO 收藏
	articleRoutes.POST("/collect/:id", middleware.AuthMiddleware(), articleController.Collect)

	// TODO 取消收藏
	articleRoutes.DELETE("/cancel/collect/:id", middleware.AuthMiddleware(), articleController.CancelCollect)

	// TODO 查看收藏状态
	articleRoutes.GET("/collect/show/:id", middleware.AuthMiddleware(), articleController.CollectShow)

	// TODO 查看收藏列表
	articleRoutes.GET("/collect/list/:id", articleController.CollectList)

	// TODO 查看收藏数量
	articleRoutes.GET("/collect/number/:id", articleController.CollectNumber)

	// TODO 查看用户收藏夹
	articleRoutes.GET("/collects/:id", articleController.Collects)

	// TODO 游览文章
	articleRoutes.POST("/visit/:id", middleware.AuthMiddleware(), articleController.Visit)

	// TODO 游览数量
	articleRoutes.GET("/visit/number/:id", articleController.VisitNumber)

	// TODO 游览列表
	articleRoutes.GET("/visit/list/:id", articleController.VisitList)

	// TODO 游览历史
	articleRoutes.GET("/visits/:id", articleController.Visits)

	// TODO 创建文章标签
	articleRoutes.POST("/label/:id/:label", middleware.AuthMiddleware(), articleController.LabelCreate)

	// TODO 删除文章标签
	articleRoutes.DELETE("/label/:id/:label", middleware.AuthMiddleware(), articleController.LabelDelete)

	// TODO 查看文章标签
	articleRoutes.GET("/label/:id", articleController.LabelShow)

	// TODO 按文本搜索文章
	articleRoutes.GET("/search/:text", articleController.Search)

	// TODO 按标签搜索文章
	articleRoutes.GET("/search/label", articleController.SearchLabel)

	// TODO 按文本和标签交集搜索文章
	articleRoutes.GET("/search/with/label/:text", articleController.SearchWithLabel)

	// TODO 获取文章热度排行
	articleRoutes.GET("/hot/rank", articleController.HotRanking)

	return r
}
