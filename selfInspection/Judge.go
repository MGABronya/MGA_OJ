package selfInspection

import (
	"MGA_OJ/controller"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"log"
)

func JudgeInspection() {
	log.Println("The profile is being detected...")
	if !util.FileExit("../config/application.yml") {
		log.Println("ERROR!!!" + "The configuration file does not exist! Please make sure that the configuration file application.yml exists in the config directory under the current directory, the contents of the configuration file can be found at " +
			"https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E9%83%A8%E7%BD%B2%E6%96%87%E6%A1%A3.md")
	}
	log.Println("language environment detection in ...")
	var test vo.TestRequest

	test.Input = ""
	test.TimeLimit = 500
	test.MemoryLimit = 512

	// TODO 接下来进行hello world测试
	test.Language = "C"
	test.Code = `#include <stdio.h>
	int main()
	{
	   printf("Hello, World!");
	   return 0;
	}
	`

	output, condition, _, _ := controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "C#"
	test.Code = `using System;
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
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "C++"
	test.Code = `#include <iostream>
	using namespace std;
	 
	int main() 
	{
		cout << "Hello, World!";
		return 0;
	}
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "C++11"
	test.Code = `#include <iostream>
	using namespace std;
	 
	int main() 
	{
		cout << "Hello, World!";
		return 0;
	}
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Erlang"
	test.Code = `-module(main).
	-export([start/0]).
	start() -> io:fwrite("hi,hello world!\n").
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Go"
	test.Code = `package main

	import "fmt"
	
	func main() {
		fmt.Println("Hello world") 
	}
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Java"
	test.Code = `public class Main { 
		public static void main(String[] args){ 
			System.out.println("hello world!");
		}
	}
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "JavaScript"
	test.Code = ` console.log("Hello World");
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Kotlin"
	test.Code = `fun main(args: Array<String>) {
		println("Hello world")
	}
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Pascal"
	test.Code = `begin
    writeln('Hello, World!');
end.
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "PHP"
	test.Code = `<?php
	echo "Hello World \n";
	?>
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Python"
	test.Code = `print('Hello World!')
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Racket"
	test.Code = `#lang racket
	"Hello, world!"
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Ruby"
	test.Code = `
	puts "Hello World!"
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Rust"
	test.Code = `fn main() {
		println!("Hello, world!");
	}
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Scala"
	test.Code = `object main {
		def main(args: Array[String]): Unit = {
		  println("Hello, world!")
		}
	  }
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Swift"
	test.Code = `
	print("hello world")
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
	} else {
		log.Println(test.Language + ":" + output)
	}

	log.Println("You can find containers with various programming language environments already deployed in the documentation " +
		"https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E9%83%A8%E7%BD%B2%E6%96%87%E6%A1%A3.md")

}
