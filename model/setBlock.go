// @Title  SetBlock
// @Description  定义表单黑名单
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// SetBlock			定义表单黑名单
type SetBlock struct {
	gorm.Model
	SetId   uint `json:"set_id" gorm:"type:uint;index:idx_setId;not null"`     // 表单外键
	GroupId uint `json:"group_id" gorm:"type:uint;index:idx_groupId;not null"` // 用户组外键
}
