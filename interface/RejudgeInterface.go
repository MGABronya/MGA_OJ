// @Title  RejudgeInterface
// @Description  该文件用于封装比赛重判方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

import "github.com/gin-gonic/gin"

// RejudgeInterface		定义了比赛重判方法
type RejudgeInterface interface {
	Rejudge(ctx *gin.Context)               // 进行重判
	CompetitionDataDelete(ctx *gin.Context) // 对比赛结果进行清空
}
