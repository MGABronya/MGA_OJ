// @Title  Exam
// @Description  定义测试
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Exam		定义测试
type Exam struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                 // id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                                    // 创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                                    // 更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;"`                      // 用户外键
	GroupId   uuid.UUID `json:"group_id" gorm:"type:char(36);index:idx_groupId;"`                    // 用户组外键
	StartTime Time      `json:"start_time" gorm:"type:timestamp;not null"`                           // 起始时间
	EndTime   Time      `json:"end_time" gorm:"type:timestamp;not null"`                             // 终止时间
	Title     string    `json:"title" gorm:"type:char(50);not null;index:search_idx,class:FULLTEXT"` // 标题
	Content   string    `json:"content" gorm:"type:text;not null;index:search_idx,class:FULLTEXT"`   // 内容
	Reslong   string    `json:"res_long" gorm:"type:text;index:search_idx,class:FULLTEXT"`           // 备用长文本
	Resshort  string    `json:"res_short" gorm:"type:text;index:search_idx,class:FULLTEXT"`          // 备用短文本
	Type      string    `json:"type" gorm:"type:char(36);not null"`                                  // 类型
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (exam *Exam) BeforeCreate(scope *gorm.DB) error {
	exam.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于测试删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (e *Exam) BeforDelete(tx *gorm.DB) (err error) {

	tx = tx.Where("exam_id = ?", e.ID)

	// TODO 删除选择题
	tx.Delete(&ProblemMCQs{})

	// TODO 删除文件提交
	tx.Delete(&ProblemFile{})

	// TODO 删除填空题
	tx.Delete(&ProblemCloze{})

	// TODO 删除测试分数
	tx.Delete(&ExamScore{})

	return
}
