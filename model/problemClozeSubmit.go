// @Title  problemMCQsSubmit
// @Description  定义了填空题的提交
// @Author  MGAronya
// @Update  MGAronya  2022-10-17 21:07
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// problemClozeSubmit			定义了填空题的提交
type ProblemClozeSubmit struct {
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                               // id
	CreatedAt      Time      `json:"created_at" gorm:"type:timestamp"`                                  // 创建日期
	UpdatedAt      Time      `json:"updated_at" gorm:"type:timestamp"`                                  // 更新日期
	UserId         uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`            // 用户外键
	ProblemClozeId uuid.UUID `json:"problem_cloze_id" gorm:"type:char(36);index:idx_cloze_Id;not null"` // 题目外键
	Answer         string    `json:"answer" gorm:"type:text;"`                                          // 该题答案
	Score          uint      `json:"score" gorm:"type:uint;"`                                           // 该题分数
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (problemCloze *ProblemClozeSubmit) BeforeCreate(scope *gorm.DB) error {
	problemCloze.ID = uuid.NewV4()
	return nil
}
