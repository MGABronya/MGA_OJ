// @Title  EnterInterface
// @Description  该文件用于封装报名方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// EnterInterface			定义了报名方法
type EnterInterface interface {
	Enter(ctx *gin.Context)          // 报名
	EnterCondition(ctx *gin.Context) // 报名状态
	CancelEnter(ctx *gin.Context)    // 取消报名
	EnterPage(ctx *gin.Context)      // 报名名单
}
