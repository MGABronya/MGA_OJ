// @Title  SetApply
// @Description  定义表单申请
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// SetApply			定义表单申请
type SetApply struct {
	gorm.Model
	SetId     uint   `json:"set_id" gorm:"type:uint;index:idx_setId;not null"`     // 表单外键
	GroupId   uint   `json:"group_id" gorm:"type:uint;index:idx_groupId;not null"` // 用户组外键
	Condition bool   `json:"condition" gorm:"type:boolean;not null"`               // 申请状态
	Content   string `json:"content" gorm:"type:text;not null"`                    // 内容
	Reslong   string `json:"res_long" gorm:"type:text"`                            // 备用长文本
	Resshort  string `json:"res_short" gorm:"type:text"`                           // 备用短文本
}
