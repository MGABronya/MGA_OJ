// @Title  TopicRequest
// @Description  定义主题
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

import (
	uuid "github.com/satori/go.uuid"
)

// TopicRequest			定义主题
type TopicRequest struct {
	Title    string      `json:"title"`     // 题目
	Content  string      `json:"content"`   // 内容
	ResLong  string      `json:"res_long"`  // 备用长文本
	ResShort string      `json:"res_short"` // 备用短文本
	Problems []uuid.UUID `json:"problems"`  // 主题包含的题目
}
