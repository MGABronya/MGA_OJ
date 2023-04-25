// @Title  ExamRequest
// @Description  定义测试
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package vo

import (
	"MGA_OJ/model"

	uuid "github.com/satori/go.uuid"
)

// ExamRequest		定义测试
type ExamRequest struct {
	StartTime model.Time `json:"start_time"` // 起始时间
	EndTime   model.Time `json:"end_time"`   // 终止时间
	Title     string     `json:"title"`      // 标题
	Content   string     `json:"content"`    // 内容
	Reslong   string     `json:"res_long"`   // 备用长文本
	Resshort  string     `json:"res_short"`  // 备用短文本
	GroupId   uuid.UUID  `json:"group_id"`   // 用户组
	Type      string     `json:"type"`       // 测试类型
}
