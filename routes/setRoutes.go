// @Title  setRoutes
// @Description  程序的表单管理相关路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"
	"MGA_OJ/middleware"

	"github.com/gin-gonic/gin"
)

// @title    SetRoutes
// @description   给gin引擎挂上表单相关的路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func SetRoutes(r *gin.Engine) *gin.Engine {

	// TODO 表单管理的路由分组
	setRoutes := r.Group("/set")

	// TODO 创建表单controller
	setController := controller.NewSetController()

	// TODO 创建表单
	setRoutes.POST("/create", middleware.AuthMiddleware(), setController.Create)

	// TODO 查看表单
	setRoutes.GET("/show/:id", setController.Show)

	// TODO 更新表单
	setRoutes.PUT("/update/:id", middleware.AuthMiddleware(), setController.Update)

	// TODO 删除表单
	setRoutes.DELETE("/delete/:id", middleware.AuthMiddleware(), setController.Delete)

	// TODO 查看表单列表
	setRoutes.GET("/list", setController.PageList)

	// TODO 查看表单内用户排行
	setRoutes.GET("/rank/list/:id", setController.RankList)

	// TODO 更新表单排行
	setRoutes.PUT("/rank/update/:id", middleware.AuthMiddleware(), setController.RankUpdate)

	// TODO 查看某一用户的表单列表
	setRoutes.GET("/user/list/:id", setController.UserList)

	// TODO 查看某一表单的主题列表
	setRoutes.GET("/topic/list/:id", setController.TopicList)

	// TODO 查看某一表单的用户组列表
	setRoutes.GET("/group/list/:id", setController.GroupList)

	// TODO 点赞、点踩表单
	setRoutes.POST("/like/:id", middleware.AuthMiddleware(), setController.Like)

	// TODO 取消点赞、点踩状态
	setRoutes.DELETE("/cancle/like/:id", middleware.AuthMiddleware(), setController.CancelLike)

	// TODO 查看点赞、点踩数量
	setRoutes.GET("/like/number/:id", setController.LikeNumber)

	// TODO 查看点赞、点踩列表
	setRoutes.GET("/like/list/:id", setController.LikeList)

	// TODO 查看用户当前点赞状态
	setRoutes.GET("/like/show/:id", middleware.AuthMiddleware(), setController.LikeShow)

	// TODO 查看用户点赞、点踩列表
	setRoutes.GET("/likes", middleware.AuthMiddleware(), setController.Likes)

	// TODO 收藏
	setRoutes.POST("/collect/:id", middleware.AuthMiddleware(), setController.Collect)

	// TODO 取消收藏
	setRoutes.DELETE("/cancel/collect/:id", middleware.AuthMiddleware(), setController.CancelCollect)

	// TODO 查看收藏状态
	setRoutes.GET("/collect/show/:id", middleware.AuthMiddleware(), setController.CollectShow)

	// TODO 查看收藏列表
	setRoutes.GET("/collect/list/:id", setController.CollectList)

	// TODO 查看收藏数量
	setRoutes.GET("/collect/number/:id", setController.CollectNumber)

	// TODO 查看用户收藏夹
	setRoutes.GET("/collects", middleware.AuthMiddleware(), setController.Collects)

	// TODO 游览表单
	setRoutes.POST("/visit/:id", middleware.AuthMiddleware(), setController.Visit)

	// TODO 游览数量
	setRoutes.GET("/visit/number/:id", setController.VisitNumber)

	// TODO 游览列表
	setRoutes.GET("/visit/list/:id", setController.VisitList)

	// TODO 游览历史
	setRoutes.GET("/visits", middleware.AuthMiddleware(), setController.Visits)

	// TODO 用户组申请加入某个表单
	setRoutes.POST("/apply/:id", middleware.AuthMiddleware(), setController.Apply)

	// TODO 用户组查看申请
	setRoutes.GET("/applying/list/:id", middleware.AuthMiddleware(), setController.ApplyingList)

	// TODO 表单查看申请
	setRoutes.GET("/applied/list/:id", middleware.AuthMiddleware(), setController.AppliedList)

	// TODO 表单通过申请
	setRoutes.PUT("/consent/:id", middleware.AuthMiddleware(), setController.Consent)

	// TODO 表单拒绝申请
	setRoutes.PUT("/refuse/:id", middleware.AuthMiddleware(), setController.Refuse)

	// TODO 表单拉黑某用户组
	setRoutes.POST("/block/:set/:group", middleware.AuthMiddleware(), setController.Block)

	// TODO 移除某用户组的黑名单
	setRoutes.DELETE("/remove/black/:set/:group", middleware.AuthMiddleware(), setController.RemoveBlack)

	// TODO 查看黑名单
	setRoutes.GET("/black/list/:id", middleware.AuthMiddleware(), setController.BlackList)

	// TODO 用户组退出某个表单
	setRoutes.DELETE("/quit/:set/:group", middleware.AuthMiddleware(), setController.Quit)

	// TODO 创建表单标签
	setRoutes.POST("/label/:id/:label", middleware.AuthMiddleware(), setController.LabelCreate)

	// TODO 删除表单标签
	setRoutes.DELETE("/label/:id/:label", middleware.AuthMiddleware(), setController.LabelDelete)

	// TODO 查看表单标签
	setRoutes.GET("/label/:id", setController.LabelShow)

	return r
}
