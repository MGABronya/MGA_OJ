// @Title  problem
// @Description  定义了题目
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Problem			定义了题目的各种元素
type Problem struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                   // id
	CreatedAt    Time      `json:"created_at" gorm:"type:timestamp"`                                      // 创建日期
	UpdatedAt    Time      `json:"updated_at" gorm:"type:timestamp"`                                      // 更新日期
	UserId       uuid.UUID `json:"user_id" gorm:"type:char(36);index:idx_userId;not null"`                // 用户外键
	TimeLimit    uint      `json:"time_limit" gorm:"type:uint;not null"`                                  // 时间限制
	MemoryLimit  uint      `json:"memory_limit" gorm:"type:uint;not null"`                                // 内存限制
	Title        string    `json:"title" gorm:"type:char(50);not null;index:search_idx,class:FULLTEXT"`   // 标题
	Description  string    `json:"description" gorm:"type:text;not null;index:search_idx,class:FULLTEXT"` // 内容
	ResLong      string    `json:"res_long" gorm:"type:text;index:search_idx,class:FULLTEXT"`             // 备用长文本
	ResShort     string    `json:"res_short" gorm:"type:text;index:search_idx,class:FULLTEXT"`            // 备用短文本
	Input        string    `json:"input" gorm:"type:text;not null"`                                       // 输入格式
	Output       string    `json:"output" gorm:"type:text;not null"`                                      // 输出格式
	Hint         string    `json:"hint" gorm:"type:text"`                                                 // 提示
	Source       string    `json:"source" gorm:"type:varchar(64)"`                                        // 来源
	SpecialJudge uuid.UUID `json:"specail_judge_id" gorm:"type:char(36);"`                                // 特判程序外键
	Standard     uuid.UUID `json:"standard_id" gorm:"type:char(36);"`                                     // 标准程序外键
	InputCheck   uuid.UUID `json:"input_check_id" gorm:"type:char(36);"`                                  // 输入检查程序外键
}

// @title    BeforeCreate
// @description   计算出一个uuid
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (problem *Problem) BeforeCreate(scope *gorm.DB) error {
	problem.ID = uuid.NewV4()
	return nil
}

// @title    BeforDelete
// @description   关于题目删除的一些级联操作
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    tx *gorm.DB       接收一个数据库指针
// @return   err error		   返回一个错误信息
func (p *Problem) BeforDelete(tx *gorm.DB) (err error) {
	tx = tx.Where("problem_id = ?", p.ID)

	// TODO 删除题目收藏
	tx.Delete(&ProblemCollect{})

	// TODO 删除题目点赞
	tx.Delete(&ProblemLike{})

	// TODO 删除题目访问
	tx.Delete(&ProblemVisit{})

	// TODO 删除题解
	tx.Delete(&Post{})

	// TODO 删除回复
	tx.Delete(&Comment{})

	// TODO 删除提交记录
	tx.Delete(&Record{})

	// TODO 删除主题包含
	tx.Delete(&ProblemList{})

	// TODO 删除题目标签
	tx.Delete(&ProblemLabel{})

	return
}
