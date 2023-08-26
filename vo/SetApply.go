// @Title  SetApply
// @Description  表单申请
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

import (
	uuid "github.com/satori/go.uuid"
)

// SetApplyRequest			定义表单申请
type SetApplyRequest struct {
	Content  string    `json:"content"`   // 内容
	ResLong  string    `json:"res_long"`  // 备用长文本
	ResShort string    `json:"res_short"` // 备用短文本
	GroupId  uuid.UUID `json:"group_id"`  // 用户组id
}
