// @Title  CmdInterface
// @Description  该文件用于封装命令执行方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Interface

import "os/exec"

// CmdInterface			定义了命令执行方法
type CmdInterface interface {
	Compile(string) *exec.Cmd // 编译
	Suffix() string           // 后缀
	Run(string) *exec.Cmd     // 运行
	TimeMultiplier() uint    // 语言的运行时间倍
	RunUpTime() uint         // 语言运行的启动所需时间
}
