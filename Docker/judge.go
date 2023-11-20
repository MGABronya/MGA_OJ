// @Title  judge
// @Description  该文件提供关于Judge容器的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Docker

import (
	"MGA_OJ/Interface"
	"os/exec"
)

// Judge			定义了Judge容器类
type Judge struct {
	Interface.GuardInterface
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Judge) Run(ID string, cpu string, postMap string) *exec.Cmd {
	return exec.Command("docker", "run", "-itd", "--cpuset-cpus", cpu, "--name", ID, "judge", "/bin/bash", "/autorun.sh")
}

// @title    RM
// @description   获得删除指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Judge) RM(ID string) *exec.Cmd {
	return exec.Command("docker", "rm", ID)
}

// @title    Stop
// @description   获得暂停指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Judge) Stop(ID string) *exec.Cmd {
	return exec.Command("docker", "stop", ID)
}

// @title    NewJudge
// @description   新建一个GuardInterface
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   GuardInterface		返回一个GuardInterface用于调用各种函数
func NewJudge() Interface.GuardInterface {
	return Judge{}
}
