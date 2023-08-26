// @Title  SearchInterface
// @Description  该文件用于封装搜索方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// SearchInterface			定义了搜索方法
type SearchInterface interface {
	Search(ctx *gin.Context)          // 实现文本搜索
	SearchLabel(ctx *gin.Context)     // 实现标签搜索
	SearchWithLabel(ctx *gin.Context) // 实现带标签搜索
}
