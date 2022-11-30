// @Title  postVisit
// @Description  定义题解的游览
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"gorm.io/gorm"
)

// PostVisit			定义题解的游览
type PostVisit struct {
	gorm.Model
	UserId uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"` // 用户外键
	PostId uint `json:"post_id" gorm:"type:uint;index:idx_postId;not null"` // 题解外键
}
