// @Title  CompetitionMember
// @Description  定义比赛参与者问题罚时
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// CompetitionMember		定义比赛参与者问题罚时
type CompetitionMember struct {
	ID            uuid.UUID     `json:"id" gorm:"type:char(36);primary_key"`                                  // id
	CreatedAt     Time          `json:"created_at" gorm:"type:timestamp"`                                     // 创建日期
	UpdatedAt     Time          `json:"updated_at" gorm:"type:timestamp"`                                     // 更新日期
	MemberId      uuid.UUID     `json:"member_id" gorm:"type:char(36);index:idx_memberId;not null"`           // 成员外键
	CompetitionId uuid.UUID     `json:"competition_id" gorm:"type:char(36);index:idx_competitionId;not null"` // 竞赛外键
	ProblemId     uuid.UUID     `json:"accept_num" gorm:"type:char(36);not null"`                             // 题目外键
	Penalties     time.Duration `json:"penalties" gorm:"type:timestamp;not null"`                             // 罚时
	Condition     string        `json:"condition" gorm:"type:varchar(64);not null"`                           // 记录状态
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (competitionMember *CompetitionMember) BeforeCreate(scope *gorm.DB) error {
	competitionMember.ID = uuid.NewV4()
	return nil
}
