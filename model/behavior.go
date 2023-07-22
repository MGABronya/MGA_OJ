// @Title  behavior
// @Description  定义行为
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
)

// Behavior			定义行为
type Behavior struct {
	Name      string    `json:"name" gorm:"type:char(50);not null;"`                    // 行为名称
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"` // 用户外键
	Score     float64   `json:"score" gorm:"type:double;not null"`                      // 分数
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                       // 徽章的创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                       // 徽章的更新日期
}
