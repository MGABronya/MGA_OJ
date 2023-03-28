// @Title  chat
// @Description  定义群聊
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
)

// Chat			定义群聊
type Chat struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                      // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                         // 创建日期
	GroupId   uuid.UUID `json:"group_id" gorm:"type:char(36);index:idx_groupId;not null"` // 所属组外键
	Author    uuid.UUID `json:"author_id" gorm:"type:char(36);not null"`                  // 作者用户外键
	Content   string    `json:"content" gorm:"type:text;not null"`                        // 私信的内容
	Reslong   string    `json:"res_long" gorm:"type:text"`                                // 备用长文本
	Resshort  string    `json:"res_short" gorm:"type:text"`                               // 备用短文本
}

// TODO 用于给chat排序
type ChatSlice []Chat

func (c ChatSlice) Len() int { return len(c) }
func (c ChatSlice) Less(i, j int) bool {
	return c[i].CreatedAt.After(c[j].CreatedAt)
}
func (c ChatSlice) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
