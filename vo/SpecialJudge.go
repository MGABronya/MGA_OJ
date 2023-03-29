// @Title  SpecailJudge
// @Description  定义特判
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package vo

// SpecialJudge			定义特判
type SpecialJudgeRequest struct {
	Language string `json:"language"` // 语言
	Code     string `json:"code"`     // 代码
	Input    string `json:"input"`    // 测试输入
}
