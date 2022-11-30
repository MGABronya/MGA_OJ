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

	return r
}
