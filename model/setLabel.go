// @Title  setLabel
// @Description  定义表单的标签
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// SetLabel			定义表单标签
type SetLabel struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                   // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                      // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                      // 更新日期
	Label     string    `json:"label" gorm:"type:char(36);index:label;not null"`       // 标签
	SetId     uuid.UUID `json:"set_id" gorm:"type:char(36);index:idx_set=Id;not null"` // 表单外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (setLabel *SetLabel) BeforeCreate(scope *gorm.DB) error {
	setLabel.ID = uuid.NewV4()
	return nil
}
