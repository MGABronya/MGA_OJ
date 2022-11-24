// @Title  LikeInterface
// @Description  该文件用于封装点赞方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// LikeInterface			定义了点赞方法
type LikeInterface interface {
	Like(ctx *gin.Context)   // 点赞或点踩
	CancelLike(ctx *gin.Context)     // 取消点赞、点踩状态
	LikeNumber(ctx *gin.Context)   // 点赞、点踩数量
	LikeList(ctx *gin.Context) // 查看点赞、点踩列表
	LikeShow(ctx *gin.Context) // 点赞状态查询
	Likes(ctx *gin.Context)   // 查看用户点赞、点踩列表
}
