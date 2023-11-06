// @Title  test
// @Description  该文件提供关于test容器的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Docker

import (
	"MGA_OJ/Interface"
	"os/exec"
)

// Test			定义了Test容器类
type Test struct {
	Interface.GuardInterface
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t Test) Run(ID string, cpu string, postMap string) *exec.Cmd {
	return exec.Command("docker", "run", "-itd", "--cpuset-cpus", cpu, "-p", postMap, "--name", ID, "mgaronya/test:v1", "/bin/bash", "/autorun.sh")
}

// @title    RM
// @description   获得删除指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t Test) RM(ID string) *exec.Cmd {
	return exec.Command("docker", "rm", ID)
}

// @title    Stop
// @description   获得暂停指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t Test) Stop(ID string) *exec.Cmd {
	return exec.Command("docker", "stop", ID)
}

// @title    NewTest
// @description   新建一个GuardInterface
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   GuardInterface		返回一个GuardInterface用于调用各种函数
func NewTest() Interface.GuardInterface {
	return Test{}
}
