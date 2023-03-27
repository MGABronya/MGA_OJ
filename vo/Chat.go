// @Title  chat
// @Description  定义群聊消息
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

// ChatRequest			定义私信
type ChatRequest struct {
	Content  string `json:"content"`   // 内容
	Reslong  string `json:"res_long"`  // 备用长文本
	Resshort string `json:"res_short"` // 备用短文本
}
