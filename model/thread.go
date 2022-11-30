// @Title  thread
// @Description  定义题解的回复
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"gorm.io/gorm"
)

// Thread			定义题解的回复
type Thread struct {
	gorm.Model
	UserId   uint   `json:"user_id" gorm:"type:uint;index:idx_userId;not null"` // 用户外键
	PostId   uint   `json:"post_id" gorm:"type:uint;index:idx_postId;not null"` // 题解外键
	Content  string `json:"content" gorm:"type:text;not null"`                  // 内容
	Reslong  string `json:"res_long" gorm:"type:text"`                          // 备用长文本
	Resshort string `json:"res_short" gorm:"type:text"`                         // 备用短文本
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
