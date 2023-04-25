// @Title  examScore
// @Description  测试分数
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// ExamScore			定义测试分数
type ExamScore struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                    // 测试的id
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                       // 测试的创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                       // 测试的更新日期
	UserId    uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"` // 用户外键
	ExamId    uuid.UUID `json:"exam_id" gorm:"type:char(36);index:idx_examId;not null"` // 测试外键
	Score     uint      `json:"score" gorm:"type:uint;"`                                // 测试得分
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (exam *ExamScore) BeforeCreate(scope *gorm.DB) error {
	exam.ID = uuid.NewV4()
	return nil
}
