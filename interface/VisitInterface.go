// @Title  VisitInterface
// @Description  该文件用于封装游览方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// VisitInterface			定义了游览方法
type VisitInterface interface {
	Visit(ctx *gin.Context)       // 游览
	VisitNumber(ctx *gin.Context) // 游览人数
	VisitList(ctx *gin.Context)   // 游览列表
	Visits(ctx *gin.Context)      // 用户游览历史记录
}
