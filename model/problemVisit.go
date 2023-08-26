// @Title  problemVisit
// @Description  定义题目的游览
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ProblemVisit			定义题目的游览
type ProblemVisit struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                          // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                             // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                             // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`       // 用户外键
	ProblemId uuid.UUID `json:"problem_id" gorm:"type:char(36);index:idx_problemId;not null"` // 题目外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (problemVisit *ProblemVisit) BeforeCreate(scope *gorm.DB) error {
	problemVisit.ID = uuid.NewV4()
	return nil
}
