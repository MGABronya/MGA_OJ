// @Title  articleVisit
// @Description  定义题目的游览
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"gorm.io/gorm"
)

// ArticleVisit			定义题目的游览
type ArticleVisit struct {
	gorm.Model
	UserId    uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`       // 用户外键
	ArticleId uint `json:"article_id" gorm:"type:uint;index:idx_articleId;not null"` // 文章外键
}
