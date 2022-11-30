// @Title  setCollect
// @Description  定义表单的收藏
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// SetCollect			定义表单收藏
type SetCollect struct {
	gorm.Model
	UserId uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"` // 用户外键
	SetId  uint `json:"set_id" gorm:"type:uint;index:idx_setId;not null"`   // 主题外键
}
