// @Title  postLike
// @Description  定义题解的点赞
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// PostLike			定义题解点赞
type PostLike struct {
	gorm.Model
	UserId uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"` // 用户外键
	PostId uint `json:"post_id" gorm:"type:uint;index:idx_postId;not null"` // 题解外键
	Like   bool `json:"like" gorm:"type:boolean;not null"`                  // 点赞或踩
}
