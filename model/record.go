// @Title  record
// @Description  定义提交记录
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Record			定义提交记录
type Record struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                          // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                             // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                             // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`       // 用户外键
	ProblemId uuid.UUID `json:"problem_id" gorm:"type:char(36);index:idx_problemId;not null"` // 题目外键
	Language  string    `json:"language" gorm:"type:varchar(64);index:idx_language;not null"` // 语言
	Code      string    `json:"code" gorm:"type:text;not null"`                               // 代码
	Condition string    `json:"condition" gorm:"type:varchar(64);not null"`                   // 记录状态
	Pass      uint      `json:"pass" gorm:"type:uint;not null"`                               // 测试通过数量
	HackId    uuid.UUID `json:"hack_id" gorm:"type:varchar(64);"`                             // hack外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (record *Record) BeforeCreate(scope *gorm.DB) error {
	record.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于提交记录删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (r *Record) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("record_id = (?)", r.ID)

	// TODO 删除提交记录用例相关
	tx.Delete(&CaseCondition{})

	return
}
