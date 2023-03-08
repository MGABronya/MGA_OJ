// @Title  problemRoutes
// @Description  程序的题目管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    ProblemRoutes
// @description   给gin引擎挂上题目相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ProblemRoutes(r *gin.Engine) *gin.Engine {

	// TODO 题目管理的路由分组
	problemRoutes := r.Group("/problem")

	// TODO 创建题目controller
	problemController := controller.NewProblemController()

	// TODO 创建题目
	problemRoutes.POST("/create", middleware.AuthMiddleware(), problemController.Create)

	// TODO 查看题目
	problemRoutes.GET("/show/:id", middleware.AuthMiddleware(), problemController.Show)

	// TODO 查看题目测试样例数量
	problemRoutes.GET("/test/num/:id", problemController.TestNum)

	// TODO 更新题目
	problemRoutes.PUT("/update/:id", middleware.AuthMiddleware(), problemController.Update)

	// TODO 删除题目
	problemRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), problemController.Delete)

	// TODO 查看题目列表
	problemRoutes.GET("/list", middleware.AuthMiddleware(), problemController.PageList)

	// TODO 查看指定用户上传的题目列表
	problemRoutes.GET("/user/list/:id", middleware.AuthMiddleware(), problemController.UserList)

	// TODO 点赞、点踩题目
	problemRoutes.POST("/like/:id", middleware.AuthMiddleware(), problemController.Like)

	// TODO 取消点赞、点踩状态
	problemRoutes.DELETE("/cancle/like/:id", middleware.AuthMiddleware(), problemController.CancelLike)

	// TODO 查看点赞、点踩数量
	problemRoutes.GET("/like/number/:id", problemController.LikeNumber)

	// TODO 查看点赞、点踩列表
	problemRoutes.GET("/like/list/:id", problemController.LikeList)

	// TODO 查看用户当前点赞状态
	problemRoutes.GET("/like/show/:id", middleware.AuthMiddleware(), problemController.LikeShow)

	// TODO 查看用户点赞、点踩列表
	problemRoutes.GET("/likes", middleware.AuthMiddleware(), problemController.Likes)

	// TODO 收藏
	problemRoutes.POST("/collect/:id", middleware.AuthMiddleware(), problemController.Collect)

	// TODO 取消收藏
	problemRoutes.DELETE("/cancel/collect/:id", middleware.AuthMiddleware(), problemController.CancelCollect)

	// TODO 查看收藏状态
	problemRoutes.GET("/collect/show/:id", middleware.AuthMiddleware(), problemController.CollectShow)

	// TODO 查看收藏列表
	problemRoutes.GET("/collect/list/:id", problemController.CollectList)

	// TODO 查看收藏数量
	problemRoutes.GET("/collect/number/:id", problemController.CollectNumber)

	// TODO 查看用户收藏夹
	problemRoutes.GET("/collects", middleware.AuthMiddleware(), problemController.Collects)

	// TODO 游览题目
	problemRoutes.POST("/visit/:id", middleware.AuthMiddleware(), problemController.Visit)

	// TODO 游览数量
	problemRoutes.GET("/visit/number/:id", problemController.VisitNumber)

	// TODO 游览列表
	problemRoutes.GET("/visit/list/:id", problemController.VisitList)

	// TODO 游览历史
	problemRoutes.GET("/visits", middleware.AuthMiddleware(), problemController.Visits)

	// TODO 创建题目标签
	problemRoutes.POST("/label/:id/:label", middleware.AuthMiddleware(), problemController.LabelCreate)

	// TODO 删除题目标签
	problemRoutes.DELETE("/label/:id/:label", middleware.AuthMiddleware(), problemController.LabelDelete)

	// TODO 查看题目标签
	problemRoutes.GET("/label/:id", problemController.LabelShow)

	return r
}
