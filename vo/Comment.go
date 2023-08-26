// @Title  comment
// @Description  定义讨论
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

// CommentRequest			定义讨论
type CommentRequest struct {
	Content  string `json:"content"`   // 内容
	ResLong  string `json:"res_long"`  // 备用长文本
	ResShort string `json:"res_short"` // 备用短文本
}
