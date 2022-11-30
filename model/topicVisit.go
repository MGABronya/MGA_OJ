// @Title  topicVisit
// @Description  定义主题主题的游览
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// TopicVisit			定义主题主题游览
type TopicVisit struct {
	gorm.Model
	UserId  uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`   // 用户外键
	TopicId uint `json:"topic_id" gorm:"type:uint;index:idx_topicId;not null"` // 主题外键
}
