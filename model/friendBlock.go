// @Title  FriendBlock
// @Description  定义用户黑名单
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// FriendBlock			定义用户黑名单
type FriendBlock struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                      // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                         // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                         // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`   // 被拉黑的用户外键
	OwnerId   uuid.UUID `json:"owner_id" gorm:"type:char(36);index:idx_ownerId;not null"` // 拥有者外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (friendBlock *FriendBlock) BeforeCreate(scope *gorm.DB) error {
	friendBlock.ID = uuid.NewV4()
	return nil
}
