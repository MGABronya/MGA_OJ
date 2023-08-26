// @Title  Remark
// @Description  定义文章的回复
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

// RemarkRequest			定义文章的回复
type RemarkRequest struct {
	Content  string `json:"content"`   // 内容
	ResLong  string `json:"res_long"`  // 备用长文本
	ResShort string `json:"res_short"` // 备用短文本
}
