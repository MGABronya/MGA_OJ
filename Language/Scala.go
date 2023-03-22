// @Title  Scala
// @Description  该文件提供关于Scala文件的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Handle

import (
	"MGA_OJ/Interface"
	"os/exec"
)

// Scala			定义了Scala文件类
type Scala struct{}

// @title    Compile
// @description   获得编译指令
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Scala) Compile(path string, ID string) *exec.Cmd {
	return exec.Command("scalac", path+ID+".scala", "-d", path)
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Scala) Run(path string, ID string) *exec.Cmd {
	return exec.Command("scala", "-classpath", "./"+path, ID)
}

// @title    Suffix
// @description   获得文件后缀
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Scala) Suffix() string {
	return "scala"
}

// @title    Name
// @description   获得文件名
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Scala) Name() string {
	return "main"
}

// @title    TimeMultiplier
// @description   运行时间倍率
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Scala) TimeMultiplier() uint {
	return 1
}

// @title    RunUpTime
// @description   运行启动时间
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (s Scala) RunUpTime() uint {
	return 1
}

// @title    NewSwift
// @description   新建一个CmdInterface
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   CmdInterface		返回一个CmdInterface用于调用各种函数
func NewScala() Interface.CmdInterface {
	return Scala{}
}
