// @Title  badge
// @Description  定义徽章
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import uuid "github.com/satori/go.uuid"

// UserBadge			定义用户解锁的徽章
type UserBadge struct {
	BadgeId   uuid.UUID `json:"badge_id" gorm:"type:char(36);index:idx_badgeId;not null"` // 徽章外键
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`   // 用户外键
	MaxScore  int       `json:"max_score" gorm:"type:int;not null"`                       // 用户累计最大得分
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                         // 徽章的创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                         // 徽章的更新日期
}
