// @Title  Test
// @Description  定义题解的回复
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

// TestRequest			定义题解的回复
type TestRequest struct {
	Language string `json:"language"` // 语言
	Code     string `json:"code"`     // 代码
	Input    string `json:"input"`    // 输入
}
