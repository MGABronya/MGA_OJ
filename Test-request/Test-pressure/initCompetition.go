package main

import (
	TQ "MGA_OJ/Test-request"
	"time"
)

func initCompetitions1() {
	var competition TQ.Competition
	competition.Start = time.Date(2023, 4, 7, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2023, 4, 7, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `abcdNTR`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["A"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `4 2
		1 5
		3 3
		6 6
		6 6`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["B"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `4 4
		..!.
		.@.#
		!##!
		#!!!`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		0 1 0`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["D"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		4 4
		2 4
		29 13`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["E"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `7 1
		1 2
		1 3
		1 4
		2 5
		2 6
		2 7
		6
		2 3 4 5 6 7`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["F"] = problem
	}
	competitions = append(competitions, competition)
}

func initCompetitions2() {
	var competition TQ.Competition
	competition.Start = time.Date(2023, 3, 24, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2023, 3, 24, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `20 6`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["A"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		1 2 3`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["B"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3 120 50
		500 2 150 6 1
		1000 4 300 12 2
		1500 6 450 120 3`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `4 6 7
		1 2 3
		1 3 4
		1 4 6
		2 3 2
		2 4 1
		3 4 5`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["D"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `4
		1 1
		-1 1
		-1 -1
		1 -1`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["E"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `4
		1 1
		-1 1
		-1 -1
		1 -1`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["F"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `(2x+3y)^2%5`)
		problem.Input = append(problem.Input, `(2x+3y)^2%6`)
		problem.Input = append(problem.Input, `(c+2c)^3%4`)
		problem.Input = append(problem.Input, `(2a+2b)^2%2`)
		problem.Input = append(problem.Input, `(3c+6d)%4`)
		problem.Input = append(problem.Input, `(x+2y)^3%7`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["G"] = problem
	}
	competitions = append(competitions, competition)
}

func initCompetitions3() {
	var competition TQ.Competition
	competition.Start = time.Date(2023, 3, 10, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2023, 3, 10, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		123 456
		1000000000 1000000000
		1 23`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["A"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		2 20 5
		2 30 20
		25 35
		1 15 5
		18
		2 20 5
		2 30 20
		14 18
		1 15 5
		13
		1 20 10
		1 35 30
		35`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["B"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `4
		5
		aabab
		1
		a
		6
		aaabab
		11
		teeqtqrqwwev`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2
		10 100 200 300
		10 10 10 10`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["D"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2
		4 5
		0 6`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["E"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `5
		10 1
		1 2 3 2 1 2 3 2 1 2
		10 1
		1 2 3 2 1 2 3 2 1 2
		10 1
		1 2 3 2 1 2 3 2 1 2
		10 1
		1 2 3 2 2 1 3 1 2 2
		10 1
		2 2 3 3 1 2 2 1 1 2`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["F"] = problem
	}
	competitions = append(competitions, competition)
}

func initCompetitions4() {
	var competition TQ.Competition
	competition.Start = time.Date(2023, 2, 24, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2023, 2, 24, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `8`)
		problem.Input = append(problem.Input, `0`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["A"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `10 5
		1 4
		1 10
		1 4
		2 4
		2 4`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["B"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		4 2 4
		27 3 15
		25 13 22`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `10 5
		1 4
		1 10
		1 4
		2 4
		2 4`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["D"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2
		2
		7 8
		10 8
		2
		10 1000
		100 1`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["E"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		2 2
		1 2
		8 5`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["F"] = problem
	}
	competitions = append(competitions, competition)
}

func initCompetitions5() {
	var competition TQ.Competition
	competition.Start = time.Date(2023, 2, 10, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2023, 2, 10, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		3
		1 2 3
		2
		0 1
		2
		2 1`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["A"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2
		3
		156
		123
		3
		111
		110`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["B"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `1
		500000000`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `10 3
		3 6 9`)
		problem.Input = append(problem.Input, `3 3
		1 2 3`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["D"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2
		3 2
		3 3`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["E"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		5 7 6`)
		problem.Input = append(problem.Input, `10
		729 9 81 19683 1 2187 3 27 243 6561`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["F"] = problem
	}
	competitions = append(competitions, competition)
}

func initCompetitions6() {
	var competition TQ.Competition
	competition.Start = time.Date(2023, 1, 6, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2023, 1, 6, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `7 5 3`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["A"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `25 4
		niu1niun\|olo5ve,ni+um/ei
		love`)
		problem.Input = append(problem.Input, `5 4
		liike
		like`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["B"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `5 4
		2 1
		1 3
		2 5
		2 4`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2
		1 2
		3 3`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["D"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3 1`)
		problem.Input = append(problem.Input, `1 100`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["E"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `9
		3 3 2 3 4 5 3 6 3
		1 1 2 2 2 3 3 7`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["F"] = problem
	}
	competitions = append(competitions, competition)
}

func initCompetitions7() {
	var competition TQ.Competition
	competition.Start = time.Date(2022, 12, 30, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2022, 12, 30, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `7 3 3 2`)
		problem.Input = append(problem.Input, `7 4 3 2`)
		problem.Input = append(problem.Input, `7 5 3 2`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["A"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3 3
		...
		.#.
		...`)
		problem.Input = append(problem.Input, `3 3
		#..
		...
		...`)
		problem.Input = append(problem.Input, `1 1
		#`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["B"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `5 5`)
		problem.Input = append(problem.Input, `1 2`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `9
		3 5 3 4 6 1 5 1 6
		1 2 2 2 1 6 6 8`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["D"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `5 5 1
		0 3 2 4 0
		3 0
		4 0
		3 2
		2 0
		1 5`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["E"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `16 4
		1 3
		2 7
		3 11
		2 15`)
		problem.Input = append(problem.Input, `10 1
		1 9`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["F"] = problem
	}
	competitions = append(competitions, competition)
}

func initCompetitions8() {
	var competition TQ.Competition
	competition.Start = time.Date(2022, 12, 16, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2022, 12, 16, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `6
		1 1 4 5 1 4`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["A"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `5 6
		1 2 3 4 5`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["B"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2
		0 2
		1 3`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `5 4
		1 2 3 4`)
		problem.Input = append(problem.Input, `3 1
		1`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["D"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2`)
		problem.Input = append(problem.Input, `4`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["E"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `11
		91 88
		913 344
 393
		318 345
		882 425
		428 754
		134 425
		135 454
		21 90
		49 39
		140 13`)
		problem.Input = append(problem.Input, `10 1
		1 9`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["F"] = problem
	}
	competitions = append(competitions, competition)
}

func initCompetitions9() {
	var competition TQ.Competition
	competition.Start = time.Date(2022, 11, 25, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2022, 11, 25, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2
		6
		1 10 100 200 120 230
		5 230 200
		7
		5
		1 2 3 4 5
		10 5 2
		1`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["A"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `1
		1 5
		2
		2
		3`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["B"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		4 1 1
		9 7 10`)
		problem.Input = append(problem.Input, `4
		1 3 5 7
		2 4 8 22`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2
		9 3 5
		0 8 2 1 3
		1 1 1
		0`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["D"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `7 4
		1 2
		3 5
		4 3
		2 7`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["E"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `11
		91 88
		913 344
 393
		318 345
		882 425
		428 754
		134 425
		135 454
		21 90
		49 39
		140 13`)
		problem.Input = append(problem.Input, `acaacb`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["F"] = problem
	}
	competitions = append(competitions, competition)
}

func initCompetitions10() {
	var competition TQ.Competition
	competition.Start = time.Date(2022, 11, 18, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2022, 11, 18, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `5 3
		1 2 3 3 1`)
		problem.Input = append(problem.Input, `3 3
		1 3 1`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["A"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `1 4`)
		problem.Input = append(problem.Input, `1 5`)
		problem.Input = append(problem.Input, `1 6`)
		problem.Input = append(problem.Input, `0 0`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["B"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `5 3
		2 1
		...
		.M.
		...
		...
		.P.`)
		problem.Input = append(problem.Input, `5 3
		1 2
		**.
		*M.
		**.
		*..
		*P.`)
		problem.Input = append(problem.Input, `5 3
		2 1
		**.
		*M.
		**.
		*..
		*P.`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `2 1 5 5 2
		2 5`)
		problem.Input = append(problem.Input, `5 6 2 7 6
		2 3 6 7 7`)
		problem.Input = append(problem.Input, `7 9 5 2 7
		4 5 6 8 8 10 10`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["D"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		1 2 3`)
		problem.Input = append(problem.Input, `7
		2 4 4 3 1 1 2`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["E"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3 2
		1 1
		1 2
		2 2`)
		problem.Input = append(problem.Input, `4 6
		2 2 2 2 2 3
		3 1 1 1 2 2
		2 2 2 4 3 4`)
		problem.Input = append(problem.Input, `9 3
		9 2 6
		7 5 6
		1 5 3`)
		problem.Input = append(problem.Input, `acaacb`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["F"] = problem
	}
	competitions = append(competitions, competition)
}
