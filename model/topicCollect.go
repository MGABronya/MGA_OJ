// @Title  topicCollect
// @Description  定义主题主题的收藏
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// TopicCollect			定义主题主题收藏
type TopicCollect struct {
	gorm.Model
	UserId  uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`   // 用户外键
	TopicId uint `json:"topic_id" gorm:"type:uint;index:idx_topicId;not null"` // 主题外键
}
