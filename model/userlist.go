// @Title  userlist
// @Description  定义题目列表
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

// UserList		定义题目列表
type UserList struct {
	GroupId uint `json:"group_id" gorm:"type:uint;index:idx_groupId;not null"` // 用户组外键
	UserId  uint `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`   // 用户外键
}
