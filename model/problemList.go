// @Title  problemTopic
// @Description  定义题目列表
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

// ProblemList		定义题目列表
type ProblemList struct {
	TopicId   uint `json:"topic_id" gorm:"type:uint;index:idx_topicId;not null"` // 用户外键
	ProblemId uint `json:"problem_id" gorm:"type:uint;not null"`                 // 题目外键
}
