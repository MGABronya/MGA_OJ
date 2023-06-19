// @Title  post
// @Description  定义题解
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Post			定义题解
type Post struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                 // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                                    // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                                    // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`              // 用户外键
	ProblemId uuid.UUID `json:"problem_id" gorm:"type:char(36);index:idx_problemId;not null"`        // 题目外键
	Title     string    `json:"title" gorm:"type:char(50);not null;index:search_idx,class:FULLTEXT"` // 标题
	Content   string    `json:"content" gorm:"type:text;not null;index:search_idx,class:FULLTEXT"`   // 内容
	ResLong   string    `json:"res_long" gorm:"type:text;index:search_idx,class:FULLTEXT"`           // 备用长文本
	ResShort  string    `json:"res_short" gorm:"type:text;index:search_idx,class:FULLTEXT"`          // 备用短文本
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (post *Post) BeforeCreate(scope *gorm.DB) error {
	post.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于题解删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (p *Post) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("post_id = (?)", p.ID)

	// TODO 删除题解收藏
	tx.Delete(&PostCollect{})

	// TODO 删除题解点赞
	tx.Delete(&PostLike{})

	// TODO 删除题解访问
	tx.Delete(&PostVisit{})

	// TODO 删除题解回复
	tx.Delete(&Thread{})

	// TODO 删除题解标签
	tx.Delete(&PostLabel{})

	return
}
