// @Title  testOutput
// @Description  定义了题目的输出
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package model

import (
	uuid "github.com/satori/go.uuid"
)

// TestOutput			定义了题目的输出
type TestOutput struct {
	Id        uint      `json:"id" gorm:"type:uint;not null"`                                 // 输出Id
	ProblemId uuid.UUID `json:"problem_id" gorm:"type:char(36);index:idx_problemId;not null"` // 题目外键
	Output    string    `json:"input" gorm:"type:text;not null"`                              // 输出
}
