// @Title  ApplyInterface
// @Description  该文件用于封装申请方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// ApplyInterface			定义了申请方法
type ApplyInterface interface {
	ApplyingList(ctx *gin.Context) // 发出请求方看到的请求列表
	AppliedList(ctx *gin.Context)  // 接收请求方看到的请求列表
	Apply(ctx *gin.Context)        // 发出请求
	Consent(ctx *gin.Context)      // 通过请求
	Refuse(ctx *gin.Context)       // 拒绝申请
	Quit(ctx *gin.Context)         // 退出/删除
}
