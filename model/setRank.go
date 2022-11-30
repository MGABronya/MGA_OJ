// @Title  setRank
// @Description  定义表单用户排名
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// setRank			定义表单用户排名
type SetRank struct {
	gorm.Model
	UserId uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"` // 用户外键
	SetId  uint `json:"set_id" gorm:"type:uint;index:idx_setId;not null"`   // 表单外键
	Pass   uint `json:"pass" gorm:"type:uint;not null"`                     // 表单内通过数量
}
