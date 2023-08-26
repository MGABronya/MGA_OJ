// @Title  Swift
// @Description  该文件提供关于Swift文件的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Handle

import (
	"MGA_OJ/Interface"
	"os/exec"
)

// Swift			定义了Swift文件类
type Swift struct {
	Interface.CmdInterface
}

// @title    Compile
// @description   获得编译指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Swift) Compile(path string, ID string) *exec.Cmd {
	return exec.Command("swiftc", path+ID+".swift", "-o", path+ID)
}

// @title    Chmod
// @description   获得权限
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Swift) Chmod(path string, ID string) *exec.Cmd {
	return exec.Command("chmod", "755", path+ID)
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Swift) Run(path string, ID string) *exec.Cmd {
	return exec.Command(path + ID)
}

// @title    Suffix
// @description   获得文件后缀
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Swift) Suffix() string {
	return "swift"
}

// @title    Name
// @description   获得文件名
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Swift) Name() string {
	return "main"
}

// @title    TimeMultiplier
// @description   运行时间倍率
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Swift) TimeMultiplier() uint {
	return 2
}

// @title    MemoryMultiplier
// @description   运行内存倍率
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Swift) MemoryMultiplier() uint {
	return 2
}

// @title    RunUpTime
// @description   运行启动时间
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Swift) RunUpTime() uint {
	return 4
}

// @title    NewSwift
// @description   新建一个CmdInterface
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   CmdInterface		返回一个CmdInterface用于调用各种函数
func NewSwift() Interface.CmdInterface {
	return Swift{}
}
