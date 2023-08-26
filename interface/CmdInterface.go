// @Title  CmdInterface
// @Description  该文件用于封装命令执行方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import (
	"os/exec"
)

// CmdInterface			定义了命 令执行方法
type CmdInterface interface {
	Compile(string, string) *exec.Cmd // 编译
	Chmod(string, string) *exec.Cmd   // 获取运行权限
	Suffix() string                   // 后缀
	Run(string, string) *exec.Cmd     // 运行
	TimeMultiplier() uint             // 语言的运行时间倍
	MemoryMultiplier() uint           // 语言的运行内存倍
	RunUpTime() uint                  // 语言运行的启动所需时间
	Name() string                     // 语言文件名称
}
