// @Title  routes
// @Description  程序的路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package main

import (
	"MGA_OJ/middleware"
	"MGA_OJ/routes"

	"github.com/gin-gonic/gin"
)

// @title    CollectRoute
// @description   给gin引擎挂上路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CollectRoute(r *gin.Engine) *gin.Engine {

	// TODO 添加中间件
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())

	// TODO 挂上题目路由
	r = routes.ProblemRoutes(r)

	// TODO 挂上用户路由
	r = routes.UserRoutes(r)

	// TODO 挂上讨论路由
	r = routes.CommentRoutes(r)

	// TODO 挂上讨论的回复路由
	r = routes.ReplyRoutes(r)

	// TODO 挂上题解路由
	r = routes.PostRoutes(r)

	// TODO 挂上题解的回复路由
	r = routes.ThreadRoutes(r)

	// TODO 挂上文章路由
	r = routes.ArticleRoutes(r)

	// TODO 挂上文章分类路由
	r = routes.CategoryRoutes(r)

	// TODO 挂上文章回复路由
	r = routes.RemarkRoutes(r)

	// TODO 挂上题目主题路由
	r = routes.TopicRoutes(r)

	// TODO 挂上用户组路由
	r = routes.GroupRoutes(r)

	// TODO 挂上用户好友路由
	r = routes.FriendRoutes(r)

	// TODO 挂上表单路由
	r = routes.SetRoutes(r)

	// TODO 挂上提交路由
	r = routes.RecordRoutes(r)

	// TODO 挂上比赛路由
	r = routes.CompetitionRoutes(r)

	// TODO 留言板路由
	r = routes.MessageRoutes(r)

	// TODO 私信路由
	r = routes.LetterRoutes(r)

	// TODO 群聊路由
	r = routes.ChatRoutes(r)

	// TODO 特判路由
	r = routes.SpecialRoutes(r)

	return r
}
