// @Title  PasswdInterface
// @Description  该文件用于封装密码方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// PasswdInterface			定义了密码方法
type PasswdInterface interface {
	CreatePasswd(ctx *gin.Context) // 创建密码
	DeletePasswd(ctx *gin.Context) // 删除密码
}
