package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var url = "http://43.139.97.161:1003/test/create"

//"http://test_oj.mgaronya.com/test/create"

func main() {
	// 执行17 * 100次hello world
	/*start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			Hello()
			wg.Done()
		}()
	}
	wg.Wait()
	//2m10.4980861s
	fmt.Print(time.Since(start))*/
	code := `#include <stdio.h>
	int main(){
		int n;
		while(scanf("%d",&n) != EOF);
		return 0;
	}
	`
	TestRun("C", code, "\n", 1000, 1000)
}

// TODO 进行hello world测试
func Hello() {
	var language, code, input string
	var time_limit, memory_limit uint
	input = ""
	time_limit = 500
	memory_limit = 512

	language = "C"
	code = `#include <stdio.h>
	int main()
	{
	   printf("Hello, World!");
	   return 0;
	}
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "C#"
	code = `using System;
	using System.Collections.Generic;
	using System.Linq;
	using System.Text;
	using System.Threading.Tasks;
	
	namespace HelloWorld//命名空间
	{
		class main//类
		{
			static void Main(string[] args)//main函数
			{
				string str = "Hello World !";
				Console.WriteLine(str);
			}
		}
	}
	
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "C++"
	code = `#include <iostream>
	using namespace std;
	 
	int main() 
	{
		cout << "Hello, World!";
		return 0;
	}
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "C++11"
	code = `#include <iostream>
	using namespace std;
	 
	int main() 
	{
		cout << "Hello, World!";
		return 0;
	}
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Erlang"
	code = `-module(main).
	-export([start/0]).
	start() -> io:fwrite("hi,hello world!\n").
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Go"
	code = `package main

	import "fmt"
	
	func main() {
		fmt.Println("Hello world") 
	}
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Java"
	code = `public class Main { 
		public static void main(String[] args){ 
			System.out.println("hello world!");
		}
	}
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "JavaScript"
	code = ` console.log("Hello World");
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Kotlin"
	code = `fun main(args: Array<String>) {
		println("Hello world")
	}
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Pascal"
	code = `begin
    writeln('Hello, World!');
end.
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "PHP"
	code = `<?php
	echo "Hello World \n";
	?>
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Python"
	code = `print('Hello World!')
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Racket"
	code = `#lang racket
	"Hello, world!"
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Ruby"
	code = `
	puts "Hello World!"
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Rust"
	code = `fn main() {
		println!("Hello, world!");
	}
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Scala"
	code = `object main {
		def main(args: Array[String]): Unit = {
		  println("Hello, world!")
		}
	  }
	`
	TestRun(language, code, input, memory_limit, time_limit)

	language = "Swift"
	code = `
	print("hello world")
	`
	TestRun(language, code, input, memory_limit, time_limit)
}

func TestRun(language string, code string, input string, memorylimit uint, timelimit uint) {
	data := make(map[string]interface{})
	var body []byte
	data["language"] = language
	data["code"] = code
	data["input"] = input
	data["memory_limit"] = memorylimit
	data["time_limit"] = timelimit
	bytesData, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bytesData))
	if err != nil {
		fmt.Println(language + " has problem")
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(language + " has problem")
		return
	}
	fmt.Println(language + ":" + string(body))
}

/*
初次，运行结果
C:{"code":200,"data":{"condition":"ok","memory":10,"output":"Hello, World!","time":1},"msg":"测试创建成功"}
C#:{"code":200,"data":{"condition":"ok","memory":13,"output":"Hello World !\n","time":17},"msg":"测试创建成功"}
C++:{"code":200,"data":{"condition":"ok","memory":9,"output":"Hello, World!","time":1},"msg":"测试创建成功"}
C++11:{"code":200,"data":{"condition":"ok","memory":9,"output":"Hello, World!","time":2},"msg":"测试创建成功"}
Erlang:{"code":200,"data":{"condition":"ok","memory":12,"output":"hi,hello world!\n","time":1168},"msg":"测试创建成功"}
Go:{"code":200,"data":{"condition":"ok","memory":10,"output":"Hello world\n","time":1},"msg":"测试创建成功"}
Java:{"code":200,"data":{"condition":"ok","memory":12,"output":"hello world!\n","time":52},"msg":"测试创建成功"}
JavaScript:{"code":200,"data":{"condition":"ok","memory":12,"output":"Hello World\n","time":74},"msg":"测试创建成功"}
Kotlin:{"code":200,"data":{"condition":"ok","memory":12,"output":"Hello world\n","time":56},"msg":"测试创建成功"}
Pascal:{"code":200,"data":{"condition":"ok","memory":9,"output":"Hello, World!\n","time":0},"msg":"测试创建成功"}
PHP:{"code":200,"data":{"condition":"ok","memory":12,"output":"Hello World \n\t","time":8},"msg":"测试创建成功"}
Python:{"code":200,"data":{"condition":"ok","memory":12,"output":"Hello World!\n","time":7},"msg":"测试创建成功"}
Racket:{"code":200,"data":{"condition":"ok","memory":12,"output":"\"Hello, world!\"\n","time":365},"msg":"测试创建成功"}
Ruby:{"code":200,"data":{"condition":"ok","memory":12,"output":"Hello World!\n","time":53},"msg":"测试创建成功"}
Rust:{"code":200,"data":{"condition":"ok","memory":9,"output":"Hello, world!\n","time":1},"msg":"测试创建成功"}
Scala:{"code":200,"data":{"condition":"ok","memory":13,"output":"Hello, world!\n","time":409},"msg":"测试创建成功"}
Swift:{"code":200,"data":{"condition":"ok","memory":9,"output":"hello world\n","time":2},"msg":"测试创建成功"}
*/
