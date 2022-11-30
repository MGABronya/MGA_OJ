// @Title  setLike
// @Description  定义表单的点赞
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// SetLike			定义表单点赞
type SetLike struct {
	gorm.Model
	UserId uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"` // 用户外键
	SetId  uint `json:"set_id" gorm:"type:uint;index:idx_setId;not null"`   // 表单外键
	Like   bool `json:"like" gorm:"type:boolean;not null"`                  // 点赞或踩
}
