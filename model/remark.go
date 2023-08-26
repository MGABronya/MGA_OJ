// @Title  remark
// @Description  定义文章的回复
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Remark			定义文章的回复
type Remark struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                          // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                             // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                             // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`       // 用户外键
	ArticleId uuid.UUID `json:"article_id" gorm:"type:char(36);index:idx_articleId;not null"` // 文章外键
	Content   string    `json:"content" gorm:"type:text;not null"`                            // 内容
	ResLong   string    `json:"res_long" gorm:"type:text"`                                    // 备用长文本
	ResShort  string    `json:"res_short" gorm:"type:text"`                                   // 备用短文本
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (remark *Remark) BeforeCreate(scope *gorm.DB) error {
	remark.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于回复删除的一些级联操作
// @auth      MGAronya       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (r *Remark) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("remark_id = (?)", r.ID)

	// TODO 删除回复点赞
	tx.Delete(&RemarkLike{})

	return
}
