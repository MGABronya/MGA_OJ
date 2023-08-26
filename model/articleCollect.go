// @Title  articleCollect
// @Description  定义文章的收藏
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ArticleCollect			定义文章收藏
type ArticleCollect struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                          // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                             // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                             // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`       // 用户外键
	ArticleId uuid.UUID `json:"article_id" gorm:"type:char(36);index:idx_articleId;not null"` // 文章外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (articleCollect *ArticleCollect) BeforeCreate(scope *gorm.DB) error {
	articleCollect.ID = uuid.NewV4()
	return nil
}
