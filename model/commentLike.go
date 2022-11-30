// @Title  commentLike
// @Description  定义讨论的点赞
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// CommentLike			定义讨论点赞
type CommentLike struct {
	gorm.Model
	UserId    uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`       // 用户外键
	CommentId uint `json:"comment_id" gorm:"type:uint;index:idx_commentId;not null"` // 讨论外键
	Like      bool `json:"like" gorm:"type:boolean;not null"`                        // 点赞或踩
}
