// @Title  remarkLike
// @Description  定义文章回复的点赞
// @Author  MGAronya
// @Update  MGAronya  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// RemarkLike			定义文章回复点赞
type RemarkLike struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                        // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                           // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                           // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`     // 用户外键
	RemarkId  uuid.UUID `json:"remark_id" gorm:"type:char(36);index:idx_remarkId;not null"` // 文章回复外键
	Like      bool      `json:"like" gorm:"type:boolean;not null"`                          // 点赞或踩
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (remarkLike *RemarkLike) BeforeCreate(scope *gorm.DB) error {
	remarkLike.ID = uuid.NewV4()
	return nil
}
