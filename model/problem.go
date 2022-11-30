// @Title  problem
// @Description  定义了题目
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package model

import (
	"gorm.io/gorm"
)

// Problem			定义了题目的各种元素
type Problem struct {
	gorm.Model
	UserId        uint   `json:"user_id" gorm:"type:uint;index:idx_userId;not null"`               // 用户外键
	Title         string `json:"title" gorm:"type:varchar(64);not null"`                           // 题目标题
	TimeLimit     uint   `json:"time_limit" gorm:"type:uint;not null"`                             // 时间限制
	MemoryLimit   uint   `json:"memory_limit" gorm:"type:uint;not null"`                           // 内存限制
	Description   string `json:"description" gorm:"type:text;not null"`                            // 内容描述
	Reslong       string `json:"res_long" gorm:"type:text"`                                        // 备用长文本
	Resshort      string `json:"res_short" gorm:"type:text"`                                       // 备用短文本
	Input         string `json:"input" gorm:"type:text;not null"`                                  // 输入格式
	Output        string `json:"output" gorm:"type:text;not null"`                                 // 输出格式
	SampleInput   string `json:"sample_input" gorm:"type:text;not null"`                           // 输入样例
	SampleOutput  string `json:"sample_output" gorm:"type:text;not null"`                          // 输出样例
	Hint          string `json:"hint" gorm:"type:text"`                                            // 提示
	Source        string `json:"source" gorm:"type:varchar(64)"`                                   // 来源
	CompetitionId uint   `json:"competition_id" gorm:"type:uint;index:idx_competitionId;not null"` // 比赛外键
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

	return
}
