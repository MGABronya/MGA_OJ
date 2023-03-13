// @Title  user
// @Description  定义用户
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// user			定义用户
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                          // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                                             // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                                             // 更新日期
	Name      string    `json:"name" gorm:"type:varchar(20);not null;unique;index:search_idx,class:FULLTEXT"` // 用户名称
	Email     string    `json:"email" gorm:"type:varchar(50);not null;unique"`                                // 邮箱
	Password  string    `json:"password" gorm:"size:255;not null"`                                            // 密码
	Icon      string    `json:"icon" gorm:"type:varchar(50)"`                                                 // 这里的Icon存储的是图像文件的url后缀
	Blog      string    `json:"blog" gorm:"type:varchar(25)"`                                                 // 博客
	Sex       bool      `json:"sex" gorm:"type:boolean"`                                                      // 性别
	Address   string    `json:"address" gorm:"type:varchar(20)"`                                              // 地址
	Level     int       `json:"level" gorm:"type:int;not null"`                                               // 用户管理等级
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (user *User) BeforeCreate(scope *gorm.DB) error {
	user.ID = uuid.NewV4()
	return nil
}
