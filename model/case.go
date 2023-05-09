// @Title  case
// @Description  定义了题目的用例
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package model

import uuid "github.com/satori/go.uuid"

// Case			定义了题目的用例
type Case struct {
	CID       uint      `json:"cid" gorm:"type:uint;not null"`                                // 用例Id
	ProblemId uuid.UUID `json:"problem_id" gorm:"type:char(36);index:idx_problemId;not null"` // 题目外键
	Input     string    `json:"input" gorm:"type:text;"`                                      // 输入
	Output    string    `json:"output" gorm:"type:text;"`                                     // 输出
	Score     uint      `json:"score" gorm:"type:uint;"`                                      // 输入分数
}
