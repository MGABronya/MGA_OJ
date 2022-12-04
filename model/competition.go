// @Title  Competition
// @Description  定义比赛
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	"gorm.io/gorm"
)

// Competition		定义比赛
type Competition struct {
	gorm.Model
	UserId    uint   `json:"user_id" gorm:"type:uint;index:idx_userId;not null"` // 用户外键
	SetId     uint   `json:"set_id" gorm:"type:uint;index:idx_setId;not null"`   // 表单外键
	StartTime Time   `json:"start_time" gorm:"type:timestamp;not null"`          // 起始时间
	EndTime   Time   `json:"end_time" gorm:"type:timestamp;not null"`            // 终止时间
	Type      string `json:"type" gorm:"type:varchar(20);not null"`              // 比赛类型
}

// @title    BeforDelete
// @description   关于比赛删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (c *Competition) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("competition_id = ?", c.ID)

	// TODO 删除比赛成员
	tx.Delete(&CompetitionMember{})

	// TODO 删除比赛排名
	tx.Delete(&CompetitionRank{})

	return
}
