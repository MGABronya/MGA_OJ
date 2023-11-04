// @Title  problemCode
// @Description  题目的基本信息+代码
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:47
package vo

// ProblemCodeRequest		题目信息+代码，无标签限制
type ProblemCodeRequest struct {
	Problem ProblemRequest
	Code    string `json:"code"`
}
