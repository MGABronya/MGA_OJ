// @Title  userlist
// @Description  定义题目列表
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import uuid "github.com/satori/go.uuid"

// UserList		定义题目列表
type UserList struct {
	GroupId uuid.UUID `json:"group_id" gorm:"type:char(36);index:idx_groupId;not null"` // 用户组外键
	UserId  uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`   // 用户外键
}
