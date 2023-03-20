// @Title  Python
// @Description  该文件提供关于python文件的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Handle

import (
	"MGA_OJ/Interface"
	"os/exec"
)

// Python			定义了python文件类
type Python struct{}

// @title    Compile
// @description   获得编译指令
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p Python) Compile(path string, ID string) *exec.Cmd {
	return exec.Command("g++", path+ID+".cpp", "-o", path+ID)
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p Python) Run(path string, ID string) *exec.Cmd {
	return exec.Command("./" + path + ID)
}

// @title    Suffix
// @description   获得文件后缀
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p Python) Suffix() string {
	return "cpp"
}

// @title    Name
// @description   获得文件名
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p Python) Name() string {
	return "main"
}

// @title    TimeMultiplier
// @description   运行时间倍率
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p Python) TimeMultiplier() uint {
	return 1
}

// @title    RunUpTime
// @description   运行启动时间
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p Python) RunUpTime() uint {
	return 1
}

// @title    NewCppPlusPlus
// @description   新建一个ICppPlusPlus
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ICppPlusPlus		返回一个ICppPlusPlus用于调用各种函数
func NewPython() Interface.CmdInterface {
	return Python{}
}
