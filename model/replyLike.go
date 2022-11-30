// @Title  replyLike
// @Description  定义回复的点赞
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// ReplyLike			定义回复点赞
type ReplyLike struct {
	gorm.Model
	UserId  uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`   // 用户外键
	ReplyId uint `json:"reply_id" gorm:"type:uint;index:idx_replyId;not null"` // 回复外键
	Like    bool `json:"like" gorm:"type:boolean;not null"`                    // 点赞或踩
}
