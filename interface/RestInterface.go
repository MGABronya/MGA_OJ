// @Title  RestInterface
// @Description  该文件用于封装增删查改方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// RestInterface			定义了增删查改方法
type RestInterface interface {
	Create(ctx *gin.Context)   // 增
	Update(ctx *gin.Context)   // 删
	Show(ctx *gin.Context)     // 查
	Delete(ctx *gin.Context)   // 改
	PageList(ctx *gin.Context) // 查看列表
}
