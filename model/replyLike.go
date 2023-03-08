// @Title  replyLike
// @Description  定义回复的点赞
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ReplyLike			定义回复点赞
type ReplyLike struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                      // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                         // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                         // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`   // 用户外键
	ReplyId   uuid.UUID `json:"reply_id" gorm:"type:char(36);index:idx_replyId;not null"` // 回复外键
	Like      bool      `json:"like" gorm:"type:boolean;not null"`                        // 点赞或踩
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (replyLike *ReplyLike) BeforeCreate(scope *gorm.DB) error {
	replyLike.ID = uuid.NewV4()
	return nil
}
