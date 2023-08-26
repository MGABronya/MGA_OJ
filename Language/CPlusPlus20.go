// @Title  CPlusPlus20
// @Description  该文件提供关于c++20文件的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Handle

import (
	"MGA_OJ/Interface"
	"os/exec"
)

// CppPlusPlus20			定义了c++文件类
type CppPlusPlus20 struct {
	Interface.CmdInterface
}

// @title    Compile
// @description   获得编译指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus20) Compile(path string, ID string) *exec.Cmd {
	return exec.Command("g++", "-std=c++20", path+ID+".cpp", "-o", path+ID)
}

// @title    Chmod
// @description   获得权限
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus20) Chmod(path string, ID string) *exec.Cmd {
	return exec.Command("chmod", "755", path+ID)
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus20) Run(path string, ID string) *exec.Cmd {
	return exec.Command(path + ID)
}

// @title    Suffix
// @description   获得文件后缀
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus20) Suffix() string {
	return "cpp"
}

// @title    TimeMultiplier
// @description   运行时间倍率
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus20) TimeMultiplier() uint {
	return 1
}

// @title    MemoryMultiplier
// @description   运行内存倍率
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus20) MemoryMultiplier() uint {
	return 1
}

// @title    RunUpTime
// @description   运行启动时间
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus20) RunUpTime() uint {
	return 4
}

// @title    Name
// @description   获得文件名
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus20) Name() string {
	return "main"
}

// @title    NewCppPlusPlus11
// @description   新建一个CmdInterface
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   CmdInterface		返回一个CmdInterface用于调用各种函数
func NewCppPlusPlus20() Interface.CmdInterface {
	return CppPlusPlus20{}
}
