// @Title  articleCollect
// @Description  定义文章的收藏
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// ArticleCollect			定义文章收藏
type ArticleCollect struct {
	gorm.Model
	UserId    uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`       // 用户外键
	ArticleId uint `json:"article_id" gorm:"type:uint;index:idx_articleId;not null"` // 文章外键
}
