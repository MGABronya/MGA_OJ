// @Title  Group
// @Description  定义用户组
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// Group			定义用户组
type Group struct {
	gorm.Model
	LeaderId uint   `json:"leader_id" gorm:"type:uint;index:idx_userId;not null"` // 用户外键
	Title    string `json:"title" gorm:"type:varchar(64);not null"`               // 主题
	Content  string `json:"content" gorm:"type:text;not null"`                    // 内容
	Reslong  string `json:"res_long" gorm:"type:text"`                            // 备用长文本
	Resshort string `json:"res_short" gorm:"type:text"`                           // 备用短文本
	Auto     bool   `json:"auto" gorm:"type:boolean"`                             // 是否自动通过申请
}

// @title    BeforDelete
// @description   关于用户组删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (g *Group) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("group_id = ?", g.ID)

	// TODO 删除用户列表
	tx.Delete(&UserList{})

	// TODO 删除用户组点赞
	tx.Delete(&GroupLike{})

	// TODO 删除用户组收藏
	tx.Delete(&GroupCollect{})

	// TODO 删除用户组申请
	tx.Delete(&GroupApply{})

	// TODO 删除用户组黑名单
	tx.Delete(&GroupBlock{})

	// TODO 删除表单相关
	tx.Delete(&GroupList{})

	return
}
