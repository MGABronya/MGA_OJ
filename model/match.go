// @Title  match
// @Description  定义匹配列表
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import uuid "github.com/satori/go.uuid"

// Match		定义匹配列表
type Match struct {
	CompetitionId uuid.UUID `json:"competition_id" gorm:"type:char(36);index:idx_competitionId;not null"` // 比赛外键
	UserId        uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`               // 用户外键
}
