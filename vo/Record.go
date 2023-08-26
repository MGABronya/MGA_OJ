// @Title  Record
// @Description  定义提交
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

import (
	uuid "github.com/satori/go.uuid"
)

// RecordRequest			定义提交
type RecordRequest struct {
	ProblemId uuid.UUID `json:"problem_id"` // 题目外键
	Language  string    `json:"language"`   // 语言
	Code      string    `json:"code"`       // 代码
}
