// @Title  groupRoutes
// @Description  程序的用户组管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    GroupRoutes
// @description   给gin引擎挂上用户组相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func GroupRoutes(r *gin.Engine) *gin.Engine {

	// TODO 用户组管理的路由分组
	groupRoutes := r.Group("/group")

	// TODO 创建用户组controller
	groupController := controller.NewGroupController()

	// TODO 创建用户组
	groupRoutes.POST("/create", middleware.AuthMiddleware(), groupController.Create)

	// TODO 生成标准用户组
	groupRoutes.POST("/standard/create/:id/:num", middleware.AuthMiddleware(), groupController.CreateStandard)

	// TODO 标准用户组成员信息
	groupRoutes.GET("/standard/list/:id", middleware.AuthMiddleware(), groupController.ShowStandard)

	// TODO 查看用户组
	groupRoutes.GET("/show/:id", groupController.Show)

	// TODO 更新用户组
	groupRoutes.PUT("/update/:id", middleware.AuthMiddleware(), groupController.Update)

	// TODO 删除用户组
	groupRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), groupController.Delete)

	// TODO 查看用户组列表
	groupRoutes.GET("/list", groupController.PageList)

	// TODO 查看某一用户的用户组列表
	groupRoutes.GET("/leader/list/:id", groupController.LeaderList)

	// TODO 查看某一用户的用户组列表
	groupRoutes.GET("/member/list/:id", groupController.MemberList)

	// TODO 查看某一用户组的用户列表
	groupRoutes.GET("/user/list/:id", groupController.UserList)

	// TODO 用户申请加入某个用户组
	groupRoutes.POST("/apply/:id", middleware.AuthMiddleware(), groupController.Apply)

	// TODO 用户查看申请
	groupRoutes.GET("/applying/list", middleware.AuthMiddleware(), groupController.ApplyingList)

	// TODO 用户组组长查看申请
	groupRoutes.GET("/applied/list/:id", middleware.AuthMiddleware(), groupController.AppliedList)

	// TODO 用户组组长通过申请
	groupRoutes.PUT("/consent/:id", middleware.AuthMiddleware(), groupController.Consent)

	// TODO 用户组组长拒绝申请
	groupRoutes.PUT("/refuse/:id", middleware.AuthMiddleware(), groupController.Refuse)

	// TODO 用户组组长拉黑某用户
	groupRoutes.POST("/block/:group/:user", middleware.AuthMiddleware(), groupController.Block)

	// TODO 移除某用户的黑名单
	groupRoutes.DELETE("/remove/black/:group/:user", middleware.AuthMiddleware(), groupController.RemoveBlack)

	// TODO 查看黑名单
	groupRoutes.GET("/black/list/:id", middleware.AuthMiddleware(), groupController.BlackList)

	// TODO 用户退出某个用户组
	groupRoutes.DELETE("/quit/:id", middleware.AuthMiddleware(), groupController.Quit)

	// TODO 点赞或点踩
	groupRoutes.POST("/like/:id", middleware.AuthMiddleware(), groupController.Like)

	// TODO 取消点赞、点踩状态
	groupRoutes.DELETE("/cancel/like/:id", middleware.AuthMiddleware(), groupController.CancelLike)

	// TODO 查看点赞、点踩数量
	groupRoutes.GET("/like/number/:id", groupController.LikeNumber)

	// TODO 查看点赞、点踩列表
	groupRoutes.GET("/like/list/:id", groupController.LikeList)

	// TODO 查看用户点赞、点踩列表
	groupRoutes.GET("/likes/:id", groupController.Likes)

	// TODO 查看用户当前点赞状态
	groupRoutes.GET("/like/show/:id", middleware.AuthMiddleware(), groupController.LikeShow)

	// TODO 收藏
	groupRoutes.POST("/collect/:id", middleware.AuthMiddleware(), groupController.Collect)

	// TODO 取消收藏
	groupRoutes.DELETE("/cancel/collect/:id", middleware.AuthMiddleware(), groupController.CancelCollect)

	// TODO 查看收藏状态
	groupRoutes.GET("/collect/show/:id", middleware.AuthMiddleware(), groupController.CollectShow)

	// TODO 查看收藏列表
	groupRoutes.GET("/collect/list/:id", groupController.CollectList)

	// TODO 查看收藏数量
	groupRoutes.GET("/collect/number/:id", groupController.CollectNumber)

	// TODO 查看用户收藏夹
	groupRoutes.GET("/collects/:id", groupController.Collects)

	// TODO 创建用户组标签
	groupRoutes.POST("/label/:id/:label", middleware.AuthMiddleware(), groupController.LabelCreate)

	// TODO 删除用户组标签
	groupRoutes.DELETE("/label/:id/:label", middleware.AuthMiddleware(), groupController.LabelDelete)

	// TODO 查看用户组标签
	groupRoutes.GET("/label/:id", groupController.LabelShow)

	// TODO 按文本搜索用户组
	groupRoutes.GET("/search/:text", groupController.Search)

	// TODO 按标签搜索用户组
	groupRoutes.GET("/search/label", groupController.SearchLabel)

	// TODO 按文本和标签交集搜索用户组
	groupRoutes.GET("/search/with/label/:text", groupController.SearchWithLabel)

	// TODO 获取文章热度排行
	groupRoutes.GET("/hot/rank", groupController.HotRanking)

	return r
}
