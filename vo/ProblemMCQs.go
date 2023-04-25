// @Title  problemMCQs
// @Description  定义了选择题
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package vo

// ProblemMCQsRequest			定义了选择题的各种元素
type ProblemMCQsRequest struct {
	Description string `json:"description"` // 内容
	Reslong     string `json:"res_long"`    // 备用长文本
	Resshort    string `json:"res_short"`   // 备用短文本
	Answer      string `json:"answer"`      // 该题答案
	Score       uint   `json:"score"`       // 该题分数
}
