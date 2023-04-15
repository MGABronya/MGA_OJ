// @Title  HackInterface
// @Description  该文件用于hack功能
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// HackInterface			定义了hack方法
type HackInterface interface {
	Hack(ctx *gin.Context) // hack功能
}
