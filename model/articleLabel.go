// @Title  articleLabel
// @Description  定义文章的标签
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ArticleLabel			定义文章标签
type ArticleLabel struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                          // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                             // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                             // 更新日期
	Label     string    `json:"label" gorm:"type:char(36);index:label;not null"`              // 标签
	ArticleId uuid.UUID `json:"article_id" gorm:"type:char(36);index:idx_articleId;not null"` // 文章外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (articleLable *ArticleLabel) BeforeCreate(scope *gorm.DB) error {
	articleLable.ID = uuid.NewV4()
	return nil
}
