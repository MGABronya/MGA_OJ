// @Title  Texts
// @Description  接收一组文本
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:47
package vo

// TextsRequest			接收一组文本
type TextsRequest struct {
	Texts []string `form:"texts"` // 文本
}
