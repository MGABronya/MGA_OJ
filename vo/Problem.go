// @Title  problem
// @Description  题目的基本信息
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package vo

import (
	uuid "github.com/satori/go.uuid"
)

// ProblemRequest		通过前端发送请求接收的题目信息，无标签限制
type ProblemRequest struct {
	Title         string    `json:"title"`         // 题目标题
	TimeLimit     uint      `json:"time_limit"`    // 时间限制
	TimeUnits     string    `json:"time_unit"`     // 时间单位
	MemoryLimit   uint      `json:"memory_limit"`  // 内存限制
	MemoryUnits   string    `json:"memory_unit"`   // 内存单位
	Description   string    `json:"description"`   // 内容描述
	Reslong       string    `json:"res_long"`      // 备用长文本
	Resshort      string    `json:"res_short"`     // 备用短文本
	Input         string    `json:"input"`         // 输入格式
	Output        string    `json:"output"`        // 输出格式
	SampleInput   []string  `json:"sample_input"`  // 输入样例
	SampleOutput  []string  `json:"sample_output"` // 输出样例
	TestInput     []string  `json:"test_input"`    // 输入测试
	TestOutput    []string  `json:"test_output"`   // 输出测试
	Hint          string    `json:"hint"`          // 提示
	Source        string    `json:"source"`        // 来源
	SpecialJudge  uuid.UUID `json:"special_judge"` // 特判
	Standard      uuid.UUID `json:"standard"`      // 标准程序外键
	InputCheck    uuid.UUID `json:"input_check"`   // 输入检查程序外键
	Scores        []uint    `json:"scores"`        // 分值
	CompetitionId uuid.UUID `json:"competitionId"` // 比赛外键
}
