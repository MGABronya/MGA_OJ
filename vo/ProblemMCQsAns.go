// @Title  problemMCQsAns
// @Description  定义了选择题答案
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package vo

// ProblemMCQsAnsRequest			定义了选择题答案的各种元素
type ProblemMCQsAnsRequest struct {
	Answer      string `json:"answer"`      // 该题答案
}
