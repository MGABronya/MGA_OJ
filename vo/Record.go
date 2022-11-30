// @Title  Record
// @Description  定义提交
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

// RecordRequest			定义提交
type RecordRequest struct {
	ProblemId uint   `json:"problem_id"` // 题目外键
	Language  string `json:"language"`   // 语言
	Code      string `json:"code"`       // 代码
}
