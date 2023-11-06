// @Title  GuardInterface
// @Description  该文件用于封装命令执行方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Interface

import (
	"os/exec"
)

// GuardInterface			定义了哨兵执行命令方法
type GuardInterface interface {
	Stop(string) *exec.Cmd                            // 停止
	RM(string) *exec.Cmd                              // 删除
	Run(id string, cpu string, post string) *exec.Cmd // 运行
}
