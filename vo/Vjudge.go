// @Title  vjudge
// @Description  外站题目的基本信息
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package vo

// Vjudge		通过前端发送请求接收的题目信息，无标签限制
type Vjudge struct {
	OJ          string `json:"oj"`           // 来自的平台
	ProblemId   string `json:"problem_id"`   // 题目的id
	Title       string `json:"title"`        // 题目标题
	TimeLimit   uint   `json:"time_limit"`   // 时间限制
	TimeUnits   string `json:"time_unit"`    // 时间单位
	MemoryLimit uint   `json:"memory_limit"` // 内存限制
	MemoryUnits string `json:"memory_unit"`  // 内存单位
	Description string `json:"description"`  // 内容描述
	ResLong     string `json:"res_long"`     // 备用长文本
	ResShort    string `json:"res_short"`    // 备用短文本
	Input       string `json:"input"`        // 输入格式
	Output      string `json:"output"`       // 输出格式
	SampleCase  []Case `json:"sample_case"`  // 样例
	Hint        string `json:"hint"`         // 提示
	Source      string `json:"source"`       // 来源
}
