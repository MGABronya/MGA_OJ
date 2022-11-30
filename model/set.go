// @Title  set
// @Description  定义表单
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"gorm.io/gorm"
)

// Set			定义表单
type Set struct {
	gorm.Model
	UserId     uint   `json:"user_id" gorm:"type:uint;index:idx_userId;not null"` // 用户外键
	Title      string `json:"title" gorm:"type:varchar(64);not null"`             // 题目
	Content    string `json:"content" gorm:"type:text;not null"`                  // 内容
	Reslong    string `json:"res_long" gorm:"type:text"`                          // 备用长文本
	Resshort   string `json:"res_short" gorm:"type:text"`                         // 备用短文本
	AutoUpdate bool   `json:"auto_update" gorm:"type:boolean"`                    // 是否每日自动更新排行
	AutoPass   bool   `json:"auto_pass" gorm:"type:boolean"`                      // 是否自动通过组加入申请
	PassNum    uint   `json:"pass_num" gorm:"type:uint"`                          // 每组的人数限制
	PassRe     bool   `json:"pass_re" gorm:"type:boolean"`                        // 每组人员是否可以重复
}

// @title    BeforDelete
// @description   关于表单删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (s *Set) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("set_id = ?", s.ID)

	// TODO 删除表单收藏
	tx.Delete(&SetCollect{})

	// TODO 删除表单点赞
	tx.Delete(&SetLike{})

	// TODO 删除表单访问
	tx.Delete(&SetVisit{})

	// TODO 删除表单内的用户排行
	tx.Delete(&SetRank{})

	// TODO 删除表单的申请列表
	tx.Delete(&SetApply{})

	// TODO 删除表单的黑名单
	tx.Delete(&SetBlock{})

	return
}
