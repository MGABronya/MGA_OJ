// @Title  FriendBlock
// @Description  定义用户黑名单
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// FriendBlock			定义用户黑名单
type FriendBlock struct {
	gorm.Model
	UserId  uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`   // 被拉黑的用户外键
	OwnerId uint `json:"owner_id" gorm:"type:uint;index:idx_ownerId;not null"` // 拥有者外键
}
