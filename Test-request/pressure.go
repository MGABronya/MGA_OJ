package test

import (
	"fmt"
	"time"
)

type User struct {
	Id      string
	Commits []Commit
}

type Commit struct {
	Created_at time.Time
	Code       string
	Language   string
	Problem    Problem
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
}

var Records []Record

func (user User) Do(start time.Time, end time.Time) {
	for i := range user.Commits {
		go func(i int) {
			<-time.NewTimer(user.Commits[i].Created_at.Sub(start)).C
			for j := range user.Commits[i].Problem.Input {
				now := time.Now()
				JudgeRun(user.Commits[i].Language, user.Commits[i].Code, user.Commits[i].Problem.Input[j], user.Commits[i].Problem.Memory_limit, user.Commits[i].Problem.Time_limit)
				spand := now.UnixMilli() - time.Now().UnixMilli()
				record := Record{
					spand,
					user.Id,
					now.Sub(start).Milliseconds(),
				}
				fmt.Println(record)
				Records = append(Records, record)
			}
		}(i)
	}
	<-time.NewTimer(end.Sub(start)).C
}
