// @Title  record
// @Description  定义提交记录
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"gorm.io/gorm"
)

// Record			定义提交记录
type Record struct {
	gorm.Model
	UserId        uint   `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`               // 用户外键
	ProblemId     uint   `json:"problem_id" gorm:"type:uint;index:idx_problemId;not null"`         // 题目外键
	Language      string `json:"language" gorm:"type:varchar(64);index:idx_language;not null"`     // 语言
	Code          string `json:"code" gorm:"type:text;not null"`                                   // 代码
	Condition     string `json:"condition" gorm:"type:varchar(64);not null"`                       // 记录状态
	CompetitionId uint   `json:"competition_id" gorm:"type:uint;index:idx_competitionId;not null"` // 比赛外键
	Pass          uint   `json:"pass" gorm:"type:uint;not null"`                                   // 测试通过数量
}

// @title    BeforDelete
// @description   关于提交记录删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (r *Record) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("record_id = ?", r.ID)

	// TODO 删除提交记录用例相关
	tx.Delete(&Case{})

	return
}
