// @Title  FriendApply
// @Description  定义好友申请
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// FriendApply			定义好友申请
type FriendApply struct {
	gorm.Model
	UserId    uint   `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`     // 用户外键
	FriendId  uint   `json:"friend_id" gorm:"type:uint;index:idx_friendId;not null"` // 好友外键
	Condition bool   `json:"condition" gorm:"type:boolean;not null"`                 // 申请状态
	Content   string `json:"content" gorm:"type:text;not null"`                      // 内容
	Reslong   string `json:"res_long" gorm:"type:text"`                              // 备用长文本
	Resshort  string `json:"res_short" gorm:"type:text"`                             // 备用短文本
}
