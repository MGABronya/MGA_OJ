// @Title  HotInterface
// @Description  该文件用于封装热度方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// HotInterface			定义了热度方法
type HotInterface interface {
	HotRanking(ctx *gin.Context) // 热度排行
}
