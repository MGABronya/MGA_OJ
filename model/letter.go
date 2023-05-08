// @Title  letter
// @Description  定义私信
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
)

// Letter			定义私信
type Letter struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                    // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                       // 创建日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"` // 所属用户外键
	Author    uuid.UUID `json:"author_id" gorm:"type:char(36);not null"`                // 作者用户外键
	Content   string    `json:"content" gorm:"type:text;not null"`                      // 私信的内容
	ResLong   string    `json:"res_long" gorm:"type:text"`                              // 备用长文本
	ResShort  string    `json:"res_short" gorm:"type:text"`                             // 备用短文本
	Read      bool      `json:"read" gorm:"type:boolean"`                               // 是否已读
}

// TODO 用于给letter排序
type LetterSlice []Letter

func (l LetterSlice) Len() int { return len(l) }
func (l LetterSlice) Less(i, j int) bool {
	if l[i].Read && l[j].Read || !l[i].Read && !l[j].Read {
		// TODO 最近优先
		return l[i].CreatedAt.After(l[j].CreatedAt)
	} else {
		// TODO 未读优先
		return !l[i].Read
	}
}
func (l LetterSlice) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
