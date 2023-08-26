// @Title  userScoreChange
// @Description  定义用户的分数变化
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// UserScoreChange			定义用户分数变化
type UserScoreChange struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                  // id
	CreatedAt     Time      `json:"created_at" gorm:"type:timestamp"`                                     // 创建日期
	UpdatedAt     Time      `json:"updated_at" gorm:"type:timestamp"`                                     // 更新日期
	ScoreChange   float64   `json:"score_change" gorm:"type:double;not null"`                             // 分数变化
	CompetitionId uuid.UUID `json:"competition_id" gorm:"type:char(36);index:idx_competitionId;not null"` // 竞赛外键
	UserId        uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`               // 用户外键
	Type          string    `json:"type" gorm:"type:char(20);not null"`                                   // 竞赛类型
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (userScoreChange *UserScoreChange) BeforeCreate(scope *gorm.DB) error {
	userScoreChange.ID = uuid.NewV4()
	return nil
}
