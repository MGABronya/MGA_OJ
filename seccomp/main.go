package main

import (
	"MGA_OJ/util"
	"bytes"
	"flag"
	"fmt"
	"io"

	"os/user"
	"strconv"
	"syscall"
)

func main() {
	var language string
	flag.StringVar(&language, "language", "C++", "language of the code")
	var input string
	flag.StringVar(&input, "input", "", "input of the test")
	flag.Parse()
	cmdI := util.LanguageMap[language]
	cmd := cmdI.Run("./user-code/", cmdI.Name())
	// TODO 更换系统用户
	user, err := user.Lookup("mgaoj")
	if err != nil {
		panic(err)
	}

	uid, _ := strconv.Atoi(user.Uid)
	gid, _ := strconv.Atoi(user.Gid)

	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}

	// TODO 用于记录代码的输入输出
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out

	stdinPipe, err := cmd.StdinPipe()
	// TODO 系统错误
	if err != nil {
		panic(err)
	}
	io.WriteString(stdinPipe, input+"\n")
	// TODO 关闭管道制造EOF信息
	stdinPipe.Close()

	if err := Seccomp(); err != nil {
		panic(err)
	}
	// TODO 运行代码
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	fmt.Print(out.String())
}
