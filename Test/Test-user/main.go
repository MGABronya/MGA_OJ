package main

import (
	"MGA_OJ/util"
	"bytes"
	"fmt"
)

func main() {
	for k := range util.LanguageMap {
		testforhello(k)
	}
}

func testforhello(language string) {
	// TODO 找到提交记录后，开始判题逻辑
	cmdI := util.LanguageMap[language]
	// id		定义文件名
	id := cmdI.Name()
	// TODO 运行可执行文件
	cmd := cmdI.Run("./user-code/", id)
	var out, stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	var err error

	// TODO 系统错误
	if err = cmd.Start(); err != nil {
		fmt.Println(language + ":" + "System Error")
		return
	}
	// TODO 启动routine等待结束
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	select {
	// TODO 运行超时
	case err = <-done:
	}

	// TODO 运行时错误
	if err != nil {
		fmt.Println(language + ":" + "Runtime Error")
		return
	}
	fmt.Println(language + ":" + out.String())
}
