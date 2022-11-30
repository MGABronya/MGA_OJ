// @Title  testInput
// @Description  定义了题目的输入
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package model

// TestInput			定义了题目的输入
type TestInput struct {
	Id        uint   `json:"id" gorm:"type:uint;not null"`                             // 输入Id
	ProblemId uint   `json:"problem_id" gorm:"type:uint;index:idx_problemId;not null"` // 题目外键
	Input     string `json:"input" gorm:"type:text;not null"`                          // 输入
}
