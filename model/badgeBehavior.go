// @Title  badgeBehavior
// @Description  定义徽章
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
)

// BadgeBehavior			定义徽章行为映射
type BadgeBehavior struct {
	BadgeId   uuid.UUID `json:"badge_id" gorm:"type:char(36);primary_key"` // 徽章的id
	Name      string    `json:"name" gorm:"type:char(50);not null;"`       // 行为名称
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`          // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`          // 更新日期
}
