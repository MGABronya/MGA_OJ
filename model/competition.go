// @Title  Competition
// @Description  定义比赛
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// Competition		定义比赛
type Competition struct {
	gorm.Model
	SetId     uint   `json:"set_id" gorm:"type:uint;not null"`          // 表单外键
	StartTime Time   `json:"start_time" gorm:"type:timestamp;not null"` // 起始时间
	EndTime   Time   `json:"end_time" gorm:"type:timestamp;not null"`   // 终止时间
	Type      string `json:"type" gorm:"type:varchar(20);not null"`     // 比赛类型
}
