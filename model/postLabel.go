// @Title  postLabel
// @Description  定义题解的标签
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// PostLabel			定义题解标签
type PostLabel struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                     // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                        // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                        // 更新日期
	Label     string    `json:"label" gorm:"type:char(36);index:label;not null"`         // 标签
	PostId    uuid.UUID `json:"post_id" gorm:"type:char(36);index:idx_post=Id;not null"` // 题解外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (postLabel *PostLabel) BeforeCreate(scope *gorm.DB) error {
	postLabel.ID = uuid.NewV4()
	return nil
}
