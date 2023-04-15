// @Title  HackNum
// @Description  定义HackNum
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import uuid "github.com/satori/go.uuid"

// HackNum		定义某场比赛中某用户hack提交的数量
type HackNum struct {
	CompetitionId uuid.UUID `json:"competition_id" gorm:"type:char(36);index:idx_competitionId;not null"` // 竞赛外键
	MemberId      uuid.UUID `json:"member_id" gorm:"type:char(36);index:idx_memberId;not null"`           // 成员外键
	Num           uint      `json:"num" gorm:"type:uint;"`                                                // 输入
	Score         uint      `json:"score" gorm:"type:uint;"`                                              // hack获得的分数
}
