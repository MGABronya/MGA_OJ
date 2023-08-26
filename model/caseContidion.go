// @Title  caseCondition
// @Description  定义提交记录中的样例通过情况
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package model

import uuid "github.com/satori/go.uuid"

// CaseCondition			定义提交记录中的样例
type CaseCondition struct {
	RecordId uuid.UUID `json:"record_id" gorm:"type:char(36);index:idx_recordId;not null"` // 提交记录外键
	CID      uint      `json:"cid" gorm:"type:uint;not null"`                              // 表示第几个测试
	Time     uint      `json:"time" gorm:"type:uint;not null"`                             // 表示测试使用时间
	Memory   uint      `json:"memory" gorm:"type:uint;not null"`                           // 表示测试使用空间
	Input    string    `json:"input" gorm:"type:text;"`                                    // 输入
	Output   string    `json:"output" gorm:"type:text;"`                                   // 输出
}
