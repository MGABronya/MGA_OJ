// @Title  RecordRabbit
// @Description  定义提交
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

import (
	uuid "github.com/satori/go.uuid"
)

// RecordRabbit			定义提交
type RecordRabbit struct {
	RecordId uuid.UUID `json:"record_id"` // 提交外键
	Type     string    `json:"type"`      // 提交类型
}
