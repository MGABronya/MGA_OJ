// @Title  Friend
// @Description  定义好友
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// Friend		定义好友
type Friend struct {
	gorm.Model
	UserId   uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`     // 用户外键
	FriendId uint `json:"friend_id" gorm:"type:uint;index:idx_friendId;not null"` // 好友外键
}
