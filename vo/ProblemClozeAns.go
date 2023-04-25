// @Title  problemClozeAns
// @Description  定义了填空题答案
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-17 21:07
package vo

// ProblemClozeAnsRequest			定义了填空题答案的各种元素
type ProblemClozeAnsRequest struct {
	Answer string `json:"answer"` // 该题答案
}
