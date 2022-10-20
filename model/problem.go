// @Title  problem
// @Description  定义了问题
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package model

import (
	"gorm.io/gorm"
)

// Problem			定义了问题的各种元素
type Problem struct {
	gorm.Model
	UserId       uint   `json:"user_id" gorm:"type:uint;index:idx_userId;not null"` // 用户外键
	Title        string `json:"title" gorm:"type:varchar(64);not null"`             // 问题标题
	TimeLimit    uint   `json:"time_limit" gorm:"type:uint;not null"`               // 时间限制
	MemoryLimit  uint   `json:"memory_limit" gorm:"type:uint;not null"`             // 内存限制
	Description  string `json:"description" gorm:"type:text;not null"`              // 内容描述
	Reslong      string `json:"res_long" gorm:"type:text"`                          // 备用长文本
	Resshort     string `json:"res_short" gorm:"type:text"`                         // 备用短文本
	Input        string `json:"input" gorm:"type:text;not null"`                    // 输入格式
	Output       string `json:"output" gorm:"type:text;not null"`                   // 输出格式
	SampleInput  string `json:"sample_input" gorm:"type:text;not null"`             // 输入样例
	SampleOutput string `json:"sample_output" gorm:"type:text;not null"`            // 输出样例
	Hint         string `json:"hint" gorm:"type:text"`                              // 提示
	Source       string `json:"source" gorm:"type:varchar(64)"`                     // 来源
}
