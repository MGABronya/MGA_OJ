// @Title  topicVisit
// @Description  定义主题主题的游览
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// TopicVisit			定义主题主题游览
type TopicVisit struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                      // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                         // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                         // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`   // 用户外键
	TopicId   uuid.UUID `json:"topic_id" gorm:"type:char(36);index:idx_topicId;not null"` // 主题外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (topicVisit *TopicVisit) BeforeCreate(scope *gorm.DB) error {
	topicVisit.ID = uuid.NewV4()
	return nil
}
