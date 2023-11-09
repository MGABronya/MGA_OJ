// @Title  chat
// @Description  定义群聊消息
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

import (
	"MGA_OJ/model"
)

// ChatRequest			定义私信
type ChatRequest struct {
	Content  string `json:"content"`   // 内容
	ResLong  string `json:"res_long"`  // 备用长文本
	ResShort string `json:"res_short"` // 备用短文本
}

// ChatWithUser 定义群聊与发出用户
type ChatWithUser struct {
	Chat model.Chat `json:"chat"` // 群聊消息
	User UserDto    `json:"user"` // 发出用户
}
