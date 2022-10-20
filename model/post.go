// @Title  post
// @Description  定义题解
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"gorm.io/gorm"
)

// Post			定义题解
type Post struct {
	gorm.Model
	UserId    uint   `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`       // 用户外键
	ProblemId uint   `json:"problem_id" gorm:"type:uint;index:idx_problemId;not null"` // 题目外键
	Title     string `json:"title" gorm:"type:varchar(64);not null"`                   // 题目
	Content   string `json:"content" gorm:"type:text;not null"`                        // 内容
	Reslong   string `json:"res_long" gorm:"type:text"`                                // 备用长文本
	Resshort  string `json:"res_short" gorm:"type:text"`                               // 备用短文本
}
