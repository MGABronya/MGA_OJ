// @Title  Topic
// @Description  定义主题主题
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Topic			定义主题主题
type Topic struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                 // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                                    // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                                    // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`              // 用户外键
	Title     string    `json:"title" gorm:"type:char(50);not null;index:search_idx,class:FULLTEXT"` // 标题
	Content   string    `json:"content" gorm:"type:text;not null;index:search_idx,class:FULLTEXT"`   // 内容
	ResLong   string    `json:"res_long" gorm:"type:text;index:search_idx,class:FULLTEXT"`           // 备用长文本
	ResShort  string    `json:"res_short" gorm:"type:text;index:search_idx,class:FULLTEXT"`          // 备用短文本
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (topic *Topic) BeforeCreate(scope *gorm.DB) error {
	topic.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于主题删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (t *Topic) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("topic_id = (?)", t.ID)

	// TODO 删除主题收藏
	tx.Delete(&TopicCollect{})

	// TODO 删除主题点赞
	tx.Delete(&TopicLike{})

	// TODO 删除主题访问
	tx.Delete(&TopicVisit{})

	// TODO 删除主题题目
	tx.Delete(&ProblemList{})

	// TODO 删除表单相关
	tx.Delete(&TopicList{})

	return
}
