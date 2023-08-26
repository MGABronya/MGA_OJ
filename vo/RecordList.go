// @Title  RecorkList
// @Description  提交频道
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

import uuid "github.com/satori/go.uuid"

// RecordList			定义提交频道消息
type RecordList struct {
	RecordId uuid.UUID `json:"record_id"`
}
