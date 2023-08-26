// @Title  problemFile
// @Description  定义了文件题
// @Author  MGAronya
// @Update  MGAronya  2022-10-17 21:07
package vo

// ProblemFileRequest			定义了文件题的各种元素
type ProblemFileRequest struct {
	Description string `json:"description"` // 内容
	ResLong     string `json:"res_long"`    // 备用长文本
	ResShort    string `json:"res_short"`   // 备用短文本
	Score       uint   `json:"score"`       // 该题分数
}
