// @Title  remarkLike
// @Description  定义文章回复的点赞
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// RemarkLike			定义文章回复点赞
type RemarkLike struct {
	gorm.Model
	UserId   uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`     // 用户外键
	RemarkId uint `json:"remark_id" gorm:"type:uint;index:idx_remarkId;not null"` // 文章回复外键
	Like     bool `json:"like" gorm:"type:boolean;not null"`                      // 点赞或踩
}
