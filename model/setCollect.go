// @Title  setCollect
// @Description  定义表单的收藏
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// SetCollect			定义表单收藏
type SetCollect struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                    // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                       // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                       // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"` // 用户外键
	SetId     uuid.UUID `json:"set_id" gorm:"type:char(36);index:idx_setId;not null"`   // 主题外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (setCollect *SetCollect) BeforeCreate(scope *gorm.DB) error {
	setCollect.ID = uuid.NewV4()
	return nil
}
