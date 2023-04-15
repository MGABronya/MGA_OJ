// @Title  Group
// @Description  定义用户组
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Group			定义用户组
type Group struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                 // id
	CreatedAt     Time      `json:"created_at" gorm:"type:timestamp"`                                    // 创建日期
	UpdatedAt     Time      `json:"updated_at" gorm:"type:timestamp"`                                    // 更新日期
	LeaderId      uuid.UUID `json:"leader_id" gorm:"type:char(36);index:idx_userId;not null"`            // 用户外键
	Title         string    `json:"title" gorm:"type:char(50);not null;index:search_idx,class:FULLTEXT"` // 文章的标题
	Content       string    `json:"content" gorm:"type:text;not null;index:search_idx,class:FULLTEXT"`   // 文章的内容
	Reslong       string    `json:"res_long" gorm:"type:text;index:search_idx,class:FULLTEXT"`           // 备用长文本
	Resshort      string    `json:"res_short" gorm:"type:text;index:search_idx,class:FULLTEXT"`          // 备用短文本
	Auto          bool      `json:"auto" gorm:"type:boolean"`                                            // 是否自动通过申请
	CompetitionAt Time      `json:"competition_at" gorm:"type:timestamp"`                                // 比赛结束时间
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (group *Group) BeforeCreate(scope *gorm.DB) error {
	group.ID = uuid.NewV4()
	return nil
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

	// TODO 删除用户组标签
	tx.Delete(&GroupLabel{})

	return
}
