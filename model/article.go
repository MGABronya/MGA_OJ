// @Title  article
// @Description  定义文章
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"gorm.io/gorm"
)

// Article			定义文章
type Article struct {
	gorm.Model
	UserId     uint   `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`         // 用户外键
	CategoryId uint   `json:"category_id" gorm:"type:uint;index:idx_categoryId;not null"` // 分类外键
	Title      string `json:"title" gorm:"type:varchar(64);not null"`                     // 题目
	Content    string `json:"content" gorm:"type:text;not null"`                          // 内容
	Reslong    string `json:"res_long" gorm:"type:text"`                                  // 备用长文本
	Resshort   string `json:"res_short" gorm:"type:text"`                                 // 备用短文本
}

// @title    BeforDelete
// @description   关于文章删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (a *Article) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("article_id = ?", a.ID)

	// TODO 删除文章收藏
	tx.Delete(&ArticleCollect{})

	// TODO 删除文章点赞
	tx.Delete(&ArticleLike{})

	// TODO 删除文章访问
	tx.Delete(&ArticleVisit{})

	// TODO 删除文章回复
	tx.Delete(&Remark{})

	return
}
