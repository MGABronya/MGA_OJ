// @Title  thread
// @Description  定义题解的回复
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Thread			定义题解的回复
type Thread struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                    // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                       // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                       // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"` // 用户外键
	PostId    uuid.UUID `json:"post_id" gorm:"type:char(36);index:idx_postId;not null"` // 题解外键
	Content   string    `json:"content" gorm:"type:text;not null"`                      // 内容
	Reslong   string    `json:"res_long" gorm:"type:text"`                              // 备用长文本
	Resshort  string    `json:"res_short" gorm:"type:text"`                             // 备用短文本
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (thread *Thread) BeforeCreate(scope *gorm.DB) error {
	thread.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于回复删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (t *Thread) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("thread_id = ?", t.ID)

	// TODO 删除回复点赞
	tx.Delete(&ThreadLike{})

	return
}
