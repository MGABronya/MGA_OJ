// @Title  Hack
// @Description  定义Hack
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Hack		定义Hack
type Hack struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                    // id
	RecordId  uuid.UUID `json:"record_id" gorm:"type:char(36);not null"`                // 记录外键
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"` // 用户外键
	Input     string    `json:"input" gorm:"type:text;"`                                // 输入
	Type      string    `json:"type" gorm:"type:char(20);"`                             // 类型
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                       // 创建日期
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (hack *Hack) BeforeCreate(scope *gorm.DB) error {
	hack.ID = uuid.NewV4()
	return nil
}