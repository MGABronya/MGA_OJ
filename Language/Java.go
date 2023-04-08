// @Title  Java
// @Description  该文件提供关于Java文件的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Handle

import (
	"MGA_OJ/Interface"
	"os/exec"
)

// Java			定义了Java文件类
type Java struct{}

// @title    Compile
// @description   获得编译指令
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) Compile(path string, ID string) *exec.Cmd {
	return exec.Command("javac", "-encoding", "utf-8", path+ID+".java")
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) Run(path string, ID string) *exec.Cmd {
	return exec.Command("java", "-Dfile.encoding=utf-8", "-cp", path, ID)
}

// @title    Suffix
// @description   获得文件后缀
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) Suffix() string {
	return "java"
}

// @title    Name
// @description   获得文件名
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) Name() string {
	return "Main"
}

// @title    TimeMultiplier
// @description   运行时间倍率
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) TimeMultiplier() uint {
	return 3
}

// @title    MemoryMultiplier
// @description   运行内存倍率
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) MemoryMultiplier() uint {
	return 3
}

// @title    RunUpTime
// @description   运行启动时间
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) RunUpTime() uint {
	return 52
}

// @title    NewJava
// @description   新建一个CmdInterface
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   CmdInterface		返回一个CmdInterface用于调用各种函数
func NewJava() Interface.CmdInterface {
	return Java{}
}
