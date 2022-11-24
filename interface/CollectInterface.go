// @Title  CollectInterface
// @Description  该文件用于封装收藏方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// CollectInterface			定义了收藏方法
type CollectInterface interface {
	Collect(ctx *gin.Context)       // 收藏
	CancelCollect(ctx *gin.Context) // 取消收藏
	CollectShow(ctx *gin.Context)   // 查看收藏状态
	CollectList(ctx *gin.Context)   // 查看收藏用户列表
	CollectNumber(ctx *gin.Context) // 查看收藏用户数量
	Collects(ctx *gin.Context)      // 查看用户收藏夹
}
