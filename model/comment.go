// @Title  comment
// @Description  定义讨论
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"gorm.io/gorm"
)

// Comment			定义讨论
type Comment struct {
	gorm.Model
	UserId    uint   `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`       // 用户外键
	ProblemId uint   `json:"problem_id" gorm:"type:uint;index:idx_problemId;not null"` // 题目外键
	Content   string `json:"content" gorm:"type:text;not null"`                        // 内容
	Reslong   string `json:"res_long" gorm:"type:text"`                                // 备用长文本
	Resshort  string `json:"res_short" gorm:"type:text"`                               // 备用短文本
}
