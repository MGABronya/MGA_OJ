// @Title  groupCollect
// @Description  定义用户组的收藏
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// GroupCollect			定义用户组收藏
type GroupCollect struct {
	gorm.Model
	UserId  uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`   // 用户外键
	GroupId uint `json:"group_id" gorm:"type:uint;index:idx_groupId;not null"` // 用户组外键
}
