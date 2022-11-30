// @Title  article
// @Description  定义文章
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

// ArticleRequest			定义文章
type ArticleRequest struct {
	Title    string `json:"title"`     // 标题
	Content  string `json:"content"`   // 内容
	Reslong  string `json:"res_long"`  // 备用长文本
	Resshort string `json:"res_short"` // 备用短文本
}
