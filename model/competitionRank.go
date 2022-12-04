// @Title  CompetitionRank
// @Description  定义比赛排名
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"time"

	"gorm.io/gorm"
)

// CompetitionRank		定义比赛
type CompetitionRank struct {
	gorm.Model
	MemberId      uint      `json:"member_id" gorm:"type:uint;index:idx_memberId;not null"`           // 成员外键
	CompetitionId uint      `json:"competition_id" gorm:"type:uint;index:idx_competitionId;not null"` // 竞赛外键
	AcceptNum     uint      `json:"accept_num" gorm:"type:uint;not null"`                             // 题目完成数量
	Penalties     time.Time `json:"penalties" gorm:"type:timestamp;not null"`                         // 罚时
}
