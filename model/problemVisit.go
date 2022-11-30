// @Title  problemVisit
// @Description  定义题目的游览
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"gorm.io/gorm"
)

// ProblemVisit			定义题目的游览
type ProblemVisit struct {
	gorm.Model
	UserId    uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`       // 用户外键
	ProblemId uint `json:"problem_id" gorm:"type:uint;index:idx_problemId;not null"` // 题目外键
}
