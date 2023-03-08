// @Title  problemLike
// @Description  定义题目的点赞
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ProblemLike			定义题目点赞
type ProblemLike struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                          // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                             // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                             // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`       // 用户外键
	ProblemId uuid.UUID `json:"problem_id" gorm:"type:char(36);index:idx_problemId;not null"` // 题目外键
	Like      bool      `json:"like" gorm:"type:boolean;not null"`                            // 点赞或踩
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (problemLike *ProblemLike) BeforeCreate(scope *gorm.DB) error {
	problemLike.ID = uuid.NewV4()
	return nil
}
