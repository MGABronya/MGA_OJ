// @Title  Group
// @Description  定义用户组
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-11-16 0:46
package vo

import (
	uuid "github.com/satori/go.uuid"
)

// GroupRequest			定义用户组
type GroupRequest struct {
	Title    string      `json:"title"`     // 主题
	Content  string      `json:"content"`   // 内容
	Auto     bool        `json:"auto"`      // 是否自动通过申请
	ResLong  string      `json:"res_long"`  // 备用长文本
	ResShort string      `json:"res_short"` // 备用短文本
	Users    []uuid.UUID `json:"users"`     // 用户组用户列表
}
