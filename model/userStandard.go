// @Title  UserStandard
// @Description  定义标准用户
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package model

import uuid "github.com/satori/go.uuid"

// UserStandard			定义标准用户
type UserStandard struct {
	Email    string    `json:"email" gorm:"type:varchar(50);not null;unique"` // 邮箱
	Password string    `json:"password" gorm:"size:255;not null"`             // 密码
	CID      uuid.UUID `json:"cid" gorm:"type:char(36);primary_key"`          // 外键id
}
