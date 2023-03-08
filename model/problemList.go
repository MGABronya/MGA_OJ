// @Title  problemTopic
// @Description  定义题目列表
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import uuid "github.com/satori/go.uuid"

// ProblemList		定义题目列表
type ProblemList struct {
	TopicId   uuid.UUID `json:"topic_id" gorm:"type:char(36);index:idx_topicId;not null"` // 用户外键
	ProblemId uuid.UUID `json:"problem_id" gorm:"type:char(36);not null"`                 // 题目外键
}
