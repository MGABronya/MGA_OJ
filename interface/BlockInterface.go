// @Title  BlockInterface
// @Description  该文件用于封装黑名单方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// BlockInterface			定义了黑名单方法
type BlockInterface interface {
	Block(ctx *gin.Context)       // 拉黑某用
	BlackList(ctx *gin.Context)   // 查看黑名单
	RemoveBlack(ctx *gin.Context) // 移除黑名单
}
