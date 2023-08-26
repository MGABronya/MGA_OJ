// @Title  Program
// @Description  定义程序
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:46
package vo

// Program			定义程序
type ProgramRequest struct {
	Language string `json:"language"` // 语言
	Code     string `json:"code"`     // 代码
}
