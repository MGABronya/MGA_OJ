// @Title  problemCloze
// @Description  定义了填空题
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package vo

// ProblemClozeRequest			定义了填空题的各种元素
type ProblemClozeRequest struct {
	Description string `json:"description"` // 内容
	Reslong     string `json:"res_long"`    // 备用长文本
	Resshort    string `json:"res_short"`   // 备用短文本
	Answer      string `json:"answer"`      // 该题答案
	Score       uint   `json:"score"`       // 该题分数
}
