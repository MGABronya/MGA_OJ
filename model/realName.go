// @Title  RealName
// @Description  定义实名
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
)

// RealName			定义实名
type RealName struct {
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`         // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`         // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);not null"`    // 用户id
	Name      string    `json:"name" gorm:"type:char(36);not null"`       // 名字
	StudentId string    `json:"student_id" gorm:"type:char(36);not null"` // 学号
	Major     string    `json:"major" gorm:"type:char(36);not null"`      // 专业
	Grade     string    `json:"grade" gorm:"char(36);not null"`           // 年级
}
