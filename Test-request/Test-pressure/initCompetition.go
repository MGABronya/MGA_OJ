package main

import (
	TQ "MGA_OJ/Test-request"
	"time"
)

func initCompetitions() {
	var competition TQ.Competition
	competition.Start = time.Date(2023, 4, 7, 19, 0, 0, 0, time.UTC)
	competition.End = time.Date(2023, 4, 7, 21, 0, 0, 0, time.UTC)
	competition.Problem = make(map[string]TQ.Problem)
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `abcdNTR\n`)
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
		6 6\n`)
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
		#!!!\n`)
		problem.Memory_limit = 262144
		problem.Time_limit = 1000
		competition.Problem["C"] = problem
	}
	{
		var problem TQ.Problem
		problem.Input = make([]string, 0)
		problem.Input = append(problem.Input, `3
		0 1 0\n`)
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
		29 13\n`)
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
		2 3 4 5 6 7\n`)
		problem.Memory_limit = 262144
		problem.Time_limit = 2000
		competition.Problem["F"] = problem
	}
	competitions = append(competitions, competition)
}
