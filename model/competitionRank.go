// @Title  CompetitionRank
// @Description  定义比赛成员与其分数
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// CompetitionRank		定义比赛成员与其分数
type CompetitionRank struct {
	ID            uuid.UUID     `json:"id" gorm:"type:char(36);primary_key"`                                  // id
	CreatedAt     Time          `json:"created_at" gorm:"type:timestamp"`                                     // 创建日期
	UpdatedAt     Time          `json:"updated_at" gorm:"type:timestamp"`                                     // 更新日期
	MemberId      uuid.UUID     `json:"member_id" gorm:"type:char(36);index:idx_memberId;not null"`           // 成员外键
	CompetitionId uuid.UUID     `json:"competition_id" gorm:"type:char(36);index:idx_competitionId;not null"` // 竞赛外键
	Score         uint          `json:"score" gorm:"type:uint;not null"`                                      // 获得分数
	Penalties     time.Duration `json:"penalties" gorm:"type:int;not null"`                                   // 罚时
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (competitionRank *CompetitionRank) BeforeCreate(scope *gorm.DB) error {
	competitionRank.ID = uuid.NewV4()
	return nil
}
