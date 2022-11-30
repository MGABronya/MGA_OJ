// @Title  CPlusPlus11
// @Description  该文件提供关于c++11文件的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package Handle

import (
	"MGA_OJ/Interface"
	"os/exec"
)

// CppPlusPlus11			定义了c++文件类
type CppPlusPlus11 struct{}

// @title    Compile
// @description   获得编译指令
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus11) Compile(ID string) *exec.Cmd {
	return exec.Command("g++", "-std=c++11", ID+".cpp", "-o", ID)
}

// @title    Run
// @description   获得运行指令
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus11) Run(ID string) *exec.Cmd {
	return exec.Command("./" + ID)
}

// @title    Suffix
// @description   获得文件后缀
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus11) Suffix() string {
	return "cpp"
}

// @title    TimeMultiplier
// @description   运行时间倍率
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus11) TimeMultiplier() uint {
	return 1
}

// @title    RunUpTime
// @description   运行启动时间
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (c CppPlusPlus11) RunUpTime() uint {
	return 0
}

// @title    NewCppPlusPlus11
// @description   新建一个ICppPlusPlus11
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   ICppPlusPlus		返回一个ICppPlusPlus用于调用各种函数
func NewCppPlusPlus11() Interface.CmdInterface {
	return CppPlusPlus11{}
}
