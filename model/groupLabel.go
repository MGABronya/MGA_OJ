// @Title  groupLabel
// @Description  定义用户组的标签
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// GroupLabel			定义比赛标签
type GroupLabel struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                       // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                          // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                          // 更新日期
	Label     string    `json:"label" gorm:"type:char(36);index:label;not null"`           // 标签
	GroupId   uuid.UUID `json:"group_id" gorm:"type:char(36);index:idx_group=Id;not null"` // 用户组外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (groupLabel *GroupLabel) BeforeCreate(scope *gorm.DB) error {
	groupLabel.ID = uuid.NewV4()
	return nil
}
