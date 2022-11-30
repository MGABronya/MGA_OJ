// @Title  threadLike
// @Description  定义跟帖的点赞
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// ThreadLike			定义跟帖点赞
type ThreadLike struct {
	gorm.Model
	UserId   uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`     // 用户外键
	ThreadId uint `json:"thread_id" gorm:"type:uint;index:idx_threadId;not null"` // 跟帖外键
	Like     bool `json:"like" gorm:"type:boolean;not null"`                      // 点赞或踩
}
