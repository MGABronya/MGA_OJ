// @Title  chatgptRoutes
// @Description  程序的chatgpt相关路由均集中在这个文件里
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:50
package routes

import (
	"MGA_OJ/controller"

	"github.com/gin-gonic/gin"
)

// @title    ChatGPTRoutes
// @description   给gin引擎挂上chatgpt相关的路由监听
// @auth      MGAronya             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func ChatGPTRoutes(r *gin.Engine) *gin.Engine {

	// TODO chatgpt管理的路由分组
	chatgptRoutes := r.Group("/chatgpt")

	// TODO 创建chatgptcontroller
	chatgptController := controller.NewChatGPTController()

	// TODO 按照注释生成代码
	chatgptRoutes.POST("/generate/code/:language", chatgptController.GenerateCode)

	// TODO 根据代码生成注释
	chatgptRoutes.POST("/generate/note/:language", chatgptController.GenerateNote)

	// TODO 代码转换
	chatgptRoutes.POST("/change/:language1/:language2", chatgptController.Change)

	// TODO 代码修改意见
	chatgptRoutes.POST("/opinion/:language", chatgptController.Opinion)

	return r
}
