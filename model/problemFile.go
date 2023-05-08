// @Title  problemFile
// @Description  定义了文件题
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// problemFile			定义了文件题的各种元素
type ProblemFile struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                   // id
	CreatedAt   Time      `json:"created_at" gorm:"type:timestamp"`                                      // 创建日期
	UpdatedAt   Time      `json:"updated_at" gorm:"type:timestamp"`                                      // 更新日期
	UserId      uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`                // 用户外键
	ExamId      uuid.UUID `json:"exam_id" gorm:"type:char(36);index:idx_weId;not null"`                  // 测试外键
	Description string    `json:"description" gorm:"type:text;not null;index:search_idx,class:FULLTEXT"` // 内容
	ResLong     string    `json:"res_long" gorm:"type:text;index:search_idx,class:FULLTEXT"`             // 备用长文本
	ResShort    string    `json:"res_short" gorm:"type:text;index:search_idx,class:FULLTEXT"`            // 备用短文本
	Score       uint      `json:"score" gorm:"type:uint;"`                                               // 该题分数
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (problemFile *ProblemFile) BeforeCreate(scope *gorm.DB) error {
	problemFile.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于题目删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (p *ProblemFile) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("problem_file_id = ?", p.ID)

	// TODO 删除题目提交
	tx.Delete(&ProblemFileSubmit{})

	return
}
