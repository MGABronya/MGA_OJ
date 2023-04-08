package main

import (
	TQ "MGA_OJ/Test-request"
)

// 本地测试使用
func main() {
	/*var language, code, input string
	var time_limit, memory_limit uint

	//1、一次性读取文件内容,还有一个 ReadAll的函数，也能读取
	data, err := ioutil.ReadFile("code.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	code = string(data)
	fmt.Print("Language:")
	fmt.Scan(&language)
	fmt.Print("TimeLimit:")
	fmt.Scan(&time_limit)
	fmt.Print("MemoryLimit:")
	fmt.Scan(&memory_limit)
	data, err = ioutil.ReadFile("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	input = string(data)
	TQ.TestRun(language, code, input, memory_limit, time_limit)
	time.Sleep(20 * time.Second)*/
	TQ.Hello()
}
