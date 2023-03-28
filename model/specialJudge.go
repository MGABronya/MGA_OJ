// @Title  specailJudge
// @Description  定义特判
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// SpecialJudge			定义特判
type SpecialJudge struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                          // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                             // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                             // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);not null"`                        // 作者id
	Language  string    `json:"language" gorm:"type:varchar(64);index:idx_language;not null"` // 语言
	Code      string    `json:"code" gorm:"type:text;not null"`                               // 代码
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (specialJudge *SpecialJudge) BeforeCreate(scope *gorm.DB) error {
	specialJudge.ID = uuid.NewV4()
	return nil
}
