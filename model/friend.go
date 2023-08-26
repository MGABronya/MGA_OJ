// @Title  Friend
// @Description  定义好友
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Friend		定义好友
type Friend struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                        // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                           // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                           // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`     // 用户外键
	FriendId  uuid.UUID `json:"friend_id" gorm:"type:char(36);index:idx_friendId;not null"` // 好友外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (friend *Friend) BeforeCreate(scope *gorm.DB) error {
	friend.ID = uuid.NewV4()
	return nil
}
