// @Title  topicList
// @Description  定义主题列表
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

// TopicList		定义主题列表
type TopicList struct {
	SetId   uint `json:"set_id" gorm:"type:uint;index:idx_setId;not null"` // 表单外键
	TopicId uint `json:"topic_id" gorm:"type:uint;not null"`               // 主题外键
}
