// @Title  RecordInterface
// @Description  该文件用于封装代码提交记录方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// RecordInterface		定义了代码提交记录方法
type RecordInterface interface {
	Submit(ctx *gin.Context)          // 用户进行提交操作
	ShowRecord(ctx *gin.Context)      // 用户查看指定提交
	SearchList(ctx *gin.Context)      // 用户搜索提交列表
	CaseList(ctx *gin.Context)        // 某次提交的具体测试通过情况
	Case(ctx *gin.Context)            // 某个测试的具体情况
	PublishPageList(ctx *gin.Context) // 订阅提交列表
	Publish(ctx *gin.Context)         // 订阅某个提交
}
