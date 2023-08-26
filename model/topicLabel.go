// @Title  topicLabel
// @Description  定义主题的标签
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// TopicLabel			定义主题标签
type TopicLabel struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                       // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                          // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                          // 更新日期
	Label     string    `json:"label" gorm:"type:char(36);index:label;not null"`           // 标签
	TopicId   uuid.UUID `json:"topic_id" gorm:"type:char(36);index:idx_topic=Id;not null"` // 主题外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (topicLabel *TopicLabel) BeforeCreate(scope *gorm.DB) error {
	topicLabel.ID = uuid.NewV4()
	return nil
}
