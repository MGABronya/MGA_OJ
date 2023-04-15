package test

import (
	"fmt"
	"time"
)

type Commit struct {
	UserId     string
	Created_at time.Time
	Code       string
	Language   string
	ProblemNum string
}

type Problem struct {
	Input        []string
	Time_limit   uint
	Memory_limit uint
}

type Record struct {
	Spand      int64
	UserId     string
	Created_at int64
	Condition  string
	Output     string
}

var Records []Record

type Competition struct {
	Start   time.Time
	End     time.Time
	Problem map[string]Problem
	Commits []Commit
}

func (c Competition) Do(num int) {
	s := time.Now()
	for i := range c.Commits {
		go func(i int) {
			for j := num - 1; j >= 0; j-- {
				<-time.NewTimer(c.Commits[i].Created_at.Sub(c.Start.Add(time.Second * time.Duration(j*10)))).C
				for j := range c.Problem[c.Commits[i].ProblemNum].Input {
					now := time.Now()
					condition, output := JudgeRun(c.Commits[i].Language, c.Commits[i].Code, c.Problem[c.Commits[i].ProblemNum].Input[j], c.Problem[c.Commits[i].ProblemNum].Memory_limit, c.Problem[c.Commits[i].ProblemNum].Time_limit)
					spand := time.Now().UnixMilli() - now.UnixMilli()
					record := Record{
						spand,
						c.Commits[i].UserId,
						now.Sub(s).Milliseconds(),
						condition,
						output,
					}
					if record.Created_at < 0 {
						continue
					}
					fmt.Printf("create_at:%vs, spand:%vms, user:%s, condition:%s\n", record.Created_at/1000, record.Spand, record.UserId, record.Condition)
					Records = append(Records, record)
				}
			}

		}(i)
	}
	<-time.NewTimer(c.End.Sub(c.Start)).C
}
