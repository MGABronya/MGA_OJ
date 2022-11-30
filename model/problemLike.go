// @Title  problemLike
// @Description  定义题目的点赞
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// ProblemLike			定义题目点赞
type ProblemLike struct {
	gorm.Model
	UserId    uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`       // 用户外键
	ProblemId uint `json:"problem_id" gorm:"type:uint;index:idx_problemId;not null"` // 题目外键
	Like      bool `json:"like" gorm:"type:boolean;not null"`                        // 点赞或踩
}
