// @Title  groupList
// @Description  定义组别列表
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
)

// GroupList		定义用户组列表
type GroupList struct {
	SetId   uuid.UUID `json:"set_id" gorm:"type:char(36);index:idx_setId;not null"` // 表单外键
	GroupId uuid.UUID `json:"group_id" gorm:"type:char(36);not null"`               // 用户组外键
}
