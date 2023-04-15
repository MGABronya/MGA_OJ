// @Title  caseSample
// @Description  定义了题目的样例
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package model

import uuid "github.com/satori/go.uuid"

// CaseSample			定义了题目的样例
type CaseSample struct {
	Id        uint      `json:"id" gorm:"type:uint;not null"`                                 // 样例Id
	ProblemId uuid.UUID `json:"problem_id" gorm:"type:char(36);index:idx_problemId;not null"` // 题目外键
	Input     string    `json:"input" gorm:"type:text;not null"`                              // 输入
	Output    string    `json:"output" gorm:"type:text;"`                                     // 输出
}
