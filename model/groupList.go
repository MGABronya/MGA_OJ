// @Title  groupList
// @Description  定义组别列表
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

// GroupList		定义用户组列表
type GroupList struct {
	SetId   uint `json:"set_id" gorm:"type:uint;index:idx_setId;not null"` // 表单外键
	GroupId uint `json:"group_id" gorm:"type:uint;not null"`               // 用户组外键
}
