// @Title  LabelInterface
// @Description  该文件用于封装标签方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// LabelInterface			定义了标签方法
type LabelInterface interface {
	LabelCreate(ctx *gin.Context)       // 增设标签
	LabelDelete(ctx *gin.Context)       // 删除标签
	LabelShow(ctx *gin.Context)         // 查看标签
}
