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
func (j Java) Compile(ID string) *exec.Cmd {
	return exec.Command("javac", ID+".java", "-o", ID)
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) Run(ID string) *exec.Cmd {
	return exec.Command("java", "./", ID)
}

// @title    Suffix
// @description   获得文件后缀
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) Suffix() string {
	return "java"
}

// @title    TimeMultiplier
// @description   运行时间倍率
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) TimeMultiplier() uint {
	return 2
}

// @title    RunUpTime
// @description   运行启动时间
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (j Java) RunUpTime() uint {
	return 0
}

// @title    NewJava
// @description   新建一个Java
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   Java		返回一个Java用于调用各种函数
func NewJava() Interface.CmdInterface {
	return Java{}
}
