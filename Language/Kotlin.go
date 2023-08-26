// @Title  kotlin
// @Description  该文件提供关于kotlin文件的各种方法
// @Author  MGAronya
// @Update  MGAronya  2022-9-16 0:33
package Handle

import (
	"MGA_OJ/Interface"
	"os/exec"
)

// kotlin			定义了kotlin文件类
type Kotlin struct {
	Interface.CmdInterface
}

// @title    Compile
// @description   获得编译指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (k Kotlin) Compile(path string, ID string) *exec.Cmd {
	return exec.Command("kotlinc", path+ID+".kt", "-include-runtime", "-d", path+ID+".jar")
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (k Kotlin) Run(path string, ID string) *exec.Cmd {
	return exec.Command("java", "-jar", path+ID+".jar")
}

// @title    Chmod
// @description   获得权限
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (k Kotlin) Chmod(path string, ID string) *exec.Cmd {
	return exec.Command("chmod", "755", path+ID+".jar")
}

// @title    Suffix
// @description   获得文件后缀
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (k Kotlin) Suffix() string {
	return "kt"
}

// @title    Name
// @description   获得文件名
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (k Kotlin) Name() string {
	return "main"
}

// @title    TimeMultiplier
// @description   运行时间倍率
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (k Kotlin) TimeMultiplier() uint {
	return 2
}

// @title    MemoryMultiplier
// @description   运行内存倍率
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (k Kotlin) MemoryMultiplier() uint {
	return 2
}

// @title    RunUpTime
// @description   运行启动时间
// @auth      MGAronya       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (k Kotlin) RunUpTime() uint {
	return 56
}

// @title    Newkotlin
// @description   新建一个CmdInterface
// @auth      MGAronya       2022-9-16 12:23
// @param    void
// @return   CmdInterface		返回一个CmdInterface用于调用各种函数
func NewKotlin() Interface.CmdInterface {
	return Kotlin{}
}
