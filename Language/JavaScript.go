// @Title  JavaScript
// @Description  该文件提供关于javascript文件的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Handle

import (
	"MGA_OJ/Interface"
	"os/exec"
)

// JavaScript			定义了javascript文件类
type JavaScript struct {
	Interface.CmdInterface
}

// @title    Compile
// @description   获得编译指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j JavaScript) Compile(path string, ID string) *exec.Cmd {
	return exec.Command("clear")
}

// @title    Chmod
// @description   获得权限
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j JavaScript) Chmod(path string, ID string) *exec.Cmd {
	return exec.Command("chmod", "755", path+ID+".js")
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j JavaScript) Run(path string, ID string) *exec.Cmd {
	return exec.Command("node", path+ID+".js")
}

// @title    Suffix
// @description   获得文件后缀
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j JavaScript) Suffix() string {
	return "js"
}

// @title    Name
// @description   获得文件名
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j JavaScript) Name() string {
	return "main"
}

// @title    TimeMultiplier
// @description   运行时间倍率
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j JavaScript) TimeMultiplier() uint {
	return 3
}

// @title    MemoryMultiplier
// @description   运行内存倍率
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j JavaScript) MemoryMultiplier() uint {
	return 3
}

// @title    RunUpTime
// @description   运行启动时间
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j JavaScript) RunUpTime() uint {
	return 64
}

// @title    NewJavaScript
// @description   新建一个CmdInterface
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   CmdInterface		返回一个CmdInterface用于调用各种函数
func NewJavaScript() Interface.CmdInterface {
	return JavaScript{}
}
