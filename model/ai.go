// @Title  AI
// @Description  AI模板
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:47
package model

import uuid "github.com/satori/go.uuid"

// AI			接收一个AI模板
type AI struct {
	Characters string    `json:"characters" gorm:"type:text;not null"`                   // 人设
	Reply      bool      `json:"reply" gorm:"type:boolean;not null"`                     // 是否回复自己
	Prologue   string    `json:"prologue" gorm:"type:text;not null"`                     // 开场白
	UserId     uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"` // 用户外键
}
