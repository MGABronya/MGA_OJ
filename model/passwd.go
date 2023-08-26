// @Title  Passwd
// @Description  定义密码
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Passwd			定义密码
type Passwd struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primary_key"` // 密码的id
	Password string    `json:"password" gorm:"size:255;not null"`   // 密码
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (passwd *Passwd) BeforeCreate(scope *gorm.DB) error {
	passwd.ID = uuid.NewV4()
	return nil
}
