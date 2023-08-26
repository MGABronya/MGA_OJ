// @Title  case
// @Description  题目的用例
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:47
package vo

// Case		通过前端发送请求接收的用例信息
type Case struct {
	Input  string `json:"input"`  // 输入
	Output string `json:"output"` // 输出
}
