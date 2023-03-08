// @Title  topicList
// @Description  定义主题列表
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import uuid "github.com/satori/go.uuid"

// TopicList		定义主题列表
type TopicList struct {
	SetId   uuid.UUID `json:"set_id" gorm:"type:char(36);index:idx_setId;not null"` // 表单外键
	TopicId uuid.UUID `json:"topic_id" gorm:"type:char(36);not null"`               // 主题外键
}
