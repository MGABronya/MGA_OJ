// @Title  message
// @Description  定义留言
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Message			定义留言
type Message struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                    // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                       // 创建日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"` // 所属用户外键
	Author    uuid.UUID `json:"author_id" gorm:"type:char(36);not null"`                // 作者用户外键
	Content   string    `json:"content" gorm:"type:char(72);not null"`                  // 留言的内容
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (message *Message) BeforeCreate(scope *gorm.DB) error {
	message.ID = uuid.NewV4()
	return nil
}
