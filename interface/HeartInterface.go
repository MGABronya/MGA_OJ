// @Title  HeartInterface
// @Description  该文件用于封装心跳方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// HeartInterface			定义了心跳方法
type HeartInterface interface {
	Publish(ctx *gin.Context)    // 订阅心跳长连接
	Show(ctx *gin.Context)       // 查看最近心跳情况
	Percentage(ctx *gin.Context) // 最近10s的心跳忙碌占比
}
