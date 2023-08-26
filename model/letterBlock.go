// @Title  LetterBlock
// @Description  定义表单黑名单
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// LetterBlock			定义表单黑名单
type LetterBlock struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                      // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                         // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                         // 更新日期
	UseraId   uuid.UUID `json:"usera_id" gorm:"type:char(36);index:idx_setId;not null"`   // 拉黑人外键
	UserbId   uuid.UUID `json:"userb_id" gorm:"type:char(36);index:idx_groupId;not null"` // 被拉黑人外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (letterBlock *LetterBlock) BeforeCreate(scope *gorm.DB) error {
	letterBlock.ID = uuid.NewV4()
	return nil
}
