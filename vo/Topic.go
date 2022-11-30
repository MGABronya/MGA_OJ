// @Title  TopicRequest
// @Description  定义主题
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

// TopicRequest			定义主题
type TopicRequest struct {
	Title    string `json:"title"`     // 题目
	Content  string `json:"content"`   // 内容
	Reslong  string `json:"res_long"`  // 备用长文本
	Resshort string `json:"res_short"` // 备用短文本
	Problems []uint `json:"problems"`  // 主题包含的题目
}
