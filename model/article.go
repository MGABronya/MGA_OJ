// @Title  article
// @Description  定义文章
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Article			定义文章
type Article struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                 // 文章的id
	CreatedAt  Time      `json:"created_at" gorm:"type:timestamp"`                                    // 文章的创建日期
	UpdatedAt  Time      `json:"updated_at" gorm:"type:timestamp"`                                    // 文章的更新日期
	UserId     uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`              // 用户外键
	CategoryId uuid.UUID `json:"category_id" gorm:"type:char(36);index:idx_categoryId;not null"`      // 分类外键
	Title      string    `json:"title" gorm:"type:char(50);not null;index:search_idx,class:FULLTEXT"` // 文章的标题
	Content    string    `json:"content" gorm:"type:text;not null;index:search_idx,class:FULLTEXT"`   // 文章的内容
	ResLong    string    `json:"res_long" gorm:"type:text;index:search_idx,class:FULLTEXT"`           // 备用长文本
	ResShort   string    `json:"res_short" gorm:"type:text;index:search_idx,class:FULLTEXT"`          // 备用短文本
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (article *Article) BeforeCreate(scope *gorm.DB) error {
	article.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于文章删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (a *Article) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("article_id = (?)", a.ID)

	// TODO 删除文章收藏
	tx.Delete(&ArticleCollect{})

	// TODO 删除文章点赞
	tx.Delete(&ArticleLike{})

	// TODO 删除文章访问
	tx.Delete(&ArticleVisit{})

	// TODO 删除文章回复
	tx.Delete(&Remark{})

	// TODO 删除文章标签
	tx.Delete(&ArticleLabel{})

	return
}
