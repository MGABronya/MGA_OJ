package selfInspection

import (
	"MGA_OJ/controller"
	"MGA_OJ/util"
	"MGA_OJ/vo"
	"log"
)

func TestInspection() int {
	log.Println("The profile is being detected...")
	if !util.FileExit("../config/application.yml") {
		log.Println("ERROR!!!" + "The configuration file does not exist! Please make sure that the configuration file application.yml exists in the config directory under the current directory, the contents of the configuration file can be found at " +
			"https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E9%83%A8%E7%BD%B2%E6%96%87%E6%A1%A3.md")
	}
	log.Println("language environment detection in ...")

	var test vo.TestRequest

	Err := 0
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
		Err++
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
		Err++
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
		Err++
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
		Err++
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
		Err++
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
		Err++
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
		Err++
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "JavaScript"
	test.Code = ` console.log("Hello World");
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
		Err++
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
		Err++
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
		Err++
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
		Err++
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Python"
	test.Code = `print('Hello World!')
	`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
		Err++
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
		Err++
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
		Err++
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
		Err++
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
		Err++
	} else {
		log.Println(test.Language + ":" + output)
	}

	test.Language = "Python"
	test.Code = `from time import time
from heapq import *
	
maxn = int(6e5)
	
inf = float('inf')
	
# ans
dis = [inf for _ in range(maxn)]
vis = [False for _ in range(maxn)]
	
# data structures
he = [-1 for _ in range(maxn)]
ne = [-1 for _ in range(maxn)]
w = [inf for _ in range(maxn)]
to = [-1 for _ in range(maxn)]
idx = 0
	
def add_edge(sou, des, z):
	global idx
	to[idx] = des
	w[idx] = z
	ne[idx] = he[sou]
	he[sou] = idx
	idx += 1
	
def dijkstra(sou):
	dis[sou] = 0
	q = []
	
	# {dis[y], y}
	heappush(q, (dis[sou], sou))
	
	while len(q):
		weight, now = heappop(q)
		if vis[now]:
			continue
	
		vis[now] = True
	
		i = he[now]
		while i != -1:
			y, z = to[i], w[i]
			if dis[y] > weight + z:
				dis[y] = weight + z
				heappush(q, (dis[y], y))
			i = ne[i]
	
	
# inputs
n, m = map(int, input().split())
	
# build graph
while m:
	m -= 1
	x, y, z = map(int, input().split())
	add_edge(x, y, z)
	
# excute
dijkstra(1)
	
if dis[n] != inf:
	print(dis[n])
else:
	print(-1)
	
print(int(time()))
	`
	test.Input = `3 3
	1 2 2
	2 3 1
	1 3 4`
	output, condition, _, _ = controller.Test(test)
	if condition != "ok" {
		log.Println("ERROR!!!" + condition + "," + test.Language + ":" + output)
		Err++
	} else {
		log.Println(test.Language + ":" + output)
	}

	log.Println("errors:", Err)

	log.Println("You can find containers with various programming language environments already deployed in the documentation " +
		"https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E9%83%A8%E7%BD%B2%E6%96%87%E6%A1%A3.md")
	return Err
}
